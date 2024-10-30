package server

import (
	"fmt"
	"os"

	db "github.com/AVVKavvk/LMS/DB"
	"github.com/AVVKavvk/LMS/router"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client
var DB *mongo.Database

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables.")
	}

	Client = db.MongoClient
	DB = Client.Database("LMS")
}

func Server() {
	PORT := os.Getenv("PORT")
	ClientURL := os.Getenv("ClientURL")

	route := echo.New()

	route.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom} [INFO] [<NONE>, ${path}]: {"http_method": "${method}", "http_ip": "${remote_ip}", "http_status": ${status}, "http_latency": ${latency}, "http_agent": "${user_agent}", "http_error": "${error}", "http_path": "${path}", "http_bytes_in": ${bytes_in}, "http_bytes_out": ${bytes_out}}` + "\n",
		CustomTimeFormat: "20060102150405.000",
		Output:           os.Stdout,
	}))

	route.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{ClientURL},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))

	route.Use(middleware.RequestID())
	router.RegisterRoutes(route)

	fmt.Println("Server is running on port:", PORT)
	route.Start(":" + PORT)
}
