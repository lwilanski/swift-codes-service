package main

import (
	"log"

	"github.com/your-github-name/swift-codes-service/internal/db"
	"github.com/your-github-name/swift-codes-service/internal/models"
	"github.com/your-github-name/swift-codes-service/internal/parser"
	"github.com/your-github-name/swift-codes-service/internal/repository"
	httptransport "github.com/your-github-name/swift-codes-service/internal/transport/http"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(&models.SwiftCode{}); err != nil {
		log.Fatal(err)
	}

	list, err := parser.ParseExcel("/data/Interns_2025_SWIFT_CODES.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.New(database)
	if err := repo.UpsertMany(list); err != nil {
		log.Fatal(err)
	}

	router := httptransport.Router(httptransport.New(repo))
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
