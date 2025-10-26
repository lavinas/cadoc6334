package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Conccred represents the Conccred data model.
type Conccred struct {
	Year                       int32  `fixed:"1,4" gorm:"column:ano"`
	Quarter                    int32  `fixed:"5,5" gorm:"column:trimestre"`
	Brand                      int32  `fixed:"6,7" gorm:"column:bandeira"`
	Function                   string `fixed:"8,8" gorm:"column:funcao"`
	CredentialedEstablishments int32  `fixed:"9,17" gorm:"column:quantidade_estabelecimentos_credenciados"`
	ActiveEstablishments       int32  `fixed:"18,26" gorm:"column:quantidade_estabelecimentos_ativos"`
	TransactionValue           float32 `gorm:"column:valor_transacoes"`
	TransactionValueInt        int64 `fixed:"27,41"`
	TransactionQuantity        int32 `fixed:"42,53" gorm:"column:quantidade_transacoes"`
}

// TableName returns the table name for the Conccred struct.
func (c *Conccred) TableName() string {
	return "cadoc_6334_conccred"
}

// NewConccred creates a new Conccred instance.
func NewConccred() *Conccred {
	return &Conccred{}
}

// GetKey generates a unique key for the Conccred record.
func (c *Conccred) GetKey() string {
	return fmt.Sprintf("%d-%d-%d-%s", c.Year, c.Quarter, c.Brand, c.Function)
}

// FindAll retrieves all Conccred records.
func (c *Conccred) FindAll(repo port.Repository) (map[string]port.Report, error) {
	var records []*Conccred
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

// Parse parses a line of text into a Conccred struct.
func (c *Conccred) Parse(line string) (*Conccred, error) {
	err := fixedwidth.Unmarshal([]byte(line), c)
	if err != nil {
		return nil, err
	}
	// Convert TransactionValueInt back to float32
	c.TransactionValue = float32(float64(c.TransactionValueInt) / float64(100))
	return c, nil
}

// String returns a string representation of the Conccred struct.
func (c *Conccred) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, Brand: %d, Function: %s, CredentialedEstablishments: %d, ActiveEstablishments: %d, TransactionValue: %.2f, TransactionQuantity: %d",
		c.Year, c.Quarter, c.Brand, c.Function, c.CredentialedEstablishments, c.ActiveEstablishments, c.TransactionValue, c.TransactionQuantity)
}

// GetConccred generates a list of Conccred records based on the provided parameters.
func  (c *Conccred) GetConccred(year int32, quarter int32, creden int32, actives int32, value float32) []*Conccred {
	ret := []*Conccred{}
	totCreden := int32(0)
	totActives := int32(0)
	for bi, bv := range brandValues {
		for fi, fv := range funcValues {
			valuePortion := value * brandProp[bi] * funcProp[fi]
			qty := int32(valuePortion / avgTicket)
			credentialedEstablishments := int32(float32(creden) * brandProp[bi] * funcProp[fi])
			activeEstablishments := int32(float32(actives) * brandProp[bi] * funcProp[fi])
			ret = append(ret, &Conccred{
				Year:                       year,
				Quarter:                    quarter,
				Brand:                      bv,
				Function:                   fv,
				CredentialedEstablishments: credentialedEstablishments,
				ActiveEstablishments:       activeEstablishments,
				TransactionValue:           valuePortion,
				TransactionQuantity:        qty,
			})
			totCreden += credentialedEstablishments
			totActives += activeEstablishments
		}

	}
	ret[0].ActiveEstablishments += actives - totActives
	ret[0].CredentialedEstablishments += creden - totCreden
	return ret
}

// LoadConccred loads a Conccred record from a line of text.
func (c *Conccred)LoadConccred() []*Conccred {
	ret := []*Conccred{}
	for _, y := range years {
		for _, q := range quarters {
			conc := c.GetConccred(y, q, conccredTotalEstablishments, conccredActiveEstablishments, conccredTotalValue)
			ret = append(ret, conc...)
		}
	}
	return ret
}

// ParseConccredFile parses a file containing Conccred records.
func (c *Conccred) ParseConccredFile(filename string) ([]*Conccred, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
	var records []*Conccred
	var count int32 = 0
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

// GetLoaded retrieves and maps Conccred records from mounted data.
func (c *Conccred) GetLoaded() (map[string]port.Report, error) {
	loadedConccred := c.LoadConccred()
	if loadedConccred == nil {
		return nil, fmt.Errorf("no loaded Conccred data")
	}
	mapConccred := make(map[string]port.Report)
	for _, conc := range loadedConccred {
		mapConccred[conc.GetKey()] = conc
	}
	return mapConccred, nil
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
