package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"networks/handler"
)

func main() {
	connString := "postgres://admin:admin@localhost:5432/mydb"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	r := gin.Default()

	postsHandler := handler.NewHandler(conn)
	
	r.POST("/posts", postsHandler.CreatePost)
	//r.GET("/posts/:id", postsHandler.GetPostById)
	//r.GET("/posts", postsHandler.GetPosts)
	//r.DELETE("/posts/:id", postsHandler.DeletePost)
	//r.PUT("/posts/:id", postsHandler.UpdatePost)

	r.Run(":8085")
}
