package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// checks some district files in here for common stuff

func main() {
	var stateAbbreviation string
	flag.StringVar(&stateAbbreviation, "state", "", "a state abbreviation")
	flag.Parse()

	if stateAbbreviation == "" {
		log.Fatalf("needs a state with --state")
	}
	// validateGeoJSON("cds/2016/FL-14/shape.geojson")
	// checkExistingJSON()

	cleanGeoJSONDir(stateAbbreviation)
}

func cleanGeoJSONDir(state string) {
	// takes a pile of output jsons from mapshaper and:
	// * removes the outer FeatureCollection
	// * sets the correct properties based on precident in this repo
	// * creates a directory and files with the expected names

	jsonFiles, err := ioutil.ReadDir(state)
	if err != nil {
		log.Fatal(err)
	}

	cleanedJSONs := []FC{}
	for _, file := range jsonFiles {
		if file.Name() == ".DS_Store" {
			continue
		}
		filename := fmt.Sprintf("%s/%s", state, file.Name())
		log.Printf("looking for %s", filename)
		fileBytes, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("bytes are %v", fileBytes)

		var fcGeojson FC
		err = json.Unmarshal(fileBytes, &fcGeojson)
		if err != nil {
			log.Fatalf("can't unmarshal: %s", err)
		}

		fcGeojson.Type = strings.ToLower(fcGeojson.Type)
		fcGeojson.Type = strings.ToLower(fcGeojson.Type)

		components := strings.Split(fileNameWithoutExtTrimSuffix(file.Name()), "-")
		suffixDistrict, _ := strconv.Atoi(components[len(components)-1])
		// propertiesDistrict, _ := strconv.Atoi(fcGeojson.Features[0].Properties.District)
		propertiesDistrict := fcGeojson.Features[0].Properties.District
		altPropertiesDistrict := fcGeojson.Features[0].Properties.Districtno
		alt2PropertiesDistrict := fcGeojson.Features[0].Properties.DistrictI

		// log.Printf("s: %d, p: %d a: %d", suffixDistrict, propertiesDistrict, altPropertiesDistrict)
		district := suffixDistrict
		if suffixDistrict != propertiesDistrict && propertiesDistrict != 0 {
			log.Printf("district mismatch: suffix was %d but properties said %d (defaults to a property if not 0)", suffixDistrict, propertiesDistrict)
			district = propertiesDistrict
		} else if suffixDistrict != altPropertiesDistrict && altPropertiesDistrict != 0 {
			log.Printf("district mismatch: suffix was %d but altProperties said %d (defaults to a property if not 0)", suffixDistrict, altPropertiesDistrict)
			district = altPropertiesDistrict
		} else if suffixDistrict != alt2PropertiesDistrict && alt2PropertiesDistrict != 0 {
			log.Printf("district mismatch: suffix was %d but alt2Properties said %d (defaults to a property if not 0)", suffixDistrict, alt2PropertiesDistrict)
			district = alt2PropertiesDistrict
		} else {
			log.Printf("no property district number or it was correct, using suffix: %d", suffixDistrict)
		}

		// MI
		// 1 -> 13
		// 2 -> 12
		// 3 -> 11
		// 4 -> 3
		// 5 -> 7
		// 6 -> 10
		// 7 -> 6
		// 8 -> 5
		// 9 -> 4
		// 10 -> 9
		// 11 -> 8
		// 12 -> 2
		// 13 -> 1

		// fcGeojson.Features[0].Properties.Code = fmt.Sprintf("%s-%s", state, district)
		// fcGeojson.Features[0].Properties.District = 1 //"Florida 1st"
		fcGeojson.District = district

		cleanedJSONs = append(cleanedJSONs, fcGeojson)
	}

	_ = os.Mkdir(state, 0700)
	for _, cleanJSON := range cleanedJSONs {
		_ = os.Mkdir(fmt.Sprintf("%s/%s-%d", state, state, cleanJSON.District), 0700)

		var geoJSON WriteGeoJSON
		geoJSON.Type = cleanJSON.Features[0].Type
		geoJSON.Geometry = cleanJSON.Features[0].Geometry
		geoJSON.Properties.Code = fmt.Sprintf("%s-%02d", state, cleanJSON.District)
		geoJSON.Properties.District = fmt.Sprintf("%s %d", StateNameFromAbbr(state), cleanJSON.District)

		jsonBytes, err := json.Marshal(geoJSON)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(fmt.Sprintf("%s/%s-%d/shape.geojson", state, state, cleanJSON.District), jsonBytes, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type FC struct {
	Type     string
	Features []FCGeoJSON
	District int
}

type FCGeoJSON struct {
	Type       string
	Geometry   json.RawMessage
	Properties FCProperties
}

type FCProperties struct {
	District   int // many states use this
	Districtno int // MI special case
	DistrictI  int `json:"DISTRICT_I"` // LA special case
}

type WriteGeoJSON struct {
	Type       string          `json:"type"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties GeoJSONProp     `json:"properties"`
}

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
