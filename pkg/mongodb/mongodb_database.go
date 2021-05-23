package mongodb

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/mongodb/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbDatabase struct {
	db *mongo.Collection
}

func (m mongodbDatabase) Insert(ctx context.Context, entity interface{}) error {
	if _, err := m.db.InsertOne(ctx, entity); err != nil {
		return err
	} else {
		return nil
	}
}

func (m mongodbDatabase) BulkInsert(ctx context.Context, entities []interface{}) error {
	if _, err := m.db.InsertMany(ctx, entities); err != nil {
		return err
	} else {
		return nil
	}
}

func (m mongodbDatabase) FindOneById(ctx context.Context, id string, entity interface{}) error {
	objectId, err := utilities.ToObjectID(id)
	if err != nil {
		return err
	}
	result := m.db.FindOne(ctx, bson.M{"_id": objectId})
	if err = result.Decode(&entity); err != nil {
		return err
	}

	return nil
}

func (m mongodbDatabase) FindOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error {
	result := m.db.FindOne(ctx, bson.M{condition: params})
	if err := result.Decode(&entity); err != nil {
		return err
	}
	return nil
}

func (m mongodbDatabase) FindByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error {
	if cursor, err := m.db.Find(ctx, bson.M{condition: params}); err != nil {
		return err
	} else {
		var records []interface{}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var record interface{}
			if err = cursor.Decode(&record); err != nil {
				// ignored
			}
			records = append(records, record)
		}
		entity = records
	}

	return nil
}

func (m mongodbDatabase) UpdateOneById(ctx context.Context, id string, entity interface{}) error {
	objectId, err := utilities.ToObjectID(id)
	if err != nil {
		return err
	}
	if _, err = m.db.UpdateOne(ctx, bson.M{"_id": objectId}, entity); err != nil {
		return err
	} else {
		return nil
	}
}

func (m mongodbDatabase) UpdateOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error {
	if _, err := m.db.UpdateOne(ctx, bson.M{condition: params}, entity); err != nil {
		return err
	} else {
		return nil
	}
}

func (m mongodbDatabase) DeleteOneById(ctx context.Context, id string) error {
	objectId, err := utilities.ToObjectID(id)
	if err != nil {
		return err
	}
	if _, err = m.db.DeleteOne(ctx, bson.M{"_id": objectId}); err != nil {
		return err
	} else {
		return nil
	}
}

func (m mongodbDatabase) DeleteOneByFilter(ctx context.Context, condition string, params ...interface{}) error {
	if _, err := m.db.DeleteOne(ctx, bson.M{condition: params}); err != nil {
		return err
	} else {
		return nil
	}
}

func (m mongodbDatabase) DeleteAllByFilter(ctx context.Context, condition string, params ...interface{}) error {
	if _, err := m.db.DeleteMany(ctx, bson.M{condition: params}); err != nil {
		return err
	} else {
		return nil
	}
}

func New(connectionString string, database string, collection string) (interfaces.MongoDbDatabase, error) {
	if client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString)); err != nil {
		return nil, err
	} else {
		return &mongodbDatabase{
			db: client.Database(database).Collection(collection),
		}, nil
	}
}
