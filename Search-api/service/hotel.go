package service

import (
	"Search/client"
	"Search/dto"
	log "github.com/sirupsen/logrus"
)

type hotelService struct{}

type hotelServiceInterface interface {
	InsertUpdateHotel(hotelDto dto.HotelDto) error
	GetHotels() dto.HotelsDto
	GetHotelById(id string) dto.HotelDto
	DeleteHotelById(id string) error
}

var HotelService hotelServiceInterface

func init() {
	HotelService = &hotelService{}
}

func (s hotelService) InsertUpdateHotel(hotelDto dto.HotelDto) error {

	document := map[string]interface{}{
		"add": []interface{}{
			map[string]interface{}{
				"id":            hotelDto.Id,
				"name":          hotelDto.Name,
				"room_amount":   hotelDto.RoomAmount,
				"description":   hotelDto.Description,
				"city":          hotelDto.City,
				"street_name":   hotelDto.StreetName,
				"street_number": hotelDto.StreetNumber,
				"rate":          hotelDto.Rate,
				"amenities":     hotelDto.Amenities,
				"images":        hotelDto.Images,
			},
		},
	}

	err := client.SolrHotelClient.InsertUpdateHotel(document)

	if err != nil {
		log.Info("Error inserting hotel", err)
		return err
	}

	return nil
}

func (s hotelService) GetHotels() dto.HotelsDto {
	return dto.HotelsDto{}
}

func (s hotelService) GetHotelById(id string) dto.HotelDto {
	return dto.HotelDto{}
}

func (s hotelService) DeleteHotelById(id string) error {
	return nil
}
