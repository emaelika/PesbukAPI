package data

import (
	"PesbukAPI/features/comment"
	"PesbukAPI/features/post"
	"errors"
	"log"
	"os"

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

	return post.Post{Picture: inputProcess.Picture, Content: inputProcess.Content, CreatedAt: inputProcess.CreatedAt, ID: inputProcess.ID}, nil
}

func (pm *model) UpdatePost(userid uint, postID uint, data post.Post) (post.Post, error) {
	// Cek apakah ada gambar baru yang diunggah
	if data.Picture != "" {
		// Dapatkan postingan lama
		var oldPost post.Post
		if err := pm.connection.Where("id = ?", postID).First(&oldPost).Error; err != nil {
			return post.Post{}, err
		}

		// Hapus gambar lama jika ada
		if oldPost.Picture != "" {
			err := os.Remove("image/picture/" + oldPost.Picture)
			if err != nil {
				// Tangani kesalahan saat menghapus gambar lama
				log.Println("error deleting old picture:", err.Error())
				// Anda dapat memutuskan apakah ingin melanjutkan atau mengembalikan kesalahan di sini
			}
		}
	}

	// Lanjutkan dengan pembaruan postingan seperti biasa
	var qry = pm.connection.Where("user_id = ? AND id = ?", userid, postID).Updates(data)
	if err := qry.Error; err != nil {
		return post.Post{}, err
	}

	if qry.RowsAffected < 1 {
		return post.Post{}, errors.New("no data affected")
	}
	var oldPost post.Post
	if err := pm.connection.Where("id = ?", postID).First(&oldPost).Error; err != nil {
		return post.Post{}, err
	}
	data.CreatedAt = oldPost.CreatedAt
	return data, nil
}

func (pm *model) DeletePost(postID uint) error {
	result := pm.connection.Where("post_id = ?", postID).Delete(&comment.Comment{})
	if result.Error != nil {
		return result.Error
	}

	result = pm.connection.Unscoped().Delete(&Post{}, postID)
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
