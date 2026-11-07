package repository

import (
	"errors"

	"storage-manager-api/internal/models"

	"gorm.io/gorm"
)

// UserRepository mendefinisikan kontrak akses data untuk entitas User.
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create menyimpan user baru ke dalam database.
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail mencari user berdasarkan alamat email.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID mencari user berdasarkan UUID.
func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
