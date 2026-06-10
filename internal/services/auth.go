package services

import (
	"errors"
	"strings"

	"task-manager-api/internal/auth"
	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
)

var (
	ErrEmailTaken         = errors.New("an account with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	users *repository.UserRepository
	jwt   *auth.JWTManager
	// adminEmails are promoted to the admin role at signup/login.
	adminEmails map[string]bool
}

func NewAuthService(users *repository.UserRepository, jwt *auth.JWTManager, adminEmails []string) *AuthService {
	emails := make(map[string]bool, len(adminEmails))
	for _, email := range adminEmails {
		email = strings.ToLower(strings.TrimSpace(email))
		if email != "" {
			emails[email] = true
		}
	}
	return &AuthService{users: users, jwt: jwt, adminEmails: emails}
}

func (s *AuthService) Signup(name, email, password string) (*models.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if _, err := s.users.FindByEmail(email); err == nil {
		return nil, "", ErrEmailTaken
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, "", err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	role := models.RoleUser
	if s.adminEmails[email] {
		role = models.RoleAdmin
	}

	user := &models.User{
		Name:         strings.TrimSpace(name),
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	}
	if err := s.users.Create(user); err != nil {
		return nil, "", err
	}

	token, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.users.FindByEmail(email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, "", ErrInvalidCredentials
	}
	if err != nil {
		return nil, "", err
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, "", ErrInvalidCredentials
	}

	// Promote existing accounts that were added to ADMIN_EMAILS later.
	if s.adminEmails[email] && user.Role != models.RoleAdmin {
		user.Role = models.RoleAdmin
		if err := s.users.Update(user); err != nil {
			return nil, "", err
		}
	}

	token, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *AuthService) GetUser(id uint) (*models.User, error) {
	return s.users.FindByID(id)
}
