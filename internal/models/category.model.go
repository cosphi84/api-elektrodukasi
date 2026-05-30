// models/category.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string     `gorm:"type:text;not null"                             json:"name"`
	Description *string    `gorm:"type:text"                                      json:"description,omitempty"`
	CreatedAt   time.Time  `gorm:"not null;default:now()"                         json:"created_at"`
	UpdatedAt   *time.Time `gorm:"default:null"                                   json:"updated_at,omitempty"`

	// Relations
	Articles []Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
