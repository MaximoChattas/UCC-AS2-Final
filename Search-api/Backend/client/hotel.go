package client

import (
	solrClient "Search/solr"
	"errors"
	"github.com/rtt/Go-Solr"
	log "github.com/sirupsen/logrus"
)

type hotelClient struct{}

type hotelClientInterface interface {
	UpdateHotel(document map[string]interface{}) error
	GetHotels() (*solr.DocumentCollection, error)
	GetHotelById(id string) (*solr.DocumentCollection, error)
	GetHotelsByCity(city string) (*solr.DocumentCollection, error)
}

var SolrHotelClient hotelClientInterface

func init() {
	SolrHotelClient = &hotelClient{}
}

func (c hotelClient) UpdateHotel(document map[string]interface{}) error {
	resp, err := solrClient.SolrClient.Update(document, true)
	if err != nil {
		return err
	}
	log.Printf("Solr Response: %s", resp.String())

	return nil
}

func (c hotelClient) GetHotels() (*solr.DocumentCollection, error) {

	q := "q=*:*"

	resp, err := solrClient.SolrClient.SelectRaw(q)

	if err != nil {
		log.Info(err)
		return &solr.DocumentCollection{}, err
	}

	result := resp.Results

	return result, nil

}

func (c hotelClient) GetHotelById(id string) (*solr.DocumentCollection, error) {

	q := solr.Query{
		Params: solr.URLParamMap{
			"q": []string{"id:" + id},
		},
	}

	resp, err := solrClient.SolrClient.Select(&q)

	if err != nil {
		return &solr.DocumentCollection{}, err
	}

	result := resp.Results

	if result.Len() == 0 {
		return &solr.DocumentCollection{}, errors.New("hotel not found")
	}

	return result, nil
}

func (c hotelClient) GetHotelsByCity(city string) (*solr.DocumentCollection, error) {

	q := solr.Query{
		Params: solr.URLParamMap{
			"q": []string{"city:" + city},
		},
	}

	resp, err := solrClient.SolrClient.Select(&q)

	if err != nil {
		return &solr.DocumentCollection{}, err
	}

	result := resp.Results

	return result, nil

}
