package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func main() {
	db, err := gorm.Open(mysql.Open("db1_dsn"), &gorm.Config{})

	db.Use(dbresolver.Register(dbresolver.Config{
		// use `db2` as sources, `db3`, `db4` as replicas
		Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
		Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}).Register(dbresolver.Config{
		// use `db1` as sources (DB's default connection), `db5` as replicas for `User`, `Address`
		Replicas: []gorm.Dialector{mysql.Open("db5_dsn")},
	}, &User{}, &Address{}).Register(dbresolver.Config{
		// use `db6`, `db7` as sources, `db8` as replicas for `orders`, `Product`
		Sources:  []gorm.Dialector{mysql.Open("db6_dsn"), mysql.Open("db7_dsn")},
		Replicas: []gorm.Dialector{mysql.Open("db8_dsn")},
	}, "orders", &Product{}, "secondary"))

	// `User` Resolver Examples
	db.Table("users").Rows()                           // replicas `db5`
	db.Model(&User{}).Find(&AdvancedUser{})            // replicas `db5`
	db.Exec("update users set name = ?", "jinzhu")     // sources `db1`
	db.Raw("select name from users").Row().Scan(&name) // replicas `db5`
	db.Create(&user)                                   // sources `db1`
	db.Delete(&User{}, "name = ?", "jinzhu")           // sources `db1`
	db.Table("users").Update("name", "jinzhu")         // sources `db1`

	// Global Resolver Examples
	db.Find(&Pet{}) // replicas `db3`/`db4`
	db.Save(&Pet{}) // sources `db2`

	// Orders Resolver Examples
	db.Find(&Order{})                  // replicas `db8`
	db.Table("orders").Find(&Report{}) // replicas `db8`

}
