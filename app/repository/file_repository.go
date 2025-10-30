package repository

import (
	"context"
	"errors"
	"fmt"
	"golanjutan/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DIHAPUS: Struct Counter dihapus dari sini karena sudah ada di counters.go
// type Counter struct {
//  ID  string `bson:"_id"`
//  Seq int64  `bson:"seq"`
// }

type FileRepository interface {
	Create(file *model.File) error
	FindAll() ([]model.File, error)
	FindByID(id int64) (*model.File, error) // DIUBAH: id string -> int64
	Delete(id int64) error                    // DIUBAH: id string -> int64
}

type fileRepository struct {
	collection *mongo.Collection
	db         *mongo.Database // Ditambahkan untuk counter
}

func NewFileRepository(db *mongo.Database) FileRepository {
	return &fileRepository{
		collection: db.Collection("files"),
		db:         db, //
	}
}

// Fungsi ini akan menggunakan struct 'Counter' dari file counters.go
func (r *fileRepository) GetNextSequence(sequenceName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var counter Counter // Ini akan merujuk ke Counter dari counters.go

	fmt.Println("[DEBUG] 2a-i. Siap memanggil FindOneAndUpdate ke DB...") // TAMBAHAN

	err := r.db.Collection("counters").FindOneAndUpdate(
		ctx,
		bson.M{"_id": sequenceName},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true),
	).Decode(&counter)

	if err != nil {
		fmt.Println("[DEBUG] 2a-ii. GAGAL FindOneAndUpdate:", err) // TAMBAHAN
		return 0, err
	}
	return counter.Seq, nil
}

func (r *fileRepository) Create(file *model.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Dapatkan ID auto-increment
	fmt.Println("[DEBUG] 2a. Memanggil GetNextSequence...") // TAMBAHAN
	nextID, err := r.GetNextSequence("file_id")
	if err != nil {
		fmt.Println("[DEBUG] GAGAL saat GetNextSequence:", err) // TAMBAHAN
		return errors.New("failed to get next sequence for file: " + err.Error())
	}
	fmt.Println("[DEBUG] 2b. Berhasil GetNextSequence, ID:", nextID) // TAMBAHAN

	// 2. Set ID dan waktu
	file.ID = nextID
	file.UploadedAt = time.Now()

	// 3. Insert
	fmt.Println("[DEBUG] 2c. Siap memanggil InsertOne ke DB...") // TAMBAHAN
	_, err = r.collection.InsertOne(ctx, file)
	if err != nil {
		fmt.Println("[DEBUG] GAGAL saat InsertOne:", err) // TAMBAHAN
		return err
	}
	return nil
}

func (r *fileRepository) FindAll() ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var files []model.File
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) FindByID(id int64) (*model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var file model.File
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}