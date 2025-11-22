package usecase

import (
	"fmt"

	"github.com/lavinas/cadoc6334/internal/domain"
	"github.com/lavinas/cadoc6334/internal/port"
	"golang.org/x/text/encoding/charmap"
)

const (
	inPath = "./files/in"
)

// ReconciliateCase represents the use case for checking or validating data
type ReconciliateCase struct {
	repo port.Repository
	// Add any dependencies or configurations needed for the use case
}

// NewReconciliateCase creates a new instance of ReconciliateCase
func NewReconciliateCase(repo port.Repository) *ReconciliateCase {
	return &ReconciliateCase{
		repo: repo,
	}
}

// Execute2 executes the check use case
func (uc *ReconciliateCase) ExecuteAll() {
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
	}
	for i, file := range files {
		filename := fmt.Sprintf("%s/%s", inPath, file)
		uc.ExecuteReport(reports[i], filename)
	}
}

// ExecuteReport executes the check use case for a specific report
func (uc *ReconciliateCase) ExecuteReport(report port.Report, filename string) {
	fmt.Printf("Reconciliating %s\n", filename)
	defer fmt.Println("---------------------------------------------------------------------------------------------------------")
	// Get db data
	loaded, err := report.GetDB(uc.repo)
	if err != nil {
		fmt.Printf("Error loading report data: %v\n", err)
		return
	}
	// Get file data
	filed, err := report.GetParsedFile(filename)
	if err != nil {
		fmt.Printf("Error parsing report file: %v\n", err)
		return
	}
	// validate DB
	if err := uc.validateReport(loaded); err != nil {
		fmt.Println("DB validation errors found:")
		for _, e := range err {
			fmt.Println(e)
		}
		return
	}
	// validate File
	if err := uc.validateReport(filed); err != nil {
		fmt.Println("File validation errors found:")
		for _, e := range err {
			fmt.Println(e)
		}
		return
	}
	// Match and report discrepancies
	errs := uc.match(loaded, filed)
	if len(errs) > 0 {
		fmt.Printf("Discrepancies found in %s:\n", filename)
		for _, e := range errs {
			fmt.Println(e)
		}
	} else {
		fmt.Printf("No discrepancies found in %s\n", filename)
	}
}

// validate validates records from both sources.
func (uc *ReconciliateCase) validateReport(report map[string]port.Report) []error {
	var errs []error
	for key, dbRecord := range report {
		if err := dbRecord.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("validation error for DB record with key %s: %v", key, err))
		}
	}
	return errs
}

// match compares two maps of records and returns a slice of errors for any discrepancies found.
func (uc *ReconciliateCase) match(db map[string]port.Report, file map[string]port.Report) []error {
	var errs []error
	// compare lengths
	if len(db) != len(file) {
		return []error{fmt.Errorf("length mismatch: DB has %d records, File has %d records", len(db), len(file))}
	}
	encoder := charmap.ISO8859_1.NewEncoder()

	for key, dbRecord := range db {
		fileRecord, exists := file[key]
		if !exists {
			errs = append(errs, fmt.Errorf("db record with key %s exists in DB but not in file", key))
			continue
		}
		encodedDBBytes, err := encoder.Bytes([]byte(dbRecord.String()))
		if err != nil {
			errs = append(errs, fmt.Errorf("error encoding DB record with key %s: %v", key, err))
			continue
		}
		encoderFileBytes, err := encoder.Bytes([]byte(fileRecord.String()))
		if err != nil {
			errs = append(errs, fmt.Errorf("error encoding File record with key %s: %v", key, err))
			continue
		}
		encodedDBString := string(encodedDBBytes)
		encodedFileString := string(encoderFileBytes)
		if encodedDBString != encodedFileString {
			errs = append(errs, fmt.Errorf("mismatch for key %s:\nDB: %s\nFile: %s", key, encodedDBString, encodedFileString))
		}
	}

	for key, fileRecord := range file {
		dbRecord, exists := db[key]
		if !exists {
			errs = append(errs, fmt.Errorf("filerecord with key %s exists in File but not in DB", key))
			continue
		}
		encodedDBBytes, err := encoder.Bytes([]byte(dbRecord.String()))
		if err != nil {
			errs = append(errs, fmt.Errorf("error encoding DB record with key %s: %v", key, err))
			continue
		}
		encoderFileBytes, err := encoder.Bytes([]byte(fileRecord.String()))
		if err != nil {
			errs = append(errs, fmt.Errorf("error encoding File record with key %s: %v", key, err))
			continue
		}
		encodedDBString := string(encodedDBBytes)
		encodedFileString := string(encoderFileBytes)
		if encodedDBString != encodedFileString {
			errs = append(errs, fmt.Errorf("mismatch for key %s:\nDB: %s\nFile: %s", key, encodedDBString, encodedFileString))
		}
	}

	return errs
}
