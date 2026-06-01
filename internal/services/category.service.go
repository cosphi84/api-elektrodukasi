package services

import (
	"errors"
	"elektrod/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	GetAllCategories(page, pageSize int) ([]models.Category, int64, error)
	UpdateCategory(id uuid.UUID, updates map[string]interface{}) (*models.Category, error)
	DeleteCategory(id uuid.UUID) error
}

type categoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{db: db}
}

func (s *categoryService) CreateCategory(category *models.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return s.db.Create(category).Error
}

func (s *categoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := s.db.Where("id = ?", id).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("category not found")
	}
	return &category, err
}

func (s *categoryService) GetAllCategories(page, pageSize int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.Category{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&categories).Error
	return categories, total, err
}

func (s *categoryService) UpdateCategory(id uuid.UUID, updates map[string]interface{}) (*models.Category, error) {
	category, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(category).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	return s.db.Where("id = ?", id).Delete(&models.Category{}).Error
}
