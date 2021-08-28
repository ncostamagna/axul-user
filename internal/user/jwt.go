package user

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key = os.Getenv("TOKEN")

type UserClaims struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

func CreateJWT(id, username string) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
		ID:       id,
		UserName: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("couldn't SignedString %w", err)
	}
	return ss, nil
}

func AccessJWT(token string) (*UserClaims, error) {

	verificationToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(beforeVeritificationToken *jwt.Token) (interface{}, error) {
		// validamos que el algoritmo de encriptacion sea el mismo
		if beforeVeritificationToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("SOMEONE TRIED TO HACK changed signing method")
		}
		return []byte(key), nil
	})

	if err != nil || !verificationToken.Valid {
		return nil, InvalidAuthentication
	}

	return verificationToken.Claims.(*UserClaims), nil

}
