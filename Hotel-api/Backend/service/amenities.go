package service

import (
	"Hotel/client"
	"Hotel/dto"
	"Hotel/model"
	"errors"
)

type amenityService struct{}

type amenityServiceInterface interface {
	InsertAmenity(amenityDto dto.AmenityDto) (dto.AmenityDto, error)
	GetAmenities() (dto.AmenitiesDto, error)
	DeleteAmenityById(id string) error
}

var AmenityService amenityServiceInterface

func init() {
	AmenityService = &amenityService{}
}

func (s *amenityService) InsertAmenity(amenityDto dto.AmenityDto) (dto.AmenityDto, error) {

	//Check if amenity already exists before inserting
	checkAmenity := client.AmenityClient.GetAmenityByName(amenityDto.Name)

	if checkAmenity.Id.Hex() == "000000000000000000000000" {
		var amenity model.Amenity

		amenity.Name = amenityDto.Name

		amenity = client.AmenityClient.InsertAmenity(amenity)

		if amenity.Id.Hex() == "000000000000000000000000" {
			return amenityDto, errors.New("error creating amenity")
		}

		amenityDto.Id = amenity.Id.Hex()

		return amenityDto, nil

	}

	return amenityDto, errors.New("amenity already exists")

}

func (s *amenityService) GetAmenities() (dto.AmenitiesDto, error) {
	var amenities model.Amenities = client.AmenityClient.GetAmenities()
	var amenitiesDto dto.AmenitiesDto

	for _, amenity := range amenities {
		var amenityDto dto.AmenityDto
		amenityDto.Id = amenity.Id.Hex()
		amenityDto.Name = amenity.Name

		amenitiesDto = append(amenitiesDto, amenityDto)
	}

	return amenitiesDto, nil
}

func (s *amenityService) DeleteAmenityById(id string) error {

	err := client.AmenityClient.DeleteAmenityById(id)

	if err != nil {
		return errors.New("failed to delete amenity")
	}

	return nil
}
