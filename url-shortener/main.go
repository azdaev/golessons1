package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"url-shortener-1/cache"
	"url-shortener-1/handler"
	"url-shortener-1/repo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// проверяем соединение
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	connString := "postgres://postgres:postgres@localhost:5432/shortener-db"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	linksRepository := repo.NewRepository(conn)
	linksCache := cache.NewLinksCache(rdb)
	linksHandler := handler.NewHandler(linksRepository, linksCache)

	go func() {
		err := cacheMostPopularLinks(context.Background(), conn, linksCache)
		if err != nil {
			log.Println("error cacheMostPopularLinks: ", err)
		}

		c := time.Tick(1 * time.Hour)
		for range c {
			err := cacheMostPopularLinks(context.Background(), conn, linksCache)
			if err != nil {
				log.Println("error cacheMostPopularLinks: ", err)
			}
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/shorten", linksHandler.CreateLink)
	r.GET("/analytics/:path", linksHandler.GetAnalytics)
	r.GET("/:path", linksHandler.Redirect)
	r.Run(":8080")
}

func cacheMostPopularLinks(ctx context.Context, repo *pgx.Conn, cache *cache.LinksCache) error {
	rows, err := repo.Query(
		ctx,
		`select
		        short_link,
				long_link
		   from redirects
	   group by short_link, long_link
	   order by count(id) desc
		  limit 10;`,
	)
	if err != nil {
		return fmt.Errorf("error cacheMostPopularLinks: %w", err)
	}

	for rows.Next() {
		var (
			shortLink string
			longLink  string
		)
		err := rows.Scan(&shortLink, &longLink)
		if err != nil {
			return fmt.Errorf("error scan in cacheMostPopularLinks: %w", err)
		}

		err = cache.StoreLink(shortLink, longLink)
		if err != nil {
			return fmt.Errorf("error StoreLink in cacheMostPopularLinks: %w", err)
		}
	}

	return nil
}
