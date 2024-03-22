package data

import (
	comment "PesbukAPI/features/comment/data"
	post "PesbukAPI/features/post/data"
)

type User struct {
	ID 			 uint 					`gorm:"primary_key;auto_increment"`
	Fullname     string                 `validate:"required"`
	Email        string                 `gorm:"unique"`
	Password	 string
	Birthday  	 string
	Avatar		 string
	Posts	 	 []post.Post			`gorm:"foreignKey:UserID"`
	Comments	 []comment.Comment		`gorm:"foreignKey:UserID"`
}
