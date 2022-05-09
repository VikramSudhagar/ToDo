package database

import (
	"log"
	"todo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB
var err error

func SetUp() {
	dsn := "host=127.0.0.1 user=postgres password=new@co-op222! dbname=todoDB port=5432"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Println("DB Connection was successful")
	} else {
		log.Println("The error is: ", err)
	}
}

func GetTasks() []models.Task {
	var tasks []models.Task
	db := DBConn
	db.Find(&tasks)
	for _, value := range tasks {
		log.Println(value)
	}
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
