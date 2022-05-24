package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	TaskName string `json:"taskname"`
}

type DTO_Task struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	TaskName  string         `json:"taskname"`
}
