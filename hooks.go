package main

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Creating an object

//// begin transaction
//BeforeSave
//BeforeCreate
//// save before associations
//// insert into database
//// save after associations
//AfterCreate
//AfterSave
//// commit or rollback transaction

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()

	if !u.IsValid() {
		err = errors.New("can't save invalid data")
	}
	// Modify current operation through tx.Statement, e.g:
	//tx.Statement.Select("Name", "Age")
	//tx.Statement.AddClause(clause.OnConflict{DoNothing: true})
	// operations based on it will run inside same transaction but without any current conditions
	//var role Role
	//err := tx.First(&role, "name = ?", user.Role).Error
	//// SELECT * FROM roles WHERE name = "admin"
	// ...
	return err
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("role", "admin")
	}
	// if you return any error in your hooks, the change will be rollbacked
	if !u.IsValid() {
		return errors.New("rollback invalid user")
	}
	return
}

//Updating an object
//// begin transaction
//BeforeSave
//BeforeUpdate
//// save before associations
//// update database
//// save after associations
//AfterUpdate
//AfterSave
//// commit or rollback transaction

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.readonly() {
		err = errors.New("read only user")
	}
	return
}

// Updating data in same transaction
func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	if u.Confirmed {
		tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("verfied", true)
	}
	return
}

//Deleting an object

// // begin transaction
// BeforeDelete
// // delete from database
// AfterDelete
// // commit or rollback transaction
// Updating data in same transaction
func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	if u.Confirmed {
		tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	}
	return
}

// load data from database
// Preloading (eager loading)
// AfterFind
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	if u.MemberShip == "" {
		u.MemberShip = "user"
	}
	return
}
