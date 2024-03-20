package handler

import (
	"PesbukAPI/features/user"
	"PesbukAPI/helper"
	"PesbukAPI/middlewares"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service user.UserService
}

func NewUserHandler(s user.UserService) user.UserController {
	return &controller{
		service: s,
	}
}

func (ct *controller) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input user.User
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirimkan tidak sesuai", nil))
		}
		err = ct.service.Register(input)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "selamat data sudah terdaftar", nil))
	}
}
func (ct *controller) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirimkan tidak sesuai", nil))
		}

		var processData user.User
		processData.Email = input.Email
		processData.Password = input.Password

		result, token, err := ct.service.Login(processData)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		var responseData LoginResponse
		responseData.Email = result.Email
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil login", responseData))
	}
}
func (ct *controller) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}
		result, err := ct.service.Profile(token)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

func (ct *controller) Update() echo.HandlerFunc {
    return func(c echo.Context) error {
        // Mendapatkan token dari konteks
        token, ok := c.Get("user").(*jwt.Token)
        if !ok {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        // Mendapatkan ID pengguna dari token JWT
        decodedID := middlewares.DecodeToken(token)

        // Mendapatkan ID pengguna dari parameter URL
        idStr := c.Param("id")
        idFromParam, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            log.Println("error parsing ID:", err.Error())
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        // Membandingkan ID dari parameter dengan ID yang di-decode dari token
        if uint(idFromParam) != decodedID {
            return c.JSON(http.StatusForbidden,
                helper.ResponseFormat(http.StatusForbidden, "you can only update your own account", nil))
        }

        // Mendapatkan data yang akan diperbarui dari body permintaan
        var inputData user.User
        if err := c.Bind(&inputData); err != nil {
            log.Println("error binding data:", err.Error())
            if strings.Contains(err.Error(), "unsupported") {
                return c.JSON(http.StatusUnsupportedMediaType,
                    helper.ResponseFormat(http.StatusUnsupportedMediaType, helper.UserInputFormatError, nil))
            }
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        // Mengecek apakah ada file gambar yang diunggah
        file, _, err := c.Request().FormFile("image")
        if err == nil {
            // Baca file gambar menjadi byte slice
            imageData, err := io.ReadAll(file)
            if err != nil {
                log.Println("error reading image data:", err.Error())
                return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, "failed to read image data", nil))
            }
            // Set imageData ke inputData.Image
            inputData.Image = imageData
        }

        // Validasi bahwa setidaknya satu bidang yang akan diperbarui diisi
        if inputData.Username == "" && inputData.Password == "" && inputData.Name == "" &&
            inputData.Email == "" && inputData.Placeofbirth == "" && inputData.Dateofbirth == "" &&
            len(inputData.Image) == 0 {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, "at least one field must be provided for update", nil))
        }


        // Memanggil service untuk melakukan pembaruan data pengguna
        updatedUser, err := ct.service.Update(token, uint(idFromParam), inputData)
        if err != nil {
            log.Println("failed to update user:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.CannotUpdate, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "user successfully updated", updatedUser))
    }
}

func (ct *controller) Delete() echo.HandlerFunc {
    return func(c echo.Context) error {
        token, ok := c.Get("user").(*jwt.Token)
        if !ok {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        // Mendapatkan ID pengguna dari URL
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            log.Println("error parsing ID:", err.Error())
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        // Memanggil metode Delete dengan token dan id pengguna
        err = ct.service.Delete(token, uint(id))
        if err != nil {
            log.Println("gagal menghapus user:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusForbidden, helper.CannotDelete, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "pengguna berhasil dihapus", nil))
    }
}

func (ct *controller) GetUserByIDParam() echo.HandlerFunc {
    return func(c echo.Context) error {
        idStr := c.Param("id")
        idFromParam, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            log.Println("error parsing ID:", err.Error())
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        token, ok := c.Get("user").(*jwt.Token)
        if !ok {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        result, err := ct.service.GetUserByIDParam(token, uint(idFromParam))
        if err != nil {
            var code = http.StatusInternalServerError
            if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
                code = http.StatusBadRequest
            }
            return c.JSON(code,
                helper.ResponseFormat(code, err.Error(), nil))
        }

        // Mengirimkan respons dengan tipe MIME image/jpeg
        c.Response().Header().Set("Content-Type", "image/jpeg")
        return c.Blob(http.StatusOK, "image/jpeg", result.Image)
    }
}
