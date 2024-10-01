package process

import (
	"strings"

	"github.com/jasontconnell/geocodecsv/data"
	geonames "github.com/jasontconnell/geonames/data"
)

func MergeCities(cities []geonames.City, mcities []geonames.City) []geonames.City {
	mlookup := make(map[string]geonames.City)
	found := make(map[string]bool)
	for _, c := range mcities {
		key := getKey(c.Name, c.State, c.Country)
		mlookup[key] = c
		found[key] = false
	}

	merged := []geonames.City{}
	for _, c := range cities {
		key := getKey(c.Name, c.State, c.Country)
		if mc, ok := mlookup[key]; ok {
			found[key] = true
			c.AlternateNames = append(c.AlternateNames, mc.AlternateNames...)
		}
		merged = append(merged, c)
	}

	for _, c := range mcities {
		k := getKey(c.Name, c.State, c.Country)
		if fnd, ok := found[k]; ok && !fnd {
			merged = append(merged, c)
		}
	}
	return merged
}

func Convert(list []data.Location, countries []geonames.Country) []geonames.City {
	converted := []geonames.City{}

	clookup := make(map[string]string)
	for _, c := range countries {
		clookup[c.Abbr3] = c.Abbr
	}

	for _, loc := range list {
		c2 := loc.Country
		if tc, ok := clookup[c2]; len(c2) == 3 && ok {
			c2 = tc
		}
		c := geonames.City{
			Name:    loc.City,
			State:   loc.State,
			Country: c2,
		}
		converted = append(converted, c)
	}

	return converted
}

type filteredCity struct {
	c      *geonames.City
	parent *geonames.City
}

func Filter(list []geonames.City, find []geonames.City) []geonames.City {
	allKeys := make(map[string]filteredCity)
	for _, c := range list {
		k := getKey(c.Name, c.State, c.Country)
		cp := geonames.City{Name: c.Name, Latitude: c.Latitude, Longitude: c.Longitude, Country: c.Country, State: c.State, TimeZone: c.TimeZone, AlternateNames: []string{}}
		allKeys[k] = filteredCity{c: &cp}
		for _, alt := range c.AlternateNames {
			k2 := getKey(alt, c.State, c.Country)
			if _, ok := allKeys[k2]; !ok {
				allKeys[k2] = filteredCity{c: &geonames.City{Name: alt, State: c.State, Country: c.Country}, parent: &cp}
			}
		}
	}

	findKeys := make(map[string]string)
	for _, c := range find {
		k := getKey(c.Name, c.State, c.Country)
		findKeys[k] = k
	}

	pfiltered := []*geonames.City{}
	for _, c := range find {
		k := getKey(c.Name, c.State, c.Country)
		if fc, ok := allKeys[k]; ok {
			if fc.parent == nil {
				pfiltered = append(pfiltered, fc.c)
			} else {
				p := fc.parent
				p.AlternateNames = append(p.AlternateNames, fc.c.Name)
			}
		}
	}
	filtered := []geonames.City{}
	for _, p := range pfiltered {
		for i := len(p.AlternateNames) - 1; i >= 0; i-- {
			a := p.AlternateNames[i]
			k2 := getKey(a, p.State, p.Country)
			if _, ok := findKeys[k2]; !ok {
				p.AlternateNames = append(p.AlternateNames[:i], p.AlternateNames[i+1:]...)
			}

		}
		filtered = append(filtered, *p)
	}
	return filtered
}

func MergeCountries(list []geonames.Country, mlist []geonames.Country) []geonames.Country {
	mlookup := make(map[string]geonames.Country)
	found := make(map[string]bool)
	for _, c := range mlist {
		mlookup[c.Abbr] = c
		found[c.Abbr] = false
	}

	merged := []geonames.Country{}
	for _, c := range list {
		if _, ok := mlookup[c.Abbr]; ok {
			found[c.Abbr] = true
		}
		merged = append(merged, c)
	}
	for _, c := range mlist {
		if fnd, ok := found[c.Abbr]; ok && !fnd {
			merged = append(merged, c)
		}
	}
	return merged
}

func getKey(city, state, country string) string {
	if country != "US" && country != "USA" {
		state = ""
	}
	key := strings.ToLower(city + "_" + state + "_" + country)
	return key
}
