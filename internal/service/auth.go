package service

import (
	"fmt"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       repository.Authorization
	jwtService *JWTService
}

func NewAuthService(repo repository.Authorization, jwtService *JWTService) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
	}
}

// CreateUser implements the Authorization interface
func (s *AuthService) CreateUser(user models.User) (int, error) {
	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

// Login implements the Authorization interface
func (s *AuthService) Login(username, password string) (AuthResponse, error) {
	// Try to get user by username first
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		// Return the same error for security (don't reveal if user exists)
		return AuthResponse{}, fmt.Errorf("invalid credentials")
	}

	// Compare hashed password with provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return AuthResponse{}, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateJWT(user)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return AuthResponse{Token: token}, nil
}

// Register implements the Authorization interface
func (s *AuthService) Register(user models.User) (AuthResponse, error) {
	// Create the user in the database
	id, err := s.CreateUser(user)
	if err != nil {
		return AuthResponse{}, err
	}

	// Set the ID for JWT generation
	user.ID = id

	// Generate JWT token
	token, err := s.jwtService.GenerateJWT(user)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return AuthResponse{Token: token}, nil
}

// GetUser implements the Authorization interface
func (s *AuthService) GetUser(username, password string) (models.User, error) {
	// This is a legacy method that's kept for backward compatibility
	// It's recommended to use Login instead which returns a proper token
	_, err := s.Login(username, password)
	if err != nil {
		return models.User{}, err
	}
	
	// For backward compatibility, return a minimal user object
	return models.User{
		UserName: username,
	}, nil
}
