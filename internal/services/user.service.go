package services

import (
	"errors"
	"elektrod/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers(page, pageSize int) ([]models.User, int64, error)
	UpdateUser(id uuid.UUID, updates map[string]interface{}) (*models.User, error)
	DeleteUser(id uuid.UUID) error
	DeactivateUser(id uuid.UUID) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) CreateUser(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return s.db.Create(user).Error
}

func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (s *userService) GetAllUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (s *userService) UpdateUser(id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(user).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (s *userService) DeactivateUser(id uuid.UUID) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", false).Error
}
