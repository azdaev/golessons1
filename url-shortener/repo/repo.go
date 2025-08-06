package repo

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

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
	var longLink string
	err := r.db.QueryRow(c, "SELECT long_link FROM links WHERE short_link=$1", shortLink).Scan(&longLink)
	if err != nil {
		return "", err
	}

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
