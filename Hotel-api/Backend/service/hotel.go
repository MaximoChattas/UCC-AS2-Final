package service

import (
	"Hotel/client"
	"Hotel/dto"
	"Hotel/model"
	"errors"
)

type hotelService struct{}

type hotelServiceInterface interface {
	GetHotelById(id string) (dto.HotelDto, error)
	GetHotels() (dto.HotelsDto, error)
	InsertHotel(hotelDto dto.HotelDto) (dto.HotelDto, error)
	//DeleteHotel(id int) error
	//UpdateHotel(hotelDto dto.HotelDto) (dto.HotelDto, error)
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
	hotel.StreetName = hotelDto.StreetName
	hotel.StreetNumber = hotelDto.StreetNumber
	hotel.Rate = hotelDto.Rate

	for _, amenityName := range hotelDto.Amenities {
		amenity := client.AmenityClient.GetAmenityByName(amenityName)

		if amenity.Id.Hex() == "000000000000000000000000" {
			return hotelDto, errors.New("amenity not found")
		}

		hotel.Amenities = append(hotel.Amenities, amenity)

	}

	hotel = client.HotelClient.InsertHotel(hotel)

	hotelDto.Id = hotel.Id.Hex()

	if hotel.Id.Hex() == "000000000000000000000000" {
		return hotelDto, errors.New("error creating hotel")
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
		hotelDto.Description = hotel.Description
		hotelDto.StreetName = hotel.StreetName
		hotelDto.StreetNumber = hotel.StreetNumber
		hotelDto.Rate = hotel.Rate

		//if len(hotel.Images) > 0 {
		//	var imageDto dto.ImageDto
		//	imageDto.Id = hotel.Images[0].Id
		//	imageDto.Path = hotel.Images[0].Path
		//	imageDto.HotelId = hotel.Images[0].HotelId
		//
		//	hotelDto.Images = append(hotelDto.Images, imageDto)
		//}

		hotelsDto = append(hotelsDto, hotelDto)
	}

	return hotelsDto, nil
}

func (s *hotelService) GetHotelById(id string) (dto.HotelDto, error) {

	var hotel model.Hotel = client.HotelClient.GetHotelById(id)
	var hotelDto dto.HotelDto

	if hotel.Id.Hex() == "000000000000000000000000" {
		return hotelDto, errors.New("hotel not found")
	}
	hotelDto.Id = hotel.Id.Hex()
	hotelDto.Name = hotel.Name
	hotelDto.RoomAmount = hotel.RoomAmount
	hotelDto.Description = hotel.Description
	hotelDto.StreetName = hotel.StreetName
	hotelDto.StreetNumber = hotel.StreetNumber
	hotelDto.Rate = hotel.Rate

	for _, amenity := range hotel.Amenities {
		hotelDto.Amenities = append(hotelDto.Amenities, amenity.Name)
	}

	//for _, image := range hotel.Images {
	//	var imageDto dto.ImageDto
	//	imageDto.Id = image.Id
	//	imageDto.Path = image.Path
	//	imageDto.HotelId = image.HotelId
	//
	//	hotelDto.Images = append(hotelDto.Images, imageDto)
	//}

	return hotelDto, nil
}

//func (s *hotelService) DeleteHotel(id int) error {
//
//	hotel := client.AmenityClient.GetHotelById(id)
//
//	if hotel.Id == 0 {
//		return errors.New("hotel not found")
//	}
//
//	err := client.AmenityClient.DeleteHotel(hotel)
//
//	return err
//}
//
//func (s *hotelService) UpdateHotel(hotelDto dto.HotelDto) (dto.HotelDto, error) {
//
//	hotel := client.AmenityClient.GetHotelById(hotelDto.Id)
//
//	if hotel.Id == 0 {
//		return hotelDto, errors.New("hotel not found")
//	}
//
//	hotel.Name = hotelDto.Name
//	hotel.StreetName = hotelDto.StreetName
//	hotel.StreetNumber = hotelDto.StreetNumber
//	hotel.Rate = hotelDto.Rate
//	hotel.Description = hotelDto.Description
//	hotel.RoomAmount = hotelDto.RoomAmount
//	hotel.Amenities = model.Amenities{}
//
//	for _, amenityName := range hotelDto.Amenities {
//		amenity := client.GetAmenityByName(amenityName)
//
//		if amenity.Id == 0 {
//			return hotelDto, errors.New("amenity not found")
//		}
//
//		hotel.Amenities = append(hotel.Amenities, amenity)
//	}
//
//	hotel = client.AmenityClient.UpdateHotel(hotel)
//
//	if hotel.Id == 0 {
//		return hotelDto, errors.New("error updating hotel")
//	}
//
//	return hotelDto, nil
//
//}
