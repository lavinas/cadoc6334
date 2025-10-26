package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Infresta represents the infresta data model
type Infresta struct {
	Year              int32  `fixed:"1,4" gorm:"column:ano"`
	Quarter           int32  `fixed:"5,5" gorm:"column:trimestre"`
	UF                string `fixed:"6,7" gorm:"column:uf"`
	TotalCli          int32  `fixed:"8,15" gorm:"column:quantidade_estabelecimentos_totais"`
	TotalCliManual    int32  `fixed:"16,23" gorm:"column:quantidade_estabelecimentos_captura_manual"`
	TotalCliEletronic int32  `fixed:"24,31" gorm:"column:quantidade_estabelecimentos_captura_eletronica"`
	TotalCliRemote    int32  `fixed:"32,39" gorm:"column:quantidade_estabelecimentos_captura_remota"`
}

// NewInfresta creates a new Infresta instance
func NewInfresta() *Infresta {
	return &Infresta{}
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
func (r *Infresta) FindAll(repo port.Repository) (map[string]port.Report, error) {
	var records []*Infresta
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

// GetInfresta returns the infresta data for a given year and quarter
func (r *Infresta) GetInfresta(year int32, quarter int32, totalCli int32) []*Infresta {
	ret := []*Infresta{}
	countCli := int32(0)
	for ui, uv := range ufValues {
		cli := int32(float32(totalCli) * ufProp[ui])
		inf := &Infresta{
			Year:              year,
			Quarter:           quarter,
			UF:                uv,
			TotalCli:          cli,
			TotalCliManual:    int32(float32(cli) * infrestaProp[0]),
			TotalCliEletronic: int32(float32(cli) * infrestaProp[1]),
			TotalCliRemote:    int32(float32(cli) * infrestaProp[2]),
		}
		inf.TotalCliEletronic = cli - inf.TotalCliManual - inf.TotalCliRemote
		ret = append(ret, inf)
		countCli += cli
	}
	ret[1].TotalCli += totalCli - countCli
	ret[1].TotalCliEletronic = ret[1].TotalCli - ret[1].TotalCliManual - ret[1].TotalCliRemote
	return ret
}

// LoadInfresta loads infresta data by year and quarter
func (r *Infresta) LoadInfresta() []*Infresta {
	ret := []*Infresta{}
	for _, y := range years {
		for _, q := range quarters {
			inf := r.GetInfresta(y, q, infrestaTotalEstablishments)
			ret = append(ret, inf...)
		}
	}
	return ret
}
	
// LoadInfrestaFile loads infresta data from a file
func (r *Infresta) LoadInfrestaFile(filename string) ([]*Infresta, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []*Infresta{}
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

// GetLoaded retrieves and maps Infresta records from mounted data.
func (r *Infresta) GetLoaded() (map[string]port.Report, error) {
	loadedInfresta := r.LoadInfresta()
	if loadedInfresta == nil {
		return nil, fmt.Errorf("no infresta data loaded")
	}
	mapInfresta := make(map[string]port.Report)
	for _, i := range loadedInfresta {
		mapInfresta[i.GetKey()] = i
	}
	return mapInfresta, nil
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
