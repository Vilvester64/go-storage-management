package repository

import (
	"errors"

	"storage-manager-api/internal/models"

	"gorm.io/gorm"
)

// FileRepository mendefinisikan kontrak akses database untuk entitas File.
type FileRepository interface {
	Create(file *models.File) error
	FindByUserID(userID string) ([]models.File, error)
	FindByID(id string) (*models.File, error)
	SoftDelete(file *models.File) error
}

type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository membuat instance baru FileRepository.
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

// Create menyimpan metadata file baru ke database.
func (r *fileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// FindByUserID mengambil semua file aktif milik pengguna tertentu (tidak termasuk yang ada di Trash).
func (r *fileRepository) FindByUserID(userID string) ([]models.File, error) {
	var files []models.File
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// FindByID mencari satu file aktif berdasarkan UUID.
func (r *fileRepository) FindByID(id string) (*models.File, error) {
	var file models.File
	if err := r.db.Where("id = ?", id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// SoftDelete memindahkan file ke Trash Bin dengan memperbarui kolom deleted_at via GORM.
func (r *fileRepository) SoftDelete(file *models.File) error {
	return r.db.Delete(file).Error
}
