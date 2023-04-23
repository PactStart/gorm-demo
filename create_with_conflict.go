package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User8 struct {
	ID   uint
	Name string
	Age  uint
	Role string
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User8{})
	user := &User8{
		ID:   1,
		Name: "zhangsan",
		Age:  18,
		Role: "Admin",
	}
	//db.Create(&user)
	// Do nothing on conflict
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
	// Update columns to default value on `id` conflict
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"role": "user"}),
	}).Create(&user)
	// Use SQL expression
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"name": gorm.Expr("UPPER(name)")}),
	}).Create(&user)
	// Update columns to new value on `id` conflict
	user.Age = 20
	user.Name = "lisi"
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
	}).Create(&user)
	// Update all columns to new value on conflict except primary keys and those columns having default values from sql func
	user.Age = 30
	user.Name = "wangwu"
	user.Role = "boss"
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user)

}
