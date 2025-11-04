package main

import (
	"os"

	"github.com/djwhocodes/ecom_cart_golang/controllers"
	"github.com/djwhocodes/ecom_cart_golang/middleware"
	"github.com/djwhocodes/ecom_cart_golang/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database))

	router := gin.Default()

	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("instantbuy", app.InstantBuy())
}
