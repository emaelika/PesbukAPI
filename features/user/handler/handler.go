package handler

import (
	"PesbukAPI/features/user"
	"PesbukAPI/helper"
	"PesbukAPI/middlewares"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"

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
				helper.ResponseFormat(http.StatusBadRequest, "terdapat kesalahan pada data input", nil))
		}
		err = ct.service.Register(input)
		if err != nil {
			if strings.Contains(err.Error(), "validation") {
				return c.JSON(http.StatusBadRequest,
                    helper.ResponseFormat(http.StatusBadRequest, "terdapat kesalahan pada data input", nil))
			} else if strings.Contains(err.Error(), "database") {
				return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
			}
			return c.JSON(http.StatusConflict,
				helper.ResponseFormat(http.StatusConflict, "data yang dimasukkan sudah terdaftar", nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "data anda berhasil mendaftar", nil))
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
			if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
				return c.JSON(http.StatusInternalServerError,
					helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
			} else if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "terdapat kesalahan pada data input", nil))
		}

		var responseData LoginResponse
		responseData.Fullname = result.Fullname
		responseData.Token = token
        responseData.Avatar = result.Avatar

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
            return c.JSON(http.StatusNotFound,
                helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
        }

        // Mendapatkan ID pengguna dari token JWT
        decodedID := middlewares.DecodeToken(token)

        token = c.Get("user").(*jwt.Token)

        existingUser, err := ct.service.GetUserByIDParam(token, uint(decodedID))
        if err != nil {
            log.Println("error getting user:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
            }


        // Mendapatkan data yang akan diperbarui dari body permintaan
        var inputData user.User
        if err := c.Bind(&inputData); err != nil {
            log.Println("error binding data:", err.Error())
            if strings.Contains(err.Error(), "unsupported") {
                return c.JSON(http.StatusUnauthorized,
                    helper.ResponseFormat(http.StatusUnauthorized, "anda tidak bisa mengakses perintah ini", nil))
            }
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, "terdapat kesalahan pada data input", nil))
        }
 
        // Handle avatar upload
        avatar, err := c.FormFile("avatar")
        if err != nil && err != http.ErrMissingFile {
            log.Println("error uploading avatar:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
        }

        // Validasi bahwa setidaknya satu bidang yang akan diperbarui diisi
        if  inputData.Password == "" && inputData.Fullname == "" &&
            inputData.Email == ""  && inputData.Birthday == "" && inputData.Avatar == "" && avatar == nil {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, "terdapat kesalahan pada data input", nil))
        }

        // If avatar is provided, save it and update inputData
        if avatar != nil {
            if existingUser.Avatar != "" {
            if err := os.Remove("image/avatar/" + existingUser.Avatar); err != nil {
            log.Println("error deleting old avatar:", err.Error())
            // Anda bisa memilih untuk mengabaikan kesalahan penghapusan
            // atau mengembalikan kesalahan jika penghapusan gagal
            }
        }

            // Generate a unique filename
            filename := uuid.New().String() + filepath.Ext(avatar.Filename)
            // Open the uploaded file
            src, err := avatar.Open()
            if err != nil {
                log.Println("error opening uploaded file:", err.Error())
                return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
            }
            defer src.Close()

            // Create the destination file
            dst, err := os.Create("image/avatar/" + filename)
            if err != nil {
                log.Println("error creating destination file:", err.Error())
                return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
            }
            defer dst.Close()

            // Copy the file
            if _, err := io.Copy(dst, src); err != nil {
                log.Println("error copying file:", err.Error())
                return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
            }

            // Update the avatar path in inputData
            inputData.Avatar = filename
        }

        // Memanggil service untuk melakukan pembaruan data pengguna
        updatedUser, err := ct.service.Update(token, inputData)
        if err != nil {
            log.Println("failed to update user:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.CannotUpdate, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "berhasil mengubah data", updatedUser))
    }
}



func (ct *controller) Delete() echo.HandlerFunc {
    return func(c echo.Context) error {
        // Mendapatkan token dari konteks
        token, ok := c.Get("user").(*jwt.Token)
        if !ok {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, "terdapat kesalahan pada data input", nil))
        }

        // Mendapatkan avatar pengguna sebelum penghapusan
        existingUser, err := ct.service.Profile(token)
        if err != nil {
            log.Println("gagal mendapatkan avatar pengguna:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, "error pada server", nil))
        }

        // Memanggil metode Delete dengan token dan id pengguna
        err = ct.service.Delete(token)
        if err != nil {
            log.Println("gagal menghapus user:", err.Error())
            return c.JSON(http.StatusUnauthorized,
                helper.ResponseFormat(http.StatusUnauthorized, "anda tidak bisa mengakses perintah ini", nil))
        }

        // Hapus avatar pengguna dari direktori destinasi jika ada
        if existingUser.Avatar != "" {
            if err := os.Remove("image/avatar/" + existingUser.Avatar); err != nil {
                log.Println("gagal menghapus avatar:", err.Error())
                // Anda bisa memilih untuk mengabaikan kesalahan penghapusan
                // atau mengembalikan kesalahan jika penghapusan gagal
            }
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "berhasil menghapus data", nil))
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

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
    }
}

func (ct *controller) Avatar() echo.HandlerFunc {
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

        // Ambil nama file avatar dari respons JSON
        avatarFileName := result.Avatar

        // Buat path lengkap ke file avatar
        avatarFilePath := filepath.Join("image/avatar/", avatarFileName)

        // Kirimkan file avatar sebagai respons HTTP
        return c.File(avatarFilePath)
    }
}


