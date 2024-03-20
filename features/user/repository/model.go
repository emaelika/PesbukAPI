package repository

import (
	"21-api/features/todo/repository"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Name     	 string                 `validate:"required"`
	Email        string                 `gorm:"unique"`
	Username 	 string                 `validate:"required"`
	Placeofbirth string
	Dateofbirth  time.Time
	Password	 string
	Image		 []byte 				`gorm:"type:longblob"`
	Posts    	 []repository.TodoModel `gorm:"foreignKey:UserID"`
}
