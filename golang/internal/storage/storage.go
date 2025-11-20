package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Storage struct {
	BasePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{BasePath: basePath}
}

func (s *Storage) SavePDF(file io.Reader, originalName string) (string, error) {
	id := uuid.NewString()
	filename := id + ".pdf"

	path := filepath.Join(s.BasePath, filename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) GetPDFPath(id string) string {
	return filepath.Join(s.BasePath, id+".pdf")
}
