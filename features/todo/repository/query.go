package repository

import (
	"PesbukAPI/features/todo"
	"log"

	"gorm.io/gorm"
)

type query struct {
	connection *gorm.DB
}

func NewTodoQuery(db *gorm.DB) todo.TodoQuery {
	return &query{
		connection: db,
	}
}

func (tq query) AddTodo(newTodo todo.Todo) (todo.Todo, error) {
	var newData = TodoModel{
		Kegiatan:  newTodo.Kegiatan,
		Deskripsi: newTodo.Deskripsi,
		Deadline:  newTodo.Deadline,
		UserID:    newTodo.UserID,
	}
	err := tq.connection.Create(&newData).Error
	if err != nil {
		log.Println(err.Error())
		return todo.Todo{}, err
	}
	var result = todo.Todo{
		ID:        newData.ID,
		UserID:    newData.UserID,
		Kegiatan:  newData.Kegiatan,
		Deskripsi: newData.Deskripsi,
		Deadline:  newData.Deadline,
	}

	return result, nil
}

func (tq query) GetTodos(id uint) ([]todo.Todo, error) {
	var data []TodoModel

	if err := tq.connection.Where("user_id = ?", id).Find(&data).Error; err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var results []todo.Todo
	for _, val := range data {
		var result = todo.Todo{
			ID:        val.ID,
			UserID:    val.UserID,
			Kegiatan:  val.Kegiatan,
			Deskripsi: val.Deskripsi,
			Deadline:  val.Deadline,
		}
		results = append(results, result)
	}

	return results, nil
}

func (tq query) GetTodo(id uint) (todo.Todo, error) {
	var data TodoModel

	if err := tq.connection.Where("id = ?", id).Find(&data).Error; err != nil {
		log.Println(err.Error())
		return todo.Todo{}, err
	}

	var result = todo.Todo{
		ID:        data.ID,
		UserID:    data.UserID,
		Kegiatan:  data.Kegiatan,
		Deskripsi: data.Deskripsi,
		Deadline:  data.Deadline,
	}

	return result, nil
}
func (tq query) UpdateTodo(update todo.Todo) (todo.Todo, error) {
	var data TodoModel
	log.Println(update)
	if err := tq.connection.First(&data, update.ID).
		Error; err != nil {
		log.Println(err.Error())
		return todo.Todo{}, err
	}
	data.Deskripsi = update.Deskripsi
	data.Deadline = update.Deadline
	data.Kegiatan = update.Kegiatan
	if err := tq.connection.Save(&data).
		Error; err != nil {
		log.Println(err.Error())
		return todo.Todo{}, err
	}

	log.Println(update)

	return update, nil
}
