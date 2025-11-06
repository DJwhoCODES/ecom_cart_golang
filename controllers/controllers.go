package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/djwhocodes/ecom_cart_golang/database"
	"github.com/djwhocodes/ecom_cart_golang/models"
	"github.com/djwhocodes/ecom_cart_golang/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var productCollection *mongo.Collection = database.ProductData(database.Client, "product")

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPassword(providedPassword, storedHashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(providedPassword))
	return err == nil
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := validator.New()
		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := database.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking email: " + err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists with this email!"})
			return
		}

		count, err = database.UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking phone: " + err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists with this phone number!"})
			return
		}

		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword

		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()

		token, refreshToken, err := tokens.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_Name, user.User_Id)
		if err != nil {
			log.Println("Error generating tokens:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		user.Access_Token = &token
		user.Refresh_Token = &refreshToken

		_, insertErr := database.UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed: " + insertErr.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User signed up successfully!",
			"user": gin.H{
				"id":        user.User_Id,
				"email":     user.Email,
				"firstName": user.First_Name,
				"lastName":  user.Last_Name,
				"token":     user.Access_Token,
				"refToken":  user.Refresh_Token,
			},
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var input models.User
		var foundUser models.User

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := database.UserCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		if !CheckPassword(*input.Password, *foundUser.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, refreshToken, err := tokens.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.First_Name,
			*foundUser.Last_Name,
			foundUser.User_Id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"access_token":  token,
				"refresh_token": refreshToken,
				"updated_at":    time.Now(),
			},
		}

		_, updateErr := database.UserCollection.UpdateOne(ctx, bson.M{"user_id": foundUser.User_Id}, update)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful!",
			"user": gin.H{
				"id":         foundUser.User_Id,
				"email":      foundUser.Email,
				"first_name": foundUser.First_Name,
				"last_name":  foundUser.Last_Name,
				"token":      token,
				"ref_token":  refreshToken,
			},
		})
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := productCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching products from database",
			})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &productList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error decoding product data",
			})
			return
		}

		if len(productList) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No products found",
			})
			return
		}

		c.JSON(http.StatusOK, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
