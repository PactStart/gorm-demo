package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	Birthday     time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func main() {

	preparedStmt()
}

func preparedStmt() {
	// globally mode, all DB operations will create prepared statements and cache them
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	var user User
	var users []User
	// session mode
	tx := db.Session(&gorm.Session{PrepareStmt: true})
	tx.First(&user, 1)
	tx.Find(&users)
	tx.Model(&user).Update("Age", 18)

	// returns prepared statements manager
	stmtManger, _ := tx.ConnPool.(*gorm.PreparedStmtDB)

	// prepared SQL for *current session*
	//stmtManger.PreparedSQL // => []string{}
	fmt.Println(stmtManger.PreparedSQL)

	// prepared statements for current database connection pool (all sessions)
	//stmtManger.Stmts // map[string]*sql.Stmt

	for sql, stmt := range stmtManger.Stmts {
		fmt.Println(sql, stmt)
		//sql          // prepared SQL
		//stmt         // prepared statement
		stmt.Close() // close the prepared statement
	}

	// close prepared statements for *current session*
	stmtManger.Close()
}

func newDb() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	tx := db.Where("name = ?", "jinzhu").Session(&gorm.Session{NewDB: true})
	var user User
	tx.First(&user)
	// SELECT * FROM users ORDER BY id LIMIT 1

	tx.First(&user, "id = ?", 10)
	// SELECT * FROM users WHERE id = 10 ORDER BY id

	// Without option `NewDB`
	tx2 := db.Where("name = ?", "jinzhu").Session(&gorm.Session{})
	tx2.First(&user)
	// SELECT * FROM users WHERE name = "jinzhu" ORDER BY id
}

func skipHooks(db *gorm.DB) {
	var user User
	var users []User
	db.Session(&gorm.Session{SkipHooks: true}).Create(&user)

	db.Session(&gorm.Session{SkipHooks: true}).Create(&users)

	db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)

	db.Session(&gorm.Session{SkipHooks: true}).Find(&user)

	db.Session(&gorm.Session{SkipHooks: true}).Delete(&user)

	db.Session(&gorm.Session{SkipHooks: true}).Model(User{}).Where("age > ?", 18).Updates(&user)

}
