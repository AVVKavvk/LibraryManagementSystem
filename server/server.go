package server

import (
	"fmt"
	"os"

	db "github.com/AVVKavvk/LMS/DB"
	"github.com/AVVKavvk/LMS/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client
var DB *mongo.Database

func init() {
	Client = db.MongoClient
	DB = Client.Database("LMS")
}

func Server() {
	PORT := os.Getenv("PORT")
	route := echo.New()

	route.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom} [INFO] [<NONE>, ${path}]: {"http_method": "${method}", "http_ip": "${remote_ip}", "http_status": ${status}, "http_latency": ${latency}, "http_agent": "${user_agent}", "http_error": "${error}", "http_path": "${path}", "http_bytes_in": ${bytes_in}, "http_bytes_out": ${bytes_out}}` + "\n",
		CustomTimeFormat: "20060102150405.000",
		Output:           os.Stdout,
	}))

	route.Use(middleware.RequestID())
	router.RegisterRoutes(route)

	fmt.Println("Server is running on port: ", PORT)
	route.Start(":" + PORT)
}
