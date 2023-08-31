package user

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID   uint64 `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string `gorm:"column:name"`
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Repository, error) {
	if !db.Migrator().HasTable(&User{}) {
		err := db.Migrator().AutoMigrate(&User{})
		if err != nil {
			fmt.Println(err)
		}
	}
	users := []User{
		{
			ID:   1,
			Name: "Bob",
		},
		{
			ID:   2,
			Name: "Kane",
		},
		{
			ID:   3,
			Name: "John",
		},
		{
			ID:   4,
			Name: "Matvey",
		},
	}
	err := db.Create(&users).Error
	fmt.Println(err)
	return &Repository{
		db: db,
	}, nil
}
