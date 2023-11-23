package service

import (
	"Search/client"
	"Search/dto"
	"errors"
	"github.com/rtt/Go-Solr"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testHotel struct{}

func init() {
	client.SolrHotelClient = testHotel{}
}

func (t testHotel) UpdateHotel(document map[string]interface{}) error {

	var id string

	add, exists := document["add"].([]interface{})

	if exists {
		id = add[0].(map[string]interface{})["id"].(string)
	} else {
		id = document["delete"].([]interface{})[0].(map[string]interface{})["id"].(string)
	}

	if id == "" {
		return errors.New("error updating hotel")
	}

	return nil
}

func (t testHotel) GetHotels() (*solr.DocumentCollection, error) {

	doc := solr.DocumentCollection{}

	doc.Collection = make([]solr.Document, 0)

	newDoc1 := solr.Document{
		Fields: map[string]interface{}{
			"id":            "1",
			"name":          []string{"Hotel Test 1"},
			"room_amount":   []int{10},
			"description":   []string{"Hotel Test Description 1"},
			"city":          []string{"Test City 1"},
			"street_name":   []string{"Test Street 1"},
			"street_number": []int{123},
			"rate":          []float64{4.5},
			"amenities":     []string{},
			"images":        []string{},
		},
	}

	newDoc2 := solr.Document{
		Fields: map[string]interface{}{
			"id":            "2",
			"name":          []string{"Hotel Test 2"},
			"room_amount":   []int{10},
			"description":   []string{"Hotel Test Description 2"},
			"city":          []string{"Test City 2"},
			"street_name":   []string{"Test Street 2"},
			"street_number": []int{123},
			"rate":          []float64{4.5},
			"amenities":     []string{},
			"images":        []string{},
		},
	}

	doc.Collection = append(doc.Collection, newDoc1)
	doc.Collection = append(doc.Collection, newDoc2)

	return &doc, nil
}

func (t testHotel) GetHotelById(id string) (*solr.DocumentCollection, error) {

	if id == "000000000000000000000000" {
		return &solr.DocumentCollection{}, errors.New("hotel not found")
	}

	doc := solr.DocumentCollection{}

	doc.Collection = make([]solr.Document, 0)

	newDoc := solr.Document{
		Fields: map[string]interface{}{
			"id":            id,
			"name":          []string{"Hotel Test"},
			"room_amount":   []int{10},
			"description":   []string{"Hotel Test Description"},
			"city":          []string{"Test City"},
			"street_name":   []string{"Test Street"},
			"street_number": []int{123},
			"rate":          []float64{4.5},
			"amenities":     []string{},
			"images":        []string{},
		},
	}

	doc.Collection = append(doc.Collection, newDoc)

	return &doc, nil
}

func (t testHotel) GetHotelsByCity(city string) (*solr.DocumentCollection, error) {

	if city == "" {
		return &solr.DocumentCollection{}, errors.New("city not found")
	}

	doc := solr.DocumentCollection{}
	doc.Collection = make([]solr.Document, 0)

	return &doc, nil
}

func TestInsertHotel_Error(t *testing.T) {

	a := assert.New(t)
	var hotelDto dto.HotelDto

	err := HotelService.InsertUpdateHotel(hotelDto)

	expectedResponse := "error updating hotel"

	a.NotNil(err)
	a.Equal(expectedResponse, err.Error())

}

func TestInsertHotel(t *testing.T) {

	a := assert.New(t)

	hotelDto := dto.HotelDto{
		Id:           "654cf68d807298d99186019f",
		Name:         "Hotel Test",
		RoomAmount:   10,
		Description:  "Hotel test description",
		City:         "Test City",
		StreetName:   "Test Street",
		StreetNumber: 123,
		Rate:         4.5,
	}

	err := HotelService.InsertUpdateHotel(hotelDto)
	a.Nil(err)

}

func TestGetHotelById_NotFound(t *testing.T) {

	a := assert.New(t)

	_, err := HotelService.GetHotelById("000000000000000000000000")

	expectedResponse := "hotel not found"
	a.NotNil(err)
	a.Equal(expectedResponse, err.Error())
}

func TestGetHotelById_Found(t *testing.T) {

	a := assert.New(t)

	id := "654cf68d807298d99186019f"

	hotelResponse, err := HotelService.GetHotelById(id)

	a.Nil(err)
	a.Equal(id, hotelResponse.Id)
}

func TestGetHotels(t *testing.T) {

	a := assert.New(t)

	hotelsResponse, err := HotelService.GetHotels()

	a.Nil(err)

	hotelsDto := dto.HotelsDto{
		dto.HotelDto{
			Id:           "1",
			Name:         "Hotel Test 1",
			RoomAmount:   10,
			Description:  "Hotel Test Description 1",
			City:         "Test City 1",
			StreetName:   "Test Street 1",
			StreetNumber: 123,
			Rate:         4.5,
			Amenities:    []string{},
			Images:       []string{},
		},

		dto.HotelDto{
			Id:           "2",
			Name:         "Hotel Test 2",
			RoomAmount:   10,
			Description:  "Hotel Test Description 2",
			City:         "Test City 2",
			StreetName:   "Test Street 2",
			StreetNumber: 123,
			Rate:         4.5,
			Amenities:    []string{},
			Images:       []string{},
		},
	}

	a.Equal(hotelsDto, hotelsResponse)
}

func TestGetHotelByCity_NotFound(t *testing.T) {

	a := assert.New(t)

	_, err := HotelService.GetHotelByCity("")

	expectedResponse := "city not found"
	a.NotNil(err)
	a.Equal(expectedResponse, err.Error())
}

func TestGetHotelByCity_Found(t *testing.T) {

	a := assert.New(t)

	city := "test"
	hotelsResponse, err := HotelService.GetHotelByCity(city)
	a.Nil(err)

	var emptyDto dto.HotelsDto
	a.Equal(emptyDto, hotelsResponse)
}

func TestDeleteHotelById_NotFound(t *testing.T) {

	a := assert.New(t)

	err := HotelService.DeleteHotelById("")

	expectedResponse := "error updating hotel"
	a.NotNil(err)
	a.Equal(expectedResponse, err.Error())
}

func TestDeleteHotelById_Found(t *testing.T) {

	a := assert.New(t)

	id := "654cf68d807298d99186019f"

	err := HotelService.DeleteHotelById(id)

	a.Nil(err)
}

func TestUnmarshallSolrResponse(t *testing.T) {

	a := assert.New(t)

	solrResponse := dto.SolrResponsesDto{
		dto.SolrResponseDto{
			Id:           "1",
			Name:         []string{"Test Hotel"},
			RoomAmount:   []int{10},
			Description:  []string{"Test Hotel description"},
			City:         []string{"City Test"},
			StreetName:   []string{"Street test"},
			StreetNumber: []int{123},
			Rate:         []float64{4.5},
		},

		dto.SolrResponseDto{
			Id:           "2",
			Name:         []string{"Test Hotel 2"},
			RoomAmount:   []int{10},
			Description:  []string{"Test Hotel description"},
			City:         []string{"City Test"},
			StreetName:   []string{"Street test"},
			StreetNumber: []int{123},
			Rate:         []float64{4.5},
		},
	}

	expectedResponse := dto.HotelsDto{
		dto.HotelDto{
			Id:           "1",
			Name:         "Test Hotel",
			RoomAmount:   10,
			Description:  "Test Hotel description",
			City:         "City Test",
			StreetName:   "Street test",
			StreetNumber: 123,
			Rate:         4.5,
		},

		dto.HotelDto{
			Id:           "2",
			Name:         "Test Hotel 2",
			RoomAmount:   10,
			Description:  "Test Hotel description",
			City:         "City Test",
			StreetName:   "Street test",
			StreetNumber: 123,
			Rate:         4.5,
		},
	}

	response := unmarshalSolrResponse(solrResponse)
	a.Equal(expectedResponse, response)
}
