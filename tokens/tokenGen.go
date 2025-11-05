package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("yoursecretkey")

func GenerateAllTokens(email, firstName, lastName, userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := jwt.MapClaims{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"user_id":    userId,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedToken, err = token.SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err = refreshToken.SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}
