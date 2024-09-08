package csv

import (
	"github.com/shvdg-coder/base-logic/pkg"
	"testing"
)

// TestGetValueFromColumnsOfCSV verifies whether values from column can be retrieved of a .csv file.
func TestGetValueFromColumnsOfCSV(t *testing.T) {
	tests := []struct {
		columnIndex int
		columnName  string
	}{
		{0, contactsColumnID},
		{1, contactsColumnName},
	}

	for _, tt := range tests {
		t.Run(tt.columnName, func(t *testing.T) {
			records, err := pkg.GetCSVRecords(contactsCSVPath, false)
			if err != nil {
				t.Fatal(err)
			}

			columnValues, err := pkg.GetCSVColumnValues(contactsCSVPath, tt.columnName)
			if err != nil {
				t.Fatal(err)
			}

			if len(records) != len(columnValues) {
				t.Fatalf("expected %d values, got %d", len(records), len(columnValues))
			}

			for i, record := range records {
				if record[tt.columnIndex] != columnValues[i] {
					t.Fatalf("expected %s, got %s", columnValues[i], record[tt.columnIndex])
				}
			}
		})
	}
}
