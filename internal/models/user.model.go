// models/user.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string     `gorm:"type:text;not null;index:idx_user_name"         json:"name"`
	Email         string     `gorm:"type:text;not null;uniqueIndex:idx_user_email"  json:"email"`
	Password      string     `gorm:"type:text;not null"                             json:"-"`
	Avatar        *string    `gorm:"type:text"                                      json:"avatar,omitempty"`
	IsActive      bool       `gorm:"not null;default:true"                          json:"is_active"`
	Role          string     `gorm:"type:text;not null;default:'user'"              json:"role"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"                         json:"created_at"`
	UpdatedAt     *time.Time `gorm:"default:null"                                   json:"updated_at,omitempty"`
	DeletedAt     *time.Time `gorm:"default:null"                                   json:"deleted_at,omitempty"`
	LastLogin     *time.Time `gorm:"default:null"                                   json:"last_login,omitempty"`
	LastLoginFrom *string    `gorm:"type:text"                                      json:"last_login_from,omitempty"`
}

func (User) TableName() string {
	return "users"
}
