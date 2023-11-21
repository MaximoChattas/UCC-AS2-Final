package client

import (
	"User-Reservation/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertUser(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	defer db.Close()

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	user := model.User{
		Name:     "User",
		LastName: "User",
		Dni:      "123456",
		Email:    "email@email.com",
		Password: "pass",
		Role:     "Customer",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` (`name`,`last_name`,`dni`,`email`,`password`,`role`) VALUES (?,?,?,?,?,?)").
		WithArgs(user.Name, user.LastName, user.Dni, user.Email, user.Password, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	response := UserClient.InsertUser(user)

	user.Id = 1
	a.NoError(mock.ExpectationsWereMet())
	a.Equal(user, response)
}

func TestGetUserById(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	defer db.Close()

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	user := model.User{
		Id:       1,
		Name:     "User",
		LastName: "User",
		Dni:      "123456",
		Email:    "email@email.com",
		Password: "pass",
		Role:     "Customer",
	}

	mock.ExpectQuery("SELECT * FROM `users` WHERE (id = ?) ORDER BY `users`.`id` ASC LIMIT 1").WithArgs(user.Id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "last_name", "dni", "email", "password", "role"}).
			AddRow(user.Id, user.Name, user.LastName, user.Dni, user.Email, user.Password, user.Role))

	result := UserClient.GetUserById(user.Id)

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(user, result)
}

func TestGetUserByEmail(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	defer db.Close()

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	user := model.User{
		Id:       1,
		Name:     "User",
		LastName: "User",
		Dni:      "123456",
		Email:    "email@email.com",
		Password: "pass",
		Role:     "Customer",
	}

	mock.ExpectQuery("SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1").WithArgs(user.Email).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "last_name", "dni", "email", "password", "role"}).
			AddRow(user.Id, user.Name, user.LastName, user.Dni, user.Email, user.Password, user.Role))

	result := UserClient.GetUserByEmail(user.Email)

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(user, result)
}

func TestGetUsers(t *testing.T) {

	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	defer db.Close()

	gormDb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	Db = gormDb

	users := model.Users{
		model.User{
			Id:       1,
			Name:     "User 1",
			LastName: "User 1",
			Dni:      "123456",
			Email:    "1@email.com",
			Password: "pass1",
			Role:     "Customer",
		},

		model.User{
			Id:       2,
			Name:     "User 2",
			LastName: "User 2",
			Dni:      "654321",
			Email:    "2@email.com",
			Password: "pass2",
			Role:     "Customer",
		},
	}

	mock.ExpectQuery("SELECT * FROM `users`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "last_name", "dni", "email", "password", "role"}).
			AddRow(users[0].Id, users[0].Name, users[0].LastName, users[0].Dni, users[0].Email, users[0].Password, users[0].Role).
			AddRow(users[1].Id, users[1].Name, users[1].LastName, users[1].Dni, users[1].Email, users[1].Password, users[1].Role))

	result := UserClient.GetUsers()

	a.NoError(mock.ExpectationsWereMet())
	a.Equal(users, result)

}
