package middleware

import (
	"log"
	"strings"
	"todo/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
)

func RetrieveSession(store *redis.Storage) (string, string, error) {
	answer, err := store.Get("session")
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
func RetrieveSessionAndVerify(store *redis.Storage, c *fiber.Ctx) error {
	email, password, e := RetrieveSession(store)
	if e != nil {
		return e
	}

	//Verify whether the information in the session storage is valid or not
	if _, err := database.VerifyUser(email, password); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"Error":   err,
		})
	}

	return nil
}
