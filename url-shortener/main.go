package main

import (
	"context"
	"log"

	"url-shortener-1/handler"
	"url-shortener-1/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	connString := "postgres://admin:admin@localhost:5432/links"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	linksRepository := repo.NewRepository(conn)
	linksHandler := handler.NewHandler(linksRepository)

	r := gin.Default()
	r.POST("/shorten", linksHandler.CreateLink)
	r.GET("/:path", linksHandler.Redirect)
	r.Run(":8080")
}
