// models/article_tag.go
package models

import "github.com/google/uuid"

// ArticleTag represents the junction table — explicit model
// needed if you want to query the pivot table directly.
type ArticleTag struct {
	ArticleID uuid.UUID `gorm:"type:uuid;primaryKey;index"                          json:"article_id"`
	TagID     uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_article_tags_tag"     json:"tag_id"`
}

func (ArticleTag) TableName() string {
	return "article_tags"
}
