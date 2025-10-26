package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Discount SQL insert statement
type Discount struct {
	Year         int32  `fixed:"1,4" gorm:"column:ano"`
	Quarter      int32  `fixed:"5,5" gorm:"column:trimestre"`
	Function     string `fixed:"6,6" gorm:"column:funcao"`
	Brand        int32  `fixed:"7,8" gorm:"column:bandeira"`
	Capture      int32  `fixed:"9,9" gorm:"column:forma_captura"`
	Installments int32  `fixed:"10,11" gorm:"column:numero_parcelas"`
	Segment      int32  `fixed:"12,14" gorm:"column:codigo_segmento"`
	AvgFee       float32 `gorm:"column:taxa_desconto_media"`
	AvgFeeInt    int32 `fixed:"15,18"`
	MinFee       float32 `gorm:"column:taxa_desconto_minima"`
	MinFeeInt    int32 `fixed:"19,22"`
	MaxFee       float32 `gorm:"column:taxa_desconto_maxima"`
	MaxFeeInt    int32 `fixed:"23,26"`
	StdDevFee    float32 `gorm:"column:desvio_padrao_taxa_desconto"`
	StdDevFeeInt int32 `fixed:"27,30"`
	Value        float32 `gorm:"column:valor_transacoes"`
	ValueInt     int32 `fixed:"31,45"`
	Qtty         int32 `fixed:"46,57" gorm:"column:quantidade_transacoes"`
}

// NewDiscount creates a new Discount instance
func NewDiscount() *Discount {
	return &Discount{}
}

// TableName returns the table name for the Discount struct
func (d *Discount) TableName() string {
	return "cadoc_6334_desconto"
}

// GetKey generates a unique key for the Discount record.
func (d *Discount) GetKey() string {
	return fmt.Sprintf("%d|%d|%s|%d|%d|%d|%d", d.Year, d.Quarter, d.Function, d.Brand, d.Capture, d.Installments, d.Segment)
}

// FindAll retrieves all Discount records.
func (d *Discount) FindAll(repo port.Repository) (map[string]port.Report, error) {
	var records []*Discount
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
func (r *Discount) GetDiscount(year int32, quarter int32, value float32, fee float32) []*Discount {
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
func (r *Discount) LoadDiscount() []*Discount {
	ret := []*Discount{}
	for _, y := range years {
		for _, q := range quarters {
			disc := r.GetDiscount(y, q, discTotalValue, discAvgFee)
			ret = append(ret, disc...)
		}
	}
	return ret
}

// ParseDiscountFile parses a discount file and returns a slice of Discount structs
func (r *Discount) ParseDiscountFile(filePath string) ([]*Discount, error) {
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

// GetLoaded retrieves and maps Discount records from mounted data.
func (r *Discount) GetLoaded() (map[string]port.Report, error) {
	loadedDiscount := r.LoadDiscount()
	if loadedDiscount == nil {
		return nil, fmt.Errorf("no discount data loaded")
	}
	discountMap := make(map[string]port.Report)
	for _, d := range loadedDiscount {
		discountMap[d.GetKey()] = d
	}
	return discountMap, nil
}

// GetParsedFile retrieves and maps Discount records from a file.
func (r *Discount) GetParsedFile(filePath string) (map[string]port.Report, error) {
	fileDiscounts, err := r.ParseDiscountFile(filePath)
	if err != nil {
		return nil, err
	}
	discountMap := make(map[string]port.Report)
	for _, d := range fileDiscounts {
		discountMap[d.GetKey()] = d
	}
	return discountMap, nil
}
