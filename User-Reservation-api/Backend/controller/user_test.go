package controller

import (
	"User-Reservation/dto"
	"User-Reservation/service"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testUser struct{}

func init() {
	service.UserService = testUser{}
}

func (t testUser) InsertUser(userDto dto.UserDto) (dto.UserDto, error) {

	if userDto.Id == 0 {
		return dto.UserDto{}, errors.New("error creating user")
	}

	return userDto, nil
}

func (t testUser) GetUserById(id int) (dto.UserDto, error) {

	if id == 0 {
		return dto.UserDto{}, errors.New("user not found")
	}

	return dto.UserDto{Id: id}, nil
}

func (t testUser) GetUsers() (dto.UsersDto, error) {

	return dto.UsersDto{}, nil
}

func (t testUser) UserLogin(loginDto dto.UserDto) (dto.UserDto, error) {

	if loginDto.Email == "" || loginDto.Password == "" {
		return dto.UserDto{}, errors.New("user not registered")
	}

	loginDto.Id = 1
	loginDto.Name = "User"
	loginDto.LastName = "User"
	loginDto.Dni = "123456"
	loginDto.Role = "Customer"

	return loginDto, nil
}

func TestInsertUser(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/user", InsertUser)

	body := `{
		"id": 1,
        "name": "User",
        "last_name": "User",
        "dni": "12345",
		"email": "email@email.com",
        "password": "pass",
        "role": "Customer"
    }`

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusCreated, w.Code)

	var response dto.UserDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal(1, response.Id)
}

func TestInsertUser_Error(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/user", InsertUser)

	body := `{
		"id": 0,
        "name": "User",
        "last_name": "User",
        "dni": "12345",
		"email": "email@email.com",
        "password": "pass",
        "role": "Customer"
    }`

	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusBadRequest, w.Code)

	var response dto.UserDto
	var emptyDto dto.UserDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	a.Equal(emptyDto, response)
}

func TestGetUserById_NotFound(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.GET("/user/:id", GetUserById)

	req, err := http.NewRequest(http.MethodGet, "/user/0", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusNotFound, w.Code)

	expectedResponse := `{"error":"user not found"}`

	a.Equal(expectedResponse, w.Body.String())
}

func TestGetUserById_Found(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.GET("/user/:id", GetUserById)

	req, err := http.NewRequest(http.MethodGet, "/user/1", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.UserDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedResponse := dto.UserDto{Id: 1}

	a.Equal(expectedResponse, response)
}

func TestGetUsers(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.GET("/user", GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/user", nil)

	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var response dto.UsersDto
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedResponse := dto.UsersDto{}

	a.Equal(expectedResponse, response)
}

func TestUserLogin_Error(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/login", UserLogin)

	body := `{
		"email": "",
        "password": ""
    }`

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusUnauthorized, w.Code)

	expectedResponse := `{"error":"user not registered"}`

	a.Equal(expectedResponse, w.Body.String())

}

func TestUserLogin(t *testing.T) {

	a := assert.New(t)

	r := gin.Default()
	r.POST("/login", UserLogin)

	body := `{
		"email": "email@email.com",
        "password": "pass"
    }`

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	if err != nil {
		t.Fatalf("New request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	r.ServeHTTP(w, req)

	a.Equal(http.StatusAccepted, w.Code)
	a.Contains(w.Body.String(), "token")

}

func TestGenerateToken(t *testing.T) {

	a := assert.New(t)

	loginDto := dto.UserDto{
		Id:       1,
		Name:     "User",
		LastName: "User",
		Dni:      "123456",
		Email:    "email@email.com",
		Password: "pass",
		Role:     "Customer",
	}

	token, err := generateToken(loginDto)

	a.Nil(err)

	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	a.Equal(float64(1), claims["id"])
	a.Equal("User", claims["name"])
	a.Equal("Customer", claims["role"])
}
