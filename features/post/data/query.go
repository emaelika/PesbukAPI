package data

import (
	"PesbukAPI/features/comment"
	"PesbukAPI/features/post"
	"PesbukAPI/helper"
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
	var user User

	if err := pm.connection.First(&user).Where(" id = ? ", userid).Error; err != nil {
		return post.Post{}, err
	}
	return post.Post{Picture: inputProcess.Picture, Content: inputProcess.Content, CreatedAt: inputProcess.CreatedAt.String(), ID: inputProcess.ID, Fullname: user.Fullname, Avatar: user.Avatar}, nil
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

func (pm *model) GetAllPosts(paginasi helper.Pagination) ([]post.Post, int, error) {
	var proses = new([]Post)
	var count int64
	offset := (paginasi.Page - 1) * paginasi.Pagesize
	if err := pm.connection.Find(&proses).Order("created at asc").Count(&count).Error; err != nil {
		log.Println("repo error: ", err.Error())
		return []post.Post{}, 0, err
	}

	var selected = new([]Post)
	if err := pm.connection.Order("created_at desc").Find(&selected).Offset(offset).
		Limit(paginasi.Pagesize).Error; err != nil {
		log.Println("repo error: ", err.Error())
		return []post.Post{}, 0, err
	}

	// parsing result
	var results []post.Post
	for _, val := range *selected {
		// nambah item
		var result = post.Post{
			ID:        val.ID,
			Content:   val.Content,
			Picture:   val.Picture,
			CreatedAt: val.CreatedAt.String(),
		}

		// nambah fullname dan avatar
		var user User
		if err := pm.connection.First(&user).Where(" id = ? ", val.UserID).Error; err != nil {
			log.Println(err.Error())
			return nil, 0, err
		}
		result.Avatar = user.Avatar
		result.Fullname = user.Fullname
		// nambah comment count
		var comments []Comment
		if err := pm.connection.
			Find(&comments).Where("post id = ?", val.ID).Error; err != nil {
			log.Println(err.Error())
			return nil, 0, err
		}
		if comments == nil {
			result.CommentCount = 0
		} else {
			result.CommentCount = len(comments)
		}
		results = append(results, result)

	}

	return results, int(count), nil

}

func (pm *model) GetPostByID(postID uint) (*post.Post, error) {
	var result post.Post
	if err := pm.connection.Where("id = ?", postID).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
