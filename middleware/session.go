package middleware

import (
	"log"
	"strings"
	"todo/database"
	"todo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
)

func RetrieveSession(store *redis.Storage, key string) (string, string, error) {
	answer, err := store.Get(key)
	if err != nil || len(string(answer)) == 0 {
		log.Println("There was a problem retrieving information from the session")
		return "", "", fiber.NewError(fiber.StatusBadRequest, "Invalid Response")
	} else {
		credentials := string(answer)
		stringArray := strings.Fields(credentials)
		return stringArray[1], stringArray[3], nil
	}
}

//TODO: Hash the passwords in the database and session with Argon2Id
//RetrieveSessionAndVerify function is currently only being used with "/task" endpoints
func RetrieveSessionAndVerify(store *redis.Storage, c *fiber.Ctx, key string) (models.User, error) {
	email, password, e := RetrieveSession(store, key)
	if e != nil {
		return models.User{}, e
	}

	//Verify whether the information in the session storage is valid or not
	user, err := database.VerifyUser(email, password)
	if err != nil {
		return user, c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"Error":   err,
		})
	}

	return user, nil
}
