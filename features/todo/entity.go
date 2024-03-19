package todo

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TodoController interface {
	AddTodo() echo.HandlerFunc
	GetTodos() echo.HandlerFunc
	GetTodo() echo.HandlerFunc
	UpdateTodo() echo.HandlerFunc
}

type TodoService interface {
	AddTodo(pemilik *jwt.Token, kegiatanBaru Todo) (Todo, error)
	GetTodos(pemilik *jwt.Token) ([]Todo, error)
	GetTodo(pemilik *jwt.Token, idTodo uint) (Todo, error)
	UpdateTodo(pemilik *jwt.Token, newTodo Todo) (Todo, error)
}

type TodoQuery interface {
	AddTodo(newData Todo) (Todo, error)
	GetTodos(id uint) ([]Todo, error)
	GetTodo(idTodo uint) (Todo, error)
	UpdateTodo(newData Todo) (Todo, error)
}

type Todo struct {
	ID        uint
	UserID    uint
	Kegiatan  string
	Deskripsi string
	Deadline  string
}
