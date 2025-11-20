package indexer

import (
	"log"

	"github.com/blevesearch/bleve/v2"
)

type Indexer struct {
	index bleve.Index
}

func NewIndexer(path string) (*Indexer, error) {
	var idx bleve.Index
	var err error

	// Try opening existing index
	idx, err = bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		// Create a new one
		mapping := bleve.NewIndexMapping()
		idx, err = bleve.New(path, mapping)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &Indexer{index: idx}, nil
}

func (i *Indexer) IndexPDF(id, text string) error {
	doc := struct {
		ID   string
		Text string
	}{ID: id, Text: text}

	return i.index.Index(id, doc)
}

func (i *Indexer) Search(queryStr string, size int) ([]string, error) {
	query := bleve.NewMatchQuery(queryStr)
	searchRequest := bleve.NewSearchRequestOptions(query, size, 0, false)

	searchResult, err := i.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, hit := range searchResult.Hits {
		results = append(results, hit.ID)
	}

	return results, nil
}

func (i *Indexer) Close() {
	if err := i.index.Close(); err != nil {
		log.Println("Failed to close index:", err)
	}
}
