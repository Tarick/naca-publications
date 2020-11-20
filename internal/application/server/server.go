package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Tarick/naca-publications/internal/entity"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

// Server defines HTTP application
type Server struct {
	httpServer        *http.Server
	logger            Logger
	repository        PublicationsRepository
	rssFeedsAPIClient RSSFeedsAPIClient
}

// PublicationsRepository represents repository for both publishers and publications
type PublicationsRepository interface {
	CreatePublication(*entity.Publication) error
	UpdatePublication(*entity.Publication) error
	DeletePublication(uuid.UUID) error
	GetPublication(uuid.UUID) (*entity.Publication, error)
	GetPublications() ([]*entity.Publication, error)
	GetPublicationsByPublisher(uuid.UUID) ([]*entity.Publication, error)
	CreatePublisher(*entity.Publisher) error
	UpdatePublisher(*entity.Publisher) error
	DeletePublisher(uuid.UUID) error
	GetPublisher(uuid.UUID) (*entity.Publisher, error)
	GetPublishers() ([]*entity.Publisher, error)
}

// RSSFeedsAPIClient is used to call RSS Feeds service
type RSSFeedsAPIClient interface {
	CreateRSSFeed(context.Context, uuid.UUID, string) error
	UpdateRSSFeed(context.Context, uuid.UUID, string) error
	DeleteRSSFeed(context.Context, uuid.UUID) error
}

// Config defines webserver configuration
type Config struct {
	Address        string `mapstructure:"address"`
	RequestTimeout int    `mapstructure:"request_timeout"`
}

// New creates new server configuration and configurates middleware
func New(serverConfig Config, logger Logger, repository PublicationsRepository, rssFeedsAPIClient RSSFeedsAPIClient) *Server {
	r := chi.NewRouter()
	srv := &Server{
		httpServer:        &http.Server{Addr: serverConfig.Address, Handler: r},
		logger:            logger,
		repository:        repository,
		rssFeedsAPIClient: rssFeedsAPIClient,
	}
	r.Use(middleware.RequestID)
	r.Use(middlewareLogger(logger))
	// Basic CORS to allow API calls from browsers (Swagger-UI)
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"},
		// Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// r.Use(middleware.Timeout(time.Duration(serverConfig.RequestTimeout) * time.Second))
	// Healthcheck
	// Could be moved back to middleware in case auth middleware meddling
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("."))
	},
	)
	// Create a route along /doc that will serve contents from
	// the ./swaggerui directory.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "swaggerui"))
	FileServer(r, "/doc", filesDir)
	r.Mount("/publishers", srv.publishersRouter())
	r.Mount("/publications", srv.publicationsRouter())
	return srv
}

// StartAndServe starts http server with signal control
func (s *Server) StartAndServe() {
	// Start server in background to unblock further signal interrupt
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal(fmt.Sprint("Server startup failed: ", err))
		}
	}()
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)
	<-signalChan
	s.logger.Warn("Interrupted by system or user, shutting down...")

	// Listen for second ctrl-c and kill unconditionally
	go func() {
		<-signalChan
		s.logger.Fatal("Killed - terminating...\n")
	}()
	// Gracefully shutdown server. (Note - without any websockets requests graceful shutdown)
	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := s.httpServer.Shutdown(gracefullCtx); err != nil {
		s.logger.Error(fmt.Sprint("Http server shutdown error: ", err))
		defer os.Exit(1)
		return
	}
	s.logger.Info("Gracefully stopped http server")
	defer os.Exit(0)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem. Used for Swagger-UI and swagger.json files.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
