package main

import (
	"flag"
	"log"
	"time"

	"github.com/jasontconnell/geocodecsv/process"
	geonames "github.com/jasontconnell/geonames/data"
)

func main() {
	start := time.Now()

	citiesfile := flag.String("cities", "", "cities file from geonames")
	countriesfile := flag.String("countries", "", "countries file from geonames")
	locationsfile := flag.String("locations", "", "csv file of locations (city, state, country)")
	modcitiesfile := flag.String("modcities", "", "modified locations file")
	addcitiesfile := flag.String("addcities", "", "added locations file")
	addcountriesfile := flag.String("addcountries", "", "added countries file")
	output := flag.String("out", "cities.json", "output filename")
	flag.Parse()

	if *citiesfile == "" || *countriesfile == "" {
		log.Fatal("can't continue, need both cities and countries files")
	}

	locs, err := process.ReadLocations(*locationsfile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("read", len(locs), "locations")

	var modlocs, addlocs []geonames.City
	var lerr error
	if *modcitiesfile != "" {
		modlocs, lerr = process.ReadJsonCities(*modcitiesfile)
		if lerr != nil {
			log.Fatal(lerr)
		}
	}

	if *addcitiesfile != "" {
		addlocs, lerr = process.ReadJsonCities(*addcitiesfile)
		if lerr != nil {
			log.Fatal(lerr)
		}
	}

	log.Println("read", len(modlocs), "modified cities")
	log.Println("read", len(addlocs), "added cities")

	cities, err := process.ReadCities(*citiesfile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("read", len(cities), "cities from geonames cities file")

	countries, err := process.ReadCountries(*countriesfile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("read", len(countries), "countries from geonames countries file")

	addcountries, err := process.ReadJsonCountries(*addcountriesfile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("read", len(addcountries), "countries from modified countries file")

	merged := process.MergeCities(cities, modlocs)
	merged = process.MergeCities(merged, addlocs)

	log.Println("total cities", len(merged))
	mergecountries := process.MergeCountries(countries, addcountries)
	for _, c := range mergecountries {
		log.Println("country", c.Name, c.Abbr, c.Abbr3)
	}

	log.Println("total countries", len(mergecountries))
	converted := process.Convert(locs, mergecountries)

	filtered := process.Filter(merged, converted)

	err = process.Write(filtered, *output)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("finished.", time.Since(start))
}
