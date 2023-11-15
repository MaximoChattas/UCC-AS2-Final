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

func TestInsertAmenity(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		db.AmenitiesCollection = mt.Coll

		amenity := AmenityClient.InsertAmenity(model.Amenity{
			Name: "Amenity Test",
		})

		a.NotNil(amenity.Id)
	})
}

func TestGetAmenityById(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		db.AmenitiesCollection = mt.Coll

		expectedAmenity := model.Amenity{
			Id:   primitive.NewObjectID(),
			Name: "Amenity Test",
		}

		cursor := mtest.CreateCursorResponse(1, "amenity.1", mtest.FirstBatch, bson.D{
			{"_id", expectedAmenity.Id},
			{"name", expectedAmenity.Name},
		})

		mt.AddMockResponses(cursor)

		amenity := AmenityClient.GetAmenityById(expectedAmenity.Id.Hex())

		a.Equal(expectedAmenity, amenity)
	})

	mt.Run("failure", func(mt *mtest.T) {

		db.AmenitiesCollection = mt.Coll

		amenity := AmenityClient.GetAmenityById(primitive.NewObjectID().Hex())

		var emptyModel model.Amenity

		a.Equal(emptyModel, amenity)
	})
}

func TestGetAmenityByName(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		db.AmenitiesCollection = mt.Coll

		expectedAmenity := model.Amenity{
			Id:   primitive.NewObjectID(),
			Name: "Amenity Test",
		}

		cursor := mtest.CreateCursorResponse(1, "amenity.1", mtest.FirstBatch, bson.D{
			{"_id", expectedAmenity.Id},
			{"name", expectedAmenity.Name},
		})

		mt.AddMockResponses(cursor)

		amenity := AmenityClient.GetAmenityByName(expectedAmenity.Name)

		a.Equal(expectedAmenity, amenity)
	})

	mt.Run("failure", func(mt *mtest.T) {

		db.AmenitiesCollection = mt.Coll

		amenity := AmenityClient.GetAmenityByName("")

		var emptyModel model.Amenity

		a.Equal(emptyModel, amenity)
	})
}

func TestGetAmenities(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		db.AmenitiesCollection = mt.Coll

		expectedAmenity1 := model.Amenity{
			Id:   primitive.NewObjectID(),
			Name: "Amenity Test 1",
		}

		expectedAmenity2 := model.Amenity{
			Id:   primitive.NewObjectID(),
			Name: "Amenity Test 2",
		}

		cursor1 := mtest.CreateCursorResponse(1, "amenity.1", mtest.FirstBatch, bson.D{
			{"_id", expectedAmenity1.Id},
			{"name", expectedAmenity1.Name},
		})

		cursor2 := mtest.CreateCursorResponse(1, "amenity.1", mtest.NextBatch, bson.D{
			{"_id", expectedAmenity2.Id},
			{"name", expectedAmenity2.Name},
		})

		killCursors := mtest.CreateCursorResponse(0, "amenity.1", mtest.NextBatch)

		mt.AddMockResponses(cursor1, cursor2, killCursors)

		amenities := AmenityClient.GetAmenities()

		var expectedAmenities model.Amenities

		expectedAmenities = append(expectedAmenities, expectedAmenity1)
		expectedAmenities = append(expectedAmenities, expectedAmenity2)

		a.Equal(expectedAmenities, amenities)

	})
}

func TestDeleteAmenityById(t *testing.T) {

	a := assert.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		db.AmenitiesCollection = mt.Coll

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})

		err := AmenityClient.DeleteAmenityById(primitive.NewObjectID().Hex())
		a.Nil(err)
	})

	mt.Run("failure", func(mt *mtest.T) {
		db.AmenitiesCollection = mt.Coll

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 0}})

		err := AmenityClient.DeleteAmenityById(primitive.NewObjectID().Hex())
		a.NotNil(err)

		expectedResponse := "amenity not found"
		a.Equal(expectedResponse, err.Error())
	})

	mt.Run("error", func(mt *mtest.T) {
		db.AmenitiesCollection = mt.Coll

		mt.AddMockResponses(bson.D{{"ok", 0}})

		err := AmenityClient.DeleteAmenityById(primitive.NewObjectID().Hex())
		a.NotNil(err)

		expectedResponse := "command failed"
		a.Equal(expectedResponse, err.Error())
	})
}
