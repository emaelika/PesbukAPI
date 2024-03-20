package handler

import (
	"PesbukAPI/features/todo"
	"PesbukAPI/helper"
	"log"
	"net/http"
	"strconv"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
)

type controller struct {
	s todo.TodoService
}

func NewHandler(service todo.TodoService) todo.TodoController {
	return &controller{
		s: service,
	}
}
func (us *controller) AddTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, _ := c.Get("user").(*golangjwt.Token)

		var input TodoRequest
		err := c.Bind(&input)
		if err != nil {
			log.Println(err.Error())
			if strings.Contains(err.Error(), "unsupport") {

				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		var processInput todo.Todo
		processInput.Kegiatan = input.Kegiatan
		processInput.Deskripsi = input.Deskripsi
		processInput.Deadline = input.Deadline

		result, err := us.s.AddTodo(token, processInput) // ini adalah fungsi yang kita buat sendiri
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		var data TodoResponse

		data.Deadline = result.Deadline
		data.Deskripsi = result.Deskripsi
		data.Kegiatan = result.Kegiatan
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "selamat data sudah terdaftar", data))
	}
}

func (us *controller) GetTodos() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*golangjwt.Token)
		if !ok {
			log.Println("kesalahan dalam mengambil token")
			return c.JSON(http.StatusUnauthorized,
				helper.ResponseFormat(http.StatusUnauthorized, "terjadi kesalahan pada token", nil))

		}
		listTodo, err := us.s.GetTodos(token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		var result []GetTodosResponse
		for _, val := range listTodo {
			var data GetTodosResponse
			data.ID = val.ID
			data.Deadline = val.Deadline
			data.Deskripsi = val.Deskripsi
			data.Kegiatan = val.Kegiatan

			result = append(result, data)
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

func (us *controller) GetTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user").(*golangjwt.Token)

		idTodo, _ := strconv.Atoi(c.Param("id"))
		if idTodo == 0 {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "parameter yang dikirimkan tidak sesuai", nil))

		}

		val, err := us.s.GetTodo(token, uint(idTodo))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "tidak ditemukan to do", nil))
			}
			if strings.Contains(err.Error(), "unauthorized") {
				return c.JSON(http.StatusUnauthorized,
					helper.ResponseFormat(http.StatusUnauthorized, "todo tersebut bukan milik anda", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		var result TodoResponse

		result.Deadline = val.Deadline
		result.Deskripsi = val.Deskripsi
		result.Kegiatan = val.Kegiatan

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

func (us *controller) UpdateTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ambil ID
		token := c.Get("user").(*golangjwt.Token)

		idTodo, _ := strconv.Atoi(c.Param("id"))
		if idTodo == 0 {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "parameter yang dikirmkan tidak sesuai", nil))

		}

		// Bind Input
		var input TodoUpdateRequest
		err := c.Bind(&input)
		if err != nil {
			log.Println(err.Error())
			if strings.Contains(err.Error(), "unsupport") {

				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		// parsing
		var processInput = todo.Todo{
			ID:        uint(idTodo),
			Kegiatan:  input.Kegiatan,
			Deskripsi: input.Deskripsi,
			Deadline:  input.Deadline,
		}

		// Input
		update, err := us.s.UpdateTodo(token, processInput) // ini adalah fungsi yang kita buat sendiri
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}

		var result TodoResponse
		result.Deadline = update.Deadline
		result.Deskripsi = update.Deskripsi
		result.Kegiatan = update.Kegiatan

		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "berhasil memperbarui data", result))
	}
}
