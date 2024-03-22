package data

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID  uint
	PostID  uint
	Content string
}
type User struct {
	ID       uint   `gorm:"primary_key;auto_increment"`
	Fullname string `validate:"required"`
	Email    string `gorm:"unique"`
	Password string
	Birthday string
	Avatar   string
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}
type Post struct {
	gorm.Model
	UserID  uint
	Picture string
	Content string
	Comment []Comment `gorm:"foreignKey:PostID"`
}
