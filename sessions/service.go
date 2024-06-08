package sessions

import (
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/shvdg-dev/base-pkg/database"
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

// Store stores a value in the session using the provided key and value.
func (s *Service) Store(request *http.Request, key string, value interface{}) {
	s.Manager.Put(request.Context(), key, value)
}

// Get retrieves the value associated with the given key from the session manager.
func (s *Service) Get(request *http.Request, key string) interface{} {
	return s.Manager.Get(request.Context(), key)
}
