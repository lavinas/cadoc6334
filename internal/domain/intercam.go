package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Intercam represents the intercam data model
type Intercam struct {
	Year         int32  `fixed:"1,4"`
	Quarter      int32  `fixed:"5,5"`
	Product      int32  `fixed:"6,7"`
	CardType     string `fixed:"8,8"`
	Function     string `fixed:"9,9"`
	Brand        int32  `fixed:"10,11"`
	Capture      int32  `fixed:"12,12"`
	Installments int32  `fixed:"13,14"`
	Segment      int32  `fixed:"15,17"`
	Fee          float32
	FeeInt       int32 `fixed:"18,21"`
	Value        float32
	ValueInt     int32 `fixed:"22,36"`
	Qtty         int32 `fixed:"37,48"`
}

// NewIntercam creates a new Intercam instance
func NewIntercam() *Intercam {
	return &Intercam{}
}

// TableName returns the table name for the Intercam struct
func (i *Intercam) TableName() string {
	return "cadoc_6334_intercam"
}

// GetKey generates a unique key for the Intercam record.
func (i *Intercam) GetKey() string {
	return fmt.Sprintf("%d|%d|%d|%s|%s|%d|%d|%d|%d", i.Year, i.Quarter, i.Product, i.CardType, i.Function, i.Brand, i.Capture, i.Installments, i.Segment)
}

// FindAll retrieves all Intercam records.
func (i *Intercam) FindAll(repo port.Repository) (map[string]port.Report, error) {
	var records []*Intercam
	err := repo.FindAll(&records)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]port.Report)
	for _, r := range records {
		ret[r.GetKey()] = r
	}
	return ret, nil
}

// Parse parses a line of text into an Intercam struct
func (i *Intercam) Parse(line string) (*Intercam, error) {
	err := fixedwidth.Unmarshal([]byte(line), i)
	if err != nil {
		return nil, err
	}
	// Convert ValueInt and FeeInt back to float32
	i.Value = float32(float64(i.ValueInt) / float64(100))
	i.Fee = float32(float64(i.FeeInt) / float64(100))
	return i, nil
}

// String returns a string representation of the Intercam struct
func (i *Intercam) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, Product: %d, CardType: %s, Function: %s, Brand: %d, Capture: %d, Installments: %d, Segment: %d, Fee: %.2f, Value: %.2f, Qtty: %d",
		i.Year, i.Quarter, i.Product, i.CardType, i.Function, i.Brand, i.Capture, i.Installments, i.Segment, i.Fee, i.Value, i.Qtty)
}

// GetIntercam returns the Intercam struct from the SQL insert statement
func (i *Intercam) GetIntercam(year int32, quarter int32, value float32, fee float32) []*Intercam {
	ret := []*Intercam{}
	totValue := float32(0)
	for si, sv := range segValues {
		for pi, pv := range prodValues {
			for ti, tv := range cardtypeValues {
				for fi, fv := range funcValues {
					for bi, bv := range brandValues {
						var insValues []int32
						var probValues []float32
						var varFee []float32
						for ci, cv := range captValues {
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
								val := value * funcProp[fi] * brandProp[bi] * captProp[ci] * probValues[ii] * segProp[si] * prodProp[pi] * cardtypeProp[ti]
								qty := int32(val / avgTicket)
								if qty == 0 {
									qty = 1
								}
								fee := fee + varFee[ii]
								inter := &Intercam{
									Year:         year,
									Quarter:      quarter,
									Product:      pv,
									CardType:     tv,
									Function:     fv,
									Brand:        bv,
									Capture:      cv,
									Installments: iv,
									Segment:      sv,
									Fee:          fee,
									Value:        val,
									Qtty:         qty,
								}
								ret = append(ret, inter)
								totValue += val

							}
						}
					}
				}
			}
		}
	}
	ret[0].Value += value - totValue
	return ret
}

// LoadIntercam loads the intercam data
func (i *Intercam) LoadIntercam() []*Intercam {
	ret := []*Intercam{}
	for _, y := range years {
		for _, q := range quarters {
			intercam := i.GetIntercam(y, q, intercamTotalValue, intercamAvgFee)
			ret = append(ret, intercam...)
		}
	}
	return ret
}

// ParseIntercamFile parses the intercam file and returns a slice of Intercam structs
func (i *Intercam) ParseIntercamFile(filename string) ([]*Intercam, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var intercams []*Intercam
	scanner := bufio.NewScanner(file)
	// read header
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	headerLine := scanner.Text()
	header := &RankingHeader{}
	_, err = header.Parse(headerLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	// read records
	var count int32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		intercam := &Intercam{}
		_, err := intercam.Parse(line)
		if err != nil {
			return nil, err
		}
		intercams = append(intercams, intercam)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("INTERCAM", count); err != nil {
		return nil, err
	}
	return intercams, nil
}

// GetLoaded retrieves and maps Intercam records from mounted data.
func (i *Intercam) GetLoaded() (map[string]port.Report, error) {
	loadedIntercam := i.LoadIntercam()
	if loadedIntercam == nil {
		return nil, fmt.Errorf("no intercam data loaded")
	}
	mapIntercam := make(map[string]port.Report)
	for _, ic := range loadedIntercam {
		mapIntercam[ic.GetKey()] = ic
	}
	return mapIntercam, nil
}


// GetParsedFile retrieves and maps Intercam records from a file.
func (i *Intercam) GetParsedFile(filename string) (map[string]port.Report, error) {
	fileIntercam, err := i.ParseIntercamFile(filename)
	if err != nil {
		return nil, err
	}
	mapIntercam := make(map[string]port.Report)
	for _, ic := range fileIntercam {
		mapIntercam[ic.GetKey()] = ic
	}
	return mapIntercam, nil
}

