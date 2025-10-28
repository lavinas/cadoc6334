package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
	"golang.org/x/text/encoding/charmap"
)

// Discount SQL insert statement
type Discount struct {
	Year         int64   `fixed:"1,4" gorm:"column:ano"`
	Quarter      int64   `fixed:"5,5" gorm:"column:trimestre"`
	Function     string  `fixed:"6,6" gorm:"column:funcao"`
	Brand        int64   `fixed:"7,8" gorm:"column:bandeira"`
	Capture      int64   `fixed:"9,9" gorm:"column:forma_captura"`
	Installments int64   `fixed:"10,11" gorm:"column:numero_parcelas"`
	Segment      int64   `fixed:"12,14" gorm:"column:codigo_segmento"`
	AvgFee       float64 `gorm:"column:taxa_desconto_media"`
	AvgFeeInt    int64   `fixed:"15,18"`
	MinFee       float64 `gorm:"column:taxa_desconto_minima"`
	MinFeeInt    int64   `fixed:"19,22"`
	MaxFee       float64 `gorm:"column:taxa_desconto_maxima"`
	MaxFeeInt    int64   `fixed:"23,26"`
	StdDevFee    float64 `gorm:"column:desvio_padrao_taxa_desconto"`
	StdDevFeeInt int64   `fixed:"27,30"`
	Value        float64 `gorm:"column:valor_transacoes"`
	ValueInt     int64   `fixed:"31,45"`
	Qtty         int64   `fixed:"46,57" gorm:"column:quantidade_transacoes"`
}

// NewDiscount creates a new Discount instance
func NewDiscount() *Discount {
	return &Discount{}
}

// GetName gets name of the report
func (d *Discount) GetName() string {
	return "DESCONTO"
}

// Format marshals the Discount struct into a fixed-width format.
func (d *Discount) Format() string {
	ret := ""
	ret += fmt.Sprintf("%04d", d.Year)
	ret += fmt.Sprintf("%01d", d.Quarter)
	ret += fmt.Sprintf("%01s", d.Function)
	ret += fmt.Sprintf("%02d", d.Brand)
	ret += fmt.Sprintf("%01d", d.Capture)
	ret += fmt.Sprintf("%02d", d.Installments)
	ret += fmt.Sprintf("%03d", d.Segment)
	// Convert AvgFee, MinFee, MaxFee, StdDevFee to int representation
	d.AvgFeeInt = int64(d.AvgFee*100 + 0.5)
	d.MinFeeInt = int64(d.MinFee*100 + 0.5)
	d.MaxFeeInt = int64(d.MaxFee*100 + 0.5)
	d.StdDevFeeInt = int64(d.StdDevFee*100 + 0.5)
	ret += fmt.Sprintf("%04d", d.AvgFeeInt)
	ret += fmt.Sprintf("%04d", d.MinFeeInt)
	ret += fmt.Sprintf("%04d", d.MaxFeeInt)
	ret += fmt.Sprintf("%04d", d.StdDevFeeInt)
	// Convert Value to int representation
	d.ValueInt = int64(d.Value*100 + 0.5)
	ret += fmt.Sprintf("%015d", d.ValueInt)
	ret += fmt.Sprintf("%012d", d.Qtty)
	return ret
}

// Validate validates the Discount header information.
func (d *Discount) Validate() error {
	if d.Year <= 0 {
		return fmt.Errorf("invalid year in header")
	}
	if d.Quarter <= 0 {
		return fmt.Errorf("invalid quarter in header")
	}
	if d.Function <= "" {
		return fmt.Errorf("invalid function in header")
	}
	if d.Brand <= 0 {
		return fmt.Errorf("invalid brand in header")
	}
	if d.Capture <= 0 {
		return fmt.Errorf("invalid capture in header")
	}
	if d.Installments <= 0 {
		return fmt.Errorf("invalid installments in header")
	}
	if d.Segment <= 0 {
		return fmt.Errorf("invalid segment in header")
	}
	if d.AvgFee <= 0 {
		return fmt.Errorf("invalid average fee in header")
	}
	if d.MinFee < 0 {
		return fmt.Errorf("invalid minimum fee in header")
	}
	if d.MaxFee <= 0 {
		return fmt.Errorf("invalid maximum fee in header")
	}
	if d.StdDevFee < 0 {
		return fmt.Errorf("invalid standard deviation fee in header")
	}
	if d.Value <= 0 {
		return fmt.Errorf("invalid transaction value in header")
	}
	if d.Qtty <= 0 {
		return fmt.Errorf("invalid transaction quantity in header")
	}
	return nil
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
func (d *Discount) GetDB(repo port.Repository) (map[string]port.Report, error) {
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
	// Convert ValueInt and DiscountInt back to float64
	r.Value = float64(float64(r.ValueInt) / float64(100))
	r.AvgFee = float64(float64(r.AvgFeeInt) / float64(100))
	r.MinFee = float64(float64(r.MinFeeInt) / float64(100))
	r.MaxFee = float64(float64(r.MaxFeeInt) / float64(100))
	r.StdDevFee = float64(float64(r.StdDevFeeInt) / float64(100))
	return r, nil
}

// String returns a string representation of the Discount struct
func (r *Discount) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, Function: %s, Brand: %d, Capture: %d, Installments: %d, Segment: %d, AvgFee: %.2f, MinFee: %.2f, MaxFee: %.2f, StdDevFee: %.2f, Value: %.2f, Qtty: %d",
		r.Year, r.Quarter, r.Function, r.Brand, r.Capture, r.Installments, r.Segment, r.AvgFee, r.MinFee, r.MaxFee, r.StdDevFee, r.Value, r.Qtty)
}

// ParseDiscountFile parses a discount file and returns a slice of Discount structs
func (r *Discount) ParseDiscountFile(filePath string) ([]*Discount, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	discounts := []*Discount{}
	decoder := charmap.ISO8859_1.NewDecoder()
	decodedReader := decoder.Reader(f)
	scanner := bufio.NewScanner(decodedReader)
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
	var count int64 = 0
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
