package repo

import (
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type StoreRedirectParams struct {
	UserAgent string
	LongLink  string
	ShortLink string
}

type Redirect struct {
	Id        int
	LongLink  string
	ShortLink string
	UserAgent string
	CreatedAt time.Time
}

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateLink(c *gin.Context, longLink string, shortLink string) error {
	_, err := r.db.Exec(c, "insert into links (long_link, short_link) VALUES ($1, $2)", longLink, shortLink)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLongByShort(c *gin.Context, shortLink string) (string, error) {
	start := time.Now()
	var longLink string
	err := r.db.QueryRow(c, "SELECT long_link FROM links WHERE short_link=$1", shortLink).Scan(&longLink)
	if err != nil {
		return "", err
	}
	log.Println("db access: ", time.Since(start).Nanoseconds())

	return longLink, err
}

func (r *Repository) GetShortByLong(c *gin.Context, longLink string) (string, error) {
	var shortLink string
	err := r.db.QueryRow(c, "SELECT short_link FROM links WHERE long_link=$1", longLink).Scan(&shortLink)
	if err != nil {
		return "", err
	}

	return shortLink, err
}

func (r *Repository) IsShortExists(c *gin.Context, shortLink string) (bool, error) {
	existingShortLink := ""
	err := r.db.QueryRow(c, "select short_link from links where short_link=$1", shortLink).Scan(&existingShortLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repository) StoreRedirect(c *gin.Context, params StoreRedirectParams) error {
	_, err := r.db.Exec(c, "insert into redirects (user_agent, short_link, long_link) values ($1, $2, $3)",
		params.UserAgent,
		params.ShortLink,
		params.LongLink,
	)
	return err
}

func (r *Repository) GetRedirectsByShortLink(c *gin.Context, shortLink string) ([]Redirect, error) {
	rows, err := r.db.Query(c, "select id, short_link, long_link, user_agent, created_at from redirects where short_link = $1", shortLink)
	if err != nil {
		return nil, err
	}

	res := make([]Redirect, 0)

	for rows.Next() {
		var redirect Redirect
		err := rows.Scan(&redirect.Id, &redirect.ShortLink, &redirect.LongLink, &redirect.UserAgent, &redirect.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, redirect)
	}

	return res, nil
}
