// models/comment.go
package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	// --- Primary key: BIGSERIAL → int64, NOT uuid ---
	ID int64 `gorm:"primaryKey;autoIncrement"  json:"id"          nestedset:"id"`

	// --- Nested Set fields (required by library) ---
	ParentID      sql.NullInt64 `gorm:"default:null"   json:"parent_id,omitempty" nestedset:"parent_id"`
	Lft           int           `gorm:"not null;default:0"                          json:"lft"           nestedset:"lft"`
	Rgt           int           `gorm:"not null;default:0"                          json:"rgt"           nestedset:"rgt"`
	Depth         int           `gorm:"not null;default:0"                          json:"depth"         nestedset:"depth"`
	ChildrenCount int           `gorm:"not null;default:0"                          json:"children_count" nestedset:"children_count"`

	// --- Scope: isolate tree per article ---
	ArticleID *uuid.UUID `gorm:"type:uuid;index:idx_comments_article" json:"article_id,omitempty" nestedset:"scope"`

	// --- Business fields ---
	UserID    *uuid.UUID `gorm:"type:uuid"      json:"user_id,omitempty"`
	Content   string     `gorm:"type:text;not null"                  json:"content"`
	Approved  bool       `gorm:"not null;default:false"              json:"approved"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// --- Relations ---
	Article *Article `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"  json:"article,omitempty"`
	User    *User    `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"     json:"user,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}
