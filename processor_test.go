package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCountries(t *testing.T) {
	tests := []struct {
		name     string
		input    []Country
		field    string
		value    string
		expected int
	}{
		{
			"Population filter",
			[]Country{
				{Name: "CountryA", Population: 100},
				{Name: "CountryB", Population: 200},
			},
			"population",
			"150",
			1,
		},
		{
			"Area filter",
			[]Country{
				{Name: "CountryA", Area: 1000},
				{Name: "CountryB", Area: 2000},
			},
			"area",
			"1500",
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterCountries(tt.input, tt.field, tt.value)
			if len(filtered) != tt.expected {
				t.Errorf("Expected %d countries, got %d", tt.expected, len(filtered))
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	path := filepath.Join("testdata", "test1.csv")
	countries, err := readCSV(path)
	if err != nil {
		t.Fatalf("readCSV failed: %v", err)
	}

	expected := []Country{
		{Name: "Russia", Population: 143400000, Area: 17075200},
		{Name: "USA", Population: 295700000, Area: 9629091},
		{Name: "France", Population: 68860000, Area: 549190},
		{Name: "Japan", Population: 127400000, Area: 377835},
	}

	assert.Equal(t, expected, countries, "readCSV failed")
}

func TestReadJSON(t *testing.T) {
	path := filepath.Join("testdata", "test2.json")
	countries, err := readJSON(path)
	if err != nil {
		t.Fatalf("readJSON failed: %v", err)
	}
	expected := []Country{
		{Name: "Russia", Population: 143400000, Area: 17075200},
		{Name: "Japan", Population: 127400000, Area: 377835},
	}

	assert.Equal(t, expected, countries, "readCSV failed")
}

func TestWriteCSV(t *testing.T) {
	countries := []Country{
		{Name: "Test", Population: 100, Area: 200},
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "output.csv")
	if err := writeCSV(tmpFile, countries); err != nil {
		t.Fatalf("writeCSV failed: %v", err)
	}

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatalf("File not created")
	}

	result, err := readCSV(tmpFile)
	if err != nil {
		t.Fatalf("readCSV failed: %v", err)
	}

	assert.Equal(t, result, countries, "writeCSV failed")
}

func TestWriteJSON(t *testing.T) {
	countries := []Country{
		{Name: "Russia", Population: 143400000, Area: 17075200},
		{Name: "Japan", Population: 127400000, Area: 377835},
	}
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "output.json")

	if err := writeJSON(tmpFile, countries); err != nil {
		t.Fatalf("writeCSV failed: %v", err)
	}
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatalf("File not created")
	}
	result, err := readJSON(tmpFile)
	if err != nil {
		t.Fatalf("readCSV failed: %v", err)
	}

	assert.Equal(t, result, countries, "writeCSV failed")
}

func TestFilterCountries_2(t *testing.T) {
	countries := []Country{
		{Name: "Russia", Population: 143400000, Area: 17075200},
		{Name: "Japan", Population: 127400000, Area: 377835},
	}
	expected := []Country{
		{Name: "Russia", Population: 143400000, Area: 17075200},
	}
	result := filterCountries(countries, "area", "10000000")

	assert.Equal(t, result, expected, "filterCountries failed")
}

func TestSortCountriesPopulation(t *testing.T) {
	countries := []Country{
		{Name: "Russia", Population: 143400000, Area: 17075200},
		{Name: "USA", Population: 295700000, Area: 9629091},
		{Name: "France", Population: 68860000, Area: 549190},
		{Name: "Japan", Population: 127400000, Area: 377835},
	}
	expected := []Country{
		{Name: "USA", Population: 295700000, Area: 9629091},
		{Name: "Russia", Population: 143400000, Area: 17075200},
		{Name: "Japan", Population: 127400000, Area: 377835},
		{Name: "France", Population: 68860000, Area: 549190},
	}
	sortCountries(countries, "population")
	assert.Equal(t, countries, expected, "sortCountries failed")
}

func TestEmptyFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "empty.csv")
	os.WriteFile(tmpFile, []byte(""), 0644)

	_, err := readCSV(tmpFile)
	if err == nil {
		t.Error("Expected error for empty file")
	}
}

func TestInvalidCSV(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "invalid.csv")
	os.WriteFile(tmpFile, []byte("Name,Population\nvalue1,not_number"), 0644)

	_, err := readCSV(tmpFile)
	if err == nil {
		t.Error("Expected error for invalid CSV")
	}
}
