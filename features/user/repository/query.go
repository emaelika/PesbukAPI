package repository

import (
	"21-api/features/user"
	"log"

	"gorm.io/gorm"
)

type query struct {
	connection *gorm.DB
}

func NewUserQuery(db *gorm.DB) user.UserQuery {
	return &query{
		connection: db,
	}
}

func (uq query) AddUser(newData user.User) error {
	var input = UserModel{
		Nama:     newData.Nama,
		Hp:       newData.Hp,
		Password: newData.Password,
	}
	err := uq.connection.Create(&input).Error
	if err != nil {
		log.Println("error di mysql, ", err.Error())
		return err
	}
	return nil

}

func (uq query) Login(hp string) (user.User, error) {
	var data UserModel
	err := uq.connection.First(&data).Where("hp = ?", hp).Error
	if err != nil {
		log.Println(err.Error())
		return user.User{}, err
	}

	var result = user.User{
		ID:       data.ID,
		Hp:       data.Hp,
		Nama:     data.Nama,
		Password: data.Password,
	}
	return result, nil

}

func (uq query) Profile(id uint) (user.User, error) {
	var data UserModel
	err := uq.connection.First(&data).Where("ID = ?", id).Error
	if err != nil {
		log.Println(err.Error())
		return user.User{}, err
	}

	var result = user.User{
		ID:       data.ID,
		Hp:       data.Hp,
		Nama:     data.Nama,
		Password: data.Password,
	}
	return result, nil

}
