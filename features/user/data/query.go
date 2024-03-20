package data

import (
	"PesbukAPI/features/user"
	"bytes"
	"errors"
	"image"
	"image/jpeg"
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

func (m *model) CekUser(username string) bool {
	var data User
	if err := m.connection.Where("username = ?", username).First(&data).Error; err != nil {
		return false
	}
	return true
}

func (m *model) Login(username string) (user.User, error) {
	var result user.User
	if err := m.connection.Where("username = ? ", username).First(&result).Error; err != nil {
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
    result := m.connection.Delete(&User{}, id)
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

    if newData.Name != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("name", newData.Name).Error; err != nil {
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

    if newData.Username != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("username", newData.Username).Error; err != nil {
            tx.Rollback()
            return user.User{}, err
        }
    }
    
    if newData.Placeofbirth != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("placeofbirth", newData.Placeofbirth).Error; err != nil {
            tx.Rollback()
            return user.User{}, err
        }
    }

    if newData.Dateofbirth != "" {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("dateofbirth", newData.Dateofbirth).Error; err != nil {
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

    if len(newData.Image) > 0 {
        if err := tx.Model(&user.User{}).Where("id = ?", id).Update("image", newData.Image).Error; err != nil {
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

func (m *model) GetUserByIDE(id uint) (user.User, error) {
    var result user.User
    if err := m.connection.Model(&User{}).Where("id = ?", id).First(&result).Error; err != nil {
        return user.User{}, err
    }

    // Konversi data gambar ke format PNG
    img, _, err := image.Decode(bytes.NewReader(result.Image))
    if err != nil {
        return user.User{}, err
    }

    var buf bytes.Buffer
    if err := jpeg.Encode(&buf, img, nil); err != nil {
        return user.User{}, err
    }

    result.Image = buf.Bytes() // Mengganti data gambar dengan data JPEG yang baru

    return result, nil
}