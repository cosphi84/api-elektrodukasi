package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"elektrod/internal/models"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	SoftDelete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.
		Where("deleted_at IS NULL").
		Order("name ASC").
		Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail is needed for login/auth — not in the interface contract
// but essential for user lookup. Checks active users only.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.
		Where("email = ? AND deleted_at IS NULL AND is_active = true", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Model(user).
		Omit("created_at", "deleted_at", "password"). // never overwrite password via Update
		Save(user).Error
}

func (r *userRepository) SoftDelete(id uuid.UUID) error {
	result := r.db.Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"is_active":  false, // deactivate on delete — belt + suspenders
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
