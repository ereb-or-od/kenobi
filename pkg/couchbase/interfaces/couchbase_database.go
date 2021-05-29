package interfaces

import "context"

type CouchbaseDatabase interface {
	Insert(ctx context.Context, id string, entity interface{}) error
	BulkInsert(ctx context.Context, entities map[string]interface{}) error
	FindOneById(ctx context.Context, id string, entity interface{}) error
	FindByOneQuery(ctx context.Context, entity interface{}, query string, params ...interface{}) error
	FindByQuery(ctx context.Context, mapToEntity func(data []byte) error, query string, params ...interface{}) error
	MutateOneById(ctx context.Context, id string, pathAndValues map[string]interface{}) error
	ReplaceOneById(ctx context.Context, id string, entity interface{}) error
	DeleteOneById(ctx context.Context, id string) error
	ExecuteNoneQuery(ctx context.Context, query string, params ...interface{}) error
}
