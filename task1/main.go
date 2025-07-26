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

	r.POST("/posts", postsHandler.CreatePostHandler)
	// r.GET("/posts/:id", handler.GetPostHandler)
	// r.GET("/posts", handler.GetPostsHandler)
	// r.DELETE("/posts/:id", handler.DeletePostHandler)
	// r.PUT("/posts/:id", handler.UpdatePostHandler)

	r.Run(":8085")
}
