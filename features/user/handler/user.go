package handler

import (
	"21-api/features/user"
	"21-api/helper"

	"log"
	"net/http"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
)

type controller struct {
	service user.UserService
}

func NewHandler(service user.UserService) user.UserController {
	return &controller{
		service: service,
	}
}

func (us *controller) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input RegisterRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				log.Println(err.Error())
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			log.Println(err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		var processInput user.User
		processInput.Hp = input.Hp
		processInput.Nama = input.Nama
		processInput.Password = input.Password

		err = us.service.AddUser(processInput) // ini adalah fungsi yang kita buat sendiri
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				return c.JSON(http.StatusConflict,
					helper.ResponseFormat(http.StatusConflict, "nomor sudah didaftarkan", nil))
			}
			if strings.Contains(err.Error(), "validasi") {
				return c.JSON(http.StatusBadRequest,
					helper.ResponseFormat(http.StatusBadRequest, "input tidak sesuai", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "selamat data sudah terdaftar", nil))
	}
}

func (us *controller) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		result, token, err := us.service.Login(input.Hp, input.Password)
		if err != nil {
			if strings.Contains(err.Error(), "validasi") {
				return c.JSON(http.StatusBadRequest,
					helper.ResponseFormat(http.StatusBadRequest, "input tidak sesuai", nil))
			}
			if strings.Contains(err.Error(), "password") {
				return c.JSON(http.StatusBadRequest,
					helper.ResponseFormat(http.StatusBadRequest, "password salah", nil))
			}
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "input tidak sesuai", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}

		var responseData LoginResponse
		responseData.Hp = result.Hp
		responseData.Nama = result.Nama
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "selamat anda berhasil login", responseData))

	}
}

func (us *controller) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {

		token, ok := c.Get("user").(*golangjwt.Token)
		if !ok {
			log.Println("error pada saat ngambil token")
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "error pada saat ngambil token", nil))
		}

		data, err := us.service.Profile(token)

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		var result = ProfilResponse{
			Nama: data.Nama,
			Hp:   data.Hp,
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

// func (us *controller) Update() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var hp = c.Param("id")
// 		var input model.User
// 		err := c.Bind(&input)
// 		if err != nil {
// 			log.Println("masalah baca input:", err.Error())
// 			if strings.Contains(err.Error(), "unsupport") {
// 				return c.JSON(http.StatusUnsupportedMediaType,
// 					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
// 			}
// 			return c.JSON(http.StatusBadRequest,
// 				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
// 		}

// 		isFound := us.Model.CekUser(hp)

// 		if !isFound {
// 			return c.JSON(http.StatusNotFound,
// 				helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
// 		}

// 		err = us.Model.Update(hp, input)

// 		if err != nil {
// 			log.Println("masalah database :", err.Error())
// 			return c.JSON(http.StatusInternalServerError,
// 				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan saat update data", nil))
// 		}

// 		return c.JSON(http.StatusOK,
// 			helper.ResponseFormat(http.StatusOK, "data berhasil di update", nil))
// 	}
// }

// func (us *controller) ListUser() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		listUser, err := us.Model.GetAllUser()
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError,
// 				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
// 		}
// 		return c.JSON(http.StatusOK,
// 			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", listUser))
// 	}
// }

// func (us *controller) Profile() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var hp = c.Param("id")
// 		result, err := us.Model.GetProfile(hp)

// 		if err != nil {
// 			if strings.Contains(err.Error(), "not found") {
// 				return c.JSON(http.StatusNotFound,
// 					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
// 			}
// 			return c.JSON(http.StatusInternalServerError,
// 				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
// 		}
// 		return c.JSON(http.StatusOK,
// 			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
// 	}
// }
