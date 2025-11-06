package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/djwhocodes/ecom_cart_golang/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewApplication(productCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "product id is empty"})
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user id is empty"})
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.productCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			log.Println("error adding product to cart:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product added to cart successfully"})
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.userCollection, productId, userQueryId)
		if err != nil {
			log.Println("error removing product from cart:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product removed from cart successfully"})
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("userID")
		if userQueryId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		if err != nil {
			log.Println("error buying cart items:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "all items purchased successfully"})
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.InstantBuy(ctx, app.productCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			log.Println("error performing instant buy:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product purchased successfully"})
	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userObjID}}}}
		unwindStage := bson.D{{Key: "$unwind", Value: "$user_cart"}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$user_cart.price"}}},
			{Key: "cart_items", Value: bson.D{{Key: "$push", Value: "$user_cart"}}},
		}}}

		cursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, groupStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching cart data"})
			return
		}
		defer cursor.Close(ctx)

		var cartResult []bson.M
		if err = cursor.All(ctx, &cartResult); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding cart data"})
			return
		}

		if len(cartResult) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Cart is empty"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"cart":    cartResult[0]["cart_items"],
			"total":   cartResult[0]["total"],
		})
	}
}
