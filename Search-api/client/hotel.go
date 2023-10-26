package client

import (
	solrClient "Search/solr"
	"github.com/rtt/Go-Solr"
	log "github.com/sirupsen/logrus"
)

type hotelClient struct{}

type hotelClientInterface interface {
	InsertUpdateHotel(document map[string]interface{}) error
	GetHotels() solr.Document
	GetHotelById(id string) solr.Document
	DeleteHotelById(id string) error
}

var SolrHotelClient hotelClientInterface

func init() {
	SolrHotelClient = &hotelClient{}
}

func (c hotelClient) InsertUpdateHotel(document map[string]interface{}) error {
	resp, err := solrClient.SolrClient.Update(document, true)
	if err != nil {
		return err
	}
	log.Printf("Solr Response: %s", resp.String())

	return nil
}

func (c hotelClient) GetHotels() solr.Document {
	return solr.Document{}

}

func (c hotelClient) GetHotelById(id string) solr.Document {
	return solr.Document{}
}

func (c hotelClient) DeleteHotelById(id string) error {
	return nil
}
