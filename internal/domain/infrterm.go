package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Infrterm represents the infrterm data model
type Infrterm struct {
	Year               int32  `fixed:"1,4" gorm:"column:ano"`
	Quarter            int32  `fixed:"5,5" gorm:"column:trimestre"`
	UF                 string `fixed:"6,7" gorm:"column:uf"`
	TotalPOSCount      int32  `fixed:"8,15" gorm:"column:quantidade_pos_totais"`
	SharedPOSCount     int32  `fixed:"16,23" gorm:"column:quantidade_pos_compartilhados"`
	ChipReaderPOSCount int32  `fixed:"24,31" gorm:"column:quantidade_pos_leitora_chip"`
	PDVCount           int32  `fixed:"32,39" gorm:"column:quantidade_pdv"`
}

// NewInfrterm creates a new Infrterm instance
func NewInfrterm() *Infrterm {
	return &Infrterm{}
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
func (r *Infrterm) FindAll(repo port.Repository) (map[string]port.Report, error) {
	var records []*Infrterm
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

// GetInfrterm returns the infrterm for a given year, quarter, state, and counts
func (r *Infrterm) GetInfrterm(year int32, quarter int32, terms int32) []*Infrterm {
	totTerms := int32(0)
	ret := []*Infrterm{}
	for ui, uv := range ufValues {
		term := int32(float32(terms) * ufProp[ui])
		inf := &Infrterm{
			Year:               year,
			Quarter:            quarter,
			UF:                 uv,
			TotalPOSCount:      term,
			SharedPOSCount:     int32(float32(term) * infretermProp[0]),
			ChipReaderPOSCount: int32(float32(term) * infretermProp[1]),
			PDVCount:           int32(float32(term) * infretermProp[2]),
		}
		totTerms += term
		inf.ChipReaderPOSCount = term - inf.SharedPOSCount - inf.PDVCount
		ret = append(ret, inf)
	}
	dif := terms - totTerms
	ret[0].SharedPOSCount += dif
	return ret
}

// LoadInfrterm loads infrterm by year and quarter
func (r *Infrterm) LoadInfrTerm() []*Infrterm {
	var infrterms []*Infrterm
	for _, y := range years {
		for _, q := range quarters {
			infrterms = append(infrterms, r.GetInfrterm(y, q, infretermTerminals)...)
		}
	}
	return infrterms
}


// LoadInfrtermFile loads infrterm data from a fixed-width file
func (i *Infrterm) LoadInfrtermFile(filename string) ([]*Infrterm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var r []*Infrterm
	scanner := bufio.NewScanner(file)
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
	var count int32 = 0
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

// GetLoaded retrieves and maps Infrterm records from mounted data.
func (r *Infrterm) GetLoaded() (map[string]port.Report, error) {
	loadedInfrterm := r.LoadInfrTerm()
	if loadedInfrterm == nil {
		return nil, fmt.Errorf("no infrterm data loaded")
	}
	mapInfrterm := make(map[string]port.Report)
	for _, i := range loadedInfrterm {
		mapInfrterm[i.GetKey()] = i
	}
	return mapInfrterm, nil
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

