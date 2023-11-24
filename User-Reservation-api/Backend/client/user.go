package client

import (
	"User-Reservation/model"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/karlseguin/ccache/v3"
	log "github.com/sirupsen/logrus"
)

type userClient struct{}

type userClientInterface interface {
	InsertUser(user model.User) model.User
	GetUserById(id int) model.User
	GetUserByEmail(email string) model.User
	GetUsers() model.Users
}

var UserClient userClientInterface

var Db *gorm.DB
var cache *ccache.Cache[model.User]

func init() {
	UserClient = &userClient{}
	cache = ccache.New(ccache.Configure[model.User]())
}

func (c *userClient) InsertUser(user model.User) model.User {

	result := Db.Create(&user)

	if result.Error != nil {
		log.Error("Failed to insert user.")
		return user
	}

	log.Debug("User created:", user.Id)

	cache.Set(user.Email, user, 24*time.Hour)
	return user
}

func (c *userClient) GetUserById(id int) model.User {
	var user model.User

	Db.Where("id = ?", id).First(&user)
	log.Debug("User: ", user)

	return user
}

func (c *userClient) GetUserByEmail(email string) model.User {

	item := cache.Get(email)

	if item != nil {
		item.Extend(24 * time.Hour)
		log.Info("Data retrieved from local cache")
		return item.Value()
	}

	var user model.User

	result := Db.Where("email = ?", email).First(&user)
	log.Debug("User: ", user)

	if result.Error == nil {
		cache.Set(email, user, 24*time.Hour)
	}

	return user
}

func (c *userClient) GetUsers() model.Users {
	var users model.Users
	Db.Find(&users)

	log.Debug("Users: ", users)

	return users
}
