package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// `Staff` belongs to `Department`, `DepartmentID` is the foreign key
type Staff struct {
	gorm.Model
	Name         string
	DepartmentID int // By default, the DepartmentID is implicitly used to create a foreign key relationship between the Staff and Department tables
	department   Department
}

type Staff2 struct {
	gorm.Model
	Name            string
	DepartmentRefer int
	department      Department `gorm:"foreignKey:DepartmentRefer"`
	// use DepartmentRefer as foreign key
}

type Staff3 struct {
	gorm.Model
	Name         string
	DepartmentID string
	Department   Department `gorm:"references:Code"` // use Code as references
}

type Staff4 struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Department struct {
	ID   int
	Name string
	Code string
}

// A belongs to association sets up a one-to-one connection with another model, such that each instance of the declaring model “belongs to” one instance of the other model.
func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.Debug()
}
