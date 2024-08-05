package pkg

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
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
			_, err := StringsToUUIDs(tt.strings...)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

// TestUUIDsToStrings tests whether UUIDs can be turned into strings.
func TestUUIDsToStrings(t *testing.T) {
	tests := []struct {
		name    string
		uuids   []uuid.UUID
		want    []string
		wantErr bool
	}{
		{
			name:    "Empty slice",
			uuids:   []uuid.UUID{},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "Single element",
			uuids:   []uuid.UUID{uuid.Must(uuid.Parse("550e8400-e29b-41d4-a716-446655440000"))},
			want:    []string{"550e8400-e29b-41d4-a716-446655440000"},
			wantErr: false,
		},
		{
			name: "Multiple elements",
			uuids: []uuid.UUID{
				uuid.Must(uuid.Parse("12345678-1234-5678-1234-567812345678")),
				uuid.Must(uuid.Parse("87654321-1234-5678-1234-567812345678")),
			},
			want: []string{
				"12345678-1234-5678-1234-567812345678",
				"87654321-1234-5678-1234-567812345678",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UUIDsToStrings(tt.uuids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("UUIDsToStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UUIDsToStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
