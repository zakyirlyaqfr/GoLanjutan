package repository

import (
	"context"
	"errors"
	"golanjutan/app/model"
	"strconv" // Ditambahkan
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive" // Tidak dipakai lagi
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IPekerjaanRepository mendefinisikan interface
type IPekerjaanRepository interface {
	GetAll(ctx context.Context) ([]model.PekerjaanAlumni, error)
	GetByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error)
	GetByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error)
	Create(ctx context.Context, p model.PekerjaanAlumni) (*model.PekerjaanAlumni, error)
	Update(ctx context.Context, id string, p model.PekerjaanAlumni) error
	GetAllWithFilter(ctx context.Context, limit, offset int, sortBy, sortOrder, search string) ([]model.PekerjaanAlumni, error)
	GetTrashed(ctx context.Context) ([]model.PekerjaanAlumni, error)
	GetTrashedByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error)
	Count(ctx context.Context, search string) (int, error)
	SoftDelete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	GetByIDIncludeDeleted(ctx context.Context, id string) (*model.PekerjaanAlumni, error)
	GetActive(ctx context.Context) ([]model.PekerjaanAlumni, error)
}

type PekerjaanRepository struct {
	DB         *mongo.Database // Ditambahkan
	Collection *mongo.Collection
}

func NewPekerjaanRepository(db *mongo.Database) IPekerjaanRepository {
	return &PekerjaanRepository{
		DB:         db, // Ditambahkan
		Collection: db.Collection("pekerjaan_alumni"),
	}
}

func (r *PekerjaanRepository) GetAll(ctx context.Context) ([]model.PekerjaanAlumni, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	return r.find(ctx, filter, opts)
}

func (r *PekerjaanRepository) GetByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error) {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID, "deleted_at": bson.M{"$exists": false}} // Diubah
	return r.findOne(ctx, filter)
}

func (r *PekerjaanRepository) GetByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error) {
	// Diubah: Konversi string ke int64
	alumniIntID, err := strconv.ParseInt(alumniID, 10, 64)
	if err != nil {
		return nil, errors.New("alumni id tidak valid")
	}
	filter := bson.M{"alumni_id": alumniIntID, "deleted_at": bson.M{"$exists": false}} // Diubah
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	return r.find(ctx, filter, opts)
}

func (r *PekerjaanRepository) Create(ctx context.Context, p model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	// PERUBAHAN UTAMA: Dapatkan ID baru dari counter
	newID, err := GetNextSequenceValue(ctx, r.DB, "pekerjaan_id")
	if err != nil {
		return nil, errors.New("gagal mendapatkan sequence id: " + err.Error())
	}

	p.ID = newID // Set ID manual

	_, err = r.Collection.InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PekerjaanRepository) Update(ctx context.Context, id string, p model.PekerjaanAlumni) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah
	update := bson.M{"$set": p}
	_, err = r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *PekerjaanRepository) GetAllWithFilter(ctx context.Context, limit, offset int, sortBy, sortOrder, search string) ([]model.PekerjaanAlumni, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	order := 1 // ASC
	if sortOrder == "DESC" {
		order = -1 // DESC
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{sortBy: order})

	return r.find(ctx, filter, opts)
}

func (r *PekerjaanRepository) GetTrashed(ctx context.Context) ([]model.PekerjaanAlumni, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": true}}
	opts := options.Find().SetSort(bson.M{"deleted_at": -1})
	return r.find(ctx, filter, opts)
}

func (r *PekerjaanRepository) GetTrashedByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error) {
	// Diubah: Konversi string ke int64
	alumniIntID, err := strconv.ParseInt(alumniID, 10, 64)
	if err != nil {
		return nil, errors.New("alumni id tidak valid")
	}
	filter := bson.M{"alumni_id": alumniIntID, "deleted_at": bson.M{"$exists": true}} // Diubah
	opts := options.Find().SetSort(bson.M{"deleted_at": -1})
	return r.find(ctx, filter, opts)
}

func (r *PekerjaanRepository) Count(ctx context.Context, search string) (int, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	count, err := r.Collection.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *PekerjaanRepository) SoftDelete(ctx context.Context, id string) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err = r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *PekerjaanRepository) HardDelete(ctx context.Context, id string) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah
	_, err = r.Collection.DeleteOne(ctx, filter)
	return err
}

func (r *PekerjaanRepository) Restore(ctx context.Context, id string) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah
	update := bson.M{"$unset": bson.M{"deleted_at": ""}}
	_, err = r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *PekerjaanRepository) GetActive(ctx context.Context) ([]model.PekerjaanAlumni, error) {
	return r.GetAll(ctx)
}

func (r *PekerjaanRepository) GetByIDIncludeDeleted(ctx context.Context, id string) (*model.PekerjaanAlumni, error) {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah
	return r.findOne(ctx, filter)
}

// Helper internal
func (r *PekerjaanRepository) findOne(ctx context.Context, filter bson.M) (*model.PekerjaanAlumni, error) {
	var p model.PekerjaanAlumni
	err := r.Collection.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("pekerjaan tidak ditemukan")
		}
		return nil, err
	}
	return &p, nil
}

// Helper internal
func (r *PekerjaanRepository) find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]model.PekerjaanAlumni, error) {
	cursor, err := r.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.PekerjaanAlumni
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}