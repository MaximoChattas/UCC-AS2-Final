package client

import (
	"Hotel/db"
	"Hotel/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestInsertHotel(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		db.HotelsCollection = mt.Coll

		hotel := HotelClient.InsertHotel(model.Hotel{
			Name:         "Hotel Test",
			RoomAmount:   10,
			Description:  "Test hotel description",
			City:         "Test City",
			StreetName:   "Test Street",
			StreetNumber: 123,
			Rate:         4.5,
		})

		a.NotNil(hotel.Id)
	})
}

func TestGetHotelById(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		db.HotelsCollection = mt.Coll

		expectedHotel := model.Hotel{
			Id:           primitive.NewObjectID(),
			Name:         "Hotel Test",
			RoomAmount:   10,
			Description:  "Hotel test Description",
			City:         "Test City",
			StreetName:   "Test Street",
			StreetNumber: 123,
			Rate:         4.5,
			Amenities:    nil,
			Images:       nil,
		}

		cursor := mtest.CreateCursorResponse(1, "hotel.1", mtest.FirstBatch, bson.D{
			{"_id", expectedHotel.Id},
			{"name", expectedHotel.Name},
			{"room_amount", expectedHotel.RoomAmount},
			{"description", expectedHotel.Description},
			{"city", expectedHotel.City},
			{"street_name", expectedHotel.StreetName},
			{"street_number", expectedHotel.StreetNumber},
			{"rate", expectedHotel.Rate},
			{"amenities", expectedHotel.Amenities},
			{"images", expectedHotel.Images},
		})

		mt.AddMockResponses(cursor)

		hotel := HotelClient.GetHotelById(expectedHotel.Id.Hex())

		a.Equal(expectedHotel, hotel)
	})

	mt.Run("failure", func(mt *mtest.T) {

		db.HotelsCollection = mt.Coll

		hotel := HotelClient.GetHotelById(primitive.NewObjectID().Hex())

		var emptyModel model.Hotel

		a.Equal(emptyModel, hotel)
	})
}

func TestGetHotels(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		db.HotelsCollection = mt.Coll

		expectedHotel1 := model.Hotel{
			Id:           primitive.NewObjectID(),
			Name:         "Hotel Test 1",
			RoomAmount:   10,
			Description:  "Hotel test Description 1",
			City:         "Test City 1",
			StreetName:   "Test Street 1",
			StreetNumber: 123,
			Rate:         4.5,
			Amenities:    nil,
			Images:       nil,
		}

		expectedHotel2 := model.Hotel{
			Id:           primitive.NewObjectID(),
			Name:         "Hotel Test 2",
			RoomAmount:   10,
			Description:  "Hotel test Description 2",
			City:         "Test City 2",
			StreetName:   "Test Street 2",
			StreetNumber: 123,
			Rate:         4.5,
			Amenities:    nil,
			Images:       nil,
		}

		cursor1 := mtest.CreateCursorResponse(1, "hotel.1", mtest.FirstBatch, bson.D{
			{"_id", expectedHotel1.Id},
			{"name", expectedHotel1.Name},
			{"room_amount", expectedHotel1.RoomAmount},
			{"description", expectedHotel1.Description},
			{"city", expectedHotel1.City},
			{"street_name", expectedHotel1.StreetName},
			{"street_number", expectedHotel1.StreetNumber},
			{"rate", expectedHotel1.Rate},
			{"amenities", expectedHotel1.Amenities},
			{"images", expectedHotel1.Images},
		})

		cursor2 := mtest.CreateCursorResponse(1, "hotel.1", mtest.NextBatch, bson.D{
			{"_id", expectedHotel2.Id},
			{"name", expectedHotel2.Name},
			{"room_amount", expectedHotel2.RoomAmount},
			{"description", expectedHotel2.Description},
			{"city", expectedHotel2.City},
			{"street_name", expectedHotel2.StreetName},
			{"street_number", expectedHotel2.StreetNumber},
			{"rate", expectedHotel2.Rate},
			{"amenities", expectedHotel2.Amenities},
			{"images", expectedHotel2.Images},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		mt.AddMockResponses(cursor1, cursor2, killCursors)

		hotels := HotelClient.GetHotels()

		var expectedHotels model.Hotels

		expectedHotels = append(expectedHotels, expectedHotel1)
		expectedHotels = append(expectedHotels, expectedHotel2)

		a.Equal(expectedHotels, hotels)

	})
}

func TestDeleteHotelById(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		db.HotelsCollection = mt.Coll

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})

		err := HotelClient.DeleteHotelById(primitive.NewObjectID().Hex())
		a.Nil(err)
	})

	mt.Run("failure", func(mt *mtest.T) {
		db.HotelsCollection = mt.Coll

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 0}})

		err := HotelClient.DeleteHotelById(primitive.NewObjectID().Hex())
		a.NotNil(err)

		expectedResponse := "hotel not found"
		a.Equal(expectedResponse, err.Error())
	})

	mt.Run("error", func(mt *mtest.T) {
		db.HotelsCollection = mt.Coll

		mt.AddMockResponses(bson.D{{"ok", 0}})

		err := HotelClient.DeleteHotelById(primitive.NewObjectID().Hex())
		a.NotNil(err)

		expectedResponse := "command failed"
		a.Equal(expectedResponse, err.Error())
	})
}

func TestUpdateHotelById(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		db.HotelsCollection = mt.Coll

		updatedHotel := model.Hotel{
			Id:           primitive.NewObjectID(),
			Name:         "Hotel Test Update",
			RoomAmount:   10,
			Description:  "Hotel test Description Update",
			City:         "Test City Update",
			StreetName:   "Test Street Update",
			StreetNumber: 123,
			Rate:         4.5,
			Amenities:    nil,
			Images:       nil,
		}

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 1},
		})

		hotelResponse := HotelClient.UpdateHotelById(updatedHotel)

		a.Equal(updatedHotel, hotelResponse)

	})

	mt.Run("failure", func(mt *mtest.T) {

		db.HotelsCollection = mt.Coll

		updatedHotel := model.Hotel{
			Id:           primitive.NewObjectID(),
			Name:         "Hotel Test Update",
			RoomAmount:   10,
			Description:  "Hotel test Description Update",
			City:         "Test City Update",
			StreetName:   "Test Street Update",
			StreetNumber: 123,
			Rate:         4.5,
			Amenities:    nil,
			Images:       nil,
		}

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 0},
		})

		hotelResponse := HotelClient.UpdateHotelById(updatedHotel)

		emptyModel := model.Hotel{}

		a.Equal(emptyModel, hotelResponse)

	})
}
