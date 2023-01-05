package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"marketapi/controller"
	"marketapi/database"
	"marketapi/database/migrate"
	_ "marketapi/docs"
	"marketapi/middleware"
)

func init() {
	loadEnv()
	database.CreateDB()
	migrate.MigrateTablesToDB()
}

// @title Market API
// @version 1.0
// @description This is a sample server Market API.

// @contact.name Seyma TUTAN GUN
// @contact.email seymatutan@gmail.com

// @host localhost:9000
// @BasePath /
func main() {

	router := gin.Default()

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)
	publicRoutes.POST("/refresh", controller.RefreshToken)

	protectedCardRoutes := router.Group("/api/card")
	protectedCardRoutes.Use(middleware.JWTAuthMiddleware())
	protectedCardRoutes.POST("/add", controller.AddCard)
	protectedCardRoutes.POST("/list", controller.ListCard)
	protectedCardRoutes.GET("/get/:id", controller.GetCardDetailsById)
	protectedCardRoutes.PUT("/update/:id", controller.UpdateCard)
	protectedCardRoutes.PUT("/deactivate/:id", controller.DeactivateCard)
	protectedCardRoutes.DELETE("/delete/:id", controller.DeleteCard)

	protectedWalletRoutes := router.Group("/api/wallet")
	protectedWalletRoutes.Use(middleware.JWTAuthMiddleware())
	protectedWalletRoutes.GET("/get", controller.GetWalletDetails)
	protectedWalletRoutes.POST("/load", controller.LoadMoneyToWallet)
	protectedWalletRoutes.GET("/qr", controller.GenerateQR)

	protectedPaymentRoutes := router.Group("/api/payment")
	protectedPaymentRoutes.Use(middleware.JWTAuthMiddleware())
	protectedPaymentRoutes.POST("/complete", controller.CompletePayment)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":9000")
	fmt.Println("Server running on port 9000")

}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
