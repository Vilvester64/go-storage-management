package services

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"storage-manager-api/internal/config"
	"storage-manager-api/internal/models"
	"storage-manager-api/internal/repository"
	"storage-manager-api/internal/utils"
)

var (
	ErrFileTooLarge        = errors.New("ukuran file melebihi batas maksimal yang diizinkan")
	ErrFileNotFound        = errors.New("file tidak ditemukan")
	ErrForbiddenFileAccess = errors.New("anda tidak memiliki akses ke file ini")
)

// FileService mendefinisikan kontrak logika bisnis untuk manajemen file.
type FileService interface {
	UploadFile(userID string, fileHeader *multipart.FileHeader) (*models.File, error)
	ListUserFiles(userID string) ([]models.File, error)
	GetFileForDownload(fileID, userID string) (*models.File, error)
	DeleteFile(fileID, userID string) error
}

type fileService struct {
	fileRepo repository.FileRepository
	cfg      *config.Config
}

// NewFileService membuat instance baru FileService.
func NewFileService(fileRepo repository.FileRepository, cfg *config.Config) FileService {
	return &fileService{
		fileRepo: fileRepo,
		cfg:      cfg,
	}
}

// UploadFile memvalidasi ukuran & MIME type, menyimpan file fisik ke direktori uploads,
// serta menyimpan metadata ke MySQL.
func (s *fileService) UploadFile(userID string, fileHeader *multipart.FileHeader) (*models.File, error) {
	if fileHeader == nil || fileHeader.Size <= 0 {
		return nil, utils.ErrEmptyFile
	}
	if s.cfg.MaxUploadSize > 0 && fileHeader.Size > s.cfg.MaxUploadSize {
		return nil, ErrFileTooLarge
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Validasi MIME type menggunakan magic bytes (h2non/filetype)
	mimeType, err := utils.DetectAndValidateMIME(src, fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	// Kembalikan pointer baca ke awal file setelah membaca header magic bytes
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// Pastikan direktori penyimpanan fisik tersedia
	if err := os.MkdirAll(s.cfg.StoragePath, os.ModePerm); err != nil {
		return nil, err
	}

	originalName := utils.SanitizeFileName(fileHeader.Filename)
	systemName := utils.GenerateSystemFileName(originalName)
	destPath := filepath.ToSlash(filepath.Join(s.cfg.StoragePath, systemName))

	dst, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}

	written, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil {
		_ = os.Remove(destPath)
		return nil, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(destPath)
		return nil, closeErr
	}

	fileModel := &models.File{
		UserID:     userID,
		Name:       originalName,
		SystemName: systemName,
		Size:       written,
		MimeType:   mimeType,
		Path:       destPath,
	}

	if err := s.fileRepo.Create(fileModel); err != nil {
		// Rollback file fisik jika gagal menyimpan ke database
		_ = os.Remove(destPath)
		return nil, err
	}

	return fileModel, nil
}

// ListUserFiles mengembalikan daftar seluruh file aktif milik pengguna.
func (s *fileService) ListUserFiles(userID string) ([]models.File, error) {
	return s.fileRepo.FindByUserID(userID)
}

// GetFileForDownload mencari file, memvalidasi kepemilikan, dan memastikan file fisik tersedia di storage.
func (s *fileService) GetFileForDownload(fileID, userID string) (*models.File, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}
	if file.UserID != userID {
		return nil, ErrForbiddenFileAccess
	}

	if _, err := os.Stat(file.Path); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return file, nil
}

// DeleteFile memvalidasi kepemilikan lalu melakukan Soft Delete (memindahkan file ke Trash).
func (s *fileService) DeleteFile(fileID, userID string) error {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return err
	}
	if file == nil {
		return ErrFileNotFound
	}
	if file.UserID != userID {
		return ErrForbiddenFileAccess
	}

	return s.fileRepo.SoftDelete(file)
}
