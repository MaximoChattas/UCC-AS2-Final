package service

import (
	"Hotel/client"
	"Hotel/dto"
	"Hotel/model"
	"Hotel/queue"
	"errors"
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type testHotel struct{}
type testQueue struct{}

var channel wabbit.Channel
var mockQueue wabbit.Queue

func init() {
	client.HotelClient = testHotel{}
	queue.QueueProducer = testQueue{}
	queue.QueueProducer.InitQueue()
}

func (t testHotel) InsertHotel(hotel model.Hotel) model.Hotel {

	if hotel.Name != "" {
		objId, _ := primitive.ObjectIDFromHex("654cf68d807298d99186019f")
		hotel.Id = objId
	}

	return hotel
}

func (t testHotel) GetHotelById(id string) model.Hotel {
	if id == "000000000000000000000000" {
		return model.Hotel{}
	}

	objId, _ := primitive.ObjectIDFromHex(id)

	return model.Hotel{
		Id: objId,
	}
}

func (t testHotel) GetHotels() model.Hotels {
	return model.Hotels{}
}

func (t testHotel) DeleteHotelById(id string) error {

	if id == "000000000000000000000000" {
		return errors.New("hotel not found")
	}

	return nil
}

func (t testHotel) UpdateHotelById(hotel model.Hotel) model.Hotel {

	if hotel.Id.Hex() == "000000000000000000000000" {
		return model.Hotel{}
	}

	return hotel
}

func (t testQueue) InitQueue() {

	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()

	if err != nil {
		log.Info("Failed to start server")
		log.Fatal(err)
	}

	connection, err := amqptest.Dial("amqp://localhost:5672/%2f")

	if err != nil {
		log.Info("Failed to connect to RabbitMQ")
		log.Fatal(err)
	} else {
		log.Info("RabbitMQ connection established")
	}

	channel, err = connection.Channel()

	if err != nil {
		log.Info("Failed to open channel")
		log.Fatal(err)
	}

	defer channel.Close()

	mockQueue, err = channel.QueueDeclare("hotel", wabbit.Option{
		"durable":    true,
		"autoDelete": false,
		"exclusive":  false,
		"noWait":     false,
	})

	if err != nil {
		log.Info("Failed to declare a queue")
		log.Fatal(err)
	} else {
		log.Info("Queue declared")
	}
}

func (t testQueue) Publish(body []byte) error {

	err := channel.Publish(
		"",
		mockQueue.Name(),
		body,
		wabbit.Option{"contentType": "application/json"})

	if err != nil {
		log.Debug("Error while publishing message", err)
		return err
	}

	return nil
}

func TestInsertHotel_Error(t *testing.T) {

	a := assert.New(t)
	var hotelDto dto.HotelDto

	_, err := HotelService.InsertHotel(hotelDto)

	expectedResponse := "error creating hotel"

	a.NotNil(err)
	a.Equal(expectedResponse, err.Error())

}

func TestInsertHotel(t *testing.T) {

	a := assert.New(t)

	hotelDto := dto.HotelDto{
		Name:         "Hotel Test",
		RoomAmount:   10,
		Description:  "Hotel test description",
		City:         "Test City",
		StreetName:   "Test Street",
		StreetNumber: 123,
		Rate:         4.5,
	}

	hotelResponse, err := HotelService.InsertHotel(hotelDto)

	hotelDto.Id = "654cf68d807298d99186019f"

	a.Nil(err)
	a.Equal(hotelDto, hotelResponse)

}

func TestGetHotelById_NotFound(t *testing.T) {

	a := assert.New(t)

	_, err := HotelService.GetHotelById("000000000000000000000000")

	expectedResponse := "hotel not found"
	a.NotNil(err)
	a.Equal(expectedResponse, err.Error())
}

func TestGetHotelById_Found(t *testing.T) {

	a := assert.New(t)

	id := "654cf68d807298d99186019f"

	hotelResponse, err := HotelService.GetHotelById(id)

	a.Nil(err)
	a.Equal(id, hotelResponse.Id)
}

func TestGetHotels(t *testing.T) {

	a := assert.New(t)

	hotelsResponse, err := HotelService.GetHotels()

	a.Nil(err)

	var emptyDto dto.HotelsDto
	a.Equal(emptyDto, hotelsResponse)
}

func TestDeleteHotel_NotFound(t *testing.T) {

	a := assert.New(t)

	id := "000000000000000000000000"
	err := HotelService.DeleteHotel(id)

	a.NotNil(err)

	expectedResponse := "hotel not found"
	a.Equal(expectedResponse, err.Error())
}

func TestDeleteHotel_Found(t *testing.T) {

	a := assert.New(t)

	id := "654cf68d807298d99186019f"
	err := HotelService.DeleteHotel(id)

	a.Nil(err)
}

func TestUpdateHotel_NotFound(t *testing.T) {

	a := assert.New(t)

	id := "000000000000000000000000"

	hotelDto := dto.HotelDto{Id: id}

	_, err := HotelService.UpdateHotel(hotelDto)

	a.NotNil(err)

	expectedResponse := "hotel not found"
	a.Equal(expectedResponse, err.Error())
}

func TestUpdateHotel_Found(t *testing.T) {

	a := assert.New(t)

	id := "654cf68d807298d99186019f"

	hotelDto := dto.HotelDto{Id: id}

	hotelResponse, err := HotelService.UpdateHotel(hotelDto)

	a.Nil(err)
	a.Equal(hotelDto, hotelResponse)
}
