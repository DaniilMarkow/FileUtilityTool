package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	inputFile, outputFile, sortField, filterValue := parseFlags()

	if err := validateInput(*inputFile, *outputFile); err != nil {
		log.Fatal(err)
	}

	countries, err := processInputFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	countries = applyFiltersAndSorting(countries, *filterValue, *sortField)

	if err := outputResults(countries, *outputFile); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() (*string, *string, *string, *string) {
	inputFile := flag.String("input", "", "Path to input file")
	outputFile := flag.String("output", "", "Path to output file")
	sortField := flag.String("sort", "", "Field to sort by (e.g. population)")
	filterValue := flag.String("filter", "", "Filter value (e.g. area)")

	flag.Parse()
	return inputFile, outputFile, sortField, filterValue
}

func validateInput(inputFile, outputFile string) error {
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}

	if format := getFileFormat(inputFile); format == "" {
		return fmt.Errorf("unsupported input file format")
	}

	if outputFile != "" && getFileFormat(outputFile) == "" {
		return fmt.Errorf("unsupported output file format")
	}

	return nil
}

func processInputFile(inputFile string) ([]Country, error) {
	format := getFileFormat(inputFile)
	switch format {
	case "csv":
		return readCSV(inputFile)
	case "json":
		return readJSON(inputFile)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", format)
	}
}

func applyFiltersAndSorting(countries []Country, filterValue, sortField string) []Country {
	if filterValue != "" {
		parts := strings.Split(filterValue, "=")
		if len(parts) == 2 {
			countries = filterCountries(countries, parts[0], parts[1])
		}
	}

	if sortField != "" {
		sortCountries(countries, sortField)
	}

	return countries
}

func outputResults(countries []Country, outputFile string) error {
	outputFormat := "console"
	if outputFile != "" {
		outputFormat = getFileFormat(outputFile)
	}

	switch outputFormat {
	case "console":
		fmt.Println("Result:")
		for i, c := range countries {
			fmt.Printf("%d. %s: %d (%d)\n", i+1, c.Name, c.Population, c.Area)
		}
		return nil
	case "csv":
		return writeCSV(outputFile, countries)
	case "json":
		return writeJSON(outputFile, countries)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func getFileFormat(filename string) string {
	if strings.HasSuffix(filename, ".csv") {
		return "csv"
	} else if strings.HasSuffix(filename, ".json") {
		return "json"
	}
	return ""
}

func validateCountry(c Country) error {
	if c.Name == "" {
		return fmt.Errorf("country name cannot be empty")
	}
	if c.Population < 0 {
		return fmt.Errorf("population cannot be negative")
	}
	if c.Area < 0 {
		return fmt.Errorf("area cannot be negative")
	}
	return nil
}
