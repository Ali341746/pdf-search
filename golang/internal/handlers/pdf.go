package handlers

import (
	"log"
	"net/http"
	"pdf-search/internal/extractor"
	"pdf-search/internal/indexer"
	"pdf-search/internal/storage"
	"strings"

	"github.com/labstack/echo/v4"
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
func (h *PDFHandler) UploadPDF(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to read file")
	}

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to open file")
	}
	defer src.Close()

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		return c.String(http.StatusBadRequest, "Only PDFs allowed")
	}

	id, err := h.Store.SavePDF(src, file.Filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save file")
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

	return c.String(http.StatusOK, id)
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
func (h *PDFHandler) GetPDF(c echo.Context) error {
	id := c.Param("id")
	path := h.Store.GetPDFPath(id)
	return c.File(path)
}

// ExtractPDFText godoc
// @Summary Extract text from PDF
// @Tags pdf
// @Param id path string true "PDF ID"
// @Produce text/plain
// @Success 200
// @Router /extract/{id} [get]
func (h *PDFHandler) ExtractPDFText(c echo.Context) error {
	id := c.Param("id")
	path := h.Store.GetPDFPath(id)

	text, err := extractor.ExtractText(path)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, text)
}

// SearchPDF godoc
// @Summary Search PDFs by text
// @Tags pdf
// @Param q query string true "Search query"
// @Produce application/json
// @Success 200 {array} string
// @Router /search [get]
func (h *PDFHandler) SearchPDF(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.String(http.StatusBadRequest, "Missing query parameter 'q'")
	}

	results, err := h.Indexer.Search(query, 20) // top 20 results
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}
