package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Member has one BankCard, MemberID is the foreign key
type Member struct {
	gorm.Model
	CreditCard CreditCard
}

type BankCard struct {
	gorm.Model
	Number   string
	memberId uint
}

// Override Foreign Key

type Member2 struct {
	gorm.Model
	BankCard BankCard2 `gorm:"foreignKey:MemberName"`
	// use MemberName as foreign key
}

type BankCard2 struct {
	gorm.Model
	Number     string
	MemberName string
}

// Override References

type Member3 struct {
	gorm.Model
	Name     string    `gorm:"index"`
	BankCard BankCard3 `gorm:"foreignKey:MemberName;references:name"`
}

type BankCard3 struct {
	gorm.Model
	Number     string
	MemberName string
}

// Polymorphism Association

type Cat struct {
	ID   int
	Name string
	Toy  Toy `gorm:"polymorphic:Owner;"`
}

type Dog struct {
	ID   int
	Name string
	Toy  Toy `gorm:"polymorphic:Owner;"` //`gorm:"polymorphic:Owner;polymorphicValue:master"`
}

type Toy struct {
	ID        int
	Name      string
	OwnerID   int
	OwnerType string
}

// Self-Referential Has One

type User13 struct {
	gorm.Model
	Name      string
	ManagerID *uint
	Manager   *User13
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

	db.Create(&Dog{Name: "dog1", Toy: Toy{Name: "toy1"}})
	// INSERT INTO `dogs` (`name`) VALUES ("dog1")
	// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy1","1","dogs")
	// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy1","1","master")
}

// Retrieve member list with eager loading bank card

func GetAll(db *gorm.DB) ([]Member, error) {
	var members []Member
	err := db.Model(&User{}).Preload("BankCard").Find(&members).Error
	return members, err
}
