package client

import (
	"User-Reservation/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertAmadeusMap(t *testing.T) {

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

	amadeusMap := model.AmadeusMap{
		HotelId:   "654cf68d807298d99186019f",
		AmadeusId: "SBMIASOF",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `amadeus_maps` (`hotel_id`,`amadeus_id`) VALUES (?,?)").
		WithArgs(amadeusMap.HotelId, amadeusMap.AmadeusId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result := AmadeusClient.InsertAmadeusMap(amadeusMap)

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(amadeusMap, result)
}

func TestGetAmadeusIdByHotelId(t *testing.T) {

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

	amadeusMap := model.AmadeusMap{
		HotelId:   "654cf68d807298d99186019f",
		AmadeusId: "SBMIASOF",
	}

	mock.ExpectQuery("SELECT * FROM `amadeus_maps` WHERE (hotel_id = ?) LIMIT 1").
		WillReturnRows(sqlmock.NewRows([]string{"hotel_id", "amadeus_id"}).
			AddRow(amadeusMap.HotelId, amadeusMap.AmadeusId))

	result := AmadeusClient.GetAmadeusIdByHotelId(amadeusMap.HotelId)

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(amadeusMap, result)
}
