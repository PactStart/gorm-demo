package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Active       bool
	Role         string
	Password     string
	Location     Location
	Birthday     time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type User3 User

type User2 struct {
	ID   uint
	Name string
}

type User4 struct {
	ID       uint
	Name     string
	Location Location
}

type Location struct {
	X, Y int
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	return nil
}

func (loc Location) GormDataType() string {
	return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
	}
}

func (u *User3) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("BeforeCreate")
	if u.Name == "admin" {
		return errors.New("invalid name")
	}
	return
}

func (u *User3) BeforeSave(tx *gorm.DB) (err error) {
	fmt.Println("BeforeSave")
	if u.Name == "admin" {
		return errors.New("invalid name")
	}
	return
}

func (u *User3) AfterSave(tx *gorm.DB) (err error) {
	fmt.Println("AfterSave")
	if u.Name == "admin" {
		return errors.New("invalid name")
	}
	return
}

func (u *User3) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("AfterCreate")
	if u.Name == "admin" {
		return errors.New("invalid name")
	}
	return
}

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&User2{})
	db.AutoMigrate(&User3{})
	db.AutoMigrate(&User4{})
	createRecord(db)
	createRecordWithSelectedFields(db)
	createRecordWithOmitFields(db)
	batchInsert(db)
	CreateInBatches(db)
	createWithHooks(db)
	createIgnoreHooks(db)
	createFromMap(db)
	createFromSQLExpr(db)
}

func createRecord(db *gorm.DB) {
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	result := db.Create(&user) // pass pointer of data to Create
	fmt.Println(user.ID, result.Error, result.RowsAffected)
	fmt.Println("Inserted user:", user)

	users := []*User{
		&User{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
		&User{Name: "Jackson", Age: 19, Birthday: time.Now()},
	}
	result2 := db.Create(users) // pass a slice to insert multiple row
	fmt.Println(result2.Error, result2.RowsAffected)

}

func createRecordWithSelectedFields(db *gorm.DB) {
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	db.Select("Name", "Age", "CreatedAt").Create(&user)
	fmt.Println(user)
}

func createRecordWithOmitFields(db *gorm.DB) {
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	db.Omit("Name", "Age", "CreatedAt").Create(&user)
}

func batchInsert(db *gorm.DB) {
	var users = []User2{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	db.Create(&users)

	for _, user := range users {
		fmt.Println(user)
	}
}

func CreateInBatches(db *gorm.DB) {
	var users []User2
	for i := 0; i < 100; i++ {
		user := User2{Name: "User" + strconv.Itoa(i)}
		users = append(users, user)
	}
	// batch size 10
	db.CreateInBatches(users, 10)
}

func createWithHooks(db *gorm.DB) {
	user := User3{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	result := db.Create(&user) // pass pointer of data to Create
	fmt.Println(user.ID, result.Error, result.RowsAffected)
	fmt.Println("Inserted user:", user)

}

func createIgnoreHooks(db *gorm.DB) {
	user := User3{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	result := db.Session(&gorm.Session{SkipHooks: true}).Create(&user) // pass pointer of data to Create
	fmt.Println(user.ID, result.Error, result.RowsAffected)
	fmt.Println("Inserted user:", user)

}

func createFromMap(db *gorm.DB) {
	db.Model(&User{}).Create(map[string]interface{}{
		"Name": "jinzhu", "Age": 18,
	})

	// batch insert from `[]map[string]interface{}{}`
	db.Model(&User{}).Create([]map[string]interface{}{
		{"Name": "jinzhu_1", "Age": 18},
		{"Name": "jinzhu_2", "Age": 20},
	})
}

func createFromSQLExpr(db *gorm.DB) {
	// INSERT INTO `users` (`name`,`location`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"))
	db.Create(&User4{
		Name:     "jinzhu",
		Location: Location{X: 100, Y: 100},
	})
	// Create from map
	// INSERT INTO `users` (`name`,`location`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"));
	db.Model(User4{}).Create(map[string]interface{}{
		"Name":     "jinzhu",
		"Location": clause.Expr{SQL: "ST_PointFromText(?)", Vars: []interface{}{"POINT(100 100)"}},
	})
}
