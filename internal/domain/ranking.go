package domain

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/ianlopshire/go-fixedwidth"
)

var (
	// insert text
	sqlRanking string = "insert into cadoc_6334_ranking(Ano, Trimestre, CodigoEstabelecimento, Funcao, Bandeira, FormaCaptura, NumeroParcelas, CodigoSegmento, ValorTransacoes, QuantidadeTransacoes, TaxaDescontoMedia) values (%d, %d, '%s', '%s', %d, %d, %d, %d, %.2f, %d, %.2f);"
)

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

// GetInsert returns the SQL insert statement for the ranking
func (r *Ranking) GetInsert() string {
	return fmt.Sprintf(sqlRanking, r.Year, r.Quarter, r.ClientCode, r.Function, r.Brand, r.Capture, r.Installments, r.Segment, r.Value, r.Qtty, r.Discount)
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
func GetRanking(year int32, quarter int32, clientCode string, clientSegment int32, amount float32, fee float32) []*Ranking {
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
func LoadRanking() ([]*Ranking, error) {
	rank := []*Ranking{}
	for _, y := range years {
		for _, q := range quarters {
			for i := int32(1); i <= maxClients; i++ {
				id := maxCode + i*maxCodeLeg
				val := maxVal + float32(i)*maxValLeg
				rank = append(rank, GetRanking(y, q, fmt.Sprintf("%08d", id), segValues[i%int32(len(segValues))], val, maxFees[i%int32(len(maxFees))])...)
			}
			for i := int32(1); i <= minClients; i++ {
				id := minCode + i*minCodeLeg
				val := minVal + float32(i)*minValLeg
				rank = append(rank, GetRanking(y, q, fmt.Sprintf("%08d", id), segValues[i%int32(len(segValues))], val, minFee[i%int32(len(minFee))])...)
			}
		}
	}
	return rank, nil
}

// ranking generate
func PrintInsertRanking() {
	rank, error := LoadRanking()
	if error != nil {
		fmt.Println("Error loading ranking:", error)
		return
	}
	for _, r := range rank {
		fmt.Println(r.GetInsert())
	}
	fmt.Printf("Total rankings generated: %d\n", len(rank))
}

// ParseRankingFile parses a file of rankings into a slice of Ranking structs
func ParseRankingFile(filename string) (*RankingHeader, []*Ranking, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// read header
	if !scanner.Scan() {
		return nil, nil, fmt.Errorf("file is empty")
	}
	headerLine := scanner.Text()
	header, err := (&RankingHeader{}).Parse(headerLine)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing header: %w", err)
	}
	// read rankings
	rankings := []*Ranking{}
	var count int32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		ranking, err := (&Ranking{}).Parse(line)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing line: %w", err)
		}
		rankings = append(rankings, ranking)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	if err := header.Validate("RANKING", count); err != nil {
		return nil, nil, err
	}
	return header, rankings, nil
}

// Reconciliate ranking
func ReconciliateRanking(filename string) {
	fmt.Println("Starting ranking reconciliation...")
	rank1, err := LoadRanking()
	if err != nil {
		fmt.Println("Error loading ranking:", err)
		return
	}
	header2, rank2, err := ParseRankingFile(filename)
	if err != nil {
		fmt.Println("Error parsing ranking file:", err)
		return
	}
	if header2.Lines != int32(len(rank1)) {
		fmt.Printf("Line count mismatch: generated %d, file %d\n", len(rank1), header2.Lines)
	} else {
		fmt.Printf("Line count match: %d lines\n", len(rank1))
	}
	mapRank1 := make(map[string]*Ranking)
	for _, r := range rank1 {
		key := fmt.Sprintf("%d|%d|%s|%s|%d|%d|%d|%d", r.Year, r.Quarter, r.ClientCode, r.Function, r.Brand, r.Capture, r.Installments, r.Segment)
		mapRank1[key] = r
	}
	mapRank2 := make(map[string]*Ranking)
	for _, r := range rank2 {
		key := fmt.Sprintf("%d|%d|%s|%s|%d|%d|%d|%d", r.Year, r.Quarter, r.ClientCode, r.Function, r.Brand, r.Capture, r.Installments, r.Segment)
		mapRank2[key] = r
	}
	if len(mapRank1) != len(mapRank2) {
		fmt.Printf("Unique ranking count mismatch: generated %d, file %d\n", len(mapRank1), len(mapRank2))
	} else {
		fmt.Printf("Unique ranking count match: %d unique rankings\n", len(mapRank1))
	}
	// Compare rankings
	mismatchCount := 0
	i := 0
	for key, r1 := range mapRank1 {
		if r2, ok := mapRank2[key]; ok {
			if r1.String() != r2.String() {
				mismatchCount++
				fmt.Printf("Mismatch at line %d:\nGenerated: %+v\nFile:      %+v\n", i+2, r1, r2)
			}
		}
	}
	if mismatchCount == 0 {
		fmt.Println("All rankings match!")
	} else {
		fmt.Printf("Total mismatches: %d\n", mismatchCount)
	}
	fmt.Println("Reconciliation complete.")
}
