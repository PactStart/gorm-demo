package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// Set field `UUID` as primary field
type Animal struct {
	ID        int64  `gorm:"column:beast_id"` // set name to `beast_id`
	UUID      string `gorm:"primaryKey"`
	Name      string
	Age       int64
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

// TableName overrides the table name used by Animal to `ani`
func (Animal) TableName() string {
	return "ani"
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.Debug()
	// Create table for `User`
	db.Migrator().CreateTable(&Animal{})
	// Append "ENGINE=InnoDB" to the creating table SQL for `User`
	db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&Animal{})
}
