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
    registerValidate.Name = newData.Name
    registerValidate.Email = newData.Email
	registerValidate.Username = newData.Username
    registerValidate.Placeofbirth = newData.Placeofbirth
    registerValidate.Dateofbirth = newData.Dateofbirth
	registerValidate.Password = newData.Password
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

func (s *service) Update(token *jwt.Token, idFromParam uint, newData user.User) (user.User, error) {
    decodedID := middlewares.DecodeToken(token)

    if idFromParam != decodedID {
        return user.User{}, errors.New("you can only update your own account")
    }

    existingUser, err := s.model.GetUserByID(decodedID)
    if err != nil {
        return user.User{}, errors.New("user not found")
    }

    if newData.Name != "" {
        existingUser.Name = newData.Name
    }

    if newData.Email != "" {
        existingUser.Email = newData.Email
    }

    if newData.Username != "" {
        existingUser.Username = newData.Username
    }

    if newData.Placeofbirth != "" {
        existingUser.Placeofbirth = newData.Placeofbirth
    }

    if newData.Dateofbirth != "" {
        existingUser.Dateofbirth = newData.Dateofbirth
    }

    if newData.Password != "" {
        newPassword, err := s.pm.HashPassword(newData.Password)
        if err != nil {
            return user.User{}, errors.New(helper.ServiceGeneralError)
        }
        existingUser.Password = newPassword
    }

    if len(newData.Image) > 0 {
        existingUser.Image = newData.Image
    }

    result, err := s.model.Update(decodedID, existingUser)
    if err != nil {
        return user.User{}, err
    }

    return result, nil
}


func (s *service) Delete(token *jwt.Token, id uint) error {
    decodedID := middlewares.DecodeToken(token)
    if decodedID == 0 {
        log.Println("error decode token:", "token tidak ditemukan")
        return errors.New("data tidak valid")
    }

    userFromToken, err := s.model.GetUserByID(decodedID)
    if err != nil {
        return err
    }

    if userFromToken.ID != id {
        return errors.New("Anda hanya dapat menghapus akun Anda sendiri")
    }

    err = s.model.Delete(id)
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

    result, err := s.model.GetUserByIDE(decodedID)
    if err != nil {
        return user.User{}, err
    }

    return result, nil
}
