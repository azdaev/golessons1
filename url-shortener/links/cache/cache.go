package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type LinksCache struct {
	rdb *redis.Client
}

func NewLinksCache(rdb *redis.Client) *LinksCache {
	return &LinksCache{
		rdb: rdb,
	}
}

func (c *LinksCache) StoreLink(shortLink string, longLink string) error {
	cmd := c.rdb.Set(shortLink, longLink, time.Hour)
	if cmd.Err() != nil {
		return fmt.Errorf("error StoreLink: %w", cmd.Err())
	}
	return nil
}

func (c *LinksCache) GetLongLink(shortLink string) (string, error) {
	cmd := c.rdb.Get(shortLink)
	if cmd.Err() != nil {
		return "", fmt.Errorf("error GetLongLink: %w", cmd.Err())
	}

	return cmd.Val(), nil
}
