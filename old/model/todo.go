package model

import (
	"log"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Kegiatan  string `json:"nama" form:"nama" validate:"required"`
	Deskripsi string `json:"hp" form:"hp" validate:"required,max=13,min=10"`
	Deadline  string `json:"password" form:"password" validate:"required"`
	UserID    uint
}

type TodoModel struct {
	Connection *gorm.DB
}

func (um *TodoModel) AddTodo(newData Todo) (Todo, error) {
	err := um.Connection.Create(&newData).Error
	if err != nil {
		return Todo{}, err
	}

	return newData, nil
}

// func (um *TodoModel) CekTodo(hp string) bool {
// 	var data Todo
// 	if err := um.Connection.Where("hp = ?", hp).First(&data).Error; err != nil {
// 		return false
// 	}
// 	return true
// }

// func (um *TodoModel) Update(hp string, data Todo) error {
// 	if err := um.Connection.Model(&data).Where("hp = ?", hp).Update("nama", data.Nama).Update("password", data.Password).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (um *TodoModel) GetTodos(id uint) ([]Todo, error) {
	var result []Todo

	if err := um.Connection.Where("user_id = ?", id).Find(&result).Error; err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func (um *TodoModel) GetTodo(id uint) (Todo, error) {
	var result Todo

	if err := um.Connection.Where("id = ?", id).First(&result).Error; err != nil {
		log.Println(err.Error())
		return Todo{}, err
	}

	return result, nil
}
func (um *TodoModel) UpdateTodo(newTodo Todo) (Todo, error) {
	var data Todo
	log.Println(newTodo)
	if err := um.Connection.First(&data, newTodo.ID).
		Error; err != nil {
		log.Println(err.Error())
		return Todo{}, err
	}
	data.Deskripsi = newTodo.Deskripsi
	data.Deadline = newTodo.Deadline
	data.Kegiatan = newTodo.Kegiatan
	if err := um.Connection.Save(&data).
		Error; err != nil {
		log.Println(err.Error())
		return Todo{}, err
	}

	log.Println(newTodo)

	return newTodo, nil
}
