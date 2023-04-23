package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
	"time"
)

type User11 struct {
	ID        uint
	Name      string
	Age       int
	Gender    string
	processed bool
}

type Order2 struct {
	UserId     int
	Amount     int
	FinishedAt *time.Time
}

type APIUser struct {
	ID   uint
	Name string
}

type Pet struct {
	ID   uint
	Name string
}

type Pizza struct {
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 100,
		//QueryFields:     true,
	})
	// 开启 GORM 的 SQL 日志
	db = db.Debug()
	//createTableIFNotExist(db)
	//smartSelectFields(db)
	//locking(db)
	//subQuery(db)
	//fromSubQuery(db)
	//groupConditions(db)
	//inWithMultipleColumns(db)
	//namedArgument(db)
	//findToMap(db)
	//firstOrInit(db)
	//firstOrCreate(db)
	//iteration(db)
	//findInBatches(db)
	//pluck(db)
	//scopes(db)
	count(db)
}

func createTableIFNotExist(db *gorm.DB) {
	db.AutoMigrate(&User11{})
	db.AutoMigrate(&Order2{})
}

func smartSelectFields(db *gorm.DB) {
	// Select `id`, `name` automatically when querying
	db.Model(&User11{}).Limit(10).Find(&APIUser{})
	// SELECT `id`, `name` FROM `users` LIMIT 10

	var user User11
	db.Find(&user)
	//SELECT * FROM `user`

	db.Session(&gorm.Session{QueryFields: true}).Find(&user)
	//SELECT `user11`.`id`,`user11`.`name`,`user11`.`age`,`user11`.`gender` FROM `user11`

}

func locking(db *gorm.DB) {
	var users []User11
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
	// SELECT * FROM `users` FOR UPDATE

	db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Find(&users)
	// SELECT * FROM `users` LOCK IN SHARE MODE

	db.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "NOWAIT",
	}).Find(&users)
	// SELECT * FROM `users` FOR UPDATE NOWAIT
}

func subQuery(db *gorm.DB) {
	var orders []Order2
	db.Where("amount > (?)", db.Table("order2").Select("AVG(amount)")).Find(&orders)
	// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

	var users []User11
	subQuery := db.Select("AVG(age)").Where("name LIKE ?", "name%").Table("user11")
	db.Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&users)
	// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")

}

func fromSubQuery(db *gorm.DB) {
	var users []User11
	db.Table("(?) as u", db.Model(&User11{}).Select("name", "age")).Where("age = ?", 18).Find(&users)
	// SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18

	subQuery1 := db.Model(&User{}).Select("name")
	subQuery2 := db.Model(&Pet{}).Select("name")
	db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&users)
	// SELECT * FROM (SELECT `name` FROM `users`) as u, (SELECT `name` FROM `pets`) as p
}

func groupConditions(db *gorm.DB) {
	db.Where(
		db.Where("pizza = ?", "pepperoni").Where(db.Where("size = ?", "small").Or("size = ?", "medium")),
	).Or(
		db.Where("pizza = ?", "hawaiian").Where("size = ?", "xlarge"),
	).Find(&Pizza{})
	// SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")

}

func inWithMultipleColumns(db *gorm.DB) {
	var users []User11
	db.Where("(name, age, role) IN ?", [][]interface{}{{"jinzhu", 18, "admin"}, {"jinzhu2", 19, "user"}}).Find(&users)
	// SELECT * FROM users WHERE (name, age, role) IN (("jinzhu", 18, "admin"), ("jinzhu 2", 19, "user"));
}

func namedArgument(db *gorm.DB) {
	var user User11
	db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

	db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu"}).First(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu" ORDER BY `users`.`id` LIMIT 1

}

func findToMap(db *gorm.DB) {
	result := map[string]interface{}{}
	db.Model(&User11{}).First(&result, "id = ?", 1)

	var results []map[string]interface{}
	db.Table("users").Find(&results)
}

func firstOrInit(db *gorm.DB) {
	var user User11
	// User not found, initialize it with give conditions
	db.FirstOrInit(&user, User11{Name: "non_existing"})
	// user -> User{Name: "non_existing"}

	// Found user with `name` = `jinzhu`
	db.Where(User11{Name: "jinzhu"}).FirstOrInit(&user)
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// Found user with `name` = `jinzhu`
	db.FirstOrInit(&user, map[string]interface{}{"name": "jinzhu"})
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// User not found, initialize it with give conditions and Attrs
	db.Where(User11{Name: "non_existing"}).Attrs(User11{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 20}

	// User not found, initialize it with give conditions and Attrs
	db.Where(User11{Name: "non_existing"}).Attrs("age", 20).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 20}

	// Found user with `name` = `jinzhu`, attributes will be ignored
	db.Where(User11{Name: "Jinzhu"}).Attrs(User11{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// User not found, initialize it with give conditions and Assign attributes
	db.Where(User11{Name: "non_existing"}).Assign(User11{Age: 20}).FirstOrInit(&user)
	// user -> User{Name: "non_existing", Age: 20}

	// Found user with `name` = `jinzhu`, update it with Assign attributes
	db.Where(User11{Name: "Jinzhu"}).Assign(User11{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "Jinzhu", Age: 20}

}

func firstOrCreate(db *gorm.DB) {
	var user User11
	// User not found, create a new record with give conditions
	result := db.FirstOrCreate(&user, User{Name: "non_existing"})
	fmt.Println(result.RowsAffected)
	// INSERT INTO "users" (name) VALUES ("non_existing");
	// user -> User{ID: 112, Name: "non_existing"}
	// result.RowsAffected // => 1

	// Found user with `name` = `jinzhu`
	result = db.Where(User11{Name: "jinzhu"}).FirstOrCreate(&user)
	// user -> User{ID: 111, Name: "jinzhu", "Age": 18}
	// result.RowsAffected // => 0

	// User not found, create it with give conditions and Attrs
	db.Where(User11{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{ID: 112, Name: "non_existing", Age: 20}

	// Found user with `name` = `jinzhu`, attributes will be ignored
	db.Where(User11{Name: "jinzhu"}).Attrs(User11{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "jinzhu", Age: 18}

	// User not found, initialize it with give conditions and Assign attributes
	db.Where(User11{Name: "non_existing"}).Assign(User11{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{ID: 112, Name: "non_existing", Age: 20}

	// Found user with `name` = `jinzhu`, update it with Assign attributes
	db.Where(User11{Name: "jinzhu"}).Assign(User11{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;
	// UPDATE users SET age=20 WHERE id = 111;
	// user -> User{ID: 111, Name: "jinzhu", Age: 20}
}

func optimizerAndIndexHints(db *gorm.DB) {
	db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&User{})
	// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`

	db.Clauses(hints.UseIndex("idx_user_name")).Find(&User11{})
	// SELECT * FROM `users` USE INDEX (`idx_user_name`)

	db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User11{})
	// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"
}

func iteration(db *gorm.DB) {
	rows, _ := db.Model(&User{}).Where("name = ?", "jinzhu").Rows()
	defer rows.Close()

	for rows.Next() {
		var user User
		// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		db.ScanRows(rows, &user)
		// do something
	}
}

func findInBatches(db *gorm.DB) {
	var users []User11
	// batch size 100
	result := db.Where("processed = ?", false).FindInBatches(&users, 100, func(tx *gorm.DB, batch int) error {
		for _, user := range users {
			// batch processing found records
			user.processed = true
		}
		tx.Save(&users)
		fmt.Println(tx.RowsAffected) // number of records in this batch
		fmt.Println(batch)           // Batch 1, 2, 3
		// returns error will stop future batches
		return nil
	})
	//result.Error // returned error
	//result.RowsAffected // processed records count in all batches
	fmt.Println(result.Error, result.RowsAffected)
}

func pluck(db *gorm.DB) {
	//Query single column from database and scan into a slice
	var ages []int64
	db.Model(&User11{}).Pluck("age", &ages)

	var names []string
	db.Model(&User11{}).Pluck("name", &names)

	db.Table("user11").Pluck("name", &names)

	// Distinct Pluck
	db.Model(&User11{}).Distinct().Pluck("Name", &names)
	// SELECT DISTINCT `name` FROM `users`

	// Requesting more than one column, use `Scan` or `Find` like this:
	var users []User11
	db.Select("name", "age").Scan(&users)
	db.Select("name", "age").Find(&users)
}

func scopes(db *gorm.DB) {
	var orders []Order2
	db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
	// Find all credit card orders and amount greater than 1000

	db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)
	// Find all COD orders and amount greater than 1000

	db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
	// Find all paid, shipped orders that amount greater than 1000
}

func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}

func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN (?)", status)
	}
}

func count(db *gorm.DB) {
	var count int64
	db.Model(&User11{}).Where("name = ?", "jinzhu").Or("name = ?", "jinzhu 2").Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'

	db.Model(&User11{}).Where("name = ?", "jinzhu").Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu';

	db.Table("user11").Count(&count)
	// SELECT count(1) FROM deleted_users;

	// Count with Distinct
	db.Model(&User11{}).Distinct("name").Count(&count)
	// SELECT COUNT(DISTINCT(`name`)) FROM `users`

	db.Table("user11").Select("count(distinct(name))").Count(&count)
	// SELECT count(distinct(name)) FROM users

	// Count with Group
	db.Model(&User11{}).Group("name").Count(&count)
}
