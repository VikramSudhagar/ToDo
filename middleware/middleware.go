package middleware

import (
	"log"
	"todo/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func middleware() {
	app := fiber.New()

	app.Use("/task", basicauth.New(basicauth.Config{
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			currentUser, err := database.VerifyUser(user, pass)
			if user == currentUser.Email && pass == currentUser.Password {
				return true
			}
			log.Println("There error was: ", err)
			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendFile("./unauthorized.html")
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}))

}
