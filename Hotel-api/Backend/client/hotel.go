package client

import (
	db "Hotel/db"
	"Hotel/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type hotelClient struct{}

type hotelClientInterface interface {
	InsertHotel(hotel model.Hotel) model.Hotel
	GetHotelById(id string) model.Hotel
	GetHotels() model.Hotels
	//DeleteHotel(hotel model.Hotel) error
	//UpdateHotel(hotel model.Hotel) model.Hotel
}

var HotelClient hotelClientInterface

func init() {
	HotelClient = &hotelClient{}
}

func (c hotelClient) InsertHotel(hotel model.Hotel) model.Hotel {

	db := db.MongoDb

	insertHotel := hotel
	insertHotel.Id = primitive.NewObjectID()

	_, err := db.Collection("hotels").InsertOne(context.TODO(), &insertHotel)

	if err != nil {
		fmt.Println(err)
		return hotel
	}

	hotel.Id = insertHotel.Id

	return hotel
}

func (c hotelClient) GetHotelById(id string) model.Hotel {
	var hotel model.Hotel

	db := db.MongoDb
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		fmt.Println(err)
		return hotel
	}

	err = db.Collection("hotels").FindOne(context.TODO(), bson.D{{"_id", objID}}).Decode(&hotel)
	if err != nil {
		fmt.Println(err)
		return hotel
	}
	return hotel
}

func (c hotelClient) GetHotels() model.Hotels {
	var hotels model.Hotels

	db := db.MongoDb

	cursor, err := db.Collection("hotels").Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Println(err)
		return hotels
	}

	err = cursor.All(context.TODO(), &hotels)

	if err != nil {
		fmt.Println(err)
		return hotels
	}

	return hotels
}

// TODO
//func (c amenityClient) DeleteHotel(hotel model.Hotel) error {
//	err := Db.Delete(&hotel).Error
//
//	if err != nil {
//		log.Debug("Failed to delete hotel")
//	} else {
//		log.Debug("Hotel deleted: ", hotel.Id)
//	}
//	return err
//}
//
//// TODO
//func (c amenityClient) UpdateHotel(hotel model.Hotel) model.Hotel {
//
//	//Db.Model(&hotel).Association("Amenities").Clear()
//
//	var newAmenities model.Amenities
//
//	for _, amenity := range hotel.Amenities {
//		newAmenities = append(newAmenities, amenity)
//	}
//	result := Db.Save(&hotel)
//
//	Db.Model(&hotel).Association("Amenities").Replace(newAmenities)
//
//	if result.Error != nil {
//		log.Debug("Failed to update hotel")
//		return model.Hotel{}
//	}
//
//	log.Debug("Updated hotel: ", hotel.Id)
//	return hotel
//}
