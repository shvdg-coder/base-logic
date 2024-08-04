package pkg

import (
	"fmt"
	"github.com/google/uuid"
)

// StringsToUUIDs function maps each string to a UUID.
func StringsToUUIDs(strings ...string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(strings))
	for i, s := range strings {
		u, err := StringToUUID(s)
		if err != nil {
			return nil, err
		}
		uuids[i] = u
	}
	return uuids, nil
}

// StringToUUID function creates a UUID from a string.
func StringToUUID(str string) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse UUID: %w", err)
	}
	return id, nil
}
