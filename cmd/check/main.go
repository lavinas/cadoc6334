package main

import (
	"github.com/lavinas/cadoc6334/internal/usecase"
	"github.com/lavinas/cadoc6334/internal/adapter"
)

// main function to run the ReconcileIntercam function
func main() {
	repo, err := adapter.NewPostgresGormAdapter(adapter.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "root",
		Password: "root",
		DBName:   "cadoc",
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}
	defer repo.Close()
	usecase.NewReconciliateCase(repo).Execute()
}


