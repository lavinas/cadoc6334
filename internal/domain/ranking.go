package domain

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Ranking represents the ranking data model
type Ranking struct {
	Year         int32  `fixed:"1,4"`
	Quarter      int32  `fixed:"5,5"`
	ClientCode   string `fixed:"6,13"`
	Function     string `fixed:"14,14"`
	Brand        int32  `fixed:"15,16"`
	Capture      int32  `fixed:"17,17"`
	Installments int32  `fixed:"18,19"`
	Segment      int32  `fixed:"20,22"`
	Value        float32
	ValueInt     int32 `fixed:"23,37"`
	Qtty         int32 `fixed:"38,49"`
	Discount     float32
	DiscountInt  int32 `fixed:"50,53"`
}

// NewRanking creates a new Ranking instance
func NewRanking() *Ranking {
	return &Ranking{}
}

// TableName returns the table name for the Ranking struct
func (r *Ranking) TableName() string {
	return "cadoc_6334_ranking"
}

// GetKey generates a unique key for the Ranking record.
func (r *Ranking) GetKey() string {
	return fmt.Sprintf("%d|%d|%s|%s|%d|%d|%d|%d", r.Year, r.Quarter, r.ClientCode, r.Function, r.Brand, r.Capture, r.Installments, r.Segment)
}

// FindAll retrieves all Ranking records.
func (r *Ranking) FindAll(repo port.Repository) (map[string]port.Report, error) {
	var records []*Ranking
	err := repo.FindAll(&records)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]port.Report)
	for _, rec := range records {
		ret[rec.GetKey()] = rec
	}
	return ret, nil
}

// ParseLine parses a line of text into a Ranking struct
func (r *Ranking) Parse(line string) (*Ranking, error) {
	err := fixedwidth.Unmarshal([]byte(line), r)
	if err != nil {
		return nil, err
	}
	// Convert ValueInt and DiscountInt back to float32
	r.Value = float32(float64(r.ValueInt) / float64(100))
	r.Discount = float32(float64(r.DiscountInt) / float64(100))
	return r, nil
}

func (r *Ranking) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, ClientCode: %s, Function: %s, Brand: %d, Capture: %d, Installments: %d, Segment: %d, Value: %.2f, Qtty: %d, Discount: %.2f",
		r.Year, r.Quarter, r.ClientCode, r.Function, r.Brand, r.Capture, r.Installments, r.Segment, r.Value, r.Qtty, r.Discount)
}

// ClientRanking returns the ranking of the client
func (r *Ranking) GetRanking(year int32, quarter int32, clientCode string, clientSegment int32, amount float32, fee float32) []*Ranking {
	ret := []*Ranking{}
	for fi, fv := range funcValues {
		for bi, bv := range brandValues {
			for ci, cv := range captValues {
				var insValues []int32
				var probValues []float32
				var varFee []float32
				if fv == "D" {
					insValues = instDebValues
					probValues = instDebProp
					varFee = instDebFee
				} else {
					insValues = instCredValues
					probValues = instCredProp
					varFee = instCredFee
				}
				for ii, iv := range insValues {
					// Calculate the value and quantity based on the amount and various proportions
					val := amount * funcProp[fi] * brandProp[bi] * probValues[ii] * captProp[ci]
					gval, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", val), 32)
					val = float32(gval)
					qty := int32(val / avgTicket)
					gdisc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", fee+varFee[ii]), 32)
					disc := float32(gdisc)
					// create and append ranking
					ranking := &Ranking{
						Year:         year,
						Quarter:      quarter,
						ClientCode:   clientCode,
						Function:     fv,
						Brand:        bv,
						Capture:      cv,
						Installments: iv,
						Segment:      clientSegment,
						Value:        val,
						ValueInt:     int32(val * 100), // store as integer in cents
						Qtty:         qty,
						Discount:     disc,
						DiscountInt:  int32(disc * 100), // store as integer in cents
					}
					ret = append(ret, ranking)
				}
			}
		}
	}
	return ret
}

// LoadRanking loads the ranking data from the file
func (r *Ranking) LoadRanking() ([]*Ranking, error) {
	rank := []*Ranking{}
	for _, y := range years {
		for _, q := range quarters {
			for i := int32(1); i <= maxClients; i++ {
				id := maxCode + i*maxCodeLeg
				val := maxVal + float32(i)*maxValLeg
				rank = append(rank,	r.GetRanking(y, q, fmt.Sprintf("%08d", id), segValues[i%int32(len(segValues))], val, maxFees[i%int32(len(maxFees))])...)
			}
			for i := int32(1); i <= minClients; i++ {
				id := minCode + i*minCodeLeg
				val := minVal + float32(i)*minValLeg
				rank = append(rank, r.GetRanking(y, q, fmt.Sprintf("%08d", id), segValues[i%int32(len(segValues))], val, minFee[i%int32(len(minFee))])...)
			}
		}
	}
	return rank, nil
}

// ParseRankingFile parses a file of rankings into a slice of Ranking structs
func (r *Ranking) ParseRankingFile(filename string) ([]*Ranking, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// read header
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	headerLine := scanner.Text()
	header, err := (&RankingHeader{}).Parse(headerLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	// read rankings
	rankings := []*Ranking{}
	var count int32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		ranking, err := (&Ranking{}).Parse(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %w", err)
		}
		rankings = append(rankings, ranking)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("RANKING", count); err != nil {
		return nil, err
	}
	return rankings, nil
}

// GetLoaded retrieves and maps Ranking records from mounted data.
func (r *Ranking) GetLoaded() (map[string]port.Report, error) {
	loadedRanking, err := r.LoadRanking()
	if err != nil {
		return nil, err
	}
	if loadedRanking == nil {
		return nil, fmt.Errorf("no ranking data loaded")
	}
	mapRanking := make(map[string]port.Report)
	for _, i := range loadedRanking {
		mapRanking[i.GetKey()] = i
	}
	return mapRanking, nil
}

// GetParsedFile retrieves and maps Ranking records from a file.
func (r *Ranking) GetParsedFile(filename string) (map[string]port.Report, error) {
	fileRankings, err := r.ParseRankingFile(filename)
	if err != nil {
		return nil, err
	}
	mapRankings := make(map[string]port.Report)
	for _, i := range fileRankings {
		mapRankings[i.GetKey()] = i
	}
	return mapRankings, nil
}
