package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User has many CreditCards, UserID is the foreign key

type User14 struct {
	gorm.Model
	CreditCards []CreditCard3
}

type CreditCard3 struct {
	gorm.Model
	Number string
	UserID uint
}

// Override Foreign Key

type User15 struct {
	gorm.Model
	CreditCards []CreditCard4 `gorm:"foreignKey:UserRefer"`
}

type CreditCard4 struct {
	gorm.Model
	Number    string
	UserRefer uint
}

//Override References

type User16 struct {
	gorm.Model
	MemberNumber string
	CreditCards  []CreditCard `gorm:"foreignKey:UserNumber;references:MemberNumber"`
}

type CreditCard5 struct {
	gorm.Model
	Number     string
	UserNumber string
}

// Polymorphism Association

type Dog2 struct {
	ID   int
	Name string
	Toys []Toy2 `gorm:"polymorphic:Owner;"`
}

type Toy2 struct {
	ID        int
	Name      string
	OwnerID   int
	OwnerType string
}

// Self-Referential Has Many

type User17 struct {
	gorm.Model
	Name      string
	ManagerID *uint
	Team      []User `gorm:"foreignkey:ManagerID"`
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
	db.Create(&Dog2{Name: "dog1", Toys: []Toy2{{Name: "toy1"}, {Name: "toy2"}}})
	// INSERT INTO `dogs` (`name`) VALUES ("dog1")
	// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy1","1","dogs"), ("toy2","1","dogs")
}

// Retrieve user list with eager loading credit cards
func getAll(db *gorm.DB) ([]User14, error) {
	var users []User14
	err := db.Model(&User{}).Preload("CreditCards").Find(&users).Error
	return users, err
}
