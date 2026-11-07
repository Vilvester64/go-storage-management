package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"storage-manager-api/internal/middleware"
	"storage-manager-api/internal/services"
	"storage-manager-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// FileHandler menangani HTTP request untuk manajemen file.
type FileHandler struct {
	fileService services.FileService
}

// NewFileHandler membuat instance baru FileHandler.
func NewFileHandler(fileService services.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

// UploadFile menangani endpoint POST /api/v1/files/upload.
func (h *FileHandler) UploadFile(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File wajib disertakan pada form-data dengan key 'file'")
		return
	}

	file, err := h.fileService.UploadFile(userID, fileHeader)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrEmptyFile), errors.Is(err, utils.ErrUnsupportedFileType):
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrFileTooLarge):
			utils.ErrorResponse(c, http.StatusRequestEntityTooLarge, err.Error())
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengunggah file")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "File uploaded successfully", gin.H{
		"id":        file.ID,
		"name":      file.Name,
		"size":      file.Size,
		"mime_type": file.MimeType,
	})
}

// ListFiles menangani endpoint GET /api/v1/files.
func (h *FileHandler) ListFiles(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}

	files, err := h.fileService.ListUserFiles(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar file")
		return
	}

	result := make([]gin.H, 0, len(files))
	for _, f := range files {
		result = append(result, gin.H{
			"id":         f.ID,
			"name":       f.Name,
			"size":       f.Size,
			"mime_type":  f.MimeType,
			"created_at": f.CreatedAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar file berhasil diambil", result)
}

// DownloadFile menangani endpoint GET /api/v1/files/download/:id menggunakan streaming Gin c.File().
func (h *FileHandler) DownloadFile(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	fileID := c.Param("id")

	file, err := h.fileService.GetFileForDownload(fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFileNotFound):
			utils.ErrorResponse(c, http.StatusNotFound, "File tidak ditemukan")
		case errors.Is(err, services.ErrForbiddenFileAccess):
			utils.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki izin untuk mengunduh file ini")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengunduh file")
		}
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Name))
	c.Header("Content-Type", file.MimeType)
	c.File(file.Path)
}

// DeleteFile menangani endpoint DELETE /api/v1/files/:id (Soft Delete ke Trash Bin).
func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	fileID := c.Param("id")

	err := h.fileService.DeleteFile(fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFileNotFound):
			utils.ErrorResponse(c, http.StatusNotFound, "File tidak ditemukan")
		case errors.Is(err, services.ErrForbiddenFileAccess):
			utils.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki izin untuk menghapus file ini")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus file")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "File moved to trash successfully", nil)
}
