package repository

import (
	"21-api/features/todo/repository"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Nama     string                 `validate:"required"`
	Hp       string                 `gorm:"unique"`
	Password string                 `validate:"required"`
	Todos    []repository.TodoModel `gorm:"foreignKey:UserID"`
}
