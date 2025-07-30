package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type Handler struct {
	db *pgx.Conn
}

func NewHandler(db *pgx.Conn) Handler {
	return Handler{
		db: db,
	}
}

func (h *Handler) CreatePost(c *gin.Context) {
	var post Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	_, err = h.db.Exec(c, "insert into posts (title, body) values ($1, $2)", post.Title, post.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Что-то пошло не так")
		log.Println("error insert into posts: ", err)
	}

	c.JSON(http.StatusOK, post)
}

// h.db.QueryRow
//func (h *Handler) GetPosts(c *gin.Context) {
//	posts := h.Posts
//	if len(posts) == 0 {
//		c.JSON(http.StatusNotFound, "Такого поста нет")
//		return
//	}
//
//	c.JSON(http.StatusOK, posts)
//}
// h.db.Query
//func (h *Handler) GetPostById(c *gin.Context) {
//	posts := h.Posts
//	if len(posts) == 0 {
//		c.JSON(http.StatusNotFound, "Такого поста нет")
//		return
//	}
//
//	idStr := c.Param("id")
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "Невалидный id")
//		log.Println("error in getPostHandler in strconv.Atoi: ", err)
//		return
//	}
//
//	post, ok := h.Posts[id]
//	if !ok {
//		c.JSON(http.StatusNotFound, "Такого поста нет")
//		return
//	}
//
//	c.JSON(http.StatusOK, post)
//}
// h.db.Exec
//func (h *Handler) DeletePost(c *gin.Context) {
//	idStr := c.Param("id")
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "Невалидный id")
//		log.Println("error in getPostHandler in strconv.Atoi: ", err)
//		return
//	}
//
//	post, ok := h.Posts[id]
//	if !ok {
//		c.JSON(http.StatusNotFound, "Такого поста нет")
//		return
//	}
//
//	delete(h.Posts, post.ID)
//	c.JSON(http.StatusOK, "Пост успешно удален")
//}
// h.db.Exec
//func (h *Handler) UpdatePost(c *gin.Context) {
//	idStr := c.Param("id")
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "Невалидный id")
//		log.Println("error in getPostHandler in strconv.Atoi: ", err)
//		return
//	}
//
//	var updatePostRequest UpdatePostRequest
//	err = c.BindJSON(&updatePostRequest)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
//		return
//	}
//
//	post, ok := h.Posts[id]
//	if !ok {
//		c.JSON(http.StatusNotFound, "Такого поста нет")
//		return
//	}
//
//	if updatePostRequest.Body != nil {
//		post.Body = *updatePostRequest.Body
//	}
//	if updatePostRequest.Title != nil {
//		post.Title = *updatePostRequest.Title
//	}
//	h.Posts[id] = post
//
//	c.JSON(http.StatusOK, h.Posts[id])
//}
