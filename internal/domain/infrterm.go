package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
)

var (
	sqlInfrterm string = "insert into cadoc_6334_infrterm(Ano, Trimestre, Uf, QuantidadePOSTotal, QuantidadePOSCompartilhados, QuantidadePOSLeitoraChip, QuantidadePDV) values (%d, %d, '%s', %d, %d, %d, %d);"
)

type Infrterm struct {
	Year               int32  `fixed:"1,4"`
	Quarter            int32  `fixed:"5,5"`
	UF                 string `fixed:"6,7"`
	TotalPOSCount      int32  `fixed:"8,15"`
	SharedPOSCount     int32  `fixed:"16,23"`
	ChipReaderPOSCount int32  `fixed:"24,31"`
	PDVCount           int32  `fixed:"32,39"`
}

// GetInsert returns the SQL insert statement for the infrterm
func (r *Infrterm) GetInsert() string {
	return fmt.Sprintf(sqlInfrterm, r.Year, r.Quarter, r.UF, r.TotalPOSCount, r.SharedPOSCount, r.ChipReaderPOSCount, r.PDVCount)
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
func GetInfrterm(year int32, quarter int32, terms int32) []*Infrterm {
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
func LoadInfrTerm() []*Infrterm {
	var infrterms []*Infrterm
	for _, y := range years {
		for _, q := range quarters {
			infrterms = append(infrterms, GetInfrterm(y, q, infretermTerminals)...)
		}
	}
	return infrterms
}

// PrintInfrterm prints the infrterm details
func PrintInfrterm() {
	totTerms := int32(0)
	infrterm := LoadInfrTerm()
	for _, i := range infrterm {
		totTerms += i.SharedPOSCount + i.ChipReaderPOSCount + i.PDVCount
		fmt.Println(i.GetInsert())
	}
	fmt.Println("---------------------------------------")
	fmt.Printf("-- total terminals: %d, expected %d\n", totTerms, infretermTerminals)
}

// LoadInfrtermFile loads infrterm data from a fixed-width file
func LoadInfrtermFile(filename string) ([]*Infrterm, error) {
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

// ReconciliateInfrterm reconciliates the infrterm data from a file
func ReconciliateInfrterm(filename string) {
	fmt.Println("Starting infrterm reconciliation...")
	fileInfrterm, err := LoadInfrtermFile(filename)
	if err != nil {
		fmt.Printf("Error loading infrterm file: %v\n", err)
		return
	}
	generatedInfrterm := LoadInfrTerm()
	map1 := make(map[string]*Infrterm)
	map2 := make(map[string]*Infrterm)
	for _, i := range generatedInfrterm {
		key := fmt.Sprintf("%d-%d-%s", i.Year, i.Quarter, i.UF)
		map1[key] = i
	}
	for _, i := range fileInfrterm {
		key := fmt.Sprintf("%d-%d-%s", i.Year, i.Quarter, i.UF)
		map2[key] = i
	}
	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok {
			fmt.Printf("Key %s not found in file data\n", k)
			continue
		}
		if v1.String() != v2.String() {
			fmt.Printf("Mismatch for %s:\nGenerated: %s\nFile:      %s\n", k, v1.String(), v2.String())
		}
	}
	fmt.Println("Reconciliation complete.")
}
