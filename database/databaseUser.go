package database

import (
	"errors"
	"log"
	"todo/models"

	"github.com/alexedwards/argon2id"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var UserDBConn *gorm.DB

func UserSetUp() {
	dsn := "host=host.docker.internal user=postgres password=mypassword dbname=todoDB port=5432"
	UserDBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Println("DB Connection was successful for UserSetUp")
		UserDBConn.AutoMigrate(models.User{})
	} else {
		log.Println("The error is: ", err)
	}
}

func VerifyUser(email string, password string) (models.User, error) {
	//Verify whether the email and password exist
	db := UserDBConn
	user := models.User{}
	dbUser := models.User{}
	//Retrieve the user information from the database based of the email information
	if err := db.Where("email = ?", email).First(&dbUser).Error; err != nil {
		//The user does not exist
		return user, err
	}
	match, e := argon2id.ComparePasswordAndHash(password, dbUser.Password)
	if e != nil {
		log.Fatal(e)
	}

	if match {
		log.Println("User credentials are valid")
		user = dbUser
	} else {
		log.Println("The user credentials are invalid")
		return user, gorm.ErrRecordNotFound
	}
	return user, nil
}

//Get User through the UserID
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
	//Hash the password with argon2id
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = hash
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
