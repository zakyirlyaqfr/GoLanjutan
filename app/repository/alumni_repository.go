package repository

import (
	"context"
	"errors"
	"golanjutan/app/model"
	"golanjutan/database"
	"strconv" // Ditambahkan
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive" // Tidak dipakai lagi
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IAlumniRepository mendefinisikan interface untuk alumni repo
type IAlumniRepository interface {
	GetAll(ctx context.Context) ([]model.Alumni, error)
	GetByID(ctx context.Context, id string) (*model.Alumni, error)
	GetByIDIncludeDeleted(ctx context.Context, id string) (*model.Alumni, error)
	GetTrashed(ctx context.Context) ([]model.Alumni, error)
	Create(ctx context.Context, alumni model.Alumni) (*model.Alumni, error)
	Update(ctx context.Context, id string, alumni model.Alumni) error
	GetAllWithFilter(ctx context.Context, limit, offset int, sortBy, sortOrder, search string) ([]model.Alumni, error)
	Count(ctx context.Context, search string) (int, error)
	SoftDelete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	GetAllActive(ctx context.Context) ([]model.Alumni, error)
}

type AlumniRepository struct {
	DB            *mongo.Database // Ditambahkan
	Collection    *mongo.Collection
	PekerjaanColl *mongo.Collection
}

func NewAlumniRepository(db *mongo.Database) IAlumniRepository {
	return &AlumniRepository{
		DB:            db, // Ditambahkan
		Collection:    db.Collection("alumni"),
		PekerjaanColl: db.Collection("pekerjaan_alumni"),
	}
}

func (r *AlumniRepository) GetAll(ctx context.Context) ([]model.Alumni, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	return r.find(ctx, filter, opts)
}

func (r *AlumniRepository) GetByID(ctx context.Context, id string) (*model.Alumni, error) {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID, "deleted_at": bson.M{"$exists": false}} // Diubah: _id adalah intID
	return r.findOne(ctx, filter)
}

func (r *AlumniRepository) GetByIDIncludeDeleted(ctx context.Context, id string) (*model.Alumni, error) {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah: _id adalah intID
	return r.findOne(ctx, filter)
}

func (r *AlumniRepository) GetTrashed(ctx context.Context) ([]model.Alumni, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": true}}
	opts := options.Find().SetSort(bson.M{"deleted_at": -1})
	return r.find(ctx, filter, opts)
}

func (r *AlumniRepository) Create(ctx context.Context, alumni model.Alumni) (*model.Alumni, error) {
	// PERUBAHAN UTAMA: Dapatkan ID baru dari counter
	newID, err := GetNextSequenceValue(ctx, r.DB, "alumni_id")
	if err != nil {
		return nil, errors.New("gagal mendapatkan sequence id: " + err.Error())
	}

	alumni.ID = newID // Set ID manual

	_, err = r.Collection.InsertOne(ctx, alumni)
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

func (r *AlumniRepository) Update(ctx context.Context, id string, alumni model.Alumni) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}
	filter := bson.M{"_id": intID} // Diubah: _id adalah intID
	update := bson.M{"$set": alumni}
	_, err = r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *AlumniRepository) GetAllWithFilter(ctx context.Context, limit, offset int, sortBy, sortOrder, search string) ([]model.Alumni, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama": bson.M{"$regex": search, "$options": "i"}},
			{"jurusan": bson.M{"$regex": search, "$options": "i"}},
			{"nim": bson.M{"$regex": search, "$options": "i"}},
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

func (r *AlumniRepository) Count(ctx context.Context, search string) (int, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama": bson.M{"$regex": search, "$options": "i"}},
			{"jurusan": bson.M{"$regex": search, "$options": "i"}},
			{"nim": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	count, err := r.Collection.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *AlumniRepository) SoftDelete(ctx context.Context, id string) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}

	session, err := database.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Soft delete alumni
		filterAlumni := bson.M{"_id": intID} // Diubah
		update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
		if _, err := r.Collection.UpdateOne(sessCtx, filterAlumni, update); err != nil {
			return nil, err
		}

		// 2. Soft delete pekerjaan terkait
		filterPekerjaan := bson.M{"alumni_id": intID} // Diubah
		if _, err := r.PekerjaanColl.UpdateMany(sessCtx, filterPekerjaan, update); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *AlumniRepository) HardDelete(ctx context.Context, id string) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}

	session, err := database.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Hapus pekerjaan terkait
		filterPekerjaan := bson.M{"alumni_id": intID} // Diubah
		if _, err := r.PekerjaanColl.DeleteMany(sessCtx, filterPekerjaan); err != nil {
			return nil, err
		}

		// 2. Hapus alumni
		filterAlumni := bson.M{"_id": intID} // Diubah
		if _, err := r.Collection.DeleteOne(sessCtx, filterAlumni); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *AlumniRepository) Restore(ctx context.Context, id string) error {
	// Diubah: Konversi string ke int64
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id tidak valid")
	}

	session, err := database.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		update := bson.M{"$unset": bson.M{"deleted_at": ""}}

		// 1. Restore alumni
		filterAlumni := bson.M{"_id": intID} // Diubah
		if _, err := r.Collection.UpdateOne(sessCtx, filterAlumni, update); err != nil {
			return nil, err
		}

		// 2. Restore pekerjaan terkait
		filterPekerjaan := bson.M{"alumni_id": intID} // Diubah
		if _, err := r.PekerjaanColl.UpdateMany(sessCtx, filterPekerjaan, update); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *AlumniRepository) GetAllActive(ctx context.Context) ([]model.Alumni, error) {
	return r.GetAll(ctx)
}

// Helper internal
func (r *AlumniRepository) findOne(ctx context.Context, filter bson.M) (*model.Alumni, error) {
	var a model.Alumni
	err := r.Collection.FindOne(ctx, filter).Decode(&a)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("alumni tidak ditemukan")
		}
		return nil, err
	}
	return &a, nil
}

// Helper internal
func (r *AlumniRepository) find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]model.Alumni, error) {
	cursor, err := r.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Alumni
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}