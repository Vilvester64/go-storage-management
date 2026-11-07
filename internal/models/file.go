package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File merepresentasikan tabel `files` di database dengan dukungan Soft Delete (Trash Bin).
type File struct {
	ID         string         `gorm:"type:char(36);primaryKey" json:"id"`
	UserID     string         `gorm:"type:char(36);index;not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	SystemName string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"system_name"`
	Size       int64          `gorm:"type:bigint;not null" json:"size"`
	MimeType   string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	Path       string         `gorm:"type:varchar(255);not null" json:"path"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName menentukan nama tabel secara eksplisit.
func (File) TableName() string {
	return "files"
}

// BeforeCreate menghasilkan UUID otomatis sebelum metadata file disimpan ke database.
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}
