package main

import (
	"log"
	"os"
	"pdf-search/internal/handlers"
	"pdf-search/internal/indexer"
	"pdf-search/internal/storage"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "pdf-search/docs"
)

// @title PDF Service API
// @version 1.0
// @description Simple service to upload and retrieve PDFs.
// @BasePath /
func main() {
	os.MkdirAll("./data/pdfs", 0755)

	store := storage.NewStorage("./data/pdfs")
	idx, err := indexer.NewIndexer("./data/index.bleve")
	if err != nil {
		log.Fatal(err)
	}
	pdfHandler := handlers.NewPDFHandler(store, idx)

	e := echo.New()

	// Routes – just pass echo.Context directly
	e.POST("/upload", pdfHandler.UploadPDF)
	e.GET("/pdf/:id", pdfHandler.GetPDF)
	e.GET("/extract/:id", pdfHandler.ExtractPDFText)
	e.GET("/search", pdfHandler.SearchPDF)

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Println("Server running on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
