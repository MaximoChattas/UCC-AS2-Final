package service

import (
	"Hotel/client"
	"Hotel/dto"
	"Hotel/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type testAmenity struct{}

func init() {
	client.AmenityClient = testAmenity{}
}

func (t testAmenity) InsertAmenity(amenity model.Amenity) model.Amenity {

	if amenity.Name != "" {
		objId, _ := primitive.ObjectIDFromHex("654cf68d807298d99186019f")
		amenity.Id = objId
	}

	return amenity
}

func (t testAmenity) GetAmenityById(id string) model.Amenity {

	if id == "000000000000000000000000" {
		return model.Amenity{}
	}

	objId, _ := primitive.ObjectIDFromHex("654cf68d807298d99186019f")

	return model.Amenity{
		Id: objId,
	}

}

func (t testAmenity) GetAmenityByName(name string) model.Amenity {

	if name == "" || name == "New Amenity" {
		return model.Amenity{}
	}

	objId, _ := primitive.ObjectIDFromHex("654cf68d807298d99186019f")

	return model.Amenity{
		Id:   objId,
		Name: name,
	}

}

func (t testAmenity) GetAmenities() model.Amenities {
	return model.Amenities{}
}

func (t testAmenity) DeleteAmenityById(id string) error {

	if id == "000000000000000000000000" {
		return errors.New("amenity not found")
	}

	return nil
}

func TestInsertAmenity_Error(t *testing.T) {

	a := assert.New(t)

	amenity := dto.AmenityDto{Id: "000000000000000000000000"}

	_, err := AmenityService.InsertAmenity(amenity)

	a.NotNil(err)

	expectedResponse := "error creating amenity"
	a.Equal(expectedResponse, err.Error())
}

func TestInsertAmenity_Exists(t *testing.T) {

	a := assert.New(t)

	amenity := dto.AmenityDto{
		Id:   "654cf68d807298d99186019f",
		Name: "Test Amenity",
	}

	_, err := AmenityService.InsertAmenity(amenity)

	a.NotNil(err)

	expectedResponse := "amenity already exists"
	a.Equal(expectedResponse, err.Error())
}

func TestInsertAmenity(t *testing.T) {

	a := assert.New(t)

	amenity := dto.AmenityDto{Name: "New Amenity"}

	amenityResponse, err := AmenityService.InsertAmenity(amenity)

	expectedId := "654cf68d807298d99186019f"

	a.Nil(err)
	a.Equal(expectedId, amenityResponse.Id)
}

func TestGetAmenities(t *testing.T) {

	a := assert.New(t)

	amenitiesResponse, err := AmenityService.GetAmenities()

	a.Nil(err)

	var emptyDto dto.AmenitiesDto
	a.Equal(emptyDto, amenitiesResponse)
}

func TestDeleteAmenityById_Error(t *testing.T) {

	a := assert.New(t)

	id := "000000000000000000000000"
	err := AmenityService.DeleteAmenityById(id)

	a.NotNil(err)

	expectedResponse := "failed to delete amenity"
	a.Equal(expectedResponse, err.Error())
}

func TestDeleteAmenityById(t *testing.T) {

	a := assert.New(t)

	id := "654cf68d807298d99186019f"
	err := AmenityService.DeleteAmenityById(id)

	a.Nil(err)
}
