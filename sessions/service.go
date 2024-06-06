package sessions

import (
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/shvdg-dev/base-pkg/database"
	"log"
	"net/http"
)

// Service is for managing sessions.
type Service struct {
	Database *database.Manager
	Manager  *scs.SessionManager
}

// NewService creates a new instance of the Service struct.
func NewService(database *database.Manager) *Service {
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(database.DB)
	return &Service{Database: database, Manager: sessionManager}
}

// CreateSessionsTable creates the sessions table in the database and adds an expiry index.
func (s *Service) CreateSessionsTable() {
	_, err := s.Database.DB.Exec(createSessionsTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	_, err = s.Database.DB.Exec(createSessionExpiryIndexQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// Store stores a value in the session using the provided key and value.
func (s *Service) Store(key string, value interface{}, request *http.Request) {
	s.Manager.Put(request.Context(), key, value)
}

// Get retrieves the value associated with the given key from the session manager.
func (s *Service) Get(key string, request *http.Request) interface{} {
	return s.Manager.Get(request.Context(), key)
}
