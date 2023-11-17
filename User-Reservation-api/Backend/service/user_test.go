package service

import (
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type testUser struct{}

func init() {
	client.UserClient = testUser{}
}

func (t testUser) InsertUser(user model.User) model.User {

	if user.Name != "" {
		user.Id = 1
	}

	return user
}

func (t testUser) GetUserById(id int) model.User {

	if id == 0 {
		return model.User{}
	}

	return model.User{Id: id}
}

func (t testUser) GetUserByEmail(email string) model.User {

	if email == "" {
		return model.User{}
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	user := model.User{
		Id:       1,
		Email:    email,
		Password: string(encryptedPassword),
	}

	return user

}

func (t testUser) GetUsers() model.Users {
	return model.Users{}
}

func TestInsertUser(t *testing.T) {

	a := assert.New(t)

	user := dto.UserDto{
		Name:     "User",
		LastName: "User",
		Email:    "user@user.com",
		Password: "password",
	}

	response, err := UserService.InsertUser(user)
	a.Nil(err)

	err = bcrypt.CompareHashAndPassword([]byte(response.Password), []byte(user.Password))
	a.Nil(err)

	user.Role = "Customer"
	user.Password = response.Password
	user.Id = response.Id

	a.Equal(user, response)
}

func TestInsertUser_Error(t *testing.T) {

	a := assert.New(t)

	user := dto.UserDto{}

	_, err := UserService.InsertUser(user)
	a.NotNil(err)

	expectedResponse := "error creating user"
	a.Equal(expectedResponse, err.Error())
}

func TestGetUserById_NotFound(t *testing.T) {

	a := assert.New(t)

	id := 0

	_, err := UserService.GetUserById(id)
	a.NotNil(err)

	expectedResponse := "user not found"
	a.Equal(expectedResponse, err.Error())
}

func TestGetUserById_Found(t *testing.T) {

	a := assert.New(t)

	id := 1

	response, err := UserService.GetUserById(id)
	a.Nil(err)

	expectedResponse := dto.UserDto{Id: id}
	a.Equal(expectedResponse, response)
}

func TestGetUsers(t *testing.T) {

	a := assert.New(t)

	response, err := UserService.GetUsers()

	var expectedResponse dto.UsersDto

	a.Nil(err)
	a.Equal(expectedResponse, response)
}

func TestUserLogin_NotRegistered(t *testing.T) {

	a := assert.New(t)

	user := dto.UserDto{}

	_, err := UserService.UserLogin(user)
	a.NotNil(err)

	expectedResponse := "user not registered"
	a.Equal(expectedResponse, err.Error())
}

func TestUserLogin_IncorrectPassword(t *testing.T) {

	a := assert.New(t)

	user := dto.UserDto{
		Email:    "email@email.com",
		Password: "incorrect",
	}

	_, err := UserService.UserLogin(user)
	a.NotNil(err)

	expectedResponse := "incorrect password"
	a.Equal(expectedResponse, err.Error())
}

func TestUserLogin_Success(t *testing.T) {

	a := assert.New(t)

	user := dto.UserDto{
		Email:    "email@email.com",
		Password: "password",
	}

	response, err := UserService.UserLogin(user)
	a.Nil(err)

	a.Equal(1, response.Id)
	a.Equal(user.Email, response.Email)
}
