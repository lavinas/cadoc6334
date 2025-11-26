package usecase

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/lavinas/cadoc6334/internal/domain"
	"github.com/lavinas/cadoc6334/internal/port"
	"golang.org/x/text/encoding/charmap"
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

// Execute all tst
func (ge *GenerateCase) ExecuteAll() {
	files := []string{
		"PIX.TXT",
	}
	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", outPath, file)
		ge.GeneratePixReport(filename)
	}
}

// Execute executes the generate use case
func (ge *GenerateCase) ExecuteAll2() {
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

// GeneratePixReport generates the PIX report
func (ge *GenerateCase) GeneratePixReport1(filename string) {
	// Implement the logic for generating data here
	fmt.Printf("Generating data for %s\n", filename)
	defer fmt.Println("---------------------------------------------------------------------------------------------------------")
	// read db data
	pix := domain.NewPix()
	lines, err := pix.GetDB(ge.repo)
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
	// Print header
	header := domain.NewPixHeader(time.Now())
	headerLine := header.Format()
	file.Write([]byte(headerLine))
	file.Write([]byte("\n"))
	// print lines
	for _, k := range order {
		r := lines[k].Format()
		// Convert to desired encoding
		file.Write([]byte(r))
		file.Write([]byte("\n"))
	}
	// Print trailer
	trailer := domain.NewPixTrailer(int64(len(lines)))
	trailerLine := trailer.Format()
	file.Write([]byte(trailerLine))
	file.Write([]byte("\n"))
}

// GeneratePixReport2 generates the PIX report
func (ge *GenerateCase) GeneratePixReport(filename string) {
	fmt.Printf("[%s]Generating data for %s\n", time.Now().Format("2006-01-02 15:04:05"), filename)
	defer fmt.Println("---------------------------------------------------------------------------------------------------------")
	// read db data
	pix := domain.NewPix()
	lines, err := pix.GetDBOrdered(ge.repo)
	if err != nil {
		fmt.Printf("Error getting data from DB: %s\n", err)
		return
	}
	fmt.Printf("[%s]Database got data successfully with %d lines.\n", time.Now().Format("2006-01-02 15:04:05"), len(lines))
	// control var
	var last_date = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	var file *os.File
	var count int64 = 0
	// loop
	for _, k := range lines {
		if k.(*domain.Pix).DataTransacao.After(last_date) {
			if file != nil {
				// Print trailer
				trailer := domain.NewPixTrailer(count)
				trailerLine := trailer.Format()
				file.Write([]byte(trailerLine))
				file.Write([]byte("\n"))
				file.Close()
			}
			// close previous file
			file.Close()
			// reset count
			count = 0
			// update date
			last_date = k.(*domain.Pix).DataTransacao
			// create new file
			filename := fmt.Sprintf("%s/PIX_%s.TXT", outPath, k.(*domain.Pix).DataTransacao.Format("20060102"))
			file, err = os.Create(filename)
			fmt.Printf("[%s]Creating file: %s\n", time.Now().Format("2006-01-02 15:04:05"), filename)
			if err != nil {
				fmt.Printf("Error creating file: %s\n", err)
				return
			}
			// Print header
			header := domain.NewPixHeader(k.(*domain.Pix).DataTransacao)
			headerLine := header.Format()
			file.Write([]byte(headerLine))
			file.Write([]byte("\n"))
		}
		// print lines
		r := k.Format()
		// Convert to desired encoding
		file.Write([]byte(r))
		file.Write([]byte("\n"))
		count += 1
	}
	if file != nil {
		// Print trailer
		trailer := domain.NewPixTrailer(count)
		trailerLine := trailer.Format()
		file.Write([]byte(trailerLine))
		file.Write([]byte("\n"))
		file.Close()
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
