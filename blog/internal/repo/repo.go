package repo

import (
	"blog/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Post struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	PostId    int       `json:"post_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePost(ctx context.Context, user_id int, title string, body string) (string, error) {
	_, err := r.db.Exec(ctx, "insert into posts (user_id, title, body) VALUES ($1, $2, $3)", user_id, title, body)
	if err != nil {
		return "", fmt.Errorf("error in CreatePost: %w", err)
	}

	return "Пост успешно создан", nil
}

func (r *Repository) GetPost(ctx context.Context, id int) (Post, error) {
	var post Post
	err := r.db.QueryRow(
		ctx,
		`select
				id,
				user_id,
				title,
				body,
				views,
				created_at,
				updated_at
		   from posts
		  where id=$1`,
		id,
	).Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Body,
		&post.Views,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return Post{}, fmt.Errorf("error in GetPost: %w", err)
	}

	return post, nil
}

func (r *Repository) RemovePost(ctx context.Context, id int) (string, error) {
	_, err := r.db.Exec(ctx, "delete from posts where id=$1", id)
	if err != nil {
		return "", fmt.Errorf("error in RemovePost: %w", err)
	}

	return "Пост удален", nil
}

func (r *Repository) UpdatePost(ctx context.Context, id int, title string, body string) (string, error) {
	_, err := r.db.Exec(ctx, "update posts set title=$1, body=$2 where id=$3", title, body, id)
	if err != nil {
		return "", fmt.Errorf("error in UpdatePost: %w", err)
	}

	return "Пост обновлен", nil
}

type CreateUser struct {
	Name           string
	HashedPassword string
	Email          string
	IsAdmin        bool
}

func (r *Repository) GetUser(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		ctx,
		`select
				id,
				name,
				email,
				is_admin,
				created_at,
				updated_at
		   from users
		  where id=$1`,
		id,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("error in GetUser: %w", err)
	}
	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user CreateUser) error {
	_, err := r.db.Exec(
		ctx,
		"insert into users (name, hashed_password, email, is_admin) VALUES ($1, $2, $3, $4)",
		user.Name,
		user.HashedPassword,
		user.Email,
		user.IsAdmin,
	)
	if err != nil {
		return fmt.Errorf("error in CreateUser: %w", err)
	}

	return nil
}

func (r *Repository) RemoveUser(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "delete from users where id=$1", id)
	if err != nil {
		return fmt.Errorf("error in RemoveUser: %w", err)
	}

	return nil
}

type UpdateUser struct {
	Name    *string
	Email   *string
	IsAdmin *bool
}

func (r *Repository) UpdateUser(ctx context.Context, id int, user UpdateUser) (model.User, error) {
	var updatedUser model.User
	err := r.db.QueryRow(
		ctx,
		`update user
			set name=COALESCE($1, name),
				email=COALESCE($2, email),
				is_admin=COALESCE($3, is_admin)б
				updated_at=now()
	      where id=$4
	  returning id,
				name,
				email,
				is_admin,
				created_at,
				updated_at;`,
		user.Name,
		user.Email,
		user.IsAdmin,
		id,
	).Scan(
		&updatedUser.Id,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.IsAdmin,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("error in UpdateUser: %w", err)
	}

	return updatedUser, nil
}

// TODO: type CreateComment struct {}
// Return only error
func (r *Repository) CreateComment(ctx context.Context, user_id int, post_id int, body string) (string, error) {
	_, err := r.db.Exec(ctx, "insert into comments (user_id, post_id, body) VALUES ($1, $2, $3)", user_id, post_id, body)
	if err != nil {
		return "", fmt.Errorf("error in CreateComment: %w", err)
	}
	return "Комментарий успешно создан", nil
}

func (r *Repository) RemoveComment(ctx context.Context, id int) (string, error) {
	_, err := r.db.Exec(ctx, "delete from comments where id=$1", id)
	if err != nil {
		return "", fmt.Errorf("error in RemoveComment: %w", err)
	}
	return "Комментарий удалён", nil
}

// TODO: type UpdateComment struct {}
func (r *Repository) UpdateComment(ctx context.Context, id int, body string) (string, error) {
	_, err := r.db.Exec(ctx, "update comments set body=$1 where id=$2", body, id)
	if err != nil {
		return "", fmt.Errorf("error in UpdateComment: %w", err)
	}
	return "Комментарий обновлён", nil
}

func (r *Repository) GetComment(ctx context.Context, id int) (Comment, error) {
	var comment Comment
	err := r.db.QueryRow(ctx, "select id, user_id, post_id, body, created_at, updated_at from comments where id=$1", id).
		Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return Comment{}, fmt.Errorf("error in GetComment: %w", err)
	}
	return comment, nil
}
