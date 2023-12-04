package service

import (
	"Hotel/client"
	"Hotel/dto"
	"Hotel/model"
	"Hotel/queue"
	"encoding/json"
	"errors"
)

type hotelService struct{}

type hotelServiceInterface interface {
	GetHotelById(id string) (dto.HotelDto, error)
	GetHotels() (dto.HotelsDto, error)
	InsertHotel(hotelDto dto.HotelDto) (dto.HotelDto, error)
	DeleteHotel(id string) error
	UpdateHotel(hotelDto dto.HotelDto) (dto.HotelDto, error)
}

var HotelService hotelServiceInterface

func init() {
	HotelService = &hotelService{}
}

func (s *hotelService) InsertHotel(hotelDto dto.HotelDto) (dto.HotelDto, error) {
	var hotel model.Hotel

	hotel.Name = hotelDto.Name
	hotel.Description = hotelDto.Description
	hotel.RoomAmount = hotelDto.RoomAmount
	hotel.City = hotelDto.City
	hotel.StreetName = hotelDto.StreetName
	hotel.StreetNumber = hotelDto.StreetNumber
	hotel.Rate = hotelDto.Rate
	hotel.Images = hotelDto.Images

	for _, amenityName := range hotelDto.Amenities {
		amenity := client.AmenityClient.GetAmenityByName(amenityName)

		if amenity.Id.Hex() == "000000000000000000000000" {
			return hotelDto, errors.New("amenity not found")
		}

		hotel.Amenities = append(hotel.Amenities, amenityName)

	}

	hotel = client.HotelClient.InsertHotel(hotel)

	hotelDto.Id = hotel.Id.Hex()

	if hotel.Id.Hex() == "000000000000000000000000" {
		return hotelDto, errors.New("error creating hotel")
	}

	body := map[string]interface{}{
		"Id":      hotel.Id.Hex(),
		"Message": "create",
	}

	jsonBody, _ := json.Marshal(body)

	err := queue.QueueProducer.Publish(jsonBody)

	if err != nil {
		return hotelDto, err
	}

	return hotelDto, nil
}

func (s *hotelService) GetHotels() (dto.HotelsDto, error) {

	var hotels model.Hotels = client.HotelClient.GetHotels()
	var hotelsDto dto.HotelsDto

	for _, hotel := range hotels {
		var hotelDto dto.HotelDto
		hotelDto.Id = hotel.Id.Hex()
		hotelDto.Name = hotel.Name
		hotelDto.RoomAmount = hotel.RoomAmount
		hotelDto.City = hotel.City
		hotelDto.Description = hotel.Description
		hotelDto.StreetName = hotel.StreetName
		hotelDto.StreetNumber = hotel.StreetNumber
		hotelDto.Rate = hotel.Rate
		hotelDto.Amenities = hotel.Amenities
		hotelDto.Images = hotel.Images

		hotelsDto = append(hotelsDto, hotelDto)
	}

	return hotelsDto, nil
}

func (s *hotelService) GetHotelById(id string) (dto.HotelDto, error) {

	var hotelDto dto.HotelDto

	hotel := client.HotelClient.GetHotelById(id)

	if hotel.Id.Hex() == "000000000000000000000000" {
		return hotelDto, errors.New("hotel not found")
	}
	hotelDto.Id = hotel.Id.Hex()
	hotelDto.Name = hotel.Name
	hotelDto.RoomAmount = hotel.RoomAmount
	hotelDto.Description = hotel.Description
	hotelDto.City = hotel.City
	hotelDto.StreetName = hotel.StreetName
	hotelDto.StreetNumber = hotel.StreetNumber
	hotelDto.Rate = hotel.Rate
	hotelDto.Images = hotel.Images

	for _, amenityName := range hotel.Amenities {
		amenity := client.AmenityClient.GetAmenityByName(amenityName)

		if amenity.Id.Hex() == "000000000000000000000000" {
			return hotelDto, errors.New("amenity not found")
		}

		hotelDto.Amenities = append(hotelDto.Amenities, amenityName)
	}

	return hotelDto, nil
}

func (s *hotelService) DeleteHotel(id string) error {

	hotel := client.HotelClient.GetHotelById(id)

	if hotel.Id.Hex() == "000000000000000000000000" {
		return errors.New("hotel not found")
	}

	err := client.HotelClient.DeleteHotelById(id)

	if err != nil {
		return err
	}

	body := map[string]interface{}{
		"id":      hotel.Id.Hex(),
		"message": "delete",
	}

	jsonBody, _ := json.Marshal(body)

	err = queue.QueueProducer.Publish(jsonBody)

	if err != nil {
		return err
	}

	return err
}

func (s *hotelService) UpdateHotel(hotelDto dto.HotelDto) (dto.HotelDto, error) {

	hotel := client.HotelClient.GetHotelById(hotelDto.Id)

	if hotel.Id.Hex() == "000000000000000000000000" {
		return hotelDto, errors.New("hotel not found")
	}

	hotel.Name = hotelDto.Name
	hotel.City = hotelDto.City
	hotel.StreetName = hotelDto.StreetName
	hotel.StreetNumber = hotelDto.StreetNumber
	hotel.Rate = hotelDto.Rate
	hotel.Description = hotelDto.Description
	hotel.RoomAmount = hotelDto.RoomAmount
	hotel.Amenities = []string{}

	for _, amenityName := range hotelDto.Amenities {
		amenity := client.AmenityClient.GetAmenityByName(amenityName)

		if amenity.Id.Hex() == "000000000000000000000000" {
			return hotelDto, errors.New("amenity not found")
		}

		hotel.Amenities = append(hotel.Amenities, amenityName)
	}

	for _, image := range hotelDto.Images {
		hotel.Images = append(hotel.Images, image)
	}

	hotelDto.Images = hotel.Images

	hotel = client.HotelClient.UpdateHotelById(hotel)

	if hotel.Id.Hex() == "000000000000000000000000" {
		return hotelDto, errors.New("error updating hotel")
	}

	body := map[string]interface{}{
		"id":      hotel.Id.Hex(),
		"message": "update",
	}

	jsonBody, _ := json.Marshal(body)

	err := queue.QueueProducer.Publish(jsonBody)

	if err != nil {
		return hotelDto, err
	}

	return hotelDto, nil

}
