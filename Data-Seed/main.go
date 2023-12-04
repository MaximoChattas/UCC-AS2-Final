package main

import (
	"Data-Seed/dto"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func main() {

	insertAmenities()
	insertHotels()
	insertAmadeusMap()
}

func insertAmenities() {

	filePath := "./amenity-data.json"
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorf("error reading file: %s", err.Error())
		return
	}

	var amenities []json.RawMessage
	err = json.Unmarshal(file, &amenities)
	if err != nil {
		log.Errorf("error unmarshaling JSON: %s", err.Error())
	}

	for _, amenity := range amenities {

		resp, err := http.Post("http://localhost:8080/amenity", "application/json", bytes.NewBuffer(amenity))
		if err != nil {
			log.Errorf("error sending request: %s", err.Error())
			return
		}

		if resp.StatusCode != http.StatusCreated {
			log.Error("error posting amenity")
		}
	}
}

func insertHotels() {

	filePath := "./hotel-data.json"
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorf("error reading file: %s", err.Error())
		return
	}

	var amenities []json.RawMessage
	err = json.Unmarshal(file, &amenities)
	if err != nil {
		log.Errorf("error unmarshaling JSON: %s", err.Error())
	}

	for _, amenity := range amenities {

		resp, err := http.Post("http://localhost:8080/hotel", "application/json", bytes.NewBuffer(amenity))
		if err != nil {
			log.Error("error sending request")
			return
		}

		if resp.StatusCode != http.StatusCreated {
			log.Error("error posting hotel")
		}
	}
}

func insertAmadeusMap() {

	resp, err := http.Get("http://localhost:8080/hotel")

	if err != nil {
		log.Error("Error in HTTP request: ", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error("Error reading response: ", err)
		return
	}

	var hotelsDto dto.HotelsDto

	err = json.Unmarshal(body, &hotelsDto)

	if err != nil {
		log.Error("Error parsing JSON: ", err)
		return
	}

	for _, hotel := range hotelsDto {

		hotelBody := fmt.Sprintf(`{"hotel_id": "%s", "amadeus_id": "SBMIASOF"}`, hotel.Id)

		resp, err = http.Post("http://localhost:8090/amadeus", "application/json", bytes.NewBuffer([]byte(hotelBody)))
		if err != nil {
			log.Errorf("error sending request: %s", err.Error())
			return
		}

		if resp.StatusCode != http.StatusCreated {
			log.Error("error posting amadeus map")
		}
	}
}
