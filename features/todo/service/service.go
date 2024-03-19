package service

import (
	"21-api/features/todo"
	"21-api/features/todo/handler"
	"21-api/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	tq todo.TodoQuery
	v  *validator.Validate
}

func NewTodoService(query todo.TodoQuery) todo.TodoService {
	return &service{
		tq: query,
		v:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (ts service) AddTodo(pemilik *jwt.Token, kegiatanBaru todo.Todo) (todo.Todo, error) {
	userID, err := middlewares.ExtractId(pemilik)
	if err != nil {
		log.Println(err.Error())
		return todo.Todo{}, err
	}

	var cekValid = handler.TodoRequest{
		Kegiatan:  kegiatanBaru.Kegiatan,
		Deadline:  kegiatanBaru.Deadline,
		Deskripsi: kegiatanBaru.Deskripsi,
	}
	err = ts.v.Struct(&cekValid)
	if err != nil {
		log.Println("error validasi", err.Error())
		return todo.Todo{}, errors.New("data tidak valid")
	}

	kegiatanBaru.UserID = userID
	result, err := ts.tq.AddTodo(kegiatanBaru)
	if err != nil {
		log.Println("service error", err.Error())
		return todo.Todo{}, err
	}
	return result, nil

}
func (ts service) GetTodos(pemilik *jwt.Token) ([]todo.Todo, error) {
	id, err := middlewares.ExtractId(pemilik)
	if err != nil {
		log.Println("todo service,", err.Error())
		return nil, err
	}
	data, err := ts.tq.GetTodos(id)
	if err != nil {
		log.Println("todo service,", err.Error())
		return nil, err
	}
	return data, nil
}

func (ts service) GetTodo(pemilik *jwt.Token, idTodo uint) (todo.Todo, error) {
	id, err := middlewares.ExtractId(pemilik)
	if err != nil {
		log.Println("todo service,", err.Error())
		return todo.Todo{}, err
	}
	data, err := ts.tq.GetTodo(idTodo)
	if err != nil {
		log.Println("todo service,", err.Error())
		return todo.Todo{}, err
	}
	if data.UserID != id {
		log.Println("todo service, todo ini bukan milik anda")
		return todo.Todo{}, errors.New("unauthorized")
	}
	return data, nil
}

func (ts service) UpdateTodo(pemilik *jwt.Token, data todo.Todo) (todo.Todo, error) {
	id, err := middlewares.ExtractId(pemilik)
	if err != nil {
		log.Println("todo service,", err.Error())
		return todo.Todo{}, err
	}
	val, err := ts.tq.GetTodo(data.ID)
	if err != nil {
		log.Println("todo service,", err.Error())
		return todo.Todo{}, err
	}
	if val.UserID != id {
		log.Println("todo service, todo ini bukan milik anda")
		return todo.Todo{}, errors.New("unauthorized")
	}

	// validasi
	var validasi = handler.TodoUpdateRequest{
		Kegiatan:  data.Kegiatan,
		Deskripsi: data.Deskripsi,
		Deadline:  data.Deadline,
	}

	err = ts.v.Struct(validasi)
	if err != nil {
		log.Println("todo service,", err.Error())
		return todo.Todo{}, err
	}

	update, err := ts.tq.UpdateTodo(data)
	if err != nil {
		log.Println("todo service,", err.Error())
		return todo.Todo{}, err
	}
	return update, nil
}
