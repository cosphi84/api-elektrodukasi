package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"elektrod/internal/models"
)

type TagRepository interface {
	FindAll() ([]models.Tag, error)
	FindByID(id uuid.UUID) (*models.Tag, error)
	Create(tag *models.Tag) error
	Update(tag *models.Tag) error
	SoftDelete(id uuid.UUID) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.
		Order("name ASC").
		Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.
		Where("id = ?", id).
		First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) Update(tag *models.Tag) error {
	return r.db.Model(tag).
		Omit("created_at").
		Save(tag).Error
}

// Tag has no deleted_at column — hard delete is correct here.
// Tags are reference data; deleting them cascades via article_tags FK.
func (r *tagRepository) SoftDelete(id uuid.UUID) error {
	result := r.db.
		Where("id = ?", id).
		Delete(&models.Tag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTagNotFound
	}
	return nil
}
