package solr

import (
	"Search/dto"
	"Search/service"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var SolrClient *solr.Connection

func InitSolr() {

	var err error

	SolrClient, err = solr.Init("solr", 8983, "hotels")
	if err != nil {
		log.Info("Failed to connect to Solr")
		log.Fatal(err)
	} else {
		log.Info("Connected to Solr successfully")
	}

	insertData()
}

func insertData() {

	resp, err := http.Get("http://hotel:8080/hotel")

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

		err = service.HotelService.InsertUpdateHotel(hotel)
	}

	if err != nil {
		log.Error(err)
		return
	}

}
