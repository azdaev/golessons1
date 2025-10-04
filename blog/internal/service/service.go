package service

import (
	"blog/internal/model"
	"blog/internal/repo"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Service struct {
	repo *repo.Repository
}

func New(repo *repo.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, user CreateUserReq) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error bcrypt.GenerateFromPassword: %w", err)
	}

	err = s.repo.CreateUser(ctx, repo.CreateUser{
		Name:           user.Name,
		HashedPassword: string(hashedPassword),
		Email:          user.Email,
	})
	if err != nil {
		return fmt.Errorf("error repo.CreateUser: %w", err)
	}

	return nil
}

func (s *Service) GetUser(ctx context.Context, id int) (model.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("error repo.GetUser: %w", err)
	}

	return user, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error repo.GetUsers: %w", err)
	}

	return users, nil
}

func (b *Service) GetPosts(ctx context.Context) ([]model.Post, error) {
	posts, err := b.repo.GetPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetPosts: %w", err)
	}

	return posts, nil
}

func (b *Service) GetPostsByUserID(ctx context.Context, userId int) ([]model.Post, error) {
	posts, err := b.repo.GetPostsByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error repo GetPostsByUserID: %w", err)
	}

	return posts, nil
}
