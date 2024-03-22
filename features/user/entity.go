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
	Avatar() echo.HandlerFunc 
}

type UserService interface {
	Register(newData User) error
	Login(loginData User) (User, string, error)
	Profile(token *jwt.Token) (User, error)
	Update(token *jwt.Token, newData User) (User, error)
	Delete(token *jwt.Token) error
	GetUserByIDParam(token *jwt.Token, idFromParam uint) (User, error)
}

type UserModel interface {
	AddUser(newData User) error
	Login(email string) (User, error)
	GetUserByID(id uint) (User, error)
	Update(id uint, newData User) (User, error)
	Delete(id uint) error
}

type User struct {
	ID	         uint	`json:"id"`
	Fullname	 string `json:"fullname" form:"fullname"`
	Email 	     string `json:"email" form:"email"`
	Password 	 string `json:"password" form:"password"`
	Birthday     string `json:"birthday" form:"birthday"`
	Avatar		 string `json:"avatar" form:"avatar"`
}

type Login struct {
	Email 		string `form:"email" validate:"required"`
	Password 	string `form:"password" validate:"required,alphanum,min=8"`
}

type Register struct {
	Fullname     string `form:"fullname"` 
	Email		 string `form:"email" validate:"required"`
	Password  	 string `form:"password" validate:"required,alphanum,min=8"`
	Birthday     string `form:"birthday"`
}