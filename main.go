package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/LMS/DB"
	"github.com/LMS/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}
	PORT := os.Getenv("PORT")

	client := db.ConnectMongoDB()
	db := client.Database("LMS")
	collection := db.Collection("User")
	log.Println("Database and collection accessed:", collection.Name())

	route := echo.New()

	route.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom} [INFO] [<NONE>, ${path}]: {"http_method": "${method}", "http_ip": "${remote_ip}", "http_status": ${status}, "http_latency": ${latency}, "http_agent": "${user_agent}", "http_error": "${error}", "http_path": "${path}", "http_bytes_in": ${bytes_in}, "http_bytes_out": ${bytes_out}}` + "\n",
		CustomTimeFormat: "20060102150405.000",
		Output:           os.Stdout,
	}))
	route.Use(middleware.RequestID())

	route.GET("/", func(ctx echo.Context) error {

		return utils.Success(ctx, http.StatusOK, "Hello World", nil)
		// return utils.Error(ctx, http.StatusBadRequest, "Hello World")
	})

	fmt.Println("Server is running on port: ", PORT)
	route.Start(":" + PORT)
}
