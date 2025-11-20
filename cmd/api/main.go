package main

import (
	"log"
	"net/http"
	"os"
	"pdf-search/internal/handlers"
	"pdf-search/internal/indexer"
	"pdf-search/internal/storage"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "pdf-search/docs"
)

// @title PDF Service API
// @version 1.0
// @description Simple service to upload and retrieve PDFs.
// @BasePath /
func main() {
	os.MkdirAll("./data/pdfs", 0755)

	store := storage.NewStorage("./data/pdfs")
	indexer, err := indexer.NewIndexer("./data/index.bleve")
	if err != nil {
		log.Fatal(err)
	}
	pdfHandler := handlers.NewPDFHandler(store, indexer)

	r := chi.NewRouter()

	r.Post("/upload", pdfHandler.UploadPDF)
	r.Get("/pdf/{id}", pdfHandler.GetPDF)
	r.Get("/extract/{id}", pdfHandler.ExtractPDFText)
	r.Get("/search", pdfHandler.SearchPDF)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
