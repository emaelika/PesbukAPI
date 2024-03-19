package user

import (
	"21-api/helper"
	"21-api/middlewares"
	"21-api/model"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	Model model.UserModel
}

func (us *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input RegisterRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)

		if err != nil {

			var message = []string{}
			for _, val := range err.(validator.ValidationErrors) {
				message = append(message, fmt.Sprint("error pada ", val.Field()))
				// 	if val.Tag() == "required" {
				// 		message = append(message, fmt.Sprint(val.Field(), " wajib diisi"))
				// 	} else if val.Tag() == "min" {
				// 		message = append(message, fmt.Sprint(val.Field(), " minimal 10 digit"))
				// 	} else {
				// 		message = append(message, fmt.Sprint(val.Field(), " ", val.Tag()))
			}
			// }
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, message, nil))
		}

		var processInput model.User
		processInput.Hp = input.Hp
		processInput.Nama = input.Nama
		processInput.Password = input.Password

		err = us.Model.AddUser(processInput) // ini adalah fungsi yang kita buat sendiri
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				return c.JSON(http.StatusConflict,
					helper.ResponseFormat(http.StatusConflict, "nomor sudah didaftarkan", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "selamat data sudah terdaftar", nil))
	}
}

func (us *UserController) Login() echo.HandlerFunc {
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

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)

		if err != nil {
			for _, val := range err.(validator.ValidationErrors) {
				fmt.Println(val.Error())
			}
		}

		result, err := us.Model.Login(input.Hp, input.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		token, err := middlewares.GenerateJWT(result.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem, gagal memproses data", nil))
		}

		var responseData LoginResponse
		responseData.Hp = result.Hp
		responseData.Nama = result.Nama
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "selamat anda berhasil login", responseData))

	}
}

func (us *UserController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hp = c.Param("id")
		var input model.User
		err := c.Bind(&input)
		if err != nil {
			log.Println("masalah baca input:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		isFound := us.Model.CekUser(hp)

		if !isFound {
			return c.JSON(http.StatusNotFound,
				helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
		}

		err = us.Model.Update(hp, input)

		if err != nil {
			log.Println("masalah database :", err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan saat update data", nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "data berhasil di update", nil))
	}
}

func (us *UserController) ListUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		listUser, err := us.Model.GetAllUser()
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", listUser))
	}
}

func (us *UserController) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hp = c.Param("id")
		result, err := us.Model.GetProfile(hp)

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}
