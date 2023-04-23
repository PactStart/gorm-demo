package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

type User5 struct {
	gorm.Model
	Name         string
	CreditCard   *CreditCard
	CreditCardID uint
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User5{}, &CreditCard{})

	// INSERT INTO `users` ...
	// INSERT INTO `credit_cards` ...
	db.Create(&User5{
		Name:       "jinzhu",
		CreditCard: &CreditCard{Number: "411111111111"},
	})

	db.Omit("CreditCard").Create(&User5{
		Name:       "jinzhu2",
		CreditCard: &CreditCard{Number: "411111111111"},
	})
	// skip all associations
	db.Omit(clause.Associations).Create(&User5{
		Name:       "jinzhu3",
		CreditCard: &CreditCard{Number: "411111111111"},
	})

}
