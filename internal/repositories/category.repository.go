package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"elektrod/internal/models"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	FindByID(id uuid.UUID) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	SoftDelete(id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.
		Where("deleted_at IS NULL").
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Model(category).
		Omit("created_at", "deleted_at").
		Save(category).Error
}

func (r *categoryRepository) SoftDelete(id uuid.UUID) error {
	result := r.db.Model(&models.Category{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
