package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/hints"
	"time"
)

type Result2 struct {
	ID   int
	Name string
	Age  int
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 100,
	})
	// 开启 GORM 的 SQL 日志
	db = db.Debug()

	rawSql(db)
	namedAgrumment(db)
	dryRun(db)
	toSql(db)
	rowAndRows(db)
	scanIntoStruct(db)
	connection(db)
	statementModifier(db)
}

func rawSql(db *gorm.DB) {

	var result Result2
	db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)

	db.Raw("SELECT id, name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)

	var age int
	db.Raw("SELECT SUM(age) FROM users WHERE role = ?", "admin").Scan(&age)

	var users []User
	db.Raw("UPDATE users SET name = ? WHERE age = ? RETURNING id, name", "jinzhu", 20).Scan(&users)

	db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})

	// Exec with SQL Expression
	db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")

	db.Exec("DROP TABLE users")

}

func namedAgrumment(db *gorm.DB) {
	var user User
	db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

	db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu2"}).First(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu2" OR name2 = "jinzhu2" ORDER BY `users`.`id` LIMIT 1

	// Named Argument with Raw SQL
	db.Raw("SELECT * FROM users WHERE name1 = @name OR name2 = @name2 OR name3 = @name",
		sql.Named("name", "jinzhu1"), sql.Named("name2", "jinzhu2")).Find(&user)
	// SELECT * FROM users WHERE name1 = "jinzhu1" OR name2 = "jinzhu2" OR name3 = "jinzhu1"

	db.Exec("UPDATE users SET name1 = @name, name2 = @name2, name3 = @name",
		sql.Named("name", "jinzhunew"), sql.Named("name2", "jinzhunew2"))
	// UPDATE users SET name1 = "jinzhunew", name2 = "jinzhunew2", name3 = "jinzhunew"

	db.Raw("SELECT * FROM users WHERE (name1 = @name AND name3 = @name) AND name2 = @name2",
		map[string]interface{}{"name": "jinzhu", "name2": "jinzhu2"}).Find(&user)
	// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"

	type NamedArgument struct {
		Name  string
		Name2 string
	}

	db.Raw("SELECT * FROM users WHERE (name1 = @Name AND name3 = @Name) AND name2 = @Name2",
		NamedArgument{Name: "jinzhu", Name2: "jinzhu2"}).Find(&user)
	// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"
}

// Generate SQL and its arguments without executing, can be used to prepare or test generated SQL
func dryRun(db *gorm.DB) {
	var user User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	fmt.Println(stmt.SQL.String()) //=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
	fmt.Println(stmt.Vars)         //=> []interface{}{1}
}

func toSql(db *gorm.DB) {
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&User{}).Where("id = ?", 100).Limit(10).Order("age desc").Find(&[]User{})
	})
	fmt.Println(sql) //=> SELECT * FROM "users" WHERE id = 100 AND "users"."deleted_at" IS NULL ORDER BY age desc LIMIT 10
}

func rowAndRows(db *gorm.DB) {
	var age int
	var name string
	var email string
	// Use GORM API build SQL
	row := db.Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row()
	row.Scan(&name, &age)

	// Use Raw SQL
	row2 := db.Raw("select name, age, email from users where name = ?", "jinzhu").Row()
	row2.Scan(&name, &age, &email)

	// Use GORM API build SQL
	rows, _ := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name, &age, &email)
		// do something
	}

	// Raw SQL
	rows2, _ := db.Raw("select name, age, email from users where name = ?", "jinzhu").Rows()
	defer rows2.Close()
	for rows2.Next() {
		rows2.Scan(&name, &age, &email)
		// do something
	}
}

func scanIntoStruct(db *gorm.DB) {
	rows, _ := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows() // (*sql.Rows, error)
	defer rows.Close()

	var user User
	for rows.Next() {
		// ScanRows scan a row into user
		db.ScanRows(rows, &user)
		// do something
	}

}

// Run mutliple SQL in same db tcp connection (not in a transaction)
func connection(db *gorm.DB) {
	db.Connection(func(tx *gorm.DB) error {
		tx.Exec("SET my.role = ?", "admin")
		tx.First(&User{})
		return nil
	})

}

func statementModifier(db *gorm.DB) {
	db.Clauses(hints.New("hint")).Find(&User{})
	// SELECT * /*+ hint */ FROM `users`
}
