package repository

import (
	"context"
	"errors"
	"golanjutan/app/model"
	"strconv" // Ditambahkan

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive" // Tidak dipakai lagi
	"go.mongodb.org/mongo-driver/mongo"
)

// IUserRepository mendefinisikan interface
type IUserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type UserRepository struct {
	DB         *mongo.Database // Ditambahkan
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) IUserRepository {
	return &UserRepository{
		DB:         db, // Ditambahkan
		Collection: db.Collection("users"),
	}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	filter := bson.M{"username": username}
	var u model.User
	err := r.Collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah
	var u model.User

	err = r.Collection.FindOne(ctx, filter).Decode(&u) // Perbaikan Anda sebelumnya sudah benar

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	// PERUBAHAN UTAMA: Dapatkan ID baru dari counter
	newID, err := GetNextSequenceValue(ctx, r.DB, "user_id")
	if err != nil {
		return nil, errors.New("gagal mendapatkan sequence id: " + err.Error())
	}

	user.ID = newID // Set ID manual

	_, err = r.Collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}