package repo

import (
	"log"
	"todo/database"
	"todo/models"
)

func GetAllTasks() []models.Task {
	db := database.DBConn
	log.Println("The value of database.DBConn is: ", database.DBConn)
	var tasks []models.Task
	if db == nil {
		log.Println("DB is null")
	}
	db.Find(&tasks)
	return tasks
}
