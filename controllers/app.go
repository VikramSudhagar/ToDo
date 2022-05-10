package main

import (
	"log"
	"todo/database"
	"todo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func main() {
	app := fiber.New()
	var todo []models.Task = make([]models.Task, 0)
	//Get all the tasks to do
	app.Get("/", func(c *fiber.Ctx) error {
		database.SetUp()
		todo = database.GetTasks()

		if c.Response().StatusCode() == 200 {
			return c.JSON(todo)
		}

		return fiber.NewError(c.Response().StatusCode(), "There was an Error")
	})

	//Adding a Task to the to do list
	app.Post("/addTask", func(c *fiber.Ctx) error {
		var body models.Task
		if err := c.BodyParser(&body); err != nil {
			log.Println("There is an error", err)
			return err
		}

		log.Println("Before Create", body.TaskName)
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
	app.Delete("/deleteTask", func(c *fiber.Ctx) error {
		var body models.Task
		if err := c.BodyParser(&body); err != nil {
			log.Println("There is an error with deletion")
			return err
		}
		tasks, err := database.DeleteTask(body)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Resource not found",
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"success": true,
			"data":    tasks,
		})
	})

	app.Listen(":3000")
}
