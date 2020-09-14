package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Tarick/naca-publications/internal/entity"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-chi/stampede"
	"github.com/gofrs/uuid"
)

func (s *Server) publishersRouter() http.Handler {
	r := chi.NewRouter()
	// Set 1 second caching and requests coalescing to avoid requests stampede. Beware of any user specific responses.
	cached := stampede.Handler(512, 1*time.Second)

	// swagger:operation GET /publishers getPublishers
	// Returns all publishers registered in db
	// ---
	// responses:
	//   '200':
	//     description: list all publishers
	//     schema:
	//       type: array
	//       items:
	//         $ref: "#/definitions/PublisherResponseBody"
	r.With(cached).Get("/", s.getPublishers)

	// swagger:operation  POST /publishers createPublisher
	// Creates publisher using supplied params from body
	// ---
	// parameters:
	//  - $ref: "#/definitions/Publisher"
	// responses:
	//    '201':
	//      $ref: "#/responses/PublisherResponse"
	//    default:
	//      $ref: "#/responses/ErrResponse"
	r.Post("/", s.createPublisher)

	r.Route("/{publisher_uuid}", func(r chi.Router) {
		r.Use(s.publisherCtx) // handle publisher_uuid

		// swagger:operation GET /publishers/{publisher_uuid} getPublisher
		// Gets single publisher using its publisher_uuid as parameter
		// ---
		// parameters:
		//  - name: publisher_uuid
		//    in: path
		//    description: publisher_uuid to get
		//    required: true
		//    type: string
		// responses:
		//    '200':
		//      $ref: "#/responses/PublisherResponse"
		//    default:
		//      $ref: "#/responses/ErrResponse"
		r.Get("/", s.getPublisher)
		// swagger:operation GET /publishers/{publisher_uuid}/publications getPublisherPublications
		// Get publisher publications
		// ---
		// parameters:
		//  - name: publisher_uuid
		//    in: path
		//    description: publisher_uuid
		//    required: true
		//    type: string
		// responses:
		//    '200':
		//      $ref: "#/responses/PublicationsResponse"
		//    default:
		//      $ref: "#/responses/ErrResponse"
		r.Get("/publications", s.getPublisherPublications)

		// swagger:operation PUT /publishers/{publication_uuid} updatePublisher
		// Modifies Publisher using supplied params from body
		// ---
		// parameters:
		//  - name: publisher_uuid
		//    in: path
		//    description: Publisher publisher_uuid to update
		//    required: true
		//    type: string
		//  - $ref: "#/definitions/Publisher"
		// responses:
		//    '200':
		//      $ref: "#/responses/PublisherResponse"
		//    default:
		//      $ref: "#/responses/ErrResponse"
		r.Put("/", s.updatePublisher)

		// swagger:operation DELETE /publishers/{publication_uuid} deletePublisher
		// Deletes publisher using its uuid
		// ---
		// parameters:
		//  - name: publisher_uuid
		//    in: path
		//    description: Publisher uuid to delete
		//    required: true
		//    type: string
		// responses:
		//  '204':
		//    description: Send success
		//  default:
		//    $ref: "#/responses/ErrResponse"
		r.Delete("/", s.deletePublisher)
	})
	return r
}

// PublisherResponse defines Feed response with Body and any additional headers
// swagger:response
type PublisherResponse struct {
	// in: body
	Body PublisherResponseBody
}

// PublisherResponseBody is returned on successfull operations to get, create publisher.
type PublisherResponseBody struct {
	// swagger:allOf
	*entity.Publisher
}

// Render converts PublisherResponseBody to json and sends it to client
func (pr *PublisherResponse) Render(w http.ResponseWriter, r *http.Request) {
	// Pre-processing before a response is marshalled and sent across the wire
	// Any instructions here
	render.JSON(w, r, pr.Body)
}

func newPublisherResponse(publisher *entity.Publisher) *PublisherResponse {
	return &PublisherResponse{Body: PublisherResponseBody{publisher}}
}

// PublisherRequest defines Publisher request with Body and any additional headers
type PublisherRequest struct {
	// in: body
	Body PublisherRequestBody
}

// PublisherRequestBody contains information on publisher creation
type PublisherRequestBody struct {
	// swagger:allOf
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Bind implements Bind interface for chi Bind to map request body to request body struct, with simple validator
func (p *PublisherRequestBody) Bind(r *http.Request) error {
	if p == nil {
		return errors.New("request body is empty")
	}
	if p.Name == "" {
		return errors.New("missing required 'name' field")
	}
	if p.URL == "" {
		return errors.New("missing required 'url' field")
	}
	return nil
}

// Used as middleware to load object from the URL parameters passed through as the request.
// If not found - 404
func (s *Server) publisherCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		publisherUUIDParam := chi.URLParam(r, "publisher_uuid")
		if publisherUUIDParam == "" {
			ErrInvalidRequest(errors.New("empty publisher uuid")).Render(w, r)
			return
		}
		publisherUUID, err := uuid.FromString(publisherUUIDParam)
		if err != nil {
			// log.Error(fmt.Sprintf("Couldn't convert publisher uuid param %s to UUID: %s", publisherUUIDParam, err))
			ErrInvalidRequest(fmt.Errorf("invalid uuid parameter %s", publisherUUIDParam)).Render(w, r)
			return
		}

		publisher, err := s.repository.GetPublisher(publisherUUID)
		if err != nil {
			ErrInternal(fmt.Errorf("Failure getting Publisher data")).Render(w, r)
			return
		}
		if publisher == nil {
			ErrNotFound.Render(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "publisher", publisher)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// // Response with single feed
func (s *Server) getPublisher(w http.ResponseWriter, r *http.Request) {
	publisher := r.Context().Value("publisher").(*entity.Publisher)
	newPublisherResponse(publisher).Render(w, r)
}
func (s *Server) updatePublisher(w http.ResponseWriter, r *http.Request) {
	publisher := r.Context().Value("publisher").(*entity.Publisher)
	data := &PublisherRequestBody{}
	if err := render.Bind(r, data); err != nil {
		ErrInvalidRequest(err).Render(w, r)
		return
	}
	publisher.Name = data.Name
	publisher.URL = data.URL
	if err := s.repository.UpdatePublisher(publisher); err != nil {
		// log.Error(fmt.Sprintf("Failure updating publisher %v: %s", publisher, err))
		ErrInternal(fmt.Errorf("Failure updating publisher")).Render(w, r)
		return
	}
	newPublisherResponse(publisher).Render(w, r)
}

func (s *Server) createPublisher(w http.ResponseWriter, r *http.Request) {
	data := &PublisherRequestBody{}
	if err := render.Bind(r, data); err != nil {
		ErrInvalidRequest(err).Render(w, r)
		return
	}
	// FIXME: validate if publisher already exists
	publisher, err := entity.NewPublisher(data.Name, data.URL)
	if err != nil {
		ErrInternal(err).Render(w, r)
		return
	}
	if err := s.repository.CreatePublisher(publisher); err != nil {
		// log.Error(fmt.Sprintf("Failure creating publisher %v in database: %s", publisher, err))
		ErrInternal(fmt.Errorf("Failure creating publisher")).Render(w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	newPublisherResponse(publisher).Render(w, r)
}

func (s *Server) deletePublisher(w http.ResponseWriter, r *http.Request) {
	publisher := r.Context().Value("publisher").(*entity.Publisher)
	if err := s.repository.DeletePublisher(publisher.UUID); err != nil {
		// log.Error(fmt.Sprintf("Failure deleting publisher %v: %s", publisher, err))
		ErrInternal(fmt.Errorf("Failure deleting publisher %v", publisher)).Render(w, r)
		return
	}
	render.NoContent(w, r)
}

// TODO: filtering
func (s *Server) getPublishers(w http.ResponseWriter, r *http.Request) {
	publishers, err := s.repository.GetPublishers()
	if err != nil {
		// log.Error(fmt.Sprint("Failure querying for publishers: ", err))
		ErrInternal(errors.New("Failure querying database for publishers")).Render(w, r)
		return
	}
	response := make([]*PublisherResponseBody, len(publishers), len(publishers))
	for i := 0; i < len(publishers); i++ {
		response[i] = &newPublisherResponse(publishers[i]).Body
	}
	render.JSON(w, r, response)
}

func (s *Server) getPublisherPublications(w http.ResponseWriter, r *http.Request) {
	publisher := r.Context().Value("publisher").(*entity.Publisher)
	publications, err := s.repository.GetPublicationsByPublisher(publisher.UUID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failure querying for publications per publisher %s: %v", publisher.UUID, err))
		ErrInternal(fmt.Errorf("Failure querying database for publications")).Render(w, r)
		return
	}
	response := make([]*PublicationResponseBody, len(publications), len(publications))
	for i := 0; i < len(publications); i++ {
		response[i] = &newPublicationResponse(publications[i]).Body
	}
	render.JSON(w, r, response)
}
