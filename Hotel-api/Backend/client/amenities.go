package client

import (
	db "Hotel/db"
	"Hotel/model"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type amenityClient struct{}

type amenityClientInterface interface {
	InsertAmenity(amenity model.Amenity) model.Amenity
	GetAmenityById(id string) model.Amenity
	GetAmenityByName(name string) model.Amenity
	GetAmenities() model.Amenities
	DeleteAmenityById(id string) error
}

var AmenityClient amenityClientInterface

func init() {
	AmenityClient = &amenityClient{}
}

func (c amenityClient) InsertAmenity(amenity model.Amenity) model.Amenity {

	insertAmenity := amenity
	insertAmenity.Id = primitive.NewObjectID()

	_, err := db.AmenitiesCollection.InsertOne(context.TODO(), &insertAmenity)

	if err != nil {
		fmt.Println(err)
		return amenity
	}

	amenity.Id = insertAmenity.Id

	return amenity
}

func (c amenityClient) GetAmenityById(id string) model.Amenity {
	var amenity model.Amenity

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		fmt.Println(err)
		return amenity
	}

	err = db.AmenitiesCollection.FindOne(context.TODO(), bson.D{{"_id", objID}}).Decode(&amenity)
	if err != nil {
		fmt.Println(err)
		return amenity
	}
	return amenity
}

func (c amenityClient) GetAmenityByName(name string) model.Amenity {
	var amenity model.Amenity

	err := db.AmenitiesCollection.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&amenity)

	if err != nil {
		fmt.Println(err)
		return amenity
	}
	return amenity
}

func (c amenityClient) GetAmenities() model.Amenities {
	var amenities model.Amenities

	cursor, err := db.AmenitiesCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Println(err)
		return amenities
	}

	err = cursor.All(context.TODO(), &amenities)

	if err != nil {
		fmt.Println(err)
		return amenities
	}

	return amenities
}

func (c amenityClient) DeleteAmenityById(id string) error {

	objID, err := primitive.ObjectIDFromHex(id)

	result, err := db.AmenitiesCollection.DeleteOne(context.TODO(), bson.D{{"_id", objID}})

	if err != nil {

		log.Debug("Failed to delete amenity")
		return err

	} else if result.DeletedCount == 0 {

		log.Debug("Amenity not found")
		return errors.New("amenity not found")
	}

	log.Debug("Amenity deleted successfully: ", id)
	return nil
}
