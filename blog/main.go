package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"blog/internal/cache"
	"blog/internal/handler"
	"blog/internal/manager"
	"blog/internal/repo"
	"blog/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
)

const cacheLinksInterval = time.Hour

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

	linksRepository := repo.New(conn)
	linksCache := cache.New(rdb)
	linksManager := manager.New(*linksCache, *linksRepository)
	linksService := service.New(*linksManager)
	linksHandler := handler.New(*linksManager, *linksService)

	go func() { // TODO: вынести из main.go в другое место
		err := cachePopularLinks(linksRepository, linksCache)
		if err != nil {
			log.Println(err)
		}

		c := time.Tick(cacheLinksInterval)
		for range c {
			err := cachePopularLinks(linksRepository, linksCache)
			if err != nil {
				log.Println(err)
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

func cachePopularLinks(linksRepository *repo.Repository, linksCache *cache.LinksCache) error { // TODO: вынести из main.go в другое место
	links, err := linksRepository.GetPopularLinks(context.Background(), 10)
	if err != nil {
		return fmt.Errorf("error updateCache GetPopularLinks: %w", err)
	}

	for _, link := range links {
		err := linksCache.StoreLink(link.Short, link.Long)
		if err != nil {
			return fmt.Errorf("error updateCache StoreLink: %w", err)
		}
	}

	return nil
}
