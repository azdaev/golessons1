package main

import (
	"context"
	"log"
	"time"

	"blog/internal/handler"
	"blog/internal/repo"
	"blog/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	connString := "postgres://postgres:postgres@localhost:5432/shortener-db"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	repository := repo.New(conn)
	svc := service.New(repository)
	h := handler.New(svc)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.Run(":8080")
}
