package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode = ?", "card")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode = ?", "cod")
}

func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(AmountGreaterThan1000).Where("status IN (?)", status)
	}
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
	// Find all credit card orders and amount greater than 1000

	db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)
	// Find all COD orders and amount greater than 1000

	db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
	// Find all paid, shipped orders that amount greater than 1000

	db.Scopes(Paginate(1, 10)).Find(&users)
	db.Scopes(Paginate(1, 10)).Find(&articles)

	db.Scopes(UserTable(user)).Create(&user)
}

func UserTable(user User) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if user.Admin {
			return tx.Table("admin_users")
		}

		return tx.Table("users")
	}
}
