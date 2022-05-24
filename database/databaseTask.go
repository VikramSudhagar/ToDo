package database

import (
	"log"
	"todo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB
var err error

func TaskSetUp() {
	dsn := "host=host.docker.internal user=postgres password=mypassword dbname=todoDB port=5432"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Println("DB Connection was successful")
		DBConn.AutoMigrate(models.Task{})
	} else {
		log.Println("The error is: ", err)
	}
}

//Get All Tasks
func GetTasks() []models.Task {
	var tasks []models.Task
	db := DBConn
	db.Find(&tasks)
	return tasks
}

func AddTask(taskname string, userID uint) (*models.Task, error) {
	db := DBConn
	task := &models.Task{
		TaskName: taskname,
		UserID:   userID,
	}

	log.Println("The memory address is: ", &task)
	if &task == nil {
		log.Println("The value is nil")
	}

	db.Create(&task)

	return task, nil
}

func GetTaskbyUserID(userID uint) ([]models.DTO_Task, error) {
	db := DBConn
	log.Println("The userID is: ", userID)
	var tasks []models.Task
	results := []models.DTO_Task{}
	if err := db.Where("user_id= ?", userID).Find(&tasks).Error; err != nil {
		return results, err
	}
	copier.Copy(&results, &tasks)
	return results, nil
}

func getTask(id int, c *fiber.Ctx) models.DTO_Task {
	db := DBConn
	task := models.Task{}
	result := models.DTO_Task{}
	db.First(&task, id)
	copier.Copy(&result, &task)
	return result
}

func DeleteTask(id int, c *fiber.Ctx) error {
	db := DBConn
	//delete the task with that specific ID. Every task has
	//a primary key, so a Batch Delete will not be triggered

	// isDeleted := c.JSON(getTask(id, c).DeletedAt)
	// log.Println("isDeleted value is: ", isDeleted)
	// if isDeleted == nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"success": false,
	// 		"Error":   "Resoruce not found",
	// 	})
	// }

	if err := db.Delete(&models.Task{}, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"Error":   "Deletion was not successful",
		})
	}
	log.Println("Deletion was successful")
	return nil
}
