package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2/log"
)

var Db *gorm.DB;

type Subscriber struct {
	gorm.Model
	Name string
	Email string
}

type UserFunctions interface {
	GetUserData(email string) Subscriber
	CreateSubscriber(payload Subscriber) bool
	GetSubscribers() []Subscriber
	Exists(email string) bool
}



func Connect() {
	var err error
	Db, err = gorm.Open(sqlite.Open("subscribers.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	Db.AutoMigrate(&Subscriber{})
}

func (s *Subscriber) GetUserData(email string) Subscriber {
	Connect()

	var sbs Subscriber
	result := Db.First(&sbs, "email = ?", email)

	if result.RowsAffected == 0  {
		log.Fatal("No data found")
	}

	return sbs

}

func (s *Subscriber) CreateSubscriber(payload Subscriber) bool {
	Connect()

	result := Db.Create(&Subscriber{
		Name: payload.Name,
		Email: payload.Email,
	})

	if result.Error != nil {
		log.Error("User not  created", result.Error)
		return false
	} else {
		log.Info("User created")
		return true
	}
}

func (s *Subscriber) GetSubscribers() []Subscriber {
	Connect()

	var subs []Subscriber
	result := Db.Find(&subs)

	if result.Error != nil {
		log.Error("No subscribers found")
	}

	return subs
}

func (s *Subscriber) Exists(email string) bool {
	Connect()
	var u Subscriber
	result := Db.First(&u, "email = ?", email)

	return (result.RowsAffected > 0)
}
