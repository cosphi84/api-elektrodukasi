package services

import (
	"errors"
	"elektrod/internal/auth"
	"elektrod/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (string, *models.User, error)
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) bool
}

type authService struct {
	userService UserService
}

func NewAuthService(userService UserService) AuthService {
	return &authService{userService: userService}
}

// Login authenticates user with email and password, returns JWT token
func (s *authService) Login(email, password string) (string, *models.User, error) {
	if email == "" || password == "" {
		return "", nil, errors.New("email and password are required")
	}

	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid email or password")
		}
		return "", nil, err
	}

	if !user.IsActive {
		return "", nil, errors.New("user account is deactivated")
	}

	if !s.VerifyPassword(user.Password, password) {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// HashPassword hashes a password using bcrypt
func (s *authService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a hashed password with a plain password
func (s *authService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
