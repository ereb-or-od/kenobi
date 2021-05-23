package interfaces

import "context"

type PostgreSqlDatabaseProvider interface {
	Insert(ctx context.Context, entity interface{}) error
	BulkInsert(ctx context.Context, entities []interface{}) error
	FindOneById(ctx context.Context, id string, entity interface{}) error
	FindOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error
	FindByFilter(ctx context.Context, entity interface{}, query string, params ...interface{}) error
	UpdateOneById(ctx context.Context, id string, entity interface{}) error
	UpdateOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error
	DeleteOneById(ctx context.Context, id string, entity interface{}) error
	DeleteOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error
	DeleteAllByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error
}


