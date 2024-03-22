package services

import (
	"PesbukAPI/features/user"
	"PesbukAPI/helper"
	"PesbukAPI/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	model 	user.UserModel
	pm 		helper.PasswordManager
	v 		*validator.Validate
}

func NewService(m user.UserModel) user.UserService {
	return &service{
		model: 	m,
		pm:		helper.NewPasswordManager(),
		v:		validator.New(),
	}
}

func (s *service) Register(newData user.User) error {
	var registerValidate user.Register
    registerValidate.Fullname = newData.Fullname
    registerValidate.Email = newData.Email
	registerValidate.Password = newData.Password
    registerValidate.Birthday = newData.Birthday
	err := s.v.Struct(&registerValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	newPassword, err := s.pm.HashPassword(newData.Password)
	if err != nil {
		return errors.New(helper.ServiceGeneralError)
	}
	newData.Password = newPassword

	err = s.model.AddUser(newData)
	if err != nil {
		return errors.New(helper.ServerGeneralError)
	}
	return nil
}
func (s *service) Login(loginData user.User) (user.User, string, error) {
	var loginValidate user.Login
	loginValidate.Email = loginData.Email
	loginValidate.Password = loginData.Password
	err := s.v.Struct(&loginValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return user.User{}, "", err
	}

	dbData, err := s.model.Login(loginValidate.Email)
	if err != nil {
		return user.User{}, "",err
	}

	err = s.pm.ComparePassword(loginValidate.Password, dbData.Password)
	if err != nil {
		return user.User{}, "", errors.New(helper.UserCredentialError)
	}

	token, err := middlewares.GenerateJWT(dbData.ID)
	if err != nil {
		return user.User{}, "", errors.New(helper.ServiceGeneralError)
	}

	return dbData, token, nil
}

func (s *service) Profile(token *jwt.Token) (user.User, error) {
	decodeId := middlewares.DecodeToken(token)
	result, err := s.model.GetUserByID(decodeId)
	if err != nil {
		return user.User{}, err
	}

	return result, nil
}

func (s *service) Update(token *jwt.Token, newData user.User) (user.User, error) {
    decodedID := middlewares.DecodeToken(token)

    existingUser, err := s.model.GetUserByID(decodedID)
    if err != nil {
        return user.User{}, errors.New("user not found")
    }

    if newData.Fullname != "" {
        existingUser.Fullname = newData.Fullname
    }

    if newData.Email != "" {
        existingUser.Email = newData.Email
    }

    if newData.Password != "" {
        newPassword, err := s.pm.HashPassword(newData.Password)
        if err != nil {
            return user.User{}, errors.New(helper.ServiceGeneralError)
        }
        existingUser.Password = newPassword
    }

    if newData.Birthday != "" {
        existingUser.Birthday = newData.Birthday
    }

    if newData.Avatar != "" {
        existingUser.Avatar = newData.Avatar
    }

    result, err := s.model.Update(decodedID, existingUser)
    if err != nil {
        return user.User{}, err
    }

    return result, nil
}


func (s *service) Delete(token *jwt.Token) error {
    decodedID := middlewares.DecodeToken(token)
    if decodedID == 0 {
        log.Println("error decode token:", "token tidak ditemukan")
        return errors.New("data tidak valid")
    }

    err := s.model.Delete(decodedID)
    if err != nil {
        return errors.New(helper.CannotDelete)
    }

    return nil
}

func (s *service) GetUserByIDParam(token *jwt.Token, idFromParam uint) (user.User, error) {
    decodedID := middlewares.DecodeToken(token)

    if idFromParam != decodedID {
        return user.User{}, errors.New("you can only view your own account")
    }

    result, err := s.model.GetUserByID(decodedID)
    if err != nil {
        return user.User{}, err
    }

    return result, nil
}

