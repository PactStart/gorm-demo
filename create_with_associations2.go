package main

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CreditCard2 struct {
	gorm.Model
	Number string
}

func (c *CreditCard2) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *CreditCard2) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type User6 struct {
	gorm.Model
	Name       string
	CreditCard *CreditCard2 `gorm:"type:JSON"`
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User6{})

	db.Create(&User6{
		Name: "jinzhu",
		CreditCard: &CreditCard2{
			Number: "411111111111",
		},
	})
}
