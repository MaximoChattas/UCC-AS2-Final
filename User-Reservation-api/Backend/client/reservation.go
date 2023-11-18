package client

import (
	"User-Reservation/model"
	log "github.com/sirupsen/logrus"
)

type reservationClient struct{}

type reservationClientInterface interface {
	InsertReservation(reservation model.Reservation) model.Reservation
	GetReservationById(id int) model.Reservation
	GetReservations() model.Reservations
	GetReservationsByUser(userId int) model.Reservations
	GetReservationsByHotel(hotelId string) model.Reservations
	DeleteReservation(reservation model.Reservation) error
}

var ReservationClient reservationClientInterface

func init() {
	ReservationClient = &reservationClient{}
}

func (c *reservationClient) InsertReservation(reservation model.Reservation) model.Reservation {

	result := Db.Create(&reservation)

	if result.Error != nil {
		log.Error("Failed to insert reservation.")
		return reservation
	}

	log.Debug("Reservation created:", reservation.Id)
	return reservation
}

func (c *reservationClient) GetReservationById(id int) model.Reservation {
	var reservation model.Reservation

	Db.Where("id = ?", id).First(&reservation)
	log.Debug("Reservation: ", reservation)

	return reservation
}

func (c *reservationClient) GetReservations() model.Reservations {
	var reservations model.Reservations
	Db.Find(&reservations)

	log.Debug("Reservations: ", reservations)

	return reservations
}

func (c *reservationClient) GetReservationsByUser(userId int) model.Reservations {
	var reservations model.Reservations

	Db.Where("user_id = ?", userId).Find(&reservations)
	log.Debug("Reservations: ", reservations)

	return reservations
}

func (c *reservationClient) GetReservationsByHotel(hotelId string) model.Reservations {
	var reservations model.Reservations

	Db.Where("hotel_id = ?", hotelId).Find(&reservations)
	log.Debug("Reservations: ", reservations)

	return reservations
}

func (c *reservationClient) DeleteReservation(reservation model.Reservation) error {
	err := Db.Delete(&reservation).Error

	if err != nil {
		log.Debug("Failed to delete reservation")
	} else {
		log.Debug("Reservation deleted: ", reservation.Id)
	}
	return err
}
