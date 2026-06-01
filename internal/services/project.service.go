package services

import (
	"errors"
	"elektrod/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectService interface {
	CreateProject(project *models.Project) error
	GetProjectByID(id uuid.UUID) (*models.Project, error)
	GetProjectBySlug(slug string) (*models.Project, error)
	GetAllProjects(page, pageSize int) ([]models.Project, int64, error)
	GetProjectsByOwner(ownerID uuid.UUID, page, pageSize int) ([]models.Project, int64, error)
	UpdateProject(id uuid.UUID, updates map[string]interface{}) (*models.Project, error)
	DeleteProject(id uuid.UUID) error
}

type projectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) ProjectService {
	return &projectService{db: db}
}

func (s *projectService) CreateProject(project *models.Project) error {
	if project.ID == uuid.Nil {
		project.ID = uuid.New()
	}
	return s.db.Create(project).Error
}

func (s *projectService) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := s.db.Preload("Owner").Where("id = ?", id).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("project not found")
	}
	return &project, err
}

func (s *projectService) GetProjectBySlug(slug string) (*models.Project, error) {
	var project models.Project
	err := s.db.Preload("Owner").Where("slug = ?", slug).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("project not found")
	}
	return &project, err
}

func (s *projectService) GetAllProjects(page, pageSize int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.Project{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Preload("Owner").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&projects).Error
	return projects, total, err
}

func (s *projectService) GetProjectsByOwner(ownerID uuid.UUID, page, pageSize int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.Project{}).Where("owner_id = ?", ownerID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Where("owner_id = ?", ownerID).
		Preload("Owner").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&projects).Error
	return projects, total, err
}

func (s *projectService) UpdateProject(id uuid.UUID, updates map[string]interface{}) (*models.Project, error) {
	project, err := s.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(project).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) DeleteProject(id uuid.UUID) error {
	return s.db.Where("id = ?", id).Delete(&models.Project{}).Error
}
