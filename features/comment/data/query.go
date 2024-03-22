package data

import (
	"PesbukAPI/features/comment"
	"errors"

	"gorm.io/gorm"
)

type model struct {
	connection *gorm.DB
}

func New(db *gorm.DB) comment.CommentModel {
	return &model{
		connection: db,
	}
}

func (cm *model) AddComment(userid uint, komentarBaru string) (comment.Comment, error) {
	var inputProcess = Comment{Komentar: komentarBaru, UserID: userid}
	if err := cm.connection.Create(&inputProcess).Error; err != nil {
		return comment.Comment{}, err
	}
	return comment.Comment{Komentar: inputProcess.Komentar},nil
}

func (cm *model) UpdateComment(userid uint, commentID uint, data comment.Comment) (comment.Comment, error) {
	var qry = cm.connection.Where("user_id = ? AND id = ?", userid, commentID).Updates(data)
	if err := qry.Error; err != nil {
		return comment.Comment{}, err
	}

	if qry.RowsAffected < 1 {
		return comment.Comment{}, errors.New("no data affected")
	}

	return data, nil
}

func (cm *model) DeleteComment(userid uint, commentID uint) error {
    result := cm.connection.Unscoped().Where("user_id = ? AND id = ?", userid, commentID).Delete(&Comment{}) // Menambahkan kondisi untuk pemilik dan ID buku
    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return errors.New("no data affected")
    }

    return nil
}

func (cm *model) GetCommentByOwner(userid uint) ([]comment.Comment, error) {
	var result []comment.Comment
	if err := cm.connection.Where("user_id = ?",userid).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}