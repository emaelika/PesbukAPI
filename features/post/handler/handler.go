package handler

import (
	"PesbukAPI/features/post"
	"PesbukAPI/helper"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
    s post.PostService
}

func NewHandler(service post.PostService) post.PostController {
    return &controller{
        s: service,
    }
}

func (ct *controller) Add() echo.HandlerFunc {
    return func(c echo.Context) error {
        // Ambil file yang diunggah dari form (jika ada)
        file, err := c.FormFile("picture")
        if err != nil && err != http.ErrMissingFile {
            log.Println("error bind data:", err.Error())
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        // Ambil nilai dari input content (jika ada)
        content := c.FormValue("content")

        // Jika tidak ada gambar atau konten yang disertakan, kembalikan error
        if file == nil && content == "" {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, "Anda harus menyertakan gambar atau konten", nil))
        }

        // Jika ada gambar yang diunggah, simpan gambar dan dapatkan path file
        var picturePath string
        if file != nil {
            // Buka file yang diunggah
            src, err := file.Open()
            if err != nil {
                log.Println("error opening uploaded file:", err.Error())
                return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
            }
            defer src.Close()

            // Simpan file ke dalam direktori upload
            dst, err := os.Create("image/picture/" + file.Filename)
            if err != nil {
                log.Println("error creating destination file:", err.Error())
                return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
            }
            defer dst.Close()

            // Salin konten file ke dalam file tujuan
            if _, err := io.Copy(dst, src); err != nil {
                log.Println("error copying file:", err.Error())
                return c.JSON(http.StatusInternalServerError,
                    helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
            }

            // Tentukan path file untuk digunakan dalam entitas Post
            picturePath = file.Filename
        }

        // Lanjutkan dengan fungsi AddPost
        token, ok := c.Get("user").(*jwt.Token)
        if !ok {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        newPost, err := ct.s.AddPost(token, picturePath, content)
        if err != nil {
            log.Println("error insert db:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
        }

        return c.JSON(http.StatusCreated,
            helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan postingan", newPost))
    }
}




func (ct *controller) Update() echo.HandlerFunc {
    return func(c echo.Context) error {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            log.Println("error parsing ID:", err.Error())
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        var input PostRequest
        if err := c.Bind(&input); err != nil {
            log.Println("error bind data:", err.Error())
            if strings.Contains(err.Error(), "unsupported") {
                return c.JSON(http.StatusUnsupportedMediaType,
                    helper.ResponseFormat(http.StatusUnsupportedMediaType, helper.UserInputFormatError, nil))
            }
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        token, ok := c.Get("user").(*jwt.Token)
        if !ok {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        updatedPost, err := ct.s.UpdatePost(token, uint(id), post.Post{
            Picture:     input.Picture,
            Content: input.Content,
        })
        if err != nil {
            log.Println("error update post:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "postingan berhasil diperbarui", updatedPost))
    }
}

func (ct *controller) Delete() echo.HandlerFunc {
    return func(c echo.Context) error {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
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

        err = ct.s.DeletePost(token, uint(id))
        if err != nil {
            log.Println("error delete post:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "postingan berhasil dihapus", nil))
    }
}

func (ct *controller) ShowAllPosts() echo.HandlerFunc {
    return func(c echo.Context) error {
        posts, err := ct.s.GetAllPosts()
        if err != nil {
            log.Println("error get all posts:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "semua postingan", posts))
    }
}


func (ct *controller) ShowPostByID() echo.HandlerFunc {
    return func(c echo.Context) error {
        // Ambil ID postingan dari path parameter
        postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil {
            log.Println("error parsing post ID:", err.Error())
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, "ID postingan tidak valid", nil))
        }

        // Panggil service untuk mendapatkan postingan berdasarkan ID postingan
        post, err := ct.s.GetPostByID(uint(postID))
        if err != nil {
            log.Println("error get post by ID:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
        }

        // Periksa jika postingan ditemukan atau tidak
        if post == nil {
            return c.JSON(http.StatusNotFound,
                helper.ResponseFormat(http.StatusNotFound, "Postingan tidak ditemukan", nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "Postingan", post))
    }
}

