package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"errors"

	"github.com/Tarick/naca-publications/internal/entity"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-chi/stampede"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
)

const PublicationTypeRSS string = "rss"

func (s *Server) publicationsRouter() http.Handler {
	r := chi.NewRouter()
	// Set 1 second caiching and requests coalescing to avoid requests stampede. Beware of any user specific responses.
	cached := stampede.Handler(512, 1*time.Second)

	// swagger:operation GET /publications getPublications
	// Returns all publications registered in db
	// ---
	// responses:
	//   '200':
	//     description: list all publications
	//     schema:
	//       type: array
	//       items:
	//         $ref: "#/definitions/PublicationResponseBody"
	r.With(cached).Get("/", s.getPublications)

	// swagger:operation  POST /publications createPublication
	// Creates publication using supplied params from body
	// ---
	// parameters:
	//  - $ref: "#/definitions/Publication"
	// responses:
	//    '201':
	//      $ref: "#/responses/PublicationResponse"
	//    default:
	//      $ref: "#/responses/ErrResponse"
	r.Post("/", s.createPublication)

	r.Route("/{publication_uuid}", func(r chi.Router) {
		r.Use(s.publicationCtx) // handle publication_uuid

		// swagger:operation GET /publications/{publication_uuid} getPublication
		// Gets single publication using its publication_uuid as parameter
		// ---
		// parameters:
		//  - name: publication_uuid
		//    in: path
		//    description: publication_uuid to get
		//    required: true
		//    type: string
		// responses:
		//    '200':
		//      $ref: "#/responses/PublicationResponse"
		//    default:
		//      $ref: "#/responses/ErrResponse"
		r.Get("/", s.getPublication)

		// swagger:operation PUT /publications/{publication_uuid} updatePublication
		// Modifies Publication using supplied params from body
		// ---
		// parameters:
		//  - name: publication_uuid
		//    in: path
		//    description: Publication publication_uuid to update
		//    required: true
		//    type: string
		//  - $ref: "#/definitions/Publication"
		// responses:
		//    '200':
		//      $ref: "#/responses/PublicationResponse"
		//    default:
		//      $ref: "#/responses/ErrResponse"
		r.Put("/", s.updatePublication)

		// swagger:operation DELETE /publications/{publication_uuid} deletePublication
		// Deletes publication using its uuid
		// ---
		// parameters:
		//  - name: publication_uuid
		//    in: path
		//    description: Publication uuid to delete
		//    required: true
		//    type: string
		// responses:
		//  '204':
		//    description: Send success
		//  default:
		//    $ref: "#/responses/ErrResponse"
		r.Delete("/", s.deletePublication)
	})
	return r
}

// PublicationResponse defines response with data body and any additional headers
// swagger:response
type PublicationResponse struct {
	// in: body
	Body PublicationResponseBody
}

// PublicationResponseBody is returned on successfull operations to get, create publication.
type PublicationResponseBody struct {
	// swagger:allOf
	*entity.Publication
}

// Render converts PublicationResponseBody to json and sends it to client
func (pr *PublicationResponse) Render(w http.ResponseWriter, r *http.Request) {
	// Pre-processing before a response is marshalled and sent across the wire
	// Any instructions here
	render.JSON(w, r, pr.Body)
}

func newPublicationResponse(publication *entity.Publication) *PublicationResponse {
	return &PublicationResponse{Body: PublicationResponseBody{publication}}
}

// PublicationRequest defines Publication create/update request with required Body and any additional headers
type PublicationRequest struct {
	// in: body
	Body PublicationRequestBody
}

// PublicationRequestBody contains information on publication creation
type PublicationRequestBody struct {
	// swagger:allOf
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	LanguageCode  string    `json:"language_code"`
	PublisherUUID uuid.UUID `json:"publisher_uuid"`
	Type          string    `json:"publication_type"`
	// Config content is different for different publication types.
	// when parsing, we decide on Type
	Config PublicationConfig `json:"config"`
}

// PublicationConfig is used to pass around different config structs
type PublicationConfig interface{}

// RSSPublicationConfig defines config for RSS Feeds
type RSSPublicationConfig struct {
	URL string `json:"url"`
}

// Only RSS is present now, but here's tested configs as well
// type APIPublicationConfig struct {
// 	URL    string `json:"url"`
// 	APIKey string `json:"api_key"`
// 	LanguageCode string `json:"language_code"`
// }

func (c *RSSPublicationConfig) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.URL, validation.Required, validation.Length(5, 100), is.URL),
	)
}

// Validate body
func (b *PublicationRequestBody) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required, validation.Length(2, 300)),
		validation.Field(&b.Description, validation.Required, validation.Length(5, 300)),
		validation.Field(&b.PublisherUUID, validation.Required, is.UUID, validation.By(checkUUIDNotNil)),
		validation.Field(&b.LanguageCode, validation.Required, validation.Length(2, 2), is.Alpha, is.LowerCase),
		validation.Field(&b.Type, validation.Required, validation.By(checkPublicationType)),
		validation.Field(&b.Config),
	)
}

// validation helper to check UUID
func checkUUIDNotNil(value interface{}) error {
	u, _ := value.(uuid.UUID)
	if u == uuid.Nil {
		return errors.New("uuid is nil")
	}
	return nil
}

// validation helper to check publication type
func checkPublicationType(value interface{}) error {
	s, _ := value.(string)
	// check can be extended with other types
	if s != PublicationTypeRSS {
		return fmt.Errorf("unknow publication type: %s", s)
	}
	return nil
}

// Used as middleware to load an feed object from the URL parameters passed through as the request.
// If not found - 404
func (s *Server) publicationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		publicationUUIDParam := chi.URLParam(r, "publication_uuid")
		if publicationUUIDParam == "" {
			ErrInvalidRequest(fmt.Errorf("empty publication uuid")).Render(w, r)
			return
		}
		publicationUUID, err := uuid.FromString(publicationUUIDParam)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Couldn't convert publication uuid param %s to UUID: %s", publicationUUIDParam, err))
			ErrInvalidRequest(fmt.Errorf("invalid uuid parameter %s", publicationUUIDParam)).Render(w, r)
			return
		}

		publication, err := s.repository.GetPublication(publicationUUID)
		if err != nil {
			ErrInternal(fmt.Errorf("Failure getting Publication data")).Render(w, r)
			return
		}
		// 404
		if publication == nil {
			ErrNotFound.Render(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "publication", publication)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Response with single feed
// TODO: with full details from subservices (RSS API)
func (s *Server) getPublication(w http.ResponseWriter, r *http.Request) {
	publication := r.Context().Value("publication").(*entity.Publication)
	newPublicationResponse(publication).Render(w, r)
}

// TODO: implement update of sub services
func (s *Server) updatePublication(w http.ResponseWriter, r *http.Request) {
	publication := r.Context().Value("publication").(*entity.Publication)
	// FIXME
	// publicationUpdated, publicationConfig, err := requestToPublication(r)
	publicationUpdated, _, err := requestToPublication(r)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failure processing request: %s", err))
		ErrInvalidRequest(err).Render(w, r)
		return
	}
	publication.Name, publication.Description, publication.LanguageCode = publicationUpdated.Name, publicationUpdated.Description, publicationUpdated.LanguageCode
	if err := s.repository.UpdatePublication(publication); err != nil {
		s.logger.Error(fmt.Sprintf("Failure updating publication %v: %s", publication, err))
		ErrInternal(fmt.Errorf("Failure updating publication")).Render(w, r)
		return
	}
	newPublicationResponse(publication).Render(w, r)
}
func requestToPublication(r *http.Request) (*entity.Publication, PublicationConfig, error) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	var (
		publicationConfigBody json.RawMessage
		publicationConfig     PublicationConfig
		publication           *entity.Publication
	)
	publicationRequestBody := &PublicationRequestBody{Config: &publicationConfigBody}
	if err := json.Unmarshal(requestBody, publicationRequestBody); err != nil {
		return nil, nil, err
	}
	if publication, err = entity.NewPublication(
		publicationRequestBody.Name,
		publicationRequestBody.Description,
		publicationRequestBody.LanguageCode,
		publicationRequestBody.PublisherUUID,
		publicationRequestBody.Type); err != nil {
		return nil, nil, err
	}
	switch publicationRequestBody.Type {
	case PublicationTypeRSS:
		config := RSSPublicationConfig{}
		if err := json.Unmarshal(publicationConfigBody, &config); err != nil {
			return nil, nil, err
		}
		publicationConfig = config
	default:
		return nil, nil, fmt.Errorf("incorrect 'publication_type' specified in request: %v", publicationRequestBody.Type)
	}
	if err := publicationRequestBody.Validate(); err != nil {
		return nil, nil, err
	}
	return publication, publicationConfig, nil
}

func (s *Server) createPublication(w http.ResponseWriter, r *http.Request) {
	publication, publicationConfig, err := requestToPublication(r)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failure processing request: %s", err))
		ErrInvalidRequest(err).Render(w, r)
		return
	}
	if err := s.repository.CreatePublication(publication); err != nil {
		s.logger.Error(fmt.Sprintf("Failure creating publication %v in database: %s", publication, err))
		ErrInternal(fmt.Errorf("Failure creating publication")).Render(w, r)
		return
	}
	switch v := publicationConfig.(type) {
	case RSSPublicationConfig:
		//FIXME: Fix context
		if err := s.rssFeedsAPIClient.CreateRSSFeed(context.Background(), publication.UUID, v.URL, publication.LanguageCode); err != nil {
			s.logger.Error("Failure creating RSS Feeds Publication in RSS service: ", err)
			errs := fmt.Errorf("failure creating publication: %w", err)
			// revert publication creation. No need for full saga patern yet.
			if err = s.repository.DeletePublication(publication.UUID); err != nil {
				s.logger.Error("Failure deleting failed RSS Publication from repository: ", err)
				errs = fmt.Errorf("%w, failure deleting created publication from repository: %w", errs, err)
			}
			s.logger.Debug("Deleted publication with UUID ", publication.UUID, "from repository")
			ErrInternal(errs).Render(w, r)
			return
		}
	}
	render.Status(r, http.StatusCreated)
	newPublicationResponse(publication).Render(w, r)
}

// TODO: implement removal from sub services (RSS)
func (s *Server) deletePublication(w http.ResponseWriter, r *http.Request) {
	publication := r.Context().Value("publication").(*entity.Publication)
	if err := s.repository.DeletePublication(publication.UUID); err != nil {
		s.logger.Error(fmt.Sprintf("Failure deleting publication %v: %s", publication, err))
		ErrInternal(fmt.Errorf("Failure deleting publication %v", publication)).Render(w, r)
		return
	}
	render.NoContent(w, r)
}

// Returns publication entries
// TODO: filtering
// TODO: get data with full details from sub services (RSS API, scraping config?)
func (s *Server) getPublications(w http.ResponseWriter, r *http.Request) {
	publications, err := s.repository.GetPublications()
	if err != nil {
		s.logger.Error(fmt.Sprint("Failure querying for publications: ", err))
		ErrInternal(fmt.Errorf("Failure querying database for publications")).Render(w, r)
		return
	}
	response := make([]*PublicationResponseBody, len(publications), len(publications))
	for i := 0; i < len(publications); i++ {
		response[i] = &newPublicationResponse(publications[i]).Body
	}
	render.JSON(w, r, response)
}
