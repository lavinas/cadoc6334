package main

import (
	"fmt"
	
	"github.com/lavinas/cadoc6334/internal/domain"

)

// main function to run the ReconcileIntercam function
func main() {
	MainReconciliate()
}


// MainFakeBins creates fake BINs and prints a separator
func MainFakeBins() {
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReplicaFakeBinsFiles("sql/insert_bin.sql")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
}

// MainReconciliate runs all the reconciliation functions and prints separators between each step
func MainReconciliate() {
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReconciliateRanking("files/RANKING.TXT")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReconciliateConcred("files/CONCCRED.TXT")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReconciliateInfresta("files/INFRESTA.TXT")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReconciliateInfrterm("files/INFRTERM.TXT")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReconciliateDiscount("files/DESCONTO.TXT")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReconciliateIntercam("files/INTERCAM.TXT")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
	domain.ReconciliateSegments("files/SEGMENTO.TXT")
	fmt.Println("---------------------------------------------------------------------------------------------------------")
}
