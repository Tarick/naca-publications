package importer

import (
	"context"
	"encoding/json"
	"fmt"

	// "github.com/Tarick/naca-publications/internal/entity"

	"github.com/Tarick/naca-publications/internal/entity"
	"github.com/gofrs/uuid"
)

type ImportError struct {
	Publisher   Publisher
	Publication Publication
	Error       error
}

func (e *ImportError) String() string {
	return fmt.Sprint("Publisher: ", e.Publisher, " Publication: ", e.Publication)
}

type PublicationsAPIClient interface {
	CreatePublisher(ctx context.Context, name string, url string) (entity.Publisher, error)
	CreatePublication(ctx context.Context, name string, description string, languageCode string, publisherUUID uuid.UUID, publicationType string, config interface{}) (entity.Publication, error)
}

type Importer struct {
	APIClient PublicationsAPIClient
}

// Actual importer
func (ip *Importer) RunImport(bytes []byte) error {
	var entries []Entrie
	if err := json.Unmarshal(bytes, &entries); err != nil {
		return fmt.Errorf("Cannot read json from file: %s", err)
	}
	var importErrors []ImportError
	fmt.Println("Starting processing of", len(entries), "entries")

	for _, entrie := range entries {
		publisher, err := ip.APIClient.CreatePublisher(context.Background(), entrie.Publisher.Name, entrie.Publisher.URL)
		if err != nil {
			importErrors = append(importErrors, ImportError{
				Publisher: entrie.Publisher,
				Error:     err,
			})
			continue
		}
		_ = publisher
		for _, publication := range entrie.Publications {
			if _, err := ip.APIClient.CreatePublication(
				context.Background(),
				publication.Name,
				publication.Description,
				publication.LanguageCode,
				publisher.UUID,
				publication.Type,
				publication.Config); err != nil {
				importErrors = append(importErrors, ImportError{
					Publisher:   entrie.Publisher,
					Publication: publication,
					Error:       err,
				})
			}
		}
	}
	if len(importErrors) > 0 {
		for _, importError := range importErrors {
			fmt.Println(importError)
		}
		return fmt.Errorf("import failed for %d entries", len(importErrors))
	}
	return nil
}
