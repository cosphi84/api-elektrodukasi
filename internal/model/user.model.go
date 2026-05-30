package model

type User struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Email string `gorm:"index" json:"email"`
}
