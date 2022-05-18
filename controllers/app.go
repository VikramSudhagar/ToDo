package main

import (
	"log"
	"strconv"
	"todo/database"
	"todo/database/cache"
	"todo/middleware"
	"todo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func main() {
	app := fiber.New()
	store := cache.CacheSetUp()
	database.TaskSetUp()
	database.UserSetUp()

	app.Get("/welcome", func(c *fiber.Ctx) error {
		return c.SendString("Welcome")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		if _, err := database.VerifyUser(user.Email, user.Password); err != nil {
			return err
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
		store.Set("session", []byte(value), 0)
		return c.SendString("Hello, World!")
	})

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
		//RetrieveSessionAndVerify we check whether the information in the
		//session is valid, and will then verify if this user exists in the DB
		err := middleware.RetrieveSessionAndVerify(store, c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"Error":   err,
			})
		}

		var body models.Task
		//Check if there was an error with parsing the body of the request
		if err := c.BodyParser(&body); err != nil {
			return err
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
		if err := middleware.RetrieveSessionAndVerify(store, c); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"Error":   err,
			})
		}

		integer := c.Params("id")

		id, e := strconv.Atoi(integer)
		if e != nil {
			return e
		}

		if err := database.DeleteTask(id); err != nil {
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

	app.Listen(":8081")
}
