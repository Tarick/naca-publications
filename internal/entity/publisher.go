package entity

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// Publisher defines minimal publisher type
// swagger:model
type Publisher struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	URL  string    `json:"url"`
}

func (p *Publisher) String() string {
	return fmt.Sprintf("{UUID: %v, Name: %v, URL: %v}", p.UUID, p.Name, p.URL)
}

// NewPublisher creates Publisher with new UUID
func NewPublisher(name string, url string) (*Publisher, error) {
	var err error
	p := &Publisher{
		Name: name,
		URL:  url,
	}
	p.UUID, err = uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return p, nil
}
