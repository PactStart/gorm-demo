package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User7 struct {
	ID        int64
	Age       int64        `gorm:"default:18"`
	Active    sql.NullBool `gorm:"default:true"`
	FirstName string       `gorm:"default:Rex" json:"first_name"`
	LastName  string       `gorm:"default:Lei" json:"last_name"`
	FullName  string       `gorm:"->;type:varchar(255);GENERATED ALWAYS AS (concat(firstname,' ',lastname))";default:"-"`
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User7{})
	db.Create(&User7{})

	var user User7
	db.Find(&user, 1)
	fmt.Println(user)
}
