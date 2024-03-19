package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Profile() echo.HandlerFunc
}

type UserService interface {
	AddUser(newData User) error
	Login(hp string, password string) (User, string, error)
	Profile(token *jwt.Token) (User, error)
}

type UserQuery interface {
	AddUser(newData User) error
	Login(hp string) (User, error)
	Profile(id uint) (User, error)
}

type User struct {
	ID       uint
	Nama     string `validate:"required"`
	Hp       string `validate:"required,max=13,min=10"`
	Password string `validate:"required"`
}
