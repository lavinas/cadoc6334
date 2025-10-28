package usecase

import (
	"fmt"
	"sort"
	"os"
	
	"golang.org/x/text/encoding/charmap"
	"github.com/lavinas/cadoc6334/internal/domain"
	"github.com/lavinas/cadoc6334/internal/port"

)

const (
	outPath = "./files/out"
)


// GenerateCase represents the use case for generating data
type GenerateCase struct {
	repo port.Repository
}


// NewGenerateCase creates a new instance of GenerateCase
func NewGenerateCase(repo port.Repository) *GenerateCase {
	return &GenerateCase{repo: repo}
}

// Execute executes the generate use case
func (ge *GenerateCase) ExecuteAll() {
	// Implement the logic for generating data here

	files := []string{
		"RANKING.TXT",
		"CONCCRED.TXT",
		"INFRESTA.TXT",
		"INFRTERM.TXT",
		"DESCONTO.TXT",
		"INTERCAM.TXT",
		"SEGMENTO.TXT",
		"LUCRCRED.TXT",
		"CONTATOS.TXT",
		"DATABASE.TXT",
	}
	reports := []port.Report{
		domain.NewRanking(),
		domain.NewConccred(),
		domain.NewInfresta(),
		domain.NewInfrterm(),
		domain.NewDiscount(),
		domain.NewIntercam(),
		domain.NewSegment(),
		domain.NewLucrCred(),
		domain.NewContact(),
		domain.NewDatabase(),
	}
	for i, file := range files {
		filename := fmt.Sprintf("%s/%s", outPath, file)
		if file == "DATABASE.TXT" {
			ge.GenerateDatabaseReport(filename)
			continue
		}
		ge.GenerateReport(reports[i], filename)
	}
}

// GenerateDatabaseReport generates the database report
func (ge *GenerateCase) GenerateDatabaseReport(filename string) {
	// Implement the logic for generating data here
	fmt.Printf("Generating data for %s\n", filename)
	defer fmt.Println("---------------------------------------------------------------------------------------------------------")
	// read db data
	db := domain.NewDatabase()
	// open file for writing
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer file.Close()
	// prepare encoder
	encoder := charmap.ISO8859_1.NewEncoder()
	// print line
	line := db.Format()
	out, err := encoder.Bytes([]byte(line))
	if err != nil {
		fmt.Printf("Error converting line to ISO-8859-1: %s\n", err)
		return
	}
	file.Write(out)
	file.Write([]byte("\n"))
}

// GenerateReport executes the generate use case for a specific report
func (ge *GenerateCase) GenerateReport(report port.Report, filename string) {
	// Implement the logic for generating data here
	fmt.Printf("Generating data for %s\n", filename)
	defer fmt.Println("---------------------------------------------------------------------------------------------------------")
	// read db data
	lines, err := report.GetDB(ge.repo)
	if err != nil {
		fmt.Printf("Error getting data from DB: %s\n", err)
		return
	}
	// sort lines
	order := make([]string, 0, len(lines))
	for k := range lines {
		order = append(order, k)
	}
	sort.Strings(order)
	// open file for writing
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer file.Close()
	// prepare encoder
	encoder := charmap.ISO8859_1.NewEncoder()
	// print header
	header := domain.NewHeader(report.GetName(), int64(len(lines)))
	headerLine := header.Format()
	out, err := encoder.Bytes([]byte(headerLine))
	if err != nil {
		fmt.Printf("Error converting header to ISO-8859-1: %s\n", err)
		return
	}
	file.Write(out)
	file.Write([]byte("\n"))
	// print lines
	for _, k := range order {
		r := lines[k].Format()
		// Convert to desired encoding
		out, err := encoder.Bytes([]byte(r))
		if err != nil {
			fmt.Printf("Error converting line to ISO-8859-1: %s\n", err)
			return
		}
		file.Write(out)
		file.Write([]byte("\n"))
	}
}
	
