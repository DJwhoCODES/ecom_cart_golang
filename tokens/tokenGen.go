package tokens

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func ValidateTokens(signedToken string) (jwt.MapClaims, string) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		return nil, err.Error()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, "invalid token"
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().After(time.Unix(int64(exp), 0)) {
			return nil, "token expired"
		}
	}

	return claims, "valid token"
}

func UpdateAllTokens(signedToken, signedRefreshToken, userId string, userCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"token":         signedToken,
			"refresh_token": signedRefreshToken,
			"updated_at":    time.Now(),
		},
	}

	_, err := userCollection.UpdateOne(
		ctx,
		bson.M{"user_id": userId},
		update,
	)

	return err
}
