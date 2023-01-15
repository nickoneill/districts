package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func checkExistingJSON() {
	shapeFiles := 0

	yearsDirs, err := ioutil.ReadDir("cds")
	if err != nil {
		log.Fatal(err)
	}

	for _, yearDir := range yearsDirs {
		if yearDir.Name() == "2012" {
			// ignore 2012, weird formats there
			continue
		}

		distDirs, err := ioutil.ReadDir(fmt.Sprintf("cds/%s", yearDir.Name()))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("starting", yearDir.Name())

		for _, distDir := range distDirs {
			if distDir.Name() == ".DS_Store" || yearDir.Name() == ".DS_Store" {
				continue
			}

			geojsonFilename := fmt.Sprintf("cds/%s/%s/shape.geojson", yearDir.Name(), distDir.Name())
			// should have some shape files in here
			if _, err := os.Stat(geojsonFilename); err != nil {
				log.Fatalf("error finding shape for %s", geojsonFilename)
			} else {
				shapeFiles++
				validateGeoJSON(geojsonFilename)
			}
		}
	}

	log.Printf("got %d shapefiles", shapeFiles)
}

func validateGeoJSON(filename string) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var geojson GeoJSON
	err = json.Unmarshal(fileBytes, &geojson)
	if err != nil {
		log.Fatal(err)
	}

	if geojson.Properties.District == "" {
		log.Fatalf("can't unmarshal %s", filename)
	}

	log.Printf("type: %s", geojson.Geometry.Type)

	if geojson.Geometry.Type == "Polygon" {
		var polygons [][][]float32
		err = json.Unmarshal(geojson.Geometry.Polygons, &polygons)
		if err != nil {
			log.Fatalf("can't unmarshal polygons at %s: %s", filename, err)
		}

		coords := 0
		for _, polygon := range polygons {
			coords = coords + len(polygon)
		}

		log.Printf("file %s has %d polys, %d total coords", geojson.Properties.District, len(polygons), coords)
	} else {
		var polygons [][][][]float32
		err = json.Unmarshal(geojson.Geometry.Polygons, &polygons)
		if err != nil {
			log.Fatalf("can't unmarshal multipolygons at %s: %s", filename, err)
		}

		coords := 0
		for _, subPolygon := range polygons {
			// yes, we have to go another level
			for _, polygon := range subPolygon {
				coords = coords + len(polygon)
			}
		}

		log.Printf("file %s has %d polys, %d total coords", geojson.Properties.District, len(polygons), coords)
	}
}

type GeoJSON struct {
	Geometry   GeoJSONGeom
	Type       string
	Properties GeoJSONProp
}

type GeoJSONGeom struct {
	Type     string
	Polygons json.RawMessage `json:"coordinates"`
}

type GeoJSONProp struct {
	Code     string
	District string
}
