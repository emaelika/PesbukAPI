package services

import (
	"PesbukAPI/features/comment"
	"PesbukAPI/helper"
	"PesbukAPI/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	m comment.CommentModel
	v *validator.Validate
}

func NewCommentService(model comment.CommentModel) comment.CommentService {
	return &service{
		m: model,
		v: validator.New(),
	}
}

func (s *service) AddComment(userid *jwt.Token, postID uint, contentBaru string) (comment.Comment, error) {
    id := middlewares.DecodeToken(userid)
    if id == 0 {
        log.Println("error decode token", "token tidak ditemukan")
        return comment.Comment{}, errors.New("data tidak valid")
    }

    err := s.v.Var(contentBaru, "required")
    if err != nil {
        log.Println("error validasi deskripsi", err.Error())
        return comment.Comment{}, err
    }

    result, err := s.m.AddComment(id, postID, contentBaru)
    if err != nil {
        return comment.Comment{}, errors.New(helper.ServerGeneralError)
    }

    return result, nil
}


func (s *service) UpdateComment(userid *jwt.Token, commentID uint, data comment.Comment) (comment.Comment, error) {
	id := middlewares.DecodeToken(userid)
	if id == 0 {
		log.Println("error decode token:", "token tidak ditemukan")
		return comment.Comment{}, errors.New("data tidak valid")
	}

	err := s.v.Struct(data)
	if err != nil {
		log.Println("error validasi aktivitas", err.Error())
		return comment.Comment{}, err
	}

	result, err := s.m.UpdateComment(id, commentID, data)
	if err != nil {
		return comment.Comment{}, errors.New(helper.CannotUpdate)
	}

	return result, nil
}

func (s *service) DeleteComment(userid *jwt.Token, commentID uint) error {
    id := middlewares.DecodeToken(userid)
    if id == 0 { // Periksa apakah id adalah nol atau tidak
        log.Println("error decode token:", "token tidak ditemukan")
        return errors.New("data tidak valid")
    }

    err := s.m.DeleteComment(id, commentID) // Memanggil DeleteBook dengan pemilik dan bookID
    if err != nil {
        return errors.New(helper.CannotDelete)
    }

    return nil
}



func (s *service) GetCommentByOwner(userid *jwt.Token) ([]comment.Comment, error) {
	id := middlewares.DecodeToken(userid)
	if id == 0 {
		log.Println("error decode token:", "token tidak ditemukan")
		return nil, errors.New("data tidak valid")
	}

	books, err := s.m.GetCommentByOwner(id)
	if err != nil {
		return nil, errors.New(helper.ServerGeneralError)
	}

	return books, nil
}

func (s *service) GetAllComments() ([]comment.Comment, error) {
	comments, err := s.m.GetAllComments()
	if err != nil {
		return nil, errors.New(helper.ServerGeneralError)
	}
	return comments, nil
}


