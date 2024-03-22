package handler

import (
	"PesbukAPI/features/comment"
	"PesbukAPI/helper"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
	s comment.CommentService
}

func NewHandler(service comment.CommentService) comment.CommentController {
	return &controller{
		s: service,
	}
}

func (ct *controller) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input CommentRequest
		err := c.Bind(&input)
		if err != nil {
			log.Println("error bind data:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
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

		newComment, err := ct.s.AddComment(token, input.Komentar)
		if err != nil {
			log.Println("error insert db:", err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan komentar", newComment))
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

        var input CommentRequest
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

        updatedComment, err := ct.s.UpdateComment(token, uint(id), comment.Comment{
            Komentar:     input.Komentar,
        })
        if err != nil {
            log.Println("gagal update komentar:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusForbidden, helper.CannotUpdate, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "komentar berhasil diperbarui", updatedComment))
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

        err = ct.s.DeleteComment(token, uint(id))
        if err != nil {
            log.Println("gagal menghapus komentar:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusForbidden, helper.CannotDelete, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "komentar berhasil dihapus", nil))
    }
}




func (ct *controller) ShowMyComments() echo.HandlerFunc {
    return func(c echo.Context) error {

        token, ok := c.Get("user").(*jwt.Token)
        if !ok {
            return c.JSON(http.StatusBadRequest,
                helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
        }

        comment, err := ct.s.GetCommentByOwner(token)
        if err != nil {
            log.Println("gagal mendapat komentar user:", err.Error())
            return c.JSON(http.StatusInternalServerError,
                helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
        }

        return c.JSON(http.StatusOK,
            helper.ResponseFormat(http.StatusOK, "komentar pengguna", comment))
    }
}