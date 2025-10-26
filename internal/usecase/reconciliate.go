package usecase

import (
	"fmt"
	
	"github.com/lavinas/cadoc6334/internal/domain"
	"github.com/lavinas/cadoc6334/internal/port"

)

const (
	path = "./files"
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

// Execute executes the reconciliate use case
func (uc *ReconciliateCase) Execute() {
	fmt.Println("Concred --------------------------------------------------")
	concred := domain.NewConccred()
	records, err := concred.FindAll(uc.repo)
	if err != nil {
		fmt.Printf("Error retrieving Concred records: %v\n", err)
		return
	}
	for key, record := range records {
		fmt.Printf("Key: %s, Record: %s\n", key, record.String())
	}
	fmt.Println("Discount --------------------------------------------------")
	discount := domain.NewDiscount()
	records, err = discount.FindAll(uc.repo)
	if err != nil {
		fmt.Printf("Error retrieving Discount records: %v\n", err)
		return
	}
	for key, record := range records {
		fmt.Printf("Key: %s, Record: %s\n", key, record.String())
	}
}


// Execute2 executes the check use case
func (uc *ReconciliateCase) Execute2() {
	files := []string{
		"RANKING.TXT",
		"CONCCRED.TXT",
		"INFRESTA.TXT",
		"INFRTERM.TXT",
		"DESCONTO.TXT",
		"INTERCAM.TXT",
		"SEGMENTO.TXT",
	}
	reports := []port.Report{
		domain.NewRanking(),
		domain.NewConccred(),
		domain.NewInfresta(),
		domain.NewInfrterm(),
		domain.NewDiscount(),
		domain.NewIntercam(),
		domain.NewSegment(),
	}
	for i, file := range files {
		fmt.Printf("Reconciliating %s\n", file)
		// Load the report data
		loaded, err := reports[i].GetLoaded()
		if err != nil {
			fmt.Printf("Error loading report data: %v\n", err)
			continue
		}
		filed, err := reports[i].GetParsedFile(fmt.Sprintf("%s/%s", path, file))
		if err != nil {
			fmt.Printf("Error parsing report file: %v\n", err)
			continue
		}
		// Compare loaded and filed data
		for key, loadedRecord := range loaded {
			filedRecord, exists := filed[key]
			if !exists {
				fmt.Printf("Record with key %s exists in loaded data but not in file data\n", key)
				continue
			}
			if loadedRecord.String() != filedRecord.String() {
				fmt.Printf("Mismatch for key %s:\nLoaded: %s\nFiled: %s\n", key, loadedRecord.String(), filedRecord.String())
			}
		}
		fmt.Println("---------------------------------------------------------------------------------------------------------")
	}
}

