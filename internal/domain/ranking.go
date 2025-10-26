package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Ranking represents the ranking data model
type Ranking struct {
	Year         int64   `fixed:"1,4" gorm:"column:ano"`
	Quarter      int64   `fixed:"5,5" gorm:"column:trimestre"`
	ClientCode   string  `fixed:"6,13" gorm:"column:codigo_estabelecimento"`
	Function     string  `fixed:"14,14" gorm:"column:funcao"`
	Brand        int64   `fixed:"15,16" gorm:"column:bandeira"`
	Capture      int64   `fixed:"17,17" gorm:"column:forma_captura"`
	Installments int64   `fixed:"18,19" gorm:"column:numero_parcelas"`
	Segment      int64   `fixed:"20,22" gorm:"column:codigo_segmento"`
	Value        float64 `gorm:"column:valor_transacoes"`
	ValueInt     int64   `fixed:"23,37"`
	Qtty         int64   `fixed:"38,49" gorm:"column:quantidade_transacoes"`
	Discount     float64 `gorm:"column:taxa_desconto_media"`
	DiscountInt  int64   `fixed:"50,53"`
}

// NewRanking creates a new Ranking instance
func NewRanking() *Ranking {
	return &Ranking{}
}

// Validate validates the Ranking header information.
func (r *Ranking) Validate() error {
	if r.Year <= 0 {
		return fmt.Errorf("invalid year in header")
	}
	if r.Quarter <= 0 {
		return fmt.Errorf("invalid quarter in header")
	}
	if r.ClientCode == "" {
		return fmt.Errorf("invalid client code in header")
	}
	if r.Function == "" {
		return fmt.Errorf("invalid function in header")
	}
	if r.Brand <= 0 {
		return fmt.Errorf("invalid brand in header")
	}
	if r.Capture <= 0 {
		return fmt.Errorf("invalid capture in header")
	}
	if r.Installments <= 0 {
		return fmt.Errorf("invalid installments in header")
	}
	if r.Segment <= 0 {
		return fmt.Errorf("invalid segment in header")
	}
	if r.Value < 0 {
		return fmt.Errorf("invalid value in header")
	}
	if r.Qtty < 0 {
		return fmt.Errorf("invalid quantity in header")
	}
	if r.Discount < 0 {
		return fmt.Errorf("invalid discount in header")
	}
	return nil
}

// TableName returns the table name for the Ranking struct
func (r *Ranking) TableName() string {
	return "cadoc_6334_ranking"
}

// GetKey generates a unique key for the Ranking record.
func (r *Ranking) GetKey() string {
	return fmt.Sprintf("%d|%d|%s|%s|%d|%d|%d|%d", r.Year, r.Quarter, r.ClientCode, r.Function, r.Brand, r.Capture, r.Installments, r.Segment)
}

// FindAll retrieves all Ranking records.
func (r *Ranking) GetDB(repo port.Repository) (map[string]port.Report, error) {
	var records []*Ranking
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

// ParseLine parses a line of text into a Ranking struct
func (r *Ranking) Parse(line string) (*Ranking, error) {
	err := fixedwidth.Unmarshal([]byte(line), r)
	if err != nil {
		return nil, err
	}
	// Convert ValueInt and DiscountInt back to float64
	r.Value = float64(float64(r.ValueInt) / float64(100))
	r.Discount = float64(float64(r.DiscountInt) / float64(100))
	return r, nil
}

func (r *Ranking) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, ClientCode: %s, Function: %s, Brand: %d, Capture: %d, Installments: %d, Segment: %d, Value: %.2f, Qtty: %d, Discount: %.2f",
		r.Year, r.Quarter, r.ClientCode, r.Function, r.Brand, r.Capture, r.Installments, r.Segment, r.Value, r.Qtty, r.Discount)
}

// ParseRankingFile parses a file of rankings into a slice of Ranking structs
func (r *Ranking) ParseRankingFile(filename string) ([]*Ranking, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// read header
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	headerLine := scanner.Text()
	header, err := (&RankingHeader{}).Parse(headerLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	// read rankings
	rankings := []*Ranking{}
	var count int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		ranking, err := (&Ranking{}).Parse(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %w", err)
		}
		rankings = append(rankings, ranking)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("RANKING", count); err != nil {
		return nil, err
	}
	return rankings, nil
}

// GetParsedFile retrieves and maps Ranking records from a file.
func (r *Ranking) GetParsedFile(filename string) (map[string]port.Report, error) {
	fileRankings, err := r.ParseRankingFile(filename)
	if err != nil {
		return nil, err
	}
	mapRankings := make(map[string]port.Report)
	for _, i := range fileRankings {
		mapRankings[i.GetKey()] = i
	}
	return mapRankings, nil
}
