/**********************************************
** @Des: JSON Web Token
** @Author: victor
** @Date:   2017-12-12 10:10:00
** @Last Modified by:   victor
** @Last Modified time: 2017-12-12 10:10:00
***********************************************/

package jwtoken

import (
	"time"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var SecretKey []byte = []byte("phoenixcore")

// 基础版

type MyCustomClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func JWTSign(uid string, expiredDays int64) (string, error) {
	// Create the Claims
	claims := MyCustomClaims{
		uid,
		jwt.StandardClaims{
				NotBefore: int64(time.Now().Unix()),
				ExpiresAt: int64(time.Now().Unix() + 86400 * expiredDays),
				Issuer:    "phoenix",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func JWTParse(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if token == nil {
		return "", errors.New("Null Token")
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return "", err
	}
}

func JWTVerify(tokenString string) (bool, string) {
	uid, err := JWTParse(tokenString)
	if err != nil {
		return false, ""
	}
	return true, uid
}

// 高级版

type MySuperCustomClaims struct {
	UserId string `json:"user_id"`
	DeviceId string `json:"device_id"`
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func SuperJWTSign(uid, deviceId, scope string, expiredDays int64) (string, error) {
	// Create the Claims
	claims := MySuperCustomClaims{
		uid,
		deviceId,
		scope,
		jwt.StandardClaims{
				NotBefore: int64(time.Now().Unix()),
				ExpiresAt: int64(time.Now().Unix() + 86400 * expiredDays),
				Issuer:    "phoenix",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func SuperJWTParse(tokenString string) (*MySuperCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MySuperCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if token == nil {
		return &MySuperCustomClaims{}, errors.New("Null Token")
	}
	if claims, ok := token.Claims.(*MySuperCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return claims, err
	}
}

func SuperJWTVerify(tokenString string) (bool, *MySuperCustomClaims) {
	claims, err := SuperJWTParse(tokenString)
	if err != nil {
		return false, claims
	}
	return true, claims
}
