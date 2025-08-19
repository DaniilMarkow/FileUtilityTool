package main 
import (
	"encoding/csv"
    "encoding/json"
    "fmt"   
    "os"
    "strconv" 
    "sort"
)

func readCSV(filePath string) ([]Country, error){
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        return nil, err
    }
    
    if fileInfo.Size() == 0 {
        return nil, fmt.Errorf("file is empty")
    }

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil{
		return nil, err
	}

	var countries []Country
	for i, record := range records {
		if i == 0{
			continue
		}

		if len(record) < 3{
			return nil, fmt.Errorf("invalid record at line %d", i+1)
		}

		name := record[0]

		population, err := strconv.Atoi(record[1])
		if err != nil{
			return nil, fmt.Errorf("invalid population at line %d: %v", i+1, err)
		}

		area, err := strconv.Atoi(record[2])
		if err != nil{
			return nil, fmt.Errorf("invalid population at line %d: %v", i+1, err)
		}

		countries = append(countries, Country{
			Name:	 name,
			Population: population,
			Area: area,
		})
	}
	return countries, nil
}

func readJSON(filePath string) ([]Country, error) {
    file, err := os.Open(filePath)
    if err != nil { return nil, err }
    defer file.Close()

    var countries []Country
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&countries)
    return countries, err
}
func writeCSV(filename string, countries []Country) error{
	file, err := os.Create(filename)
	if err != nil{
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Name", "Population", "Area"}); err!=nil{
		return err
	}

	for _, c := range countries{
		record := []string{
			c.Name, 
			strconv.Itoa(c.Population), 
			strconv.Itoa(c.Area),
		}
		if err := writer.Write(record); err != nil{
			return err
		}
	}
	return nil
}

func writeJSON(filename string, countries []Country) error{
	file, err := os.Create(filename)
	if err != nil{
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(countries); err != nil{
		return err
	}
	return nil
}

func filterCountries(countries []Country, field, value string) []Country {
    var filtered []Country
    for _, c := range countries {
        switch field {
        case "name":
            if c.Name == value {
                filtered = append(filtered, c)
            }
        case "population":
			val, err := strconv.Atoi(value)
            if err == nil && c.Population >= val {
                filtered = append(filtered, c)
            }
		case "area":
			val, err := strconv.Atoi(value)
			if err == nil && c.Area >= val{
				filtered = append(filtered, c)
			}
        }
		
    }
    return filtered
}

func sortCountries(countries []Country, field string) {
    sort.Slice(countries, func(i, j int) bool {
        switch field {
        case "name":
            return countries[i].Name > countries[j].Name
        case "population":
            return countries[i].Population > countries[j].Population
        case "area":
            return countries[i].Area > countries[j].Area
        default:
            return countries[i].Name > countries[j].Name
		}
    })
}

