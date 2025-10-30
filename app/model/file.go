// golanjutan/app/model/file.go
package model

import (
	"time"
	// Hapus "go.mongodb.org/mongo-driver/bson/primitive"
)

// File adalah struktur untuk metadata di MongoDB
type File struct {
	ID           int64     `json:"id" bson:"_id"` // DIUBAH: dari primitive.ObjectID
	UserID       int64     `json:"user_id" bson:"user_id"`
	FileName     string    `json:"file_name" bson:"file_name"`
	OriginalName string    `json:"original_name" bson:"original_name"`
	FilePath     string    `json:"file_path" bson:"file_path"`
	FileSize     int64     `json:"file_size" bson:"file_size"`
	FileType     string    `json:"file_type" bson:"file_type"`
	UploadedAt   time.Time `json:"uploaded_at" bson:"uploaded_at"`
}

// FileResponse adalah struktur untuk respon JSON
type FileResponse struct {
	ID           int64     `json:"id"` // DIUBAH: dari string
	UserID       int64     `json:"user_id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}