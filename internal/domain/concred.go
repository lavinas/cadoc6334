package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
)

var sqlConcred string = "insert into cadoc_6334_conccred(Ano, Trimestre, Bandeira, Funcao, QuantidadeEstabelecimentosCredenciados, QuantidadeEstabelecimentosAtivos, ValorTransacoes, QuantidadeTransacoes) values (%d, %d, %d, '%s', %d, %d, %.2f, %d);"

type Concred struct {
	Year                       int32  `fixed:"1,4"`
	Quarter                    int32  `fixed:"5,5"`
	Brand                      int32  `fixed:"6,7"`
	Function                   string `fixed:"8,8"`
	CredentialedEstablishments int32  `fixed:"9,17"`
	ActiveEstablishments       int32  `fixed:"18,26"`
	TransactionValue           float32
	TransactionValueInt        int64 `fixed:"27,41"`
	TransactionQuantity        int32 `fixed:"42,53"`
}

// GetInsert generates the SQL insert statement for the Concred struct.
func (c *Concred) GetInsert() string {
	return fmt.Sprintf(sqlConcred, c.Year, c.Quarter, c.Brand, c.Function, c.CredentialedEstablishments, c.ActiveEstablishments, c.TransactionValue, c.TransactionQuantity)
}

// Parse parses a line of text into a Concred struct.
func (c *Concred) Parse(line string) (*Concred, error) {
	err := fixedwidth.Unmarshal([]byte(line), c)
	if err != nil {
		return nil, err
	}
	// Convert TransactionValueInt back to float32
	c.TransactionValue = float32(float64(c.TransactionValueInt) / float64(100))
	return c, nil
}

// String returns a string representation of the Concred struct.
func (c *Concred) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, Brand: %d, Function: %s, CredentialedEstablishments: %d, ActiveEstablishments: %d, TransactionValue: %.2f, TransactionQuantity: %d",
		c.Year, c.Quarter, c.Brand, c.Function, c.CredentialedEstablishments, c.ActiveEstablishments, c.TransactionValue, c.TransactionQuantity)
}

// GetConcred generates a list of Concred records based on the provided parameters.
func GetConcred(year int32, quarter int32, creden int32, actives int32, value float32) []*Concred {
	ret := []*Concred{}
	totCreden := int32(0)
	totActives := int32(0)
	for bi, bv := range brandValues {
		for fi, fv := range funcValues {
			valuePortion := value * brandProp[bi] * funcProp[fi]
			qty := int32(valuePortion / avgTicket)
			credentialedEstablishments := int32(float32(creden) * brandProp[bi] * funcProp[fi])
			activeEstablishments := int32(float32(actives) * brandProp[bi] * funcProp[fi])
			ret = append(ret, &Concred{
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

// LoadConcred loads a Concred record from a line of text.
func LoadConcred() []*Concred {
	ret := []*Concred{}
	for _, y := range years {
		for _, q := range quarters {
			conc := GetConcred(y, q, concredTotalEstablishments, concredActiveEstablishments, concredTotalValue)
			ret = append(ret, conc...)
		}
	}
	return ret
}

// PrintConcred prints the Concred records.
func PrintConcred() {
	estabCount := int32(0)
	activeCount := int32(0)
	value := float32(0)
	qty := int32(0)
	count := int32(0)
	conc := LoadConcred()
	for _, c := range conc {
		fmt.Println(c.GetInsert())
		estabCount += c.CredentialedEstablishments
		activeCount += c.ActiveEstablishments
		value += c.TransactionValue
		qty += c.TransactionQuantity
		count++
	}
	fmt.Println("--------------------------------------")
	fmt.Printf("-- Total Credentialed Establishments: %d, expected: %d\n", estabCount, concredTotalEstablishments)
	fmt.Printf("-- Total Active Establishments: %d, expected: %d\n", activeCount, concredActiveEstablishments)
	fmt.Printf("-- Total Transaction Value: %.2f, expected: %.2f\n", value, concredTotalValue)
	fmt.Printf("-- Total Transaction Quantity: %d\n", qty)
	fmt.Printf("-- Avg ticket: %.2f expected %.2f\n", value/float32(qty), avgTicket)
	fmt.Printf("-- Total Records: %d\n", count)
}

// ParseConcredFile parses a file containing Concred records.
func ParseConcredFile(filename string) ([]*Concred, error) {
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
	var records []*Concred
	var count int32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		var c Concred
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

// ReconciliateConcred adjusts the first record to ensure totals match expected values.
func ReconciliateConcred(filename string) {
	fmt.Println("Starting concred reconciliation...")
	getConcred := LoadConcred()
	fileConcred, err := ParseConcredFile(filename)
	if err != nil {
		fmt.Printf("Error parsing concred file: %v\n", err)
		return
	}
	if len(getConcred) != len(fileConcred) {
		fmt.Printf("Length mismatch: generated %d, file %d\n", len(getConcred), len(fileConcred))
		return
	}
	map1 := make(map[string]*Concred)
	map2 := make(map[string]*Concred)
	for _, c := range getConcred {
		key := fmt.Sprintf("%d-%d-%d-%s", c.Year, c.Quarter, c.Brand, c.Function)
		map1[key] = c
	}
	for _, c := range fileConcred {
		key := fmt.Sprintf("%d-%d-%d-%s", c.Year, c.Quarter, c.Brand, c.Function)
		map2[key] = c
	}
	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok {
			fmt.Printf("Record missing in file: %s\n", k)
			continue
		}
		if v1.String() != v2.String() {
			fmt.Printf("Record mismatch for %s:\nGenerated: %s\nFile:      %s\n", k, v1.String(), v2.String())
		}
	}
	fmt.Println("Reconciliation complete.")

}
