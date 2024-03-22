package data

import (
	"PesbukAPI/features/post"
	"errors"

	"gorm.io/gorm"
)

type model struct {
    connection *gorm.DB
}

func New(db *gorm.DB) post.PostModel {
    return &model{
        connection: db,
    }
}

func (pm *model) AddPost(userid uint, pictureBaru string, contentBaru string) (post.Post, error) {
    var inputProcess = Post{Picture: pictureBaru, Content: contentBaru, UserID: userid}
    if err := pm.connection.Create(&inputProcess).Error; err != nil {
        return post.Post{}, err
    }

    return post.Post{Picture: inputProcess.Picture, Content: inputProcess.Content}, nil
}

func (pm *model) UpdatePost(userid uint, postID uint, data post.Post) (post.Post, error) {
    var qry = pm.connection.Where("user_id = ? AND id = ?", userid, postID).Updates(data)
    if err := qry.Error; err != nil {
        return post.Post{}, err
    }

    if qry.RowsAffected < 1 {
        return post.Post{}, errors.New("no data affected")
    }

    return data, nil
}

func (pm *model) DeletePost(postID uint) error {
    result := pm.connection.Unscoped().Delete(&Post{}, postID)
    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return errors.New("no data affected")
    }

    return nil
}

func (pm *model) GetAllPosts() ([]post.Post, error) {
    var result []post.Post
    if err := pm.connection.Find(&result).Error; err != nil {
        return nil, err
    }
    return result, nil
}

func (pm *model) GetPostByID(postID uint) (*post.Post, error) {
    var result post.Post
    if err := pm.connection.Where("id = ?", postID).First(&result).Error; err != nil {
        return nil, err
    }

    return &result, nil
}

