package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID       uint
	quantity uint
	price    float64
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Role == "admin" {
		return errors.New("admin user not allowed to update")
	}
	// if Role changed
	if tx.Statement.Changed("Role") {
		return errors.New("role not allowed to change")
	}

	if tx.Statement.Changed("Name", "Admin") { // if Name or Role changed
		tx.Statement.SetColumn("Age", 18)
	}

	// if any fields changed
	if tx.Statement.Changed() {
		tx.Statement.SetColumn("RefreshedAt", time.Now())
	}
	return nil
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}
	saveAllFields(db)
}

func saveAllFields(db *gorm.DB) {
	var user User
	db.First(&user)

	user.Name = "jinzhu 2"
	user.Age = 100
	db.Save(&user)
	// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;

	db.Save(&User{Name: "jinzhu", Age: 100})
	// INSERT INTO `users` (`name`,`age`,`birthday`,`update_at`) VALUES ("jinzhu",100,"0000-00-00 00:00:00","0000-00-00 00:00:00")

	db.Save(&User{ID: 1, Name: "jinzhu", Age: 100})
	// UPDATE `users` SET `name`="jinzhu",`age`=100,`birthday`="0000-00-00 00:00:00",`update_at`="0000-00-00 00:00:00" WHERE `id` = 1
}

func updateSingleColumn(db *gorm.DB) {
	// Update with conditions
	db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

	// User's ID is `111`:
	var user = User{ID: 111}
	db.Model(&user).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

	// Update with conditions and model value
	db.Model(&user).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

}

func updateMultipleColumn(db *gorm.DB) {
	var user = User{ID: 111}
	// Update attributes with `struct`, will only update non-zero fields
	db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
	// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

	// Update attributes with `map`
	db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
}

func updateSelectedFields(db *gorm.DB) {
	var user User
	// Select with Map
	// User's ID is `111`:
	db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET name='hello' WHERE id=111;

	db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

	// Select with Struct (select zero value fields)
	db.Model(&user).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})
	// UPDATE users SET name='new_name', age=0 WHERE id=111;

	// Select all fields (select all fields include zero value fields)
	db.Model(&user).Select("*").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})

	// Select all fields but omit Role (select all fields include zero value fields)
	db.Model(&user).Select("*").Omit("Role").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})
}

func batchUpdate(db *gorm.DB) {
	// Update with struct
	db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
	// UPDATE users SET name='hello', age=18 WHERE role = 'admin';

	// Update with map
	db.Table("users").Where("id IN ?", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})
	// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
}

func globalUpdate(db *gorm.DB) {
	//db.Model(&User{}).Update("name", "jinzhu").Error // gorm.ErrMissingWhereClause

	db.Model(&User{}).Where("1 = 1").Update("name", "jinzhu")
	// UPDATE users SET `name` = "jinzhu" WHERE 1=1

	db.Exec("UPDATE users SET name = ?", "jinzhu")
	// UPDATE users SET name = "jinzhu"

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("name", "jinzhu")
	// UPDATE users SET `name` = "jinzhu"

}

func updateRecordsCount(db *gorm.DB) {
	// Get updated records count with `RowsAffected`
	result := db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
	// UPDATE users SET name='hello', age=18 WHERE role = 'admin';
	fmt.Println(result.RowsAffected)
}

func updateWithSQLExpr(db *gorm.DB) {
	var product = Product{ID: 3}
	// product's ID is `3`
	db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
	// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;

	db.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})
	// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;

	db.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
	// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3;

	db.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
	// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3 AND quantity > 1;
}

func updateWithSQLExprAndCustomizedDataType(db *gorm.DB) {
	db.Model(&User{ID: 1}).Updates(User{
		Name:     "jinzhu",
		Location: Location{X: 100, Y: 100},
	})
	// UPDATE `user_with_points` SET `name`="jinzhu",`location`=ST_PointFromText("POINT(100 100)") WHERE `id` = 1
}

func updateFromSubQuery(db *gorm.DB) {
	var user User
	db.Model(&user).Update("company_name", db.Model(&Company{}).Select("name").Where("companies.id = users.company_id"))
	// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

	db.Table("users as u").Where("name = ?", "jinzhu").Update("company_name", db.Table("companies as c").Select("name").Where("c.id = u.company_id"))

	db.Table("users as u").Where("name = ?", "jinzhu").Updates(map[string]interface{}{"company_name": db.Table("companies as c").Select("name").Where("c.id = u.company_id")})

}

func updateWithoutHooksAndTimeTracing(db *gorm.DB) {
	var user User
	// Update single column
	db.Model(&user).UpdateColumn("name", "hello")
	// UPDATE users SET name='hello' WHERE id = 111;

	// Update multiple columns
	db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
	// UPDATE users SET name='hello', age=18 WHERE id = 111;

	// Update selected columns
	db.Model(&user).Select("name", "age").UpdateColumns(User{Name: "hello", Age: 0})
	// UPDATE users SET name='hello', age=0 WHERE id = 111;
}

func checkFieldHasChanged(db *gorm.DB) {
	db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(map[string]interface{}{"name": "jinzhu2"})
	// Changed("Name") => true
	db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(map[string]interface{}{"name": "jinzhu"})
	// Changed("Name") => false, `Name` not changed
	db.Model(&User{ID: 1, Name: "jinzhu"}).Select("Admin").Updates(map[string]interface{}{
		"name": "jinzhu2", "admin": false,
	})
	// Changed("Name") => false, `Name` not selected to update

	db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(User{Name: "jinzhu2"})
	// Changed("Name") => true
	db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(User{Name: "jinzhu"})
	// Changed("Name") => false, `Name` not changed
	db.Model(&User{ID: 1, Name: "jinzhu"}).Select("Admin").Updates(User{Name: "jinzhu2"})
	// Changed("Name") => false, `Name` not selected to update
}
