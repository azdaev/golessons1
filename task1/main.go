package main

import (
	"networks/handler"

	"github.com/gin-gonic/gin"
)

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

func main() {
	r := gin.Default()

	postsHandler := handler.Handler{
		LastID: 0,
		Posts:  make(map[int]handler.Post),
	}

	r.POST("/posts", postsHandler.CreatePost)
	r.GET("/posts/:id", postsHandler.GetPostById)
	r.GET("/posts", postsHandler.GetPosts)
	r.DELETE("/posts/:id", postsHandler.DeletePost)
	r.PUT("/posts/:id", postsHandler.UpdatePost)

	r.Run(":8085")
}
