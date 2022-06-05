package main

import (
	"log"
	"todo/controller"
	"todo/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {
	engine := html.New("./views", ".html")
	err := godotenv.Load("local.env")
	if err != nil {
		log.Println("Error loading .env file, the error was: ", err)
	}
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	database.TaskSetUp()
	database.UserSetUp()

	app.Static("/", "./views", fiber.Static{
		Index: "index.html",
	})

	route := app.Group("/")

	app.Static("/static", "./static")

	//User Login
	route.Post("/login", controller.Login)

	//Get the User's tasks
	route.Get("/task", controller.GetTask)

	//Adding a Task to the to do list
	route.Post("/task", controller.AddTask)

	//Deleting a Task from the to do list
	route.Delete("/task/:id", controller.DeleteTask)

	//Retrieving a User given the User ID
	route.Get("/user/:id", controller.GetUser)

	//Sign Up
	route.Post("/signup", controller.SignUp)

	//Delete User
	route.Delete("/user/delete/:id", controller.DeleteUser)

	//Logout
	route.Post("/logout", controller.Logout)

	app.Listen(":8081")
}
