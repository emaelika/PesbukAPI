package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	Add() echo.HandlerFunc
	Login() echo.HandlerFunc
	Profile() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetUserByIDParam() echo.HandlerFunc
}

type UserService interface {
	Register(newData User) error
	Login(loginData User) (User, string, error)
	Profile(token *jwt.Token) (User, error)
	Update(token *jwt.Token, idFromParam uint, newData User) (User, error)
	Delete(token *jwt.Token, id uint) error
	GetUserByIDParam(token *jwt.Token, idFromParam uint) (User, error)
}

type UserModel interface {
	AddUser(newData User) error
	Login(email string) (User, error)
	GetUserByID(id uint) (User, error)
	Update(id uint, newData User) (User, error)
	Delete(id uint) error
	GetUserByIDE(id uint) (User, error)
}

type User struct {
	ID	         uint
	Name		 string `form:"name"`
	Email 	     string `form:"email"`
	Username 	 string `form:"username"`
	Placeofbirth string `form:"placeofbirth"`
	Dateofbirth  string `form:"dateofbirth"`
	Password 	 string `form:"password"`
	Image		 []byte `gorm:"type:longblob"`
}

type Login struct {
	Email 		string `validate:"required"`
	Password 	string `validate:"required,alphanum,min=8"`
}

type Register struct {
	Name 		 string `form:"name"` 
	Email		 string `form:"email" validate:"required"`
	Username	 string `form:"username"`
	Placeofbirth string `form:"placeofbirth"`
	Dateofbirth  string `form:"dateofbirth"`
	Password  	 string `form:"password" validate:"required,alphanum,min=8"`
}