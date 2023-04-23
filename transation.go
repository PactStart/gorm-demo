package main

import (
	"errors"
	"gorm.io/gorm"
)

func makeATrans(db *gorm.DB) {
	// To perform a set of operations within a transaction, the general flow is as below.
	db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
			// return any error will rollback
			return err
		}

		if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})
}

func nestedTrans(db *gorm.DB) {
	// you can rollback a subset of operations performed within the scope of a larger transaction
	db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&user1)

		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&user2)
			return errors.New("rollback user2") // Rollback user2
		})

		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&user3)
			return nil
		})

		return nil
	})

	// Commit user1, user3
}

func manualControlTrans(db *gorm.DB) {
	// begin a transaction
	tx := db.Begin()

	// do some database operations in the transaction (use 'tx' from this point, not 'db')
	tx.Create(...)

	// ...

	// rollback the transaction in case of error
	tx.Rollback()

	// Or commit the transaction
	tx.Commit()

}

func CreateAnimals(db *gorm.DB) error {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func savePoint(db *gorm.DB)  {
	tx := db.Begin()
	tx.Create(&user1)

	tx.SavePoint("sp1")
	tx.Create(&user2)
	tx.RollbackTo("sp1") // Rollback user2

	tx.Commit() // Commit user1
}