package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/suhas-developer07/EdwinNova-Server/internals/application"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/rabbitmq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	mongoURI := getEnv("MONGO_URI", "mongodb+srv://suhas:Fordmustang1969@suhas.cbbha.mongodb.net/EdwinNova")
	dbName := getEnv("MONGO_DB", "edwinnova")
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(dbName)

	/* RabbitMq Initialization */
	RabbitMQ_URI := os.Getenv("RABBITMQ_URI")

	rabbitmqConn,err := rabbitmq.New(RabbitMQ_URI)

	if err != nil {
		log.Fatalln("RabbitMq connection failed:Error",err)
	}
	
	defer rabbitmqConn.Conn.Close()
	defer rabbitmqConn.Channel.Close()

	repo := application.NewRepository(db)
	svc := application.NewService(repo)
	handler := application.NewHandler(svc, uploadDir)

	fmt.Println("Mongo url ", os.Getenv("MONGO_URI"))

	e := echo.New()

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
