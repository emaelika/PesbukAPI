package handler

import (
	"PesbukAPI/features/post"
	"PesbukAPI/helper"
	cloudnr "PesbukAPI/utils"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
	s  post.PostService
	cl *cloudinary.Cloudinary
	ct context.Context
	up string
}

func NewHandler(service post.PostService, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) post.PostController {
	return &controller{
		s:  service,
		cl: cld,
		ct: ctx,
		up: uploadparam,
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
		if file != nil {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusBadRequest,
					helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
			}

			formFile, err := file.Open()
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError, map[string]any{
						"message": "formfile error",
					})
			}

			link, err := cloudnr.UploadImage(ct.cl, ct.ct, formFile, ct.up)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					return c.JSON(http.StatusBadRequest, map[string]any{
						"message": "harap pilih gambar",
						"data":    nil,
					})
				} else {
					return c.JSON(http.StatusInternalServerError, map[string]any{
						"message": "kesalahan pada server",
						"data":    nil,
					})
				}
			}

			newPost, err := ct.s.AddPost(token, link, content)
			if err != nil {
				log.Println("error insert db:", err.Error())
				return c.JSON(http.StatusInternalServerError,
					helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
			}

			return c.JSON(http.StatusCreated,
				helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan postingan", newPost))
		}

		// Lanjutkan dengan fungsi AddPost
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		newPost, err := ct.s.AddPost(token, "", content)
		if err != nil {
			log.Println("error insert db:", err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}
		var output = PostResponse{
			ID:        newPost.ID,
			Fullname:  newPost.Fullname,
			Avatar:    newPost.Avatar,
			Picture:   newPost.Picture,
			Content:   newPost.Content,
			CreatedAt: newPost.CreatedAt,
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan postingan", output))
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

		// Periksa apakah ada file gambar baru yang diunggah
		file, err := c.FormFile("picture")
		if err != nil && err != http.ErrMissingFile {
			log.Println("error bind data:", err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}
		if file != nil {

			formFile, err := file.Open()
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError, map[string]any{
						"message": "formfile error",
					})
			}

			link, err := cloudnr.UploadImage(ct.cl, ct.ct, formFile, ct.up)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					return c.JSON(http.StatusBadRequest, map[string]any{
						"message": "harap pilih gambar",
						"data":    nil,
					})
				} else {
					return c.JSON(http.StatusInternalServerError, map[string]any{
						"message": "kesalahan pada server",
						"data":    nil,
					})
				}
			}
			updatedPost, err := ct.s.UpdatePost(token, uint(id), post.Post{
				Picture: link,
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

		// Jika ada file gambar baru, simpan file tersebut dan dapatkan path file baru
		// upload file

		// Panggil service untuk melakukan pembaruan
		updatedPost, err := ct.s.UpdatePost(token, uint(id), post.Post{
			Picture: "",
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

		// Panggil service untuk menghapus postingan
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
