package database

import (
	"errors"
	"log"
	"todo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var UserDBConn *gorm.DB

func UserSetUp() {
	dsn := "host=127.0.0.1 user=postgres password=new@co-op222! dbname=todoDB port=5432"
	UserDBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Println("DB Connection was successful for UserSetUp")
	} else {
		log.Println("The error is: ", err)
	}
}

func VerifyUser(email string, password string) models.User {
	//Verify whether the email and password exist
	db := UserDBConn
	var user models.User = models.User{}
	err := db.Where("email = ? AND password = ?", email, password).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user
	}
	return user
}

func GetUser(id string) (error, models.User) {
	db := UserDBConn
	var user models.User
	err := db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("User not found"), user
	}
	return nil, user
}

func AddUser(user models.User) error {
	db := UserDBConn
	log.Println("The user value is:", user)
	if err := db.Create(&user); err != nil {
		return err.Error
	}
	log.Println("Adding User was successful")
	return nil
}

func DeleteUser(id string) error {
	db := UserDBConn
	if err := db.Delete(&models.User{}, id); err != nil {
		return err.Error
	}
	log.Println("Deleting User was successful")
	return nil
}
