package service

import (
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testAmadeusMap struct{}

func init() {
	client.AmadeusClient = &testAmadeusMap{}
	AmadeusService = &amadeusService{}
}

func (t testAmadeusMap) InsertAmadeusMap(mapping model.AmadeusMap) model.AmadeusMap {
	if mapping.HotelId == "" {
		return model.AmadeusMap{}
	}

	return mapping
}

func (t testAmadeusMap) GetAmadeusIdByHotelId(hotelId string) model.AmadeusMap {

	if hotelId == "" {
		return model.AmadeusMap{}
	}

	mapping := model.AmadeusMap{
		HotelId:   hotelId,
		AmadeusId: "SBMIASOF",
	}

	return mapping
}

func TestInsertAmadeusMap_Error(t *testing.T) {
	a := assert.New(t)

	body := dto.AmadeusMapDto{}

	_, err := AmadeusService.InsertAmadeusMap(body)

	a.NotNil(err)

	expectedResponse := "error creating mapping"
	a.Equal(expectedResponse, err.Error())

}

func TestInsertAmadeusMap(t *testing.T) {
	a := assert.New(t)

	hotelId := "654cf68d807298d99186019f"
	amadeusId := "SBMIASOF"

	body := dto.AmadeusMapDto{
		HotelId:   hotelId,
		AmadeusId: amadeusId,
	}

	response, err := AmadeusService.InsertAmadeusMap(body)

	a.Nil(err)
	a.Equal(body, response)
}

func TestGetAmadeusIdByHotelId_NotFound(t *testing.T) {

	a := assert.New(t)

	hotelId := ""
	_, err := AmadeusService.GetAmadeusIdByHotelId(hotelId)

	a.NotNil(err)

	expectedResponse := "no amadeus id set"
	a.Equal(expectedResponse, err.Error())

}

func TestGetAmadeusIdByHotelId_Found(t *testing.T) {

	a := assert.New(t)

	hotelId := "654cf68d807298d99186019f"
	response, err := AmadeusService.GetAmadeusIdByHotelId(hotelId)

	a.Nil(err)

	expectedResponse := dto.AmadeusMapDto{
		HotelId:   hotelId,
		AmadeusId: "SBMIASOF",
	}
	a.Equal(expectedResponse, response)

}

func TestGetAmadeusAvailability(t *testing.T) {

	a := assert.New(t)

	amadeusId := "SBMIASOF"
	startDate := time.Now().Add(72 * time.Hour)
	endDate := time.Now().Add(96 * time.Hour)

	available, err := AmadeusService.GetAmadeusAvailability(amadeusId, startDate, endDate)
	a.Nil(err)
	a.True(available)
}
