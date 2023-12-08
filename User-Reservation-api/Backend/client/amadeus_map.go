package client

import (
	"User-Reservation/model"
	log "github.com/sirupsen/logrus"
)

type amadeusClient struct{}

type amadeusClientInterface interface {
	InsertAmadeusMap(mapping model.AmadeusMap) model.AmadeusMap
	GetAmadeusIdByHotelId(hotelId string) model.AmadeusMap
	DeleteMapping(mapping model.AmadeusMap) error
}

var AmadeusClient amadeusClientInterface

func init() {
	AmadeusClient = &amadeusClient{}
}

func (c *amadeusClient) InsertAmadeusMap(mapping model.AmadeusMap) model.AmadeusMap {

	result := Db.Create(&mapping)

	if result.Error != nil {
		log.Error("Failed to insert mapping.")
		return mapping
	}

	log.Debug("Hotel mapping created:", mapping.HotelId)
	return mapping
}

func (c *amadeusClient) GetAmadeusIdByHotelId(hotelId string) model.AmadeusMap {
	var mapping model.AmadeusMap

	Db.Where("hotel_id = ?", hotelId).First(&mapping)
	log.Debug("Mapping: ", mapping)

	return mapping
}

func (c *amadeusClient) DeleteMapping(mapping model.AmadeusMap) error {
	err := Db.Where("hotel_id = ?", mapping.HotelId).Delete(&model.AmadeusMap{}).Error

	if err != nil {
		log.Debug("Failed to delete mapping")
	} else {
		log.Debug("Mapping deleted: ", mapping.HotelId)
	}

	return err
}
