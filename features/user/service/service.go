package service

import (
	"21-api/features/user"
	"21-api/helper"
	"21-api/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	uq user.UserQuery
	v  *validator.Validate
	pm helper.PasswordManager
}

func NewUserService(query user.UserQuery) user.UserService {
	return &service{
		uq: query,
		v:  validator.New(validator.WithRequiredStructEnabled()),
		pm: helper.NewPasswordManager(),
	}
}

func (us service) AddUser(data user.User) error {
	err := us.v.Struct(&data)
	if err != nil {
		log.Println("error validasi", err.Error())
		return errors.New("data tidak valid")
	}

	newPassword, err := us.pm.HashPassword(data.Password)
	if err != nil {
		return errors.New("password error")
	}
	data.Password = newPassword

	err = us.uq.AddUser(data)
	if err != nil {
		log.Println("error query", err.Error())
		return err
	}
	return nil
}

func (us service) Login(email string, password string) (user.User, string, error) {
	// validasi
	var dummyValidate = user.User{Name: "A",
		Email:       email,
		Password: password}
	err := us.v.Struct(&dummyValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return user.User{}, "", errors.New("data tidak valid")
	}

	// ngambil data
	data, err := us.uq.Login(email)
	if err != nil {
		log.Println(err.Error())
		return user.User{}, "", err
	}

	// compare hash password
	err = us.pm.ComparePassword(password, data.Password)
	if err != nil {
		log.Println("user service,", err.Error())
		return user.User{}, "", err
	}
	// token
	token, err := middlewares.GenerateJWT(data.ID)
	if err != nil {
		log.Println(err.Error())
		return user.User{}, "", err
	}

	return data, token, nil
}

func (us service) Profile(token *jwt.Token) (user.User, error) {
	id, err := middlewares.ExtractId(token)
	if err != nil {
		log.Println(err.Error())
		return user.User{}, err
	}

	data, err := us.uq.Profile(id)
	if err != nil {
		log.Println(err.Error())
		return user.User{}, err
	}
	return data, nil
}
