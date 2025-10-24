package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
)

var (
	infestaSQL = "insert into cadoc_6334_infresta(Ano, Trimestre, Uf, QuantidadeEstabelecimentosTotais, QuantidadeEstabelecimentosCapturaManual, QuantidadeEstabelecimentosCapturaEletronica, QuantidadeEstabelecimentosCapturaRemota) values (%d, %d, '%s', %d, %d, %d, %d);"
)

// Infresta represents the infresta data model
type Infresta struct {
	Year              int32  `fixed:"1,4"`
	Quarter           int32  `fixed:"5,5"`
	UF                string `fixed:"6,7"`
	TotalCli          int32  `fixed:"8,15"`
	TotalCliManual    int32  `fixed:"16,23"`
	TotalCliEletronic int32  `fixed:"24,31"`
	TotalCliRemote    int32  `fixed:"32,39"`
}

// GetInsert returns the SQL insert statement for the ranking
func (r *Infresta) GetInsert() string {
	return fmt.Sprintf(infestaSQL, r.Year, r.Quarter, r.UF, r.TotalCli, r.TotalCliManual, r.TotalCliEletronic, r.TotalCliRemote)
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
func GetInfresta(year int32, quarter int32, totalCli int32) []*Infresta {
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
func LoadInfresta() []*Infresta {
	ret := []*Infresta{}
	for _, y := range years {
		for _, q := range quarters {
			inf := GetInfresta(y, q, infrestaTotalEstablishments)
			ret = append(ret, inf...)
		}
	}
	return ret
}

// PrintInfresta prints the infresta data for a given year and quarter
func PrintInfresta() {
	totalCli := int32(0)
	totalCli2 := int32(0)
	infresta := LoadInfresta()
	for _, i := range infresta {
		fmt.Println(i.String())
		totalCli += i.TotalCli
		totalCli2 += i.TotalCliManual + i.TotalCliEletronic + i.TotalCliRemote
	}
	fmt.Println("----------------------------------------------------")
	fmt.Printf("Total Clients: %d - Total Clients2: %d, expected %d\n", totalCli, totalCli2, infrestaTotalEstablishments)
}

// LoadInfrestaFile loads infresta data from a file
func LoadInfrestaFile(filename string) ([]*Infresta, error) {
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

// ReconciliateInfresta reconciliates the infresta data from a file
func ReconciliateInfresta(filename string) {
	fmt.Println("Starting infresta reconciliation...")
	fileInfresta, err := LoadInfrestaFile(filename)
	if err != nil {
		fmt.Printf("Error loading infresta file: %v\n", err)
		return
	}
	genInfresta := LoadInfresta()
	if len(fileInfresta) != len(genInfresta) {
		fmt.Printf("Error: file infresta length %d does not match generated infresta length %d\n", len(fileInfresta), len(genInfresta))
		return
	}
	map1 := make(map[string]*Infresta)
	map2 := make(map[string]*Infresta)
	for _, i := range genInfresta {
		key := fmt.Sprintf("%d-%d-%s", i.Year, i.Quarter, i.UF)
		map1[key] = i
	}
	for _, i := range fileInfresta {
		key := fmt.Sprintf("%d-%d-%s", i.Year, i.Quarter, i.UF)
		map2[key] = i
	}
	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok {
			fmt.Printf("Missing in file: %s\n", k)
			continue
		}
		if v1.String() != v2.String() {
			fmt.Printf("Mismatch for %s:\nGenerated: %s\nFile:      %s\n", k, v1.String(), v2.String())
		}
	}
	fmt.Println("Reconciliation complete.")
}
