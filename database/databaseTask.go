package database

import (
	"log"
	"todo/models"

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

func GetTasks() []models.Task {
	var tasks []models.Task
	db := DBConn
	db.Find(&tasks)
	return tasks
}

func AddTask(taskname string) (*models.Task, error) {
	db := DBConn
	task := &models.Task{
		TaskName: taskname,
	}

	log.Println("The memory address is: ", &task)
	if &task == nil {
		log.Println("The value is nil")
	}

	db.Create(&task)

	return task, nil
}

func DeleteTask(id int) error {
	db := DBConn
	//delete the task with that specific ID. Every task has
	//a primary key, so a Batch Delete will not be triggered
	if err := db.Delete(&models.Task{}, id).Error; err != nil {
		log.Println("Deletion was not successful")
		return err
	}
	log.Println("Deletion was successful")
	return nil
}
