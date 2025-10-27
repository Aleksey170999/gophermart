package service

import (
	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUser(username, password string) (models.User, error) {
	return s.repo.GetUser(username, password)
}
