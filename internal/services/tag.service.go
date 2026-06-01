package services

import (
	"errors"
	"elektrod/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagService interface {
	CreateTag(tag *models.Tag) error
	GetTagByID(id uuid.UUID) (*models.Tag, error)
	GetAllTags(page, pageSize int) ([]models.Tag, int64, error)
	UpdateTag(id uuid.UUID, updates map[string]interface{}) (*models.Tag, error)
	DeleteTag(id uuid.UUID) error
}

type tagService struct {
	db *gorm.DB
}

func NewTagService(db *gorm.DB) TagService {
	return &tagService{db: db}
}

func (s *tagService) CreateTag(tag *models.Tag) error {
	if tag.ID == uuid.Nil {
		tag.ID = uuid.New()
	}
	return s.db.Create(tag).Error
}

func (s *tagService) GetTagByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := s.db.Where("id = ?", id).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tag not found")
	}
	return &tag, err
}

func (s *tagService) GetAllTags(page, pageSize int) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.Tag{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tags).Error
	return tags, total, err
}

func (s *tagService) UpdateTag(id uuid.UUID, updates map[string]interface{}) (*models.Tag, error) {
	tag, err := s.GetTagByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(tag).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *tagService) DeleteTag(id uuid.UUID) error {
	return s.db.Where("id = ?", id).Delete(&models.Tag{}).Error
}
