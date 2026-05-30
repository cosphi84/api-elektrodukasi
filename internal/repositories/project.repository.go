package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"elektrod/internal/models"
)

type ProjectRepository interface {
	FindAll() ([]models.Project, error)
	FindByID(id uuid.UUID) (*models.Project, error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	SoftDelete(id uuid.UUID) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) FindAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.
		Preload("Owner").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) FindByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.
		Preload("Owner").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Model(project).
		Omit("created_at", "deleted_at").
		Save(project).Error
}

func (r *projectRepository) SoftDelete(id uuid.UUID) error {
	result := r.db.Model(&models.Project{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}
	return nil
}
