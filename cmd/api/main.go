package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/suhas-developer07/EdwinNova-Server/internals/application"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mail"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mongo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dbName := getEnv("MONGO_DB", "edwinnova")
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")

	/* Database Initialization */
	client, err := mongo.InitMongo(mongo.Config{
		URI:         os.Getenv("MONGO_URI"),
		MaxPoolSize: 50,
		MinPoolSize: 5,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	db := client.Database(dbName)
	defer mongo.DisconnectMongo()

	/* SMTP Initialization */
	smtpClient, err := mail.NewSMTPClient()
	if err != nil {
		log.Fatalln("Failed to initialize SMTP client:", err)
	}

	/* Internals */
	repo := application.NewRepository(db)
	svc := application.NewService(repo, smtpClient)
	handler := application.NewHandler(svc, uploadDir)

	/* Echo Initialization */
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))

	e.POST("/applications", handler.CreateApplication)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}