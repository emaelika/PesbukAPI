package data

import (
	"PesbukAPI/features/comment"
	"PesbukAPI/features/post"
	"PesbukAPI/features/user"
	"errors"
	_ "image/png"

	"gorm.io/gorm"
)

type model struct {
	connection *gorm.DB
}

func New(db *gorm.DB) user.UserModel {
	return &model{
		connection: db,
	}
}

func (m *model) AddUser(newData user.User) error {
	err := m.connection.Create(&newData).Error
	if err != nil {
		return errors.New("terjadi masalah pada database")
	}
	return nil
}

func (m *model) Login(email string) (user.User, error) {
	var result user.User
	if err := m.connection.Where("email = ? ", email).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

func (m *model) GetUserByID(id uint) (user.User, error) {
	var result user.User
	if err := m.connection.Model(&User{}).Where("id = ?", id).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

func (m *model) Delete(id uint) error {
    result := m.connection.Where("user_id = ?", id).Delete(&post.Post{}, &comment.Comment{})
    if result.Error != nil {
        return result.Error
    }

    result = m.connection.Delete(&User{}, id)
    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return errors.New("no data affected")
    }

    return nil
}


func (m *model) Update(id uint, newData user.User) (user.User, error) {
    var updatedUser user.User

    tx := m.connection.Begin()

    if newData.Fullname != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("fullname", newData.Fullname).Error; err != nil {
            tx.Rollback()
            return user.User{}, err
        }
    }

    if newData.Email != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("email", newData.Email).Error; err != nil {
            tx.Rollback()
            return user.User{}, err
        }
    }
    
    if newData.Password != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("password", newData.Password).Error; err != nil {
            tx.Rollback()
            return user.User{}, err
        }
    }

    if newData.Birthday != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("birthday", newData.Birthday).Error; err != nil {
            tx.Rollback()
            return user.User{}, err
        }
    }

    if newData.Avatar != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("avatar", newData.Avatar).Error; err != nil {
            tx.Rollback()
            return user.User{}, err
        }
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return user.User{}, err
    }

    // Ambil data user yang telah diperbarui
    if err := m.connection.First(&updatedUser, id).Error; err != nil {
        return user.User{}, err
    }

    return updatedUser, nil
}
