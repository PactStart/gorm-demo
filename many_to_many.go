package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// User has and belongs to many languages, `user_languages` is the join table

type User18 struct {
	gorm.Model
	Languages []Language2 `gorm:"many2many:user_languages;"`
}

type Language2 struct {
	gorm.Model
	Name string
}

// Back-Reference
// User has and belongs to many languages, use `user_languages` as join table
type User19 struct {
	gorm.Model
	Languages []*Language3 `gorm:"many2many:user_languages;"`
}

type Language3 struct {
	gorm.Model
	Name  string
	Users []*User19 `gorm:"many2many:user_languages;"`
}

// Retrieve user list with eager loading languages
func GetAllUsers(db *gorm.DB) ([]User19, error) {
	var users []User19
	err := db.Model(&User19{}).Preload("Languages").Find(&users).Error
	return users, err
}

// Retrieve language list with eager loading users
func GetAllLanguages(db *gorm.DB) ([]Language3, error) {
	var languages []Language3
	err := db.Model(&User19{}).Preload("Users").Find(&languages).Error
	return languages, err
}

// Override Foreign Key
// Which creates join table: user_profiles
//
//	foreign key: user_refer_id, reference: users.refer
//	foreign key: profile_refer, reference: profiles.user_refer
type User20 struct {
	gorm.Model
	Profiles []Profile `gorm:"many2many:user_profiles;foreignKey:Refer;joinForeignKey:UserReferID;References:UserRefer;joinReferences:ProfileRefer"`
	Refer    uint      `gorm:"index:,unique"`
}

type Profile struct {
	gorm.Model
	Name      string
	UserRefer uint `gorm:"index:,unique"`
}

// Self-Referential Many2Many
// Which creates join table: user_friends
//   foreign key: user_id, reference: users.id
//   foreign key: friend_id, reference: users.id

type User21 struct {
	gorm.Model
	Friends []*User21 `gorm:"many2many:user_friends"`
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.Debug()
	// Change model Person's field Addresses' join table to PersonAddress
	// PersonAddress must defined all required foreign keys or it will raise error
	db.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{})

}

// Customize JoinTable
type Person struct {
	ID        int
	Name      string
	Addresses []Address `gorm:"many2many:person_addressses;"`
}

type Address struct {
	ID   uint
	Name string
}

type PersonAddress struct {
	PersonID  int `gorm:"primaryKey"`
	AddressID int `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (PersonAddress) BeforeCreate(db *gorm.DB) error {
	// ...
	return nil
}

// Composite Foreign Keys

type Tag struct {
	ID     uint   `gorm:"primaryKey"`
	Locale string `gorm:"primaryKey"`
	Value  string
}

type Blog struct {
	ID         uint   `gorm:"primaryKey"`
	Locale     string `gorm:"primaryKey"`
	Subject    string
	Body       string
	Tags       []Tag `gorm:"many2many:blog_tags;"`
	LocaleTags []Tag `gorm:"many2many:locale_blog_tags;ForeignKey:id,locale;References:id"`
	SharedTags []Tag `gorm:"many2many:shared_blog_tags;ForeignKey:id;References:id"`
}

// Join Table: blog_tags
//   foreign key: blog_id, reference: blogs.id
//   foreign key: blog_locale, reference: blogs.locale
//   foreign key: tag_id, reference: tags.id
//   foreign key: tag_locale, reference: tags.locale

// Join Table: locale_blog_tags
//   foreign key: blog_id, reference: blogs.id
//   foreign key: blog_locale, reference: blogs.locale
//   foreign key: tag_id, reference: tags.id

// Join Table: shared_blog_tags
//   foreign key: blog_id, reference: blogs.id
//   foreign key: tag_id, reference: tags.id
