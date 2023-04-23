package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type User9 struct {
	gorm.Model
	Name      string
	Age       int
	Email     string
	CompanyId int
}

type Order struct {
	UserId     int
	FinishedAt *time.Time
}

type Company struct {
	ID    int
	Name  string
	Alive bool
}

type User10 struct {
	Name string
}

type Language struct { // no primary key defined, results will be ordered by first field (i.e., `Code`)
	Code string
	Name string
}

type Result struct {
	Name  string
	Total int
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 100,
	})
	// 开启 GORM 的 SQL 日志
	db = db.Debug()

	if err != nil {
		panic("failed to connect database")
	}
	//createTableIfNotExist(db)
	//retrieveSingObject(db)
	//queryWithNoPrimaryKey(db)
	//retrieveObjectsWithPrimaryKey(db)
	//retrieveAllObjects(db)
	//stringConditions(db)
	//structAndMapConditions(db)
	//specStructSearchFields(db)
	//inlineCondition(db)
	//notConditions(db)
	//orConditions(db)
	//selectSpecificFields(db)
	//order(db)
	//limitAndOffset(db)
	//groupByAndHaving(db)
	//distinct(db)
	//joins(db)
	//joinPreloading(db)
	//joinADerivedTable(db)
	scan(db)

}

func init() {
	fmt.Println("init...")
}

func createTableIfNotExist(db *gorm.DB) {
	db.AutoMigrate(&User9{})
	var users []User9
	for i := 0; i < 100; i++ {
		user := User9{Name: "User" + strconv.Itoa(i+1), Age: i + 1}
		users = append(users, user)
	}
	db.Create(users)
	db.AutoMigrate(&Company{})
	db.AutoMigrate(&Order{})
}

func retrieveSingObject(db *gorm.DB) {
	var user1 User9
	// Get the first record ordered by primary key
	db.First(&user1)
	fmt.Println(user1)
	// SELECT * FROM users ORDER BY id LIMIT 1;

	// Get one record, no specified order
	var user2 User9
	db.Take(&user2)
	fmt.Println(user2)
	// SELECT * FROM users LIMIT 1;

	// Get last record, ordered by primary key desc
	var user3 User9
	db.Last(&user3)
	fmt.Println(user3)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;

	// check error ErrRecordNotFound
	var user4 User9
	result := db.First(&user4)
	fmt.Println(result.RowsAffected)
	errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func queryWithNoPrimaryKey(db *gorm.DB) {
	db.AutoMigrate(&User10{}, &Language{})
	// works because destination struct is passed in
	var user User10
	db.First(&user)
	// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1

	// works because model is specified using `db.Model()`
	result1 := map[string]interface{}{}
	db.Model(&User10{}).First(&result1)
	// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1

	// doesn't work
	result2 := map[string]interface{}{}
	db.Table("users10").First(&result2)

	// works with Take
	result3 := map[string]interface{}{}
	db.Table("users10").Take(&result3)

	db.First(&Language{})
	// SELECT * FROM `languages` ORDER BY `languages`.`code` LIMIT 1
}

func retrieveObjectsWithPrimaryKey(db *gorm.DB) {
	var user1 User9
	db.First(&user1, 10)
	fmt.Println(user1)
	// SELECT * FROM users WHERE id = 10;

	var user2 User9
	db.First(&user2, "10")
	fmt.Println(user2)
	// SELECT * FROM users WHERE id = 10;

	var user3 []User9
	db.Find(&user3, []int{1, 2, 3})
	fmt.Println(user3)
	// SELECT * FROM users WHERE id IN (1,2,3);

	//db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")
	// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";

	var user4 = User9{Model: gorm.Model{ID: 10}}
	db.First(&user4)
	fmt.Println(user4)
	// SELECT * FROM users WHERE id = 10;

	var user5 User9
	db.Model(User9{Model: gorm.Model{ID: 10}}).First(&user5)
	fmt.Println(user5)
	// SELECT * FROM users WHERE id = 10;

}

func retrieveAllObjects(db *gorm.DB) {
	// Get all records
	var users []User9
	result := db.Find(&users)
	// SELECT * FROM users;
	fmt.Println(result.RowsAffected)
}

func stringConditions(db *gorm.DB) {
	// Get first matched record
	var user User9
	db.Where("name = ?", "jinzhu").First(&user)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

	// Get all matched records
	var users1 []User9
	db.Where("name <> ?", "jinzhu").Find(&users1)
	// SELECT * FROM users WHERE name <> 'jinzhu';

	// IN
	db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users1)
	// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

	// LIKE
	db.Where("name LIKE ?", "%jin%").Find(&users1)
	// SELECT * FROM users WHERE name LIKE '%jin%';

	// AND
	db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users1)
	// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

	now := time.Now()
	offset := int(now.Weekday())
	startOfLastWeek := now.AddDate(0, 0, -offset).AddDate(0, 0, -7)
	startOfLastWeek = startOfLastWeek.Truncate(24 * time.Hour)
	fmt.Println(startOfLastWeek)

	// Time
	db.Where("updated_at > ?", startOfLastWeek).Find(&users1)
	// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';

	// BETWEEN
	db.Where("created_at BETWEEN ? AND ?", startOfLastWeek, time.Now()).Find(&users1)
	// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';

	user2 := User9{Model: gorm.Model{ID: 10}}
	db.Where("id = ?", 20).First(&user2)

}

func structAndMapConditions(db *gorm.DB) {
	var user User9
	// Struct
	db.Where(&User9{Name: "jinzhu", Age: 20}).First(&user)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

	// Map
	var users []User9
	db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

	// Slice of primary keys
	var users2 []User9
	db.Where([]int64{20, 21, 22}).Find(&users2)
	// SELECT * FROM users WHERE id IN (20, 21, 22);

	//注意：使用结构体查询时，gorm只会将非0字段构建where查询语句；如果要将零值字段构建where查询语句，使用map
	db.Where(&User9{Name: "jinzhu", Age: 0}).Find(&users)
	// SELECT * FROM users WHERE name = "jinzhu";
	db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;

}

func specStructSearchFields(db *gorm.DB) {
	var users []User9
	db.Where(&User9{Name: "jinzhu"}, "name", "Age").Find(&users)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;

	db.Where(&User9{Name: "jinzhu"}, "Age").Find(&users)
	// SELECT * FROM users WHERE age = 0;

}

func inlineCondition(db *gorm.DB) {
	var user User9
	var users []User9
	// Get by primary key if it were a non-integer type
	db.First(&user, "id = ?", "string_primary_key")
	// SELECT * FROM users WHERE id = 'string_primary_key';

	// Plain SQL
	db.Find(&user, "name = ?", "jinzhu")
	// SELECT * FROM users WHERE name = "jinzhu";

	db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)
	// SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;

	// Struct
	db.Find(&users, User9{Age: 20})
	// SELECT * FROM users WHERE age = 20;

	// Map
	db.Find(&users, map[string]interface{}{"age": 20})
	// SELECT * FROM users WHERE age = 20;
}

func notConditions(db *gorm.DB) {
	var user User9
	var users []User9
	db.Not("name = ?", "jinzhu").First(&user)
	// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;

	// Not In
	db.Not(map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}).Find(&users)
	// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");

	// Struct
	db.Not(User9{Name: "jinzhu", Age: 18}).First(&user)
	// SELECT * FROM users WHERE name <> "jinzhu" AND age <> 18 ORDER BY id LIMIT 1;

	// Not In slice of primary keys
	db.Not([]int64{1, 2, 3}).First(&user)
	// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;

}

func orConditions(db *gorm.DB) {
	var users []User9
	db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
	// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

	// Struct
	db.Where("name = 'jinzhu'").Or(User9{Name: "jinzhu 2", Age: 18}).Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);

	// Map
	db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2", "age": 18}).Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);
}

func selectSpecificFields(db *gorm.DB) {
	var users []User9
	db.Select("name", "age").Find(&users)
	// SELECT name, age FROM users;

	db.Select([]string{"name", "age"}).Find(&users)
	// SELECT name, age FROM users;

	db.Table("user9").Select("COALESCE(age,?)", 42).Rows()
	// SELECT COALESCE(age,'42') FROM users;
}

func order(db *gorm.DB) {
	var users []User9
	db.Order("age desc, name").Find(&users)
	// SELECT * FROM users ORDER BY age desc, name;

	// Multiple orders
	db.Order("age desc").Order("name").Find(&users)
	// SELECT * FROM users ORDER BY age desc, name;

	db.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{[]int{1, 2, 3}}, WithoutParentheses: true},
	}).Find(&users)
	// SELECT * FROM users ORDER BY FIELD(id,1,2,3)
}

func limitAndOffset(db *gorm.DB) {
	var users []User9
	db.Limit(3).Find(&users)
	// SELECT * FROM users LIMIT 3;

	// Cancel limit condition with -1
	var users1 []User9
	var users2 []User9
	db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
	// SELECT * FROM users LIMIT 10; (users1)
	// SELECT * FROM users; (users2)

	// sql error : 因为 MySQL 的 OFFSET 关键字必须与 LIMIT 一起使用。OFFSET 定义从结果集中的哪一行开始返回记录，而 LIMIT 定义要返回多少条记录。
	db.Offset(3).Find(&users)
	// SELECT * FROM users OFFSET 3;

	db.Limit(10).Offset(5).Find(&users)
	// SELECT * FROM users OFFSET 5 LIMIT 10;

	// Cancel offset condition with -1
	db.Offset(10).Find(&users1).Offset(-1).Find(&users2)
	// SELECT * FROM users OFFSET 10; (users1)
	// SELECT * FROM users; (users2)
}

func groupByAndHaving(db *gorm.DB) {
	var result Result
	db.Model(&User9{}).Select("name, sum(age) as total").Where("name LIKE ?", "group%").Group("name").First(&result)
	// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name` LIMIT 1

	db.Model(&User9{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "group").Find(&result)
	// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"

	rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {

	}

	rows2, err2 := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
	if err2 != nil {
		return
	}
	defer rows2.Close()
	for rows.Next() {

	}
	type Result2 struct {
		Date  time.Time
		Total int
	}
	var results []Result2
	db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
}

func distinct(db *gorm.DB) {
	var users []User9
	db.Distinct("name", "age").Order("name, age desc").Find(&users)
}

func joins(db *gorm.DB) {
	type Result2 struct {
		Name  string
		Email string
	}
	var result Result2
	db.Model(&User9{}).Select("user9.name, emails.email").Joins("left join emails on emails.user_id = user9.id").Scan(&result)
	// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

	rows, err := db.Table("user9").Select("user9.name, emails.email").Joins("left join emails on emails.user_id = user9.id").Rows()
	if err == nil {
		for rows.Next() {

		}
	}
	var results []Result2
	db.Table("user9").Select("user9.name, emails.email").Joins("left join emails on emails.user_id = user9.id").Scan(&results)

	// multiple joins with parameter
	var user User9
	db.Joins("JOIN emails ON emails.user_id = user9.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = user9.id").Where("credit_cards.number = ?", "411111111111").Find(&user)

}

func joinPreloading(db *gorm.DB) {
	var users []User9
	db.Joins("Company").Find(&users)
	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;

	// inner join
	db.InnerJoins("Company").Find(&users)
	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` INNER JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;

	db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users)
	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id` AND `Company`.`alive` = true;

}

func joinADerivedTable(db *gorm.DB) {
	var results []Order
	query := db.Table("orders").Select("MAX(order.finished_at) as latest").Joins("left join user user on orders.user_id = user.id").Where("user.age > ?", 18).Group("orders.user_id")
	db.Model(&Order{}).Joins("join (?) q on order.finished_at = q.latest", query).Scan(&results)
	// SELECT `order`.`user_id`,`order`.`finished_at` FROM `order` join (SELECT MAX(order.finished_at) as latest FROM `order` left join user user on order.user_id = user.id WHERE user.age > 18 GROUP BY `order`.`user_id`) q on order.finished_at = q.latest

}

func scan(db *gorm.DB) {
	var result Result
	db.Table("users").Select("name", "age").Where("name = ?", "Antonio").Scan(&result)
	// Raw SQL
	db.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&result)
}
