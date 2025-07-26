package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Handler struct {
	LastID int
	Posts  map[int]Post
}

func (h *Handler) CreatePostHandler(c *gin.Context) {
	var post Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	h.LastID++
	post.ID = h.LastID
	h.Posts[h.LastID] = post
	c.JSON(http.StatusOK, post)
}

// func (h *Handler) GetPostHandler(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, "Невалидный id")
// 		log.Println("error in getPostHandler in strconv.Atoi: ", err)
// 		return
// 	}

// 	post, ok := posts[id]
// 	if !ok {
// 		c.JSON(http.StatusNotFound, "Такого поста нет")
// 		return
// 	}

// 	c.JSON(http.StatusOK, post)
// }

// func GetPostsHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, posts)
// }

// func DeletePostHandler(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, "Невалидный id")
// 		log.Println("error in getPostHandler in strconv.Atoi: ", err)
// 		return
// 	}

// 	delete(posts, id)

// 	c.JSON(http.StatusOK, "Пост удален")
// }

// func UpdatePostHandler(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, "Невалидный id")
// 		log.Println("error in getPostHandler in strconv.Atoi: ", err)
// 		return
// 	}

// 	var updatePostRequest UpdatePostRequest
// 	err = c.BindJSON(&updatePostRequest)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
// 		return
// 	}

// 	post, ok := posts[id]
// 	if !ok {
// 		c.JSON(http.StatusNotFound, "Такого поста нет")
// 		return
// 	}

// 	if updatePostRequest.Body != nil {
// 		post.Body = *updatePostRequest.Body
// 	}
// 	if updatePostRequest.Title != nil {
// 		post.Title = *updatePostRequest.Title
// 	}
// 	posts[id] = post

// 	c.JSON(http.StatusOK, posts[id])
// }
