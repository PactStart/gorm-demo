package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Role string
}

type Company struct {
	gorm.Model
	Name string
}
