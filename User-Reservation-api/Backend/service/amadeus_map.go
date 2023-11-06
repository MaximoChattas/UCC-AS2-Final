package service

import (
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"errors"
)

type amadeusService struct{}

type amadeusServiceInterface interface {
	InsertAmadeusMap(amadeusMapDto dto.AmadeusMapDto) (dto.AmadeusMapDto, error)
	GetAmadeusIdByHotelId(hotelId string) (dto.AmadeusMapDto, error)
}

var AmadeusService amadeusServiceInterface

func init() {
	AmadeusService = &amadeusService{}
}

func (s *amadeusService) InsertAmadeusMap(amadeusMapDto dto.AmadeusMapDto) (dto.AmadeusMapDto, error) {
	var mapping model.AmadeusMap

	mapping.HotelId = amadeusMapDto.HotelId
	mapping.AmadeusId = amadeusMapDto.AmadeusId

	mapping = client.InsertAmadeusMap(mapping)

	if mapping.HotelId == "" {
		return amadeusMapDto, errors.New("error creating mapping")
	}

	return amadeusMapDto, nil
}

func (s *amadeusService) GetAmadeusIdByHotelId(hotelId string) (dto.AmadeusMapDto, error) {

	var mapping model.AmadeusMap = client.GetAmadeusIdByHotelId(hotelId)

	var amadeusMapDto dto.AmadeusMapDto

	if mapping.HotelId == "" {
		return amadeusMapDto, errors.New("hotel not found")
	}

	amadeusMapDto.HotelId = mapping.HotelId
	amadeusMapDto.AmadeusId = mapping.AmadeusId

	return amadeusMapDto, nil
}
