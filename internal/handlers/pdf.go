package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"pdf-search/internal/extractor"
	"pdf-search/internal/indexer"
	"pdf-search/internal/storage"
	"strings"

	"github.com/go-chi/chi/v5"
)

type PDFHandler struct {
	Store   *storage.Storage
	Indexer *indexer.Indexer
}

func NewPDFHandler(s *storage.Storage, idx *indexer.Indexer) *PDFHandler {
	return &PDFHandler{Store: s, Indexer: idx}
}

// UploadPDF godoc
// @Summary Upload a PDF file
// @Description Upload, store, extract text, and index
// @Tags pdf
// @Accept multipart/form-data
// @Produce text/plain
// @Param file formData file true "PDF file"
// @Success 200 {string} string "ID of uploaded PDF"
// @Failure 400 {string} string "Bad request"
// @Router /upload [post]
func (h *PDFHandler) UploadPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		http.Error(w, "Only PDFs allowed", http.StatusBadRequest)
		return
	}

	id, err := h.Store.SavePDF(file, header.Filename)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Extract text immediately and index
	path := h.Store.GetPDFPath(id)
	text, err := extractor.ExtractText(path)
	if err != nil {
		log.Println("Failed to extract text for indexing:", err)
	} else {
		err = h.Indexer.IndexPDF(id, text)
		if err != nil {
			log.Println("Failed to index PDF:", err)
		}
	}

	w.Write([]byte(id))
}

// GetPDF godoc
// @Summary Retrieve a PDF
// @Description Returns the raw PDF file
// @Tags pdf
// @Produce application/pdf
// @Param id path string true "PDF ID"
// @Success 200
// @Failure 404
// @Router /pdf/{id} [get]
func (h *PDFHandler) GetPDF(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	path := h.Store.GetPDFPath(id)
	http.ServeFile(w, r, path)
}

// ExtractPDFText godoc
// @Summary Extract text from PDF
// @Tags pdf
// @Param id path string true "PDF ID"
// @Produce text/plain
// @Success 200
// @Router /extract/{id} [get]
func (h *PDFHandler) ExtractPDFText(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	path := h.Store.GetPDFPath(id)

	text, err := extractor.ExtractText(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}

// SearchPDF godoc
// @Summary Search PDFs by text
// @Tags pdf
// @Param q query string true "Search query"
// @Produce application/json
// @Success 200 {array} string
// @Router /search [get]
func (h *PDFHandler) SearchPDF(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	results, err := h.Indexer.Search(query, 20) // top 20 results
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
