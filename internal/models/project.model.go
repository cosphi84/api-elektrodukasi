// models/project.go
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OwnerID   *uuid.UUID      `gorm:"type:uuid"                                       json:"owner_id,omitempty"`
	Owner     *User           `gorm:"foreignKey:OwnerID;constraint:OnDelete:SET NULL" json:"owner,omitempty"`
	Title     string          `gorm:"type:text;not null"                              json:"title"`
	Slug      string          `gorm:"type:text;not null;uniqueIndex"                  json:"slug"`
	Summary   *string         `gorm:"type:text"                                       json:"summary,omitempty"`
	Link      *string         `gorm:"type:text"                                       json:"link,omitempty"`
	Metadata  json.RawMessage `gorm:"type:jsonb"                                      json:"metadata,omitempty"`
	CreatedAt time.Time       `gorm:"type:timestamptz;not null;default:now()"         json:"created_at"`
	UpdatedAt time.Time       `gorm:"type:timestamptz;not null;default:now()"         json:"updated_at"`
}

func (Project) TableName() string {
	return "projects"
}
