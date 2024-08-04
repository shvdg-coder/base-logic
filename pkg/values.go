package pkg

import (
	"fmt"
	"github.com/google/uuid"
)

// stringsToUUIDs function maps each string to a UUID.
func stringsToUUIDs(strings ...string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(strings))
	for i, s := range strings {
		u, err := stringToUUID(s)
		if err != nil {
			return nil, err
		}
		uuids[i] = u
	}
	return uuids, nil
}

// stringToUUID function creates a UUID from a string.
func stringToUUID(str string) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse UUID: %w", err)
	}
	return id, nil
}
