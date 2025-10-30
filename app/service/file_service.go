package service

import (
	"fmt"
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FileService interface {
	HandleUpload(c *fiber.Ctx) error // Diganti menjadi HandleUpload umum
	GetAllFiles(c *fiber.Ctx) error
	GetFileByID(c *fiber.Ctx) error
	HandleDeleteFile(c *fiber.Ctx) error
}

type fileService struct {
	repo       repository.FileRepository
	uploadPath string
}

func NewFileService(repo repository.FileRepository, uploadPath string) FileService {
	// Semua metode di interface FileService diimplementasikan oleh *fileService.
	return &fileService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

// BARU: Fungsi utama untuk upload yang menangani otorisasi dan penentuan target User ID
func (s *fileService) HandleUpload(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, "message": "User context not found or unauthorized",
		})
	}

	// 1. Tentukan target User ID
	var targetUserID int64
	paramID := c.Params("id") // Cek apakah ada ID di URL (rute /upload/:id)

	if paramID != "" {
		// Rute POST /api/files/upload/:id (Admin upload untuk user lain)
		if user.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false, "message": "Access denied. Only admin can upload for other users.",
			})
		}

		parsedID, err := strconv.ParseInt(paramID, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false, "message": "Invalid user ID format in path",
			})
		}
		targetUserID = parsedID
	} else {
		// Rute POST /api/files/upload (User/Admin upload untuk dirinya sendiri)
		targetUserID = user.ID
	}

	// 2. Dapatkan jenis upload dari form
	uploadType := c.FormValue("type")
	if uploadType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Missing required field 'type' (e.g., foto or sertifikat) in form data",
		})
	}

	// 3. Panggil proses upload
	return s.processUpload(c, targetUserID, uploadType)
}

// DIUBAH: Mengimplementasikan batasan berdasarkan uploadType
func (s *fileService) processUpload(c *fiber.Ctx, userIDToSave int64, uploadType string) error {
	var allowedTypes map[string]bool
	var maxSize int64

	switch uploadType {
	case "foto":
		allowedTypes = map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
		}
		maxSize = int64(1 * 1024 * 1024) // 1MB
	case "sertifikat":
		allowedTypes = map[string]bool{
			"application/pdf": true,
		}
		maxSize = int64(2 * 1024 * 1024) // 2MB
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Invalid upload type. Must be 'foto' or 'sertifikat'",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "No file uploaded", "error": err.Error(),
		})
	}

	if fileHeader.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": fmt.Sprintf("File size exceeds %dMB for %s upload", maxSize/1024/1024, uploadType),
		})
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "File type not allowed for " + uploadType,
		})
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadPath, newFileName)

	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to create upload directory", "error": err.Error(),
		})
	}

	// --- LOG DEBUG DIMULAI DARI SINI ---

	fmt.Println("[DEBUG] 1. Validasi file selesai. Siap menyimpan ke disk.") // TAMBAHAN

	if err := c.SaveFile(fileHeader, filePath); err != nil {
		fmt.Println("[DEBUG] GAGAL saat c.SaveFile:", err) // TAMBAHAN
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to save file to disk", "error": err.Error(),
		})
	}

	fmt.Println("[DEBUG] 2. Berhasil simpan file ke disk. Siap memanggil repo.Create.") // TAMBAHAN

	fileModel := &model.File{
		UserID:       userIDToSave,
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
	}

	if err := s.repo.Create(fileModel); err != nil {
		fmt.Println("[DEBUG] GAGAL saat s.repo.Create:", err) // TAMBAHAN
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to save file metadata", "error": err.Error(),
		})
	}

	fmt.Println("[DEBUG] 3. Berhasil s.repo.Create. Mengirim respons.") // TAMBAHAN

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("File uploaded successfully for User ID: %d (Type: %s)", userIDToSave, uploadType),
		"data":    s.toFileResponse(fileModel),
	})
}

// --- Implementasi Metode Interface ---

func (s *fileService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to get files", "error": err.Error(),
		})
	}
	var responses []model.FileResponse
	for _, file := range files {
		responses = append(responses, *s.toFileResponse(&file))
	}
	return c.JSON(fiber.Map{
		"success": true, "data": responses,
	})
}

func (s *fileService) GetFileByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Invalid ID format",
		})
	}

	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "File not found", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "data": s.toFileResponse(file),
	})
}

func (s *fileService) HandleDeleteFile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Invalid ID format",
		})
	}

	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "File not found", "error": err.Error(),
		})
	}
	if err := os.Remove(file.FilePath); err != nil {
		fmt.Println("Warning: Failed to delete file from storage:", err)
	}
	if err := s.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to delete file metadata", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "File deleted successfully",
	})
}

// PERBAIKAN: Fungsi ini sekarang memiliki receiver (s *fileService)
func (s *fileService) toFileResponse(file *model.File) *model.FileResponse {
	return &model.FileResponse{
		ID:           file.ID,
		UserID:       file.UserID,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		UploadedAt:   file.UploadedAt,
	}
}
