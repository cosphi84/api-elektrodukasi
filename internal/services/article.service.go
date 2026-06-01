package services

import (
	"errors"
	"elektrod/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleService interface {
	CreateArticle(article *models.Article) error
	GetArticleByID(id uuid.UUID) (*models.Article, error)
	GetArticleBySlug(slug string) (*models.Article, error)
	GetAllArticles(page, pageSize int, published *bool) ([]models.Article, int64, error)
	GetArticlesByCategory(categoryID uuid.UUID, page, pageSize int) ([]models.Article, int64, error)
	UpdateArticle(id uuid.UUID, updates map[string]interface{}) (*models.Article, error)
	DeleteArticle(id uuid.UUID) error
	PublishArticle(id uuid.UUID) error
	AddTagToArticle(articleID, tagID uuid.UUID) error
	RemoveTagFromArticle(articleID, tagID uuid.UUID) error
	IncrementViewCount(id uuid.UUID) error
}

type articleService struct {
	db *gorm.DB
}

func NewArticleService(db *gorm.DB) ArticleService {
	return &articleService{db: db}
}

func (s *articleService) CreateArticle(article *models.Article) error {
	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}
	return s.db.Create(article).Error
}

func (s *articleService) GetArticleByID(id uuid.UUID) (*models.Article, error) {
	var article models.Article
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("id = ?", id).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("article not found")
	}
	return &article, err
}

func (s *articleService) GetArticleBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("slug = ?", slug).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("article not found")
	}
	return &article, err
}

func (s *articleService) GetAllArticles(page, pageSize int, published *bool) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := s.db.Model(&models.Article{})

	if published != nil {
		query = query.Where("published = ?", *published)
	}

	offset := (page - 1) * pageSize
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Author").Preload("Category").Preload("Tags").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles).Error
	return articles, total, err
}

func (s *articleService) GetArticlesByCategory(categoryID uuid.UUID, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.Article{}).Where("category_id = ?", categoryID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Where("category_id = ?", categoryID).
		Preload("Author").Preload("Category").Preload("Tags").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles).Error
	return articles, total, err
}

func (s *articleService) UpdateArticle(id uuid.UUID, updates map[string]interface{}) (*models.Article, error) {
	article, err := s.GetArticleByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(article).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (s *articleService) DeleteArticle(id uuid.UUID) error {
	return s.db.Where("id = ?", id).Delete(&models.Article{}).Error
}

func (s *articleService) PublishArticle(id uuid.UUID) error {
	article, err := s.GetArticleByID(id)
	if err != nil {
		return err
	}

	now := article.CreatedAt
	return s.db.Model(article).Updates(map[string]interface{}{
		"published":   true,
		"published_at": now,
	}).Error
}

func (s *articleService) AddTagToArticle(articleID, tagID uuid.UUID) error {
	return s.db.Model(&models.Article{}).Where("id = ?", articleID).
		Association("Tags").Append(&models.Tag{ID: tagID})
}

func (s *articleService) RemoveTagFromArticle(articleID, tagID uuid.UUID) error {
	return s.db.Model(&models.Article{}).Where("id = ?", articleID).
		Association("Tags").Delete(&models.Tag{ID: tagID})
}

func (s *articleService) IncrementViewCount(id uuid.UUID) error {
	return s.db.Model(&models.Article{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}
