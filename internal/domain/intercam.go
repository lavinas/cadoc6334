package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
)

var (
	intercamSQL = "INSERT INTO intercam (Ano, Trimestre, Produto, ModalidadeCartao, Funcao, Bandeira, FormaCaptura, NumeroParcelas, CodigoSegmento, TarifaIntercambio, ValorTransacoes, QuantidadeTransacoes) VALUES (%d, %d, %d, '%s', '%s', %d, %d, %d, %d, %.2f, %.2f, %d);"
)

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

// GetInsert returns the SQL insert statement for the Intercam struct
func (i *Intercam) GetInsert() string {
	return fmt.Sprintf(intercamSQL, i.Year, i.Quarter, i.Product, i.CardType, i.Function, i.Brand, i.Capture, i.Installments, i.Segment, i.Fee, i.Value, i.Qtty)
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
func GetIntercam(year int32, quarter int32, value float32, fee float32) []*Intercam {
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
func LoadIntercam() []*Intercam {
	ret := []*Intercam{}
	for _, y := range years {
		for _, q := range quarters {
			intercam := GetIntercam(y, q, intercamTotalValue, intercamAvgFee)
			ret = append(ret, intercam...)
		}
	}
	return ret
}

// PrintIntercam
func PrintIntercam() {
	tot := float32(0)
	intercam := LoadIntercam()
	for _, i := range intercam {
		fmt.Println(i.GetInsert())
		tot += i.Value
	}
	fmt.Println("---------------------------------------")
	fmt.Printf("Value: %.2f, expected: %.2f\n", tot, intercamTotalValue)
}

// ParseIntercamFile parses the intercam file and returns a slice of Intercam structs
func ParseIntercamFile(filename string) ([]*Intercam, error) {
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

// ReconcileIntercam reconciles the intercam data
func ReconciliateIntercam(filePath string) {
	fmt.Println("Starting intercam reconciliation...")
	genIntercan := LoadIntercam()
	fileIntercan, err := ParseIntercamFile(filePath)
	if err != nil {
		fmt.Println("Error parsing discount file:", err)
		return
	}
	if len(genIntercan) != len(fileIntercan) {
		fmt.Printf("Record count mismatch: generated %d, file %d\n", len(genIntercan), len(fileIntercan))
	}
	map1 := make(map[string]*Intercam)
	map2 := make(map[string]*Intercam)
	for _, i := range genIntercan {
		key := fmt.Sprintf("%d|%d|%d|%s|%s|%d|%d|%d|%d", i.Year, i.Quarter, i.Product, i.CardType, i.Function, i.Brand, i.Capture, i.Installments, i.Segment)
		map1[key] = i
	}
	for _, i := range fileIntercan {
		key := fmt.Sprintf("%d|%d|%d|%s|%s|%d|%d|%d|%d", i.Year, i.Quarter, i.Product, i.CardType, i.Function, i.Brand, i.Capture, i.Installments, i.Segment)
		map2[key] = i
	}
	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok {
			fmt.Printf("Missing in file: %s\n", k)
			continue
		}
		if v1.String() != v2.String() {
			fmt.Printf("Record mismatch for key %s:\nGenerated: %s\nFile:      %s\n", k, v1.String(), v2.String())
		}
	}
	fmt.Println("Intercam data reconciled successfully")
}
