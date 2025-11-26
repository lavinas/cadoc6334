package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
	"golang.org/x/text/encoding/charmap"
)

// Conccred represents the Conccred data model.
type Conccred struct {
	Year                       int64   `fixed:"1,4" gorm:"column:ano"`
	Quarter                    int64   `fixed:"5,5" gorm:"column:trimestre"`
	Brand                      int64   `fixed:"6,7" gorm:"column:bandeira"`
	Function                   string  `fixed:"8,8" gorm:"column:funcao"`
	CredentialedEstablishments int64   `fixed:"9,17" gorm:"column:quantidade_estabelecimentos_credenciados"`
	ActiveEstablishments       int64   `fixed:"18,26" gorm:"column:quantidade_estabelecimentos_ativos"`
	TransactionValue           float64 `gorm:"column:valor_transacoes"`
	TransactionValueInt        int64   `fixed:"27,41"`
	TransactionQuantity        int64   `fixed:"42,53" gorm:"column:quantidade_transacoes"`
}

// NewConccred creates a new Conccred instance.
func NewConccred() *Conccred {
	return &Conccred{}
}

// GetName gets name of the report
func (c *Conccred) GetName() string {
	return "CONCCRED"
}

// MarshalFixedWidth marshals the Conccred struct into a fixed-width format.
func (c *Conccred) Format() string {
	ret := ""
	ret += fmt.Sprintf("%04d", c.Year)
	ret += fmt.Sprintf("%01d", c.Quarter)
	ret += fmt.Sprintf("%02d", c.Brand)
	ret += fmt.Sprintf("%01s", c.Function)
	ret += fmt.Sprintf("%09d", c.CredentialedEstablishments)
	ret += fmt.Sprintf("%09d", c.ActiveEstablishments)
	// Convert TransactionValue to int representation
	c.TransactionValueInt = int64(c.TransactionValue*100 + 0.5)
	ret += fmt.Sprintf("%015d", c.TransactionValueInt)
	ret += fmt.Sprintf("%012d", c.TransactionQuantity)
	return ret
}

// Validate validates the Conccred header information.
func (c *Conccred) Validate() error {
	if c.Year <= 0 {
		return fmt.Errorf("invalid year in header")
	}
	if c.Quarter <= 0 {
		return fmt.Errorf("invalid quarter in header")
	}
	if c.Brand <= 0 {
		return fmt.Errorf("invalid brand in header")
	}
	if c.Function <= "" {
		return fmt.Errorf("invalid function in header")
	}
	if c.CredentialedEstablishments <= 0 {
		return fmt.Errorf("invalid number of credentialed establishments in header")
	}
	if c.ActiveEstablishments <= 0 {
		return fmt.Errorf("invalid number of active establishments in header")
	}
	if c.TransactionValue <= 0 {
		return fmt.Errorf("invalid transaction value in header")
	}
	if c.TransactionQuantity <= 0 {
		return fmt.Errorf("invalid transaction quantity in header")
	}
	return nil
}

// TableName returns the table name for the Conccred struct.
func (c *Conccred) TableName() string {
	return "cadoc_6334_conccred"
}

// GetKey generates a unique key for the Conccred record.
func (c *Conccred) GetKey() string {
	return fmt.Sprintf("%d-%d-%d-%s", c.Year, c.Quarter, c.Brand, c.Function)
}

// FindAll retrieves all Conccred records.
func (c *Conccred) GetDB(repo port.Repository) (map[string]port.Report, error) {
	var records []*Conccred
	err := repo.FindAll(&records, 0, 0, "")
	if err != nil {
		return nil, err
	}
	ret := make(map[string]port.Report)
	for _, r := range records {
		ret[r.GetKey()] = r
	}
	return ret, nil
}

// Parse parses a line of text into a Conccred struct.
func (c *Conccred) Parse(line string) (*Conccred, error) {
	err := fixedwidth.Unmarshal([]byte(line), c)
	if err != nil {
		return nil, err
	}
	// Convert TransactionValueInt back to float64
	c.TransactionValue = float64(float64(c.TransactionValueInt) / float64(100))
	return c, nil
}

// String returns a string representation of the Conccred struct.
func (c *Conccred) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, Brand: %d, Function: %s, CredentialedEstablishments: %d, ActiveEstablishments: %d, TransactionValue: %.2f, TransactionQuantity: %d",
		c.Year, c.Quarter, c.Brand, c.Function, c.CredentialedEstablishments, c.ActiveEstablishments, c.TransactionValue, c.TransactionQuantity)
}

// ParseConccredFile parses a file containing Conccred records.
func (c *Conccred) ParseConccredFile(filename string) ([]*Conccred, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := charmap.ISO8859_1.NewDecoder()
	decodedReader := decoder.Reader(file)
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
	// read records
	var records []*Conccred
	var count int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		var c Conccred
		record, err := c.Parse(line)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("CONCCRED", count); err != nil {
		return nil, err
	}
	return records, nil
}

// GetParsedFile retrieves and maps Conccred records from a file.
func (c *Conccred) GetParsedFile(filename string) (map[string]port.Report, error) {
	fileConccred, err := c.ParseConccredFile(filename)
	if err != nil {
		return nil, err
	}
	mapConccred := make(map[string]port.Report)
	for _, conc := range fileConccred {
		mapConccred[conc.GetKey()] = conc
	}
	return mapConccred, nil
}
