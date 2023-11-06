package client

import (
	"User-Reservation/model"
	log "github.com/sirupsen/logrus"
)

func InsertAmadeusMap(mapping model.AmadeusMap) model.AmadeusMap {

	result := Db.Create(&mapping)

	if result.Error != nil {
		log.Error("Failed to insert mapping.")
		return mapping
	}

	log.Debug("Hotel mapping created:", mapping.HotelId)
	return mapping
}

func GetAmadeusIdByHotelId(hotelId string) model.AmadeusMap {
	var mapping model.AmadeusMap

	Db.Where("hotel_id = ?", hotelId).First(&mapping)
	log.Debug("Mapping: ", mapping)

	return mapping
}
