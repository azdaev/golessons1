package main

import (
	"context"
	"log"
	"os"
	"time"

	"blog/internal/handler"
	"blog/internal/repo"
	"blog/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connString := os.Getenv("DBSTRING")
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

	r.GET("/users/:id", h.GetUser)
	r.GET("/users/:id/posts", h.GetPostsByUserID)
	r.GET("/users", h.GetUsers)
	r.POST("/users", h.CreateUser)

	r.GET("/posts", h.GetPosts)
	// r.GET("/post/:id", h.GetPost)
	// r.GET("/post/:id/comments", h.GetCommentsByPostId)
	// r.POST("/posts/:id", h.CreatePost)

	// r.POST("/posts/:id/comments", h.CreateComment)

	r.Run(":8080")
}
