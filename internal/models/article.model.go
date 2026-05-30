// models/article.go
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"  json:"id"`
	AuthorID    *uuid.UUID      `gorm:"type:uuid"                                       json:"author_id,omitempty"`
	CategoryID  *uuid.UUID      `gorm:"type:uuid;index:idx_articles_category"           json:"category_id,omitempty"`
	Title       string          `gorm:"type:text;not null;index:idx_articles_title"     json:"title"`
	Slug        string          `gorm:"type:text;not null;uniqueIndex"                  json:"slug"`
	Summary     *string         `gorm:"type:text"                                       json:"summary,omitempty"`
	ContentHTML *string         `gorm:"type:text;column:content_html"                   json:"content_html,omitempty"`
	ContentJSON json.RawMessage `gorm:"type:jsonb;column:content_json"                  json:"content_json,omitempty"`
	Image       *string         `gorm:"type:text"                                       json:"image,omitempty"`
	Published   bool            `gorm:"not null;default:false"                          json:"published"`
	PublishedAt *time.Time      `gorm:"index:idx_articles_published_at"                 json:"published_at,omitempty"`
	ViewCount   int64           `gorm:"not null;default:0"                              json:"view_count"`
	CreatedAt   time.Time       `gorm:"not null;default:now()"                          json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"default:null"                                    json:"updated_at,omitempty"`
	DeletedAt   *time.Time      `gorm:"default:null"                                    json:"deleted_at,omitempty"`

	// Relations
	Author   *User     `gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL"   json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	Tags     []Tag     `gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:ArticleID;References:ID;joinReferences:TagID" json:"tags,omitempty"`
}

func (Article) TableName() string {
	return "articles"
}
