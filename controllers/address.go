package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/djwhocodes/ecom_cart_golang/database"
	"github.com/djwhocodes/ecom_cart_golang/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "user")

func AddAddress() gin.HandlerFunc {
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

		var newAddress models.Address
		if err := c.BindJSON(&newAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address data"})
			return
		}

		newAddress.ID = primitive.NewObjectID()

		update := bson.M{
			"$push": bson.M{
				"address_details": newAddress,
			},
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userObjID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add address"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Address added successfully",
			"address": newAddress,
		})
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		addressID := c.Query("address_id")
		if addressID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		addressObjID, err := primitive.ObjectIDFromHex(addressID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		filter := bson.M{"_id": userObjID}
		update := bson.M{
			"$pull": bson.M{
				"address_details": bson.M{"_id": addressObjID},
			},
		}

		result, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting address"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No matching address found for this user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
	}
}
