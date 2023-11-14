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

type hotelClient struct{}

type hotelClientInterface interface {
	InsertHotel(hotel model.Hotel) model.Hotel
	GetHotelById(id string) model.Hotel
	GetHotels() model.Hotels
	DeleteHotelById(id string) error
	UpdateHotelById(hotel model.Hotel) model.Hotel
}

var HotelClient hotelClientInterface

func init() {
	HotelClient = &hotelClient{}
}

func (c hotelClient) InsertHotel(hotel model.Hotel) model.Hotel {

	insertHotel := hotel
	insertHotel.Id = primitive.NewObjectID()

	_, err := db.HotelsCollection.InsertOne(context.TODO(), &insertHotel)

	if err != nil {
		fmt.Println(err)
		return hotel
	}

	hotel.Id = insertHotel.Id

	return hotel
}

func (c hotelClient) GetHotelById(id string) model.Hotel {
	var hotel model.Hotel

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		fmt.Println(err)
		return hotel
	}

	err = db.HotelsCollection.FindOne(context.TODO(), bson.D{{"_id", objID}}).Decode(&hotel)
	if err != nil {
		fmt.Println(err)
		return hotel
	}
	return hotel
}

func (c hotelClient) GetHotels() model.Hotels {
	var hotels model.Hotels

	cursor, err := db.HotelsCollection.Find(context.TODO(), bson.D{})

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

func (c hotelClient) DeleteHotelById(id string) error {

	objID, _ := primitive.ObjectIDFromHex(id)

	result, err := db.HotelsCollection.DeleteOne(context.TODO(), bson.D{{"_id", objID}})

	if err != nil {

		log.Debug("Failed to delete hotel")
		return err

	} else if result.DeletedCount == 0 {

		log.Debug("Hotel not found")
		return errors.New("hotel not found")
	}

	log.Debug("Hotel deleted successfully: ", id)
	return nil
}

func (c hotelClient) UpdateHotelById(hotel model.Hotel) model.Hotel {

	update := bson.D{{"$set",
		bson.D{
			{"name", hotel.Name},
			{"room_amount", hotel.RoomAmount},
			{"description", hotel.Description},
			{"city", hotel.City},
			{"street_name", hotel.StreetName},
			{"street_number", hotel.StreetNumber},
			{"rate", hotel.Rate},
			{"amenities", hotel.Amenities},
			{"images", hotel.Images},
		},
	}}

	result, err := db.HotelsCollection.UpdateOne(context.TODO(), bson.D{{"_id", hotel.Id}}, update)

	if result.MatchedCount != 0 {
		log.Debug("Updated hotel successfully")
		return hotel

	} else if err != nil {
		log.Debug("Failed to update hotel")
	}

	return model.Hotel{}
}
