package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
	"golang.org/x/text/encoding/charmap"
)

// Infrterm represents the infrterm data model
type Infrterm struct {
	Year               int64  `fixed:"1,4" gorm:"column:ano"`
	Quarter            int64  `fixed:"5,5" gorm:"column:trimestre"`
	UF                 string `fixed:"6,7" gorm:"column:uf"`
	TotalPOSCount      int64  `fixed:"8,15" gorm:"column:quantidade_pos_totais"`
	SharedPOSCount     int64  `fixed:"16,23" gorm:"column:quantidade_pos_compartilhados"`
	ChipReaderPOSCount int64  `fixed:"24,31" gorm:"column:quantidade_pos_leitora_chip"`
	PDVCount           int64  `fixed:"32,39" gorm:"column:quantidade_pdv"`
}

// NewInfrterm creates a new Infrterm instance
func NewInfrterm() *Infrterm {
	return &Infrterm{}
}

// GetName gets name of the report
func (r *Infrterm) GetName() string {
	return "INFRTERM"
}

// Format marshals the Infrterm struct into a fixed-width format.
func (r *Infrterm) Format() string {
	ret := ""
	ret += fmt.Sprintf("%04d", r.Year)
	ret += fmt.Sprintf("%01d", r.Quarter)
	ret += fmt.Sprintf("%-2s", r.UF)
	ret += fmt.Sprintf("%08d", r.TotalPOSCount)
	ret += fmt.Sprintf("%08d", r.SharedPOSCount)
	ret += fmt.Sprintf("%08d", r.ChipReaderPOSCount)
	ret += fmt.Sprintf("%08d", r.PDVCount)
	return ret
}

// Validate validates the Infrterm header information.
func (r *Infrterm) Validate() error {
	if r.Year <= 0 {
		return fmt.Errorf("invalid year in header")
	}
	if r.Quarter <= 0 {
		return fmt.Errorf("invalid quarter in header")
	}
	if r.UF <= "" {
		return fmt.Errorf("invalid UF in header")
	}
	if r.TotalPOSCount < 0 {
		return fmt.Errorf("invalid total POS count in header")
	}
	if r.SharedPOSCount < 0 {
		return fmt.Errorf("invalid shared POS count in header")
	}
	if r.ChipReaderPOSCount < 0 {
		return fmt.Errorf("invalid chip reader POS count in header")
	}
	if r.PDVCount < 0 {
		return fmt.Errorf("invalid PDV count in header")
	}
	return nil
}

// TableName returns the table name for the Infrterm struct
func (r *Infrterm) TableName() string {
	return "cadoc_6334_infrterm"
}

// GetKey generates a unique key for the Infrterm record.
func (r *Infrterm) GetKey() string {
	return fmt.Sprintf("%d-%d-%s", r.Year, r.Quarter, r.UF)
}

// FindAll retrieves all Infrterm records.
func (r *Infrterm) GetDB(repo port.Repository) (map[string]port.Report, error) {
	var records []*Infrterm
	err := repo.FindAll(&records, 0, 0)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]port.Report)
	for _, rec := range records {
		ret[rec.GetKey()] = rec
	}
	return ret, nil
}

// Parse parses a fixed-width string into an Infrterm struct
func (r *Infrterm) Parse(line string) error {
	err := fixedwidth.Unmarshal([]byte(line), r)
	if err != nil {
		return err
	}
	return nil
}

// String returns a string representation of the Infrterm struct
func (r *Infrterm) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, UF: %s, TotalPOSCount: %d, SharedPOSCount: %d, ChipReaderPOSCount: %d, PDVCount: %d",
		r.Year, r.Quarter, r.UF, r.TotalPOSCount, r.SharedPOSCount, r.ChipReaderPOSCount, r.PDVCount)
}

// LoadInfrtermFile loads infrterm data from a fixed-width file
func (i *Infrterm) LoadInfrtermFile(filename string) ([]*Infrterm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var r []*Infrterm
	decoder := charmap.ISO8859_1.NewDecoder()
	decodedReader := decoder.Reader(file)
	scanner := bufio.NewScanner(decodedReader)
	// header line
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	headerLine := scanner.Text()
	header := &RankingHeader{}
	_, err = header.Parse(headerLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	// data lines
	var count int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		inf := &Infrterm{}
		err := inf.Parse(line)
		if err != nil {
			return nil, err
		}
		r = append(r, inf)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("INFRTERM", count); err != nil {
		return nil, err
	}
	return r, nil
}

// GetParsedFile retrieves and maps Infrterm records from a file.
func (r *Infrterm) GetParsedFile(filename string) (map[string]port.Report, error) {
	fileInfrterm, err := r.LoadInfrtermFile(filename)
	if err != nil {
		return nil, err
	}
	mapInfrterm := make(map[string]port.Report)
	for _, i := range fileInfrterm {
		mapInfrterm[i.GetKey()] = i
	}
	return mapInfrterm, nil
}
