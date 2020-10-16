package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Tarick/naca-publications/internal/application/server"
	"github.com/Tarick/naca-publications/internal/entity"
	"github.com/gofrs/uuid"
)

const publicationsPath string = "publications"
const publishersPath string = "publishers"

// TODO: WithTimeout?
// New creates API http client
func New(serviceAPIURL string) *client {
	return &client{
		publishersURL:   fmt.Sprintf("%s/%s", serviceAPIURL, publishersPath),
		publicationsURL: fmt.Sprintf("%s/%s", serviceAPIURL, publicationsPath),
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// TODO: add logger
type client struct {
	publishersURL   string
	publicationsURL string
	httpClient      *http.Client
}

func (c *client) CreatePublisher(ctx context.Context, name string, url string) (entity.Publisher, error) {
	publisher := &entity.Publisher{
		Name: name,
		URL:  url,
	}
	body, err := json.Marshal(publisher)
	if err != nil {
		return entity.Publisher{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", c.publishersURL), bytes.NewReader(body))
	if err != nil {
		return entity.Publisher{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return entity.Publisher{}, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated {
		// Create new publisher from response
		responsePublisher := &entity.Publisher{}
		if err = json.NewDecoder(res.Body).Decode(responsePublisher); err == nil {
			return entity.Publisher{}, err
		}
		return *responsePublisher, nil
	}
	// handle error
	var errRes server.ErrResponseBody
	if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
		return entity.Publisher{}, errors.New(errRes.ErrorText)
	}
	return entity.Publisher{}, fmt.Errorf("unknown error, status code: %d, message: %v", res.StatusCode, res.Status)
}
func (c *client) CreatePublication(
	ctx context.Context,
	name string,
	description string,
	languageCode string,
	publisherUUID uuid.UUID,
	publicationType string,
	config interface{}) (entity.Publication, error) {
	publicationRequest := &server.PublicationRequestBody{
		Name:          name,
		Description:   description,
		LanguageCode:  languageCode,
		PublisherUUID: publisherUUID,
		Type:          publicationType,
		Config:        config,
	}
	body, err := json.Marshal(publicationRequest)
	if err != nil {
		return entity.Publication{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", c.publicationsURL), bytes.NewReader(body))
	if err != nil {
		return entity.Publication{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return entity.Publication{}, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated {
		// Create new publisher from response
		responsePublication := &entity.Publication{}
		if err = json.NewDecoder(res.Body).Decode(responsePublication); err == nil {
			return entity.Publication{}, err
		}
		return *responsePublication, nil
	}
	// handle error
	var errRes server.ErrResponseBody
	if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
		return entity.Publication{}, errors.New(errRes.ErrorText)
	}
	return entity.Publication{}, fmt.Errorf("unknown error, status code: %d, message: %v", res.StatusCode, res.Status)
}
