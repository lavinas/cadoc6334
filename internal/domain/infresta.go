package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
	"golang.org/x/text/encoding/charmap"
)

// Infresta represents the infresta data model
type Infresta struct {
	Year              int64  `fixed:"1,4" gorm:"column:ano"`
	Quarter           int64  `fixed:"5,5" gorm:"column:trimestre"`
	UF                string `fixed:"6,7" gorm:"column:uf"`
	TotalCli          int64  `fixed:"8,15" gorm:"column:quantidade_estabelecimentos_totais"`
	TotalCliManual    int64  `fixed:"16,23" gorm:"column:quantidade_estabelecimentos_captura_manual"`
	TotalCliEletronic int64  `fixed:"24,31" gorm:"column:quantidade_estabelecimentos_captura_eletronica"`
	TotalCliRemote    int64  `fixed:"32,39" gorm:"column:quantidade_estabelecimentos_captura_remota"`
}

// NewInfresta creates a new Infresta instance
func NewInfresta() *Infresta {
	return &Infresta{}
}

// GetName gets name of the report
func (r *Infresta) GetName() string {
	return "INFRESTA"
}

// Format marshals the Infresta struct into a fixed-width format.
func (r *Infresta) Format() string {
	ret := ""
	ret += fmt.Sprintf("%04d", r.Year)
	ret += fmt.Sprintf("%01d", r.Quarter)
	ret += fmt.Sprintf("%-2s", r.UF)
	ret += fmt.Sprintf("%08d", r.TotalCli)
	ret += fmt.Sprintf("%08d", r.TotalCliManual)
	ret += fmt.Sprintf("%08d", r.TotalCliEletronic)
	ret += fmt.Sprintf("%08d", r.TotalCliRemote)
	return ret
}

// Validate validates the Infresta header information.
func (r *Infresta) Validate() error {
	if r.Year <= 0 {
		return fmt.Errorf("invalid year in header")
	}
	if r.Quarter <= 0 {
		return fmt.Errorf("invalid quarter in header")
	}
	if r.UF <= "" {
		return fmt.Errorf("invalid UF in header")
	}
	if r.TotalCli <= 0 {
		return fmt.Errorf("invalid total clients in header")
	}
	if r.TotalCliManual < 0 {
		return fmt.Errorf("invalid total manual clients in header")
	}
	if r.TotalCliEletronic < 0 {
		return fmt.Errorf("invalid total electronic clients in header")
	}
	if r.TotalCliRemote < 0 {
		return fmt.Errorf("invalid total remote clients in header")
	}
	return nil
}

// TableName returns the table name for the Infresta struct
func (r *Infresta) TableName() string {
	return "cadoc_6334_infresta"
}

// GetKey generates a unique key for the Infresta record.
func (r *Infresta) GetKey() string {
	return fmt.Sprintf("%d-%d-%s", r.Year, r.Quarter, r.UF)
}

// FindAll retrieves all Infresta records.
func (r *Infresta) GetDB(repo port.Repository) (map[string]port.Report, error) {
	var records []*Infresta
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

// Parse parses a line of text into an Infresta struct
func (r *Infresta) Parse(line string) (*Infresta, error) {
	err := fixedwidth.Unmarshal([]byte(line), r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// String returns a string representation of the Infresta struct
func (r *Infresta) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, UF: %s, TotalCli: %d, TotalCliManual: %d, TotalCliEletronic: %d, TotalCliRemote: %d",
		r.Year, r.Quarter, r.UF, r.TotalCli, r.TotalCliManual, r.TotalCliEletronic, r.TotalCliRemote)
}

// LoadInfrestaFile loads infresta data from a file
func (r *Infresta) LoadInfrestaFile(filename string) ([]*Infresta, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []*Infresta{}
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
		inf := &Infresta{}
		parsedInf, err := inf.Parse(line)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parsedInf)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("INFRESTA", count); err != nil {
		return nil, err
	}
	return ret, nil
}

// GetParsedFile retrieves and maps Infresta records from a file.
func (r *Infresta) GetParsedFile(filename string) (map[string]port.Report, error) {
	fileInfresta, err := r.LoadInfrestaFile(filename)
	if err != nil {
		return nil, err
	}
	mapInfresta := make(map[string]port.Report)
	for _, i := range fileInfresta {
		mapInfresta[i.GetKey()] = i
	}
	return mapInfresta, nil
}
