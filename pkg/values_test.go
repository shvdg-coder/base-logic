package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStringsToUUIDs tests whether strings can be turned into UUIDs.
func TestStringsToUUIDs(t *testing.T) {
	tests := []struct {
		name        string
		strings     []string
		expectedErr bool
	}{
		{
			name:        "Valid UUID strings",
			strings:     []string{"3b1f8bab-5c5b-468c-8a2a-0270f1db6e7e", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
			expectedErr: false,
		},
		{
			name:        "Invalid UUID string",
			strings:     []string{"invalid-uuid-string"},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := stringsToUUIDs(tt.strings...)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}
