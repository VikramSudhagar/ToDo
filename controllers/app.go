package main

import (
	"log"
	"todo/database"
	"todo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/jinzhu/copier"
)

func main() {
	app := fiber.New()
	database.TaskSetUp()
	database.UserSetUp()
	log.Println("test")
	app.Use("/task", basicauth.New(basicauth.Config{
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			log.Println("In the middleware")
			currentUser := database.VerifyUser(user, pass)
			if user == currentUser.Email && pass == currentUser.Password {
				return true
			}
			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendFile("./unauthorized.html")
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}))

	var todo []models.Task = make([]models.Task, 0)
	//Get all the tasks to do

	app.Get("/task", func(c *fiber.Ctx) error {
		todo = database.GetTasks()

		if c.Response().StatusCode() == 200 {
			return c.JSON(todo)
		}

		return fiber.NewError(c.Response().StatusCode(), "There was an Error")
	})

	//Adding a Task to the to do list
	app.Post("/task", func(c *fiber.Ctx) error {
		var body models.Task
		if err := c.BodyParser(&body); err != nil {
			log.Println("There is an error", err)
			return err
		}
		if &body == nil {
			log.Println("The value is nil")
		}
		addTask, err := database.AddTask(body.TaskName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		result := &models.Task{}
		if err := copier.Copy(&result, &addTask); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "cannot map results",
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    result,
		})
	})

	//Deleting a Task from the to do list
	app.Delete("/task/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		err := database.DeleteTask(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Resource not found",
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"success": true,
			"message": "Task was deleted",
		})
	})

	//Below are the User Endpoints
	app.Get("/user/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")
		e, User := database.GetUser(id)
		if e != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   e,
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"success": true,
			"user":    User,
		})
	})

	app.Post("/addUser", func(c *fiber.Ctx) error {
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
	})

	app.Delete("/user/delete/:id", func(c *fiber.Ctx) error {
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

	})

	app.Listen(":3000")
}
