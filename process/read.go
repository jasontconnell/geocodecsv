package process

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jasontconnell/geocodecsv/data"
	geonames "github.com/jasontconnell/geonames/data"
)

func ReadLocations(location string) ([]data.Location, error) {
	f, err := os.Open(location)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file %s. %w", location, err)
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	lines, err := rdr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing csv %s. %w", location, err)
	}
	locs := []data.Location{}
	for _, line := range lines {
		loc := data.Location{
			City:    line[0],
			State:   line[1],
			Country: line[2],
		}
		locs = append(locs, loc)
	}
	return locs, nil
}

func ReadCities(filename string) ([]geonames.City, error) {
	return geonames.ReadCities(filename)
}

func ReadJsonCities(filename string) ([]geonames.City, error) {
	return geonames.ReadJsonCities(filename)
}

func ReadCountries(filename string) ([]geonames.Country, error) {
	return geonames.ReadCountries(filename)
}

func ReadJsonCountries(filename string) ([]geonames.Country, error) {
	return geonames.ReadJsonCountries(filename)
}
