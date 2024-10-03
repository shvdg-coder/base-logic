package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
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

// UUIDsToStrings function maps each UUID to a string.
func UUIDsToStrings(uuids ...uuid.UUID) ([]string, error) {
	str := make([]string, len(uuids))
	for i, u := range uuids {
		s, err := UUIDToString(u)
		if err != nil {
			return nil, err
		}
		str[i] = s
	}
	return str, nil
}

// UUIDToString function creates a string from a UUID.
func UUIDToString(id uuid.UUID) (string, error) {
	if id == uuid.Nil {
		return "", fmt.Errorf("UUID is nil")
	}
	return id.String(), nil
}

// GetFields returns a slice of interfaces containing values of any struct.
func GetFields(s interface{}) []interface{} {
	v := reflect.ValueOf(s).Elem()
	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	return values
}

// GetFieldNames returns the field names of any struct.
func GetFieldNames(tag string, s interface{}) []string {
	t := reflect.TypeOf(s)
	fieldNames := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = t.Field(i).Tag.Get(tag)
	}
	return fieldNames
}
