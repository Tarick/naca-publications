package entity

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// Publication defines minimal publication type
// swagger:model
type Publication struct {
	UUID          uuid.UUID `json:"uuid"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	LanguageCode  string    `json:"language_code"`
	PublisherUUID uuid.UUID `json:"publisher_uuid"`
	Type          string    `json:"publication_type"`
}

func (p *Publication) String() string {
	return fmt.Sprintf("{UUID: %v, Name: %v, Description: %v, LanguageCode: %v, PublisherUUID: %v, Type: %v}",
		p.UUID, p.Name, p.Description, p.LanguageCode, p.PublisherUUID, p.Type)
}

// NewPublication creates Publication with new UUID
func NewPublication(name string, description string, languageCode string, publisherUUID uuid.UUID, publicationType string) (*Publication, error) {
	var err error
	p := &Publication{
		Name:          name,
		Description:   description,
		LanguageCode:  languageCode,
		PublisherUUID: publisherUUID,
		Type:          publicationType,
	}
	if p.UUID, err = uuid.NewV4(); err != nil {
		return nil, err
	}
	return p, nil
}
