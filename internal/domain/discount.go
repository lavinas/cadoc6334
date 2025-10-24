package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
)

var (
	sqlDiscount string = "insert into cadoc_6334_desconto(Ano, Trimestre, Funcao, Bandeira, FormaCaptura, NumeroParcelas, CodigoSegmento, TaxaDescontoMedia, TaxaDescontoMinima, TaxaDescontoMaxima, DesvioPadraoTaxaDesconto, ValorTransacoes, QuantidadeTransacoes) values (%d, %d,'%s', %d, %d, %d, %d, %.2f, %.2f, %.2f, %.2f, %.2f, %d);"
)

type Discount struct {
	Year         int32  `fixed:"1,4"`
	Quarter      int32  `fixed:"5,5"`
	Function     string `fixed:"6,6"`
	Brand        int32  `fixed:"7,8"`
	Capture      int32  `fixed:"9,9"`
	Installments int32  `fixed:"10,11"`
	Segment      int32  `fixed:"12,14"`
	AvgFee       float32
	AvgFeeInt    int32 `fixed:"15,18"`
	MinFee       float32
	MinFeeInt    int32 `fixed:"19,22"`
	MaxFee       float32
	MaxFeeInt    int32 `fixed:"23,26"`
	StdDevFee    float32
	StdDevFeeInt int32 `fixed:"27,30"`
	Value        float32
	ValueInt     int32 `fixed:"31,45"`
	Qtty         int32 `fixed:"46,57"`
}

// GetInsert returns the SQL insert statement for the ranking
func (r *Discount) GetInsert() string {
	return fmt.Sprintf(sqlDiscount, r.Year, r.Quarter, r.Function, r.Brand, r.Capture, r.Installments, r.Segment, r.AvgFee, r.MinFee, r.MaxFee, r.StdDevFee, r.Value, r.Qtty)
}

// ParseLine parses a line of text into a Discount struct
func (r *Discount) Parse(line string) (*Discount, error) {
	err := fixedwidth.Unmarshal([]byte(line), r)
	if err != nil {
		return nil, err
	}
	// Convert ValueInt and DiscountInt back to float32
	r.Value = float32(float64(r.ValueInt) / float64(100))
	r.AvgFee = float32(float64(r.AvgFeeInt) / float64(100))
	r.MinFee = float32(float64(r.MinFeeInt) / float64(100))
	r.MaxFee = float32(float64(r.MaxFeeInt) / float64(100))
	r.StdDevFee = float32(float64(r.StdDevFeeInt) / float64(100))
	return r, nil
}

// String returns a string representation of the Discount struct
func (r *Discount) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, Function: %s, Brand: %d, Capture: %d, Installments: %d, Segment: %d, AvgFee: %.2f, MinFee: %.2f, MaxFee: %.2f, StdDevFee: %.2f, Value: %.2f, Qtty: %d",
		r.Year, r.Quarter, r.Function, r.Brand, r.Capture, r.Installments, r.Segment, r.AvgFee, r.MinFee, r.MaxFee, r.StdDevFee, r.Value, r.Qtty)
}

// GetDiscount returns the discount for a given year, quarter, value, and fee
func GetDiscount(year int32, quarter int32, value float32, fee float32) []*Discount {
	ret := []*Discount{}
	totValue := float32(0)
	for si, sv := range segValues {
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
						val := value * funcProp[fi] * brandProp[bi] * captProp[ci] * probValues[ii] * segProp[si]
						qty := int32(val / avgTicket)
						fee := fee + varFee[ii]
						minFee := fee + subFee
						maxFee := fee + sumFee
						stdDevFee := fee * stdDevProp
						discount := &Discount{
							Year:         year,
							Quarter:      quarter,
							Function:     fv,
							Brand:        bv,
							Capture:      cv,
							Installments: iv,
							Segment:      sv,
							AvgFee:       fee,
							MinFee:       minFee,
							MaxFee:       maxFee,
							StdDevFee:    stdDevFee,
							Value:        val,
							Qtty:         qty,
						}
						totValue += val
						ret = append(ret, discount)
					}
				}
			}
		}
	}
	ret[0].Value += value - totValue
	return ret
}

// LoadDiscount loads the discount data
func LoadDiscount() []*Discount {
	ret := []*Discount{}
	for _, y := range years {
		for _, q := range quarters {
			disc := GetDiscount(y, q, discTotalValue, discAvgFee)
			ret = append(ret, disc...)
		}
	}
	return ret
}

// PrintDiscount prints the discount data
func PrintDiscount() {
	value := float32(0)
	qty := int32(0)
	avgfee := float32(0)
	minfee := float32(0)
	maxfee := float32(0)
	stddev := float32(0)
	count := int32(0)
	disc := LoadDiscount()
	for _, d := range disc {
		fmt.Println(d.GetInsert())
		value += d.Value
		qty += d.Qtty
		avgfee += d.AvgFee
		minfee += d.MinFee
		maxfee += d.MaxFee
		stddev += d.StdDevFee
		count++
	}
	fmt.Println("--------------------------------------")
	fmt.Printf("-- total value: %.2f, expected %.2f\n", value, discTotalValue)
	fmt.Printf("-- total quantity: %d\n", qty)
	fmt.Printf("-- avg fee: %.2f, expected %.2f\n", avgfee/float32(count), discAvgFee)
	fmt.Printf("-- min fee: %.2f\n", minfee/float32(count))
	fmt.Printf("-- max fee: %.2f\n", maxfee/float32(count))
	fmt.Printf("-- stddev fee: %.2f\n", stddev/float32(count))
	fmt.Printf("-- total records: %d\n", count)
}

// ParseDiscountFile parses a discount file and returns a slice of Discount structs
func ParseDiscountFile(filePath string) ([]*Discount, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	discounts := []*Discount{}
	scanner := bufio.NewScanner(f)
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
	// read discounts
	var count int32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		disc := &Discount{}
		parsedDisc, err := disc.Parse(line)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, parsedDisc)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("DESCONTO", count); err != nil {
		return nil, err
	}
	return discounts, nil
}

// ReconciliateDiscount reconciliates the discount data from the file with the generated data
func ReconciliateDiscount(filePath string) {
	fmt.Println("Starting discount reconciliation...")
	genDiscounts := LoadDiscount()
	fileDiscounts, err := ParseDiscountFile(filePath)
	if err != nil {
		fmt.Println("Error parsing discount file:", err)
		return
	}
	if len(genDiscounts) != len(fileDiscounts) {
		fmt.Printf("Line count mismatch: generated %d, file %d\n", len(genDiscounts), len(fileDiscounts))
		return
	}
	map1 := make(map[string]*Discount)
	map2 := make(map[string]*Discount)
	for _, d := range genDiscounts {
		key := fmt.Sprintf("%d|%d|%s|%d|%d|%d|%d", d.Year, d.Quarter, d.Function, d.Brand, d.Capture, d.Installments, d.Segment)
		map1[key] = d
	}
	for _, d := range fileDiscounts {
		key := fmt.Sprintf("%d|%d|%s|%d|%d|%d|%d", d.Year, d.Quarter, d.Function, d.Brand, d.Capture, d.Installments, d.Segment)
		map2[key] = d
	}
	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok {
			fmt.Printf("Missing in file: %s\n", k)
			continue
		}
		if v1.String() != v2.String() {
			fmt.Printf("Mismatch for %s:\nGenerated: %s\nFile:      %s\n", k, v1.String(), v2.String())
		}
	}
	fmt.Println("Reconciliation complete.")
}
