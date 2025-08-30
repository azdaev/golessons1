package manager

import (
	"context"
	"errors"
	"fmt"
	"url-shortener-1/cache"
	"url-shortener-1/model"
	"url-shortener-1/repo"

	"github.com/jackc/pgx/v5"
)

type LinksManager struct {
	cache cache.LinksCache
	repo  repo.Repository
}

func New(cache cache.LinksCache, repo repo.Repository) *LinksManager {
	return &LinksManager{
		cache: cache,
		repo:  repo,
	}
}

func (m *LinksManager) GetLongByShort(ctx context.Context, shortLink string) (string, error) {
	// сначала ищем в кэше
	longLink, err := m.cache.GetLink(shortLink)
	if err != nil {
		return "", fmt.Errorf("error LinksCache.GetLink: %w", err)
	}

	// если нашли, возвращаем
	if longLink != "" {
		return longLink, nil
	}

	// идем в бд
	longLink, err = m.repo.GetLongByShort(ctx, shortLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("error link not found")
		}
		return "", fmt.Errorf("error repo.GetLongByShort: %w", err)
	}

	return longLink, nil
}

func (m *LinksManager) IsShortExists(ctx context.Context, shortLink string) (bool, error) {
	return m.repo.IsShortExists(ctx, shortLink)
}

func (m *LinksManager) CreateLink(ctx context.Context, longLink string, shortLink string) error {
	return m.repo.CreateLink(ctx, longLink, shortLink)
}

func (m *LinksManager) StoreRedirect(ctx context.Context, params model.StoreRedirectParams) error {
	return m.repo.StoreRedirect(ctx, params)
}

func (m *LinksManager) GetRedirectsByShortLink(ctx context.Context, shortLink string) ([]model.Redirect, error) {
	return m.repo.GetRedirectsByShortLink(ctx, shortLink)
}

func (m *LinksManager) GetShortByLong(ctx context.Context, longLink string) (string, error) {
	return m.repo.GetShortByLong(ctx, longLink)
}

func (m *LinksManager) GetPopularLinks(ctx context.Context, n int) ([]model.LinkPair, error) {
	return m.repo.GetPopularLinks(ctx, n)
}
