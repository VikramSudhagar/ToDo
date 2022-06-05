package controller

import (
	"strconv"
	"todo/database"
	"todo/database/cache"
	"todo/middleware"
	"todo/models"

	"github.com/gofiber/fiber/v2"
)

var store = cache.CacheSetUp()

func GetTask(c *fiber.Ctx) error {
	user, err := middleware.RetrieveSessionAndVerify(store, c, c.Cookies("sessionID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}
	var DTO_Response_Array []models.DTO_Task = make([]models.DTO_Task, 0)
	DTO_Response_Array, e := database.GetTaskbyUserID(user.ID)

	if e != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   e,
		})
	}

	return c.Render("todo", fiber.Map{
		"DTO_Response": DTO_Response_Array,
	})
}

func AddTask(c *fiber.Ctx) error {
	//RetrieveSessionAndVerify we check whether the information in the
	//session is valid, and will then verify if this user exists in the DB
	user, err := middleware.RetrieveSessionAndVerify(store, c, c.Cookies("sessionID"))
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

	addTask, err := database.AddTask(body.TaskName, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    addTask,
	})
}

func DeleteTask(c *fiber.Ctx) error {
	if _, err := middleware.RetrieveSessionAndVerify(store, c, c.Cookies("sessionID")); err != nil {
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

	if err := database.DeleteTask(id, c); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": "Task was deleted",
	})
}
