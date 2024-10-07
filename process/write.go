package process

import (
	"encoding/json"
	"fmt"
	"os"

	geonames "github.com/jasontconnell/geonames/data"
)

func WriteCities(cities []geonames.City, outfile string) error {
	return writeJsonFile(cities, outfile)
}

func WriteCountries(countries []geonames.Country, outfile string) error {
	return writeJsonFile(countries, outfile)
}

func writeJsonFile(obj interface{}, outfile string) error {
	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't open file for writing %s. %w", outfile, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(obj)
}
