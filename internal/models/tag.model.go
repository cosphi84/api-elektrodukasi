// models/tag.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:text;not null"                              json:"name"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"         json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"         json:"updated_at"`

	// Relations
	Articles []Article `gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:TagID;References:ID;joinReferences:ArticleID" json:"articles,omitempty"`
}

func (Tag) TableName() string {
	return "tags"
}
