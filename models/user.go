package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	email    string
	password string
}
