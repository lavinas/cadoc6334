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
	Year         int64   `fixed:"1,4" gorm:"column:ano"`
	Quarter      int64   `fixed:"5,5" gorm:"column:trimestre"`
	Product      int64   `fixed:"6,7" gorm:"column:produto"`
	CardType     string  `fixed:"8,8" gorm:"column:modalidade_cartao"`
	Function     string  `fixed:"9,9" gorm:"column:funcao"`
	Brand        int64   `fixed:"10,11" gorm:"column:bandeira"`
	Capture      int64   `fixed:"12,12" gorm:"column:forma_captura"`
	Installments int64   `fixed:"13,14" gorm:"column:numero_parcelas"`
	Segment      int64   `fixed:"15,17" gorm:"column:codigo_segmento"`
	Fee          float64 `gorm:"column:tarifa_intercambio"`
	FeeInt       int64   `fixed:"18,21"`
	Value        float64 `gorm:"column:valor_transacoes"`
	ValueInt     int64   `fixed:"22,36"`
	Qtty         int64   `fixed:"37,48" gorm:"column:quantidade_transacoes"`
}

// NewIntercam creates a new Intercam instance
func NewIntercam() *Intercam {
	return &Intercam{}
}

// GetName gets name of the report
func (i *Intercam) GetName() string {
	return "INTERCAM"
}

// Format marshals the Intercam struct into a fixed-width format.
func (i *Intercam) Format() string {
	ret := ""
	ret += fmt.Sprintf("%04d", i.Year)
	ret += fmt.Sprintf("%01d", i.Quarter)
	ret += fmt.Sprintf("%02d", i.Product)
	ret += fmt.Sprintf("%01s", i.CardType)
	ret += fmt.Sprintf("%01s", i.Function)
	ret += fmt.Sprintf("%02d", i.Brand)
	ret += fmt.Sprintf("%01d", i.Capture)
	ret += fmt.Sprintf("%02d", i.Installments)
	ret += fmt.Sprintf("%03d", i.Segment)
	// Convert Fee to int representation
	i.FeeInt = int64(i.Fee * 100)
	ret += fmt.Sprintf("%04d", i.FeeInt)
	// Convert Value to int representation
	i.ValueInt = int64(i.Value * 100)
	ret += fmt.Sprintf("%015d", i.ValueInt)
	ret += fmt.Sprintf("%012d", i.Qtty)
	return ret
}

// Validate validates the Intercam header information.
func (i *Intercam) Validate() error {
	if i.Year <= 0 {
		return fmt.Errorf("invalid year in header")
	}
	if i.Quarter <= 0 {
		return fmt.Errorf("invalid quarter in header")
	}
	if i.Product <= 0 {
		return fmt.Errorf("invalid product in header")
	}
	if i.CardType <= "" {
		return fmt.Errorf("invalid card type in header")
	}
	if i.Function <= "" {
		return fmt.Errorf("invalid function in header")
	}
	if i.Brand <= 0 {
		return fmt.Errorf("invalid brand in header")
	}
	if i.Capture <= 0 {
		return fmt.Errorf("invalid capture in header")
	}
	if i.Installments <= 0 {
		return fmt.Errorf("invalid installments in header")
	}
	if i.Segment <= 0 {
		return fmt.Errorf("invalid segment in header")
	}
	if i.Fee < 0 {
		return fmt.Errorf("invalid fee in header")
	}
	if i.Value < 0 {
		return fmt.Errorf("invalid value in header")
	}
	if i.Qtty < 0 {
		return fmt.Errorf("invalid quantity in header")
	}
	return nil
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
func (i *Intercam) GetDB(repo port.Repository) (map[string]port.Report, error) {
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
	// Convert ValueInt and FeeInt back to float64
	i.Value = float64(float64(i.ValueInt) / float64(100))
	i.Fee = float64(float64(i.FeeInt) / float64(100))
	return i, nil
}

// String returns a string representation of the Intercam struct
func (i *Intercam) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, Product: %d, CardType: %s, Function: %s, Brand: %d, Capture: %d, Installments: %d, Segment: %d, Fee: %.2f, Value: %.2f, Qtty: %d",
		i.Year, i.Quarter, i.Product, i.CardType, i.Function, i.Brand, i.Capture, i.Installments, i.Segment, i.Fee, i.Value, i.Qtty)
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
	var count int64 = 0
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
