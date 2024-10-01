package process

import (
	"encoding/json"
	"fmt"
	"os"

	geonames "github.com/jasontconnell/geonames/data"
)

func Write(cities []geonames.City, outfile string) error {
	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't open file for writing %s. %w", outfile, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(cities)
}
