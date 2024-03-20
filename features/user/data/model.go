package data

import (
	"PesbukAPI/features/todo/repository"
)

type User struct {
	ID 			 uint 					`gorm:"primary_key;auto_increment"`
	Name     	 string                 `validate:"required"`
	Email        string                 `gorm:"unique"`
	Username 	 string                 `validate:"required"`
	Placeofbirth string
	Dateofbirth  string
	Password	 string
	Image		 []byte 				`gorm:"type:longblob"`
	Posts    	 []repository.TodoModel `gorm:"foreignKey:UserID"`
}
