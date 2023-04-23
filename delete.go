package main

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Email struct {
	ID    uint
	Email string
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role == "admin" {
		return errors.New("admin user not allowed to delete")
	}
	return
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

func delete(db *gorm.DB) {
	email := Email{ID: 1}
	// Email's ID is `10`
	db.Delete(&email)
	// DELETE from emails where id = 10;

	// Delete with additional conditions
	db.Where("name = ?", "jinzhu").Delete(&email)
	// DELETE from emails where id = 10 AND name = "jinzhu";
}

func deleteWithPrimaryKey(db *gorm.DB) {
	db.Delete(&User{}, 10)
	// DELETE FROM users WHERE id = 10;

	db.Delete(&User{}, "10")
	// DELETE FROM users WHERE id = 10;

	db.Delete(&User{}, []int{1, 2, 3})
	// DELETE FROM users WHERE id IN (1,2,3);

}

func batchDelete(db *gorm.DB) {
	db.Where("email LIKE ?", "%jinzhu%").Delete(&Email{})
	// DELETE from emails where email LIKE "%jinzhu%";

	db.Delete(&Email{}, "email LIKE ?", "%jinzhu%")
	// DELETE from emails where email LIKE "%jinzhu%";

	//To efficiently delete large number of records, pass a slice with primary keys to the Delete method.
	var users = []User{{ID: 1}, {ID: 2}, {ID: 3}}
	db.Delete(&users)
	// DELETE FROM users WHERE id IN (1,2,3);

	db.Delete(&users, "name LIKE ?", "%jinzhu%")
	// DELETE FROM users WHERE name LIKE "%jinzhu%" AND id IN (1,2,3);

}

func globalDelete(db *gorm.DB) {
	//db.Delete(&User{}).Error // gorm.ErrMissingWhereClause
	//db.Delete(&[]User{{Name: "jinzhu1"}, {Name: "jinzhu2"}}).Error // gorm.ErrMissingWhereClause

	db.Where("1 = 1").Delete(&User{})
	// DELETE FROM `users` WHERE 1=1

	db.Exec("DELETE FROM users")
	// DELETE FROM users

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
	// DELETE FROM users
}

// If your model includes a gorm.DeletedAt field (which is included in gorm.Model), it will get soft delete ability automatically!
// If you don’t want to include gorm.Model, you can enable the soft delete feature like:
//
//	type User struct {
//		ID      int
//		Deleted gorm.DeletedAt
//		Name    string
//	}
func softDelete(db *gorm.DB) {
	user := User{ID: 111}
	// user's ID is `111`
	db.Delete(&user)
	// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

	// Batch Delete
	db.Where("age = ?", 20).Delete(&User{})
	// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

	// Soft deleted records will be ignored when querying
	db.Where("age = 20").Find(&user)
	// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
}

func findSoftDeletedRecord(db *gorm.DB) {
	var users []User
	db.Unscoped().Where("age = 20").Find(&users)
	// SELECT * FROM users WHERE age = 20;

}

func deletePermanently(db *gorm.DB) {
	user := User{ID: 111}
	db.Unscoped().Delete(&user)
	// DELETE FROM users WHERE id=111;

}

// // Query
// SELECT * FROM users WHERE deleted_at = 0;
// // Delete
// UPDATE users SET deleted_at = /* current unix second */ WHERE ID = 1;
type Customer struct {
	ID        uint
	Name      string                `gorm:"uniqueIndex:udx_name"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
}

// // Query
// SELECT * FROM users WHERE is_del = 0;
// // Delete
// UPDATE users SET is_del = 1 WHERE ID = 1;
type Customer2 struct {
	ID    uint
	Name  string
	IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

// // Query
// SELECT * FROM users WHERE is_del = 0;
// // Delete
// UPDATE users SET is_del = 1, deleted_at = /* current unix second */ WHERE ID = 1;
type Customer3 struct {
	ID        uint
	Name      string
	DeletedAt time.Time
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	// IsDel     soft_delete.DeletedAt `gorm:"softDelete:,DeletedAtField:DeletedAt"` // use `unix second`
	// IsDel     soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"` // use `unix nano second`
}
