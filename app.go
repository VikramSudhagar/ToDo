package main

import (
	"log"
	"strconv"
	"time"
	"todo/database"
	"todo/database/cache"
	"todo/middleware"
	"todo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/jinzhu/copier"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	store := cache.CacheSetUp()
	session := session.New()
	database.TaskSetUp()
	database.UserSetUp()

	app.Static("/", "./views", fiber.Static{
		Index: "index.html",
	})

	app.Static("/static", "./static")

	app.Post("/login", func(c *fiber.Ctx) error {
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
	})
	//Get the User's tasks
	app.Get("/task", func(c *fiber.Ctx) error {
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
	})

	//Adding a Task to the to do list
	app.Post("/task", func(c *fiber.Ctx) error {
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
	})

	//Deleting a Task from the to do list
	app.Delete("/task/:id", func(c *fiber.Ctx) error {
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
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
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

	app.Post("/logout", func(c *fiber.Ctx) error {
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
	})

	app.Listen(":8081")
}
