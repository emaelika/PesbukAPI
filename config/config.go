package config

import (
	comment "PesbukAPI/features/comment/data"
	post "PesbukAPI/features/post/data"
	user "PesbukAPI/features/user/data"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var JWTSECRET = ""

type AppConfig struct {
	DBUsername string
	DBPassword string
	DBPort     string
	DBHost     string
	DBName     string
}

func assignEnv(c AppConfig) (AppConfig, bool) {
	var missing = false
	if val, found := os.LookupEnv("DBUsername"); found {
		c.DBUsername = val
	} else {
		missing = true
	}
	if val, found := os.LookupEnv("DBPassword"); found {
		c.DBPassword = val
	} else {
		missing = true
	}
	if val, found := os.LookupEnv("DBPort"); found {
		c.DBPort = val
	} else {
		missing = true
	}
	if val, found := os.LookupEnv("DBHost"); found {
		c.DBHost = val
	} else {
		missing = true
	}
	if val, found := os.LookupEnv("DBName"); found {
		c.DBName = val
	} else {
		missing = true
	}
	if val, found := os.LookupEnv("JWT_SECRET"); found {
		JWTSECRET = val
	} else {
		missing = true
	}

	return c, missing
}

func InitConfig() AppConfig {
	var result AppConfig
	var missing = false
	result, missing = assignEnv(result)
	if missing {
		godotenv.Load(".env")
		result, _ = assignEnv(result)
	}

	return result
}

func InitSQL(c AppConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DBUsername, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("terjadi error", err.Error())
		return nil
	}

	db.AutoMigrate(&user.User{}, &comment.Comment{}, &post.Post{})

	return db
}

func InitDir(c AppConfig) error {
	// init dir upload avatar
	dir := "uploads"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(dir, "does not exist")
		err = os.Mkdir(dir, 0644)
		if err != nil {
			return err
		}
	}
	dir = "uploads/avatars"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(dir, "does not exist")
		err = os.Mkdir(dir, 0644)
		if err != nil {
			return err
		}
	}

	dir = "uploads/pictures"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(dir, "does not exist")
		err = os.Mkdir(dir, 0644)
		if err != nil {
			return err
		}
	}
	return nil
	// init dir upload gambar

}
