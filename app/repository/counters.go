package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Counter adalah model untuk collection 'counters'
type Counter struct {
	ID  string `bson:"_id"`
	Seq int64  `bson:"seq"`
}

// GetNextSequenceValue mengambil dan menaikkan nilai counter secara atomik.
// Ini adalah fungsi kunci untuk auto-increment.
func GetNextSequenceValue(ctx context.Context, db *mongo.Database, sequenceName string) (int64, error) {
	collection := db.Collection("counters")

	// Tentukan filter untuk menemukan counter berdasarkan namanya
	filter := bson.M{"_id": sequenceName}

	// Tentukan update untuk menaikkan (increment) 'seq' sebesar 1
	update := bson.M{"$inc": bson.M{"seq": 1}}

	// Opsi untuk mengembalikan dokumen *setelah* di-update (ReturnDocument: After)
	// Upsert: true berarti jika counter belum ada, buat baru.
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)

	var updatedCounter Counter
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCounter)
	if err != nil {
		return 0, err
	}

	return updatedCounter.Seq, nil
}