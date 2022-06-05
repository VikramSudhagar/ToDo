package controller

import (
	"log"
	"time"
	"todo/database"
	"todo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jinzhu/copier"
)

func Login(c *fiber.Ctx) error {
	session := session.New()
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	if _, err := database.VerifyUser(user.Email, user.Password); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"Error":   "User not found",
		})
	}
	//When logging in, the username and password will be stored in redis
	value := "email: " + user.Email + " password: " + user.Password
	//TODO: Add a validator function to determine whether a valid email was passed
	if len(user.Email) == 0 && len(user.Password) == 0 {
		//The user did not pass any credentials
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"Error":   "Invalid Input",
		})
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "sessionID"
	cookie.Value = string(session.KeyGenerator())
	//Session Token will expire in 24 hours
	cookie.Expires = time.Now().Add(24 * time.Hour)

	//Setting the Cookie
	c.Cookie(cookie)
	store.Set(cookie.Value, []byte(value), 0)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"Message": "Login Successful",
	})
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	e, User := database.GetUser(id)
	if e != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   e,
		})
	}

	DTO_User := &models.DTO_User{}

	if err := copier.Copy(&DTO_User, &User); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"user":    DTO_User,
	})
}

func SignUp(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Println("Could not add the user")
	}

	if err := database.AddUser(user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	} else {
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"success": true,
			"message": "User created successfully",
		})
	}
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	delete := database.DeleteUser(id)
	if delete != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Deletion was not successful",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": "Deletion was successful",
	})

}

func Logout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "sessionID"
	cookie.Value = c.Cookies("sessionID")
	cookie.Expires = time.Now().Add(-100 * time.Hour)
	c.Cookie(cookie)
	store.Delete(c.Cookies("sessionID"))
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"Message": "Logout was successful",
	})
}
