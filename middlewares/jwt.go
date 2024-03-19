package middlewares

import (
	"21-api/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(id uint) (string, error) {
	var data = jwt.MapClaims{}
	// custom data
	data["id"] = id
	// mandatory data
	data["iat"] = time.Now().Unix()
	data["exp"] = time.Now().Add(time.Hour * 3).Unix()

	var proccessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	result, err := proccessToken.SignedString([]byte(config.JWTSECRET))

	if err != nil {
		return "", err
	}

	return result, nil
}

func ExtractId(t *jwt.Token) (uint, error) {
	var userID uint

	expiredTime, err := t.Claims.GetExpirationTime()
	if err != nil {
		return 0, err
	}

	var eTime = *expiredTime

	if t.Valid && eTime.Compare(time.Now()) > 0 {
		var tokenClaims = t.Claims.(jwt.MapClaims)
		userID = uint(tokenClaims["id"].(float64))

		return userID, nil
	}

	return 0, errors.New("token tidak valid")
}
