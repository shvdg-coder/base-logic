package sessions

import (
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	consts "tabs/internal/constants"
	"tabs/pkg/base/database"
)

// Service is for managing sessions.
type Service struct {
	Database *database.Connection
	Manager  *scs.SessionManager
}

// NewService creates a new instance of the Service struct.
func NewService(database *database.Connection) *Service {
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

// IsAuthenticated checks whether a request is authenticated.
func (s *Service) IsAuthenticated(request *http.Request) bool {
	isAuthenticated, ok := s.Get(consts.ValueIsAuthenticated, request).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
