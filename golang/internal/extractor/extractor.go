package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

type PDFPath struct {
	Path string `json:"path"`
}

type ExtractResponse struct {
	Filename string `json:"filename"`
	Text     string `json:"text"`
}

func ExtractText(path string) (string, error) {
	url := "http://localhost:8081/extract"

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	// Replace backslashes with forward slashes for Python
	absPath = strings.ReplaceAll(absPath, `\`, `/`)

	payload, err := json.Marshal(PDFPath{Path: absPath})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error from API: %s", string(body))
	}

	var result ExtractResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}
