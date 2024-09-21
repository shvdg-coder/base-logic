package pkg

import (
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"net/http"
)

// SessionManager is for managing sessions.
type SessionManager struct {
	Database *DbService
	Manager  *scs.SessionManager
}

// NewSessionManager creates a new instance of the SessionManager struct.
func NewSessionManager(database *DbService) *SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(database.DB())
	return &SessionManager{Database: database, Manager: sessionManager}
}

// Store stores a value in the session using the provided key and value.
func (s *SessionManager) Store(request *http.Request, key string, value interface{}) {
	s.Manager.Put(request.Context(), key, value)
}

// Get retrieves the value associated with the given key from the session manager.
func (s *SessionManager) Get(request *http.Request, key string) interface{} {
	return s.Manager.Get(request.Context(), key)
}
