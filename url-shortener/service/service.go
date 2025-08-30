package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"unicode"
	"url-shortener-1/manager"
)

const ShortLinkLength = 6

type LinksService struct {
	linksManager manager.LinksManager
}

func New(linksManager manager.LinksManager) *LinksService {
	return &LinksService{
		linksManager: linksManager,
	}
}

func (s *LinksService) CreateShortLink(ctx context.Context, longLink string, customShortLink *string) (string, error) {
	existingShortLink, err := s.linksManager.GetShortByLong(ctx, longLink)
	if err == nil {
		return existingShortLink, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("linksManager.GetShortByLong: %w", err) // 500
	}

	// создание по кастомной ссылке
	if customShortLink != nil {
		err := validateShortLink(*customShortLink)
		if err != nil {
			return "", err
		}

		isExists, err := s.linksManager.IsShortExists(ctx, *customShortLink)
		if err != nil {
			return "", fmt.Errorf("linksManager.IsShortExists: %w", err) // 500
		}

		if isExists {
			return "", ErrorLinkAlreadyExists
		}

		err = s.linksManager.CreateLink(ctx, longLink, *customShortLink)
		if err != nil {
			return "", fmt.Errorf("linksManager.CreateLink: %w", err) // 500
		}

		return *customShortLink, nil
	}

	shortLink := ""
	for {
		b := make([]byte, ShortLinkLength)
		rand.Read(b)
		shortLink = base64.URLEncoding.EncodeToString(b)[:ShortLinkLength]

		isExists, err := s.linksManager.IsShortExists(ctx, shortLink)
		if err != nil {
			return "", fmt.Errorf("linksManager.IsShortExists: %w", err) // 500
		}

		if !isExists {
			break
		}
	}

	err = s.linksManager.CreateLink(ctx, longLink, shortLink)
	if err != nil {
		return "", fmt.Errorf("linksManager.CreateLink: %w", err) // 500
	}

	return shortLink, nil
}

func validateShortLink(link string) error {
	if len(link) < ShortLinkLength {
		return ErrorLinkTooShort
	}

	for _, r := range link {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrorInvalidSymbolInLink
		}
	}

	return nil
}
