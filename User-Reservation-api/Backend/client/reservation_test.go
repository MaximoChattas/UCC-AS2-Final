package client

import (
	"User-Reservation/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertReservation(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database %v", err)
	}

	defer db.Close()

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	reservation := model.Reservation{
		StartDate: "20-11-2023",
		EndDate:   "21-11-2023",
		UserId:    1,
		HotelId:   "654cf68d807298d99186019f",
		Amount:    10980.5,
	}

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO `reservations` (`start_date`,`end_date`,`user_id`,`hotel_id`,`amount`) VALUES (?,?,?,?,?)").
		WithArgs(reservation.StartDate, reservation.EndDate, reservation.UserId, reservation.HotelId, reservation.Amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result := ReservationClient.InsertReservation(reservation)

	reservation.Id = 1
	a.NoError(mock.ExpectationsWereMet())
	a.Equal(reservation, result)

}

func TestGetReservationById(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	reservation := model.Reservation{
		Id:        1,
		StartDate: "20-11-2023",
		EndDate:   "21-11-2023",
		UserId:    1,
		HotelId:   "654cf68d807298d99186019f",
		Amount:    10980.5,
	}

	mock.ExpectQuery("SELECT * FROM `reservations` WHERE (id = ?) ORDER BY `reservations`.`id` ASC LIMIT 1").
		WithArgs(reservation.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "start_date", "end_date", "user_id", "hotel_id", "amount"}).
		AddRow(reservation.Id, reservation.StartDate, reservation.EndDate, reservation.UserId, reservation.HotelId, reservation.Amount))

	result := ReservationClient.GetReservationById(reservation.Id)

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(reservation, result)
}

func TestGetReservationsByUser(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	reservation := model.Reservations{
		model.Reservation{
			Id:        1,
			StartDate: "20-11-2023",
			EndDate:   "21-11-2023",
			UserId:    1,
			HotelId:   "654cf68d807298d99186019f",
			Amount:    10980.5,
		},

		model.Reservation{
			Id:        2,
			StartDate: "28-11-2023",
			EndDate:   "30-11-2023",
			UserId:    1,
			HotelId:   "654cf68d807298d99186019f",
			Amount:    4590,
		},
	}

	mock.ExpectQuery("SELECT * FROM `reservations` WHERE (user_id = ?)").
		WithArgs(reservation[0].UserId).WillReturnRows(sqlmock.NewRows([]string{"id", "start_date", "end_date", "user_id", "hotel_id", "amount"}).
		AddRow(reservation[0].Id, reservation[0].StartDate, reservation[0].EndDate, reservation[0].UserId, reservation[0].HotelId, reservation[0].Amount).
		AddRow(reservation[1].Id, reservation[1].StartDate, reservation[1].EndDate, reservation[1].UserId, reservation[1].HotelId, reservation[1].Amount))

	result := ReservationClient.GetReservationsByUser(reservation[0].UserId)

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(reservation, result)
}

func TestGetReservationsByHotel(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	reservation := model.Reservations{
		model.Reservation{
			Id:        1,
			StartDate: "20-11-2023",
			EndDate:   "21-11-2023",
			UserId:    1,
			HotelId:   "654cf68d807298d99186019f",
			Amount:    10980.5,
		},

		model.Reservation{
			Id:        2,
			StartDate: "28-11-2023",
			EndDate:   "30-11-2023",
			UserId:    1,
			HotelId:   "654cf68d807298d99186019f",
			Amount:    4590,
		},
	}

	mock.ExpectQuery("SELECT * FROM `reservations` WHERE (hotel_id = ?)").
		WithArgs(reservation[0].HotelId).WillReturnRows(sqlmock.NewRows([]string{"id", "start_date", "end_date", "user_id", "hotel_id", "amount"}).
		AddRow(reservation[0].Id, reservation[0].StartDate, reservation[0].EndDate, reservation[0].UserId, reservation[0].HotelId, reservation[0].Amount).
		AddRow(reservation[1].Id, reservation[1].StartDate, reservation[1].EndDate, reservation[1].UserId, reservation[1].HotelId, reservation[1].Amount))

	result := ReservationClient.GetReservationsByHotel(reservation[0].HotelId)

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(reservation, result)
}

func TestGetReservations(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	reservation := model.Reservations{
		model.Reservation{
			Id:        1,
			StartDate: "20-11-2023",
			EndDate:   "21-11-2023",
			UserId:    1,
			HotelId:   "654cf68d807298d99186019f",
			Amount:    10980.5,
		},

		model.Reservation{
			Id:        2,
			StartDate: "28-11-2023",
			EndDate:   "30-11-2023",
			UserId:    1,
			HotelId:   "654cf68d807298d99186019f",
			Amount:    4590,
		},
	}

	mock.ExpectQuery("SELECT * FROM `reservations`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "start_date", "end_date", "user_id", "hotel_id", "amount"}).
			AddRow(reservation[0].Id, reservation[0].StartDate, reservation[0].EndDate, reservation[0].UserId, reservation[0].HotelId, reservation[0].Amount).
			AddRow(reservation[1].Id, reservation[1].StartDate, reservation[1].EndDate, reservation[1].UserId, reservation[1].HotelId, reservation[1].Amount))

	result := ReservationClient.GetReservations()

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(reservation, result)
}

func TestDeleteReservation(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	reservation := model.Reservation{
		Id:        1,
		StartDate: "20-11-2023",
		EndDate:   "21-11-2023",
		UserId:    1,
		HotelId:   "654cf68d807298d99186019f",
		Amount:    10980.5,
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `reservations` WHERE `reservations`.`id` = ?").WithArgs(reservation.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = ReservationClient.DeleteReservation(reservation)

	a.Nil(err)
	a.NoError(mock.ExpectationsWereMet())
}
