package postgresql

import (
	"context"
	"github.com/ereb-or-od/kenobi/pkg/postgresql/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/postgresql/options"
	"github.com/go-pg/pg/v10"
)

type standalonePostgresqlDatabase struct {
	db *pg.DB
}

func (s standalonePostgresqlDatabase) FindByFilter(ctx context.Context, entity interface{}, query string, params ...interface{}) error {
	if result, err := s.db.Query(entity, query, params); err != nil {
		return err
	} else {
		result.RowsAffected()
		return nil
	}
}

func (s standalonePostgresqlDatabase) DeleteOneById(ctx context.Context, id string, entity interface{}) error {
	if _, err := s.db.Model(entity).Where("id = ? ", id).Delete(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s standalonePostgresqlDatabase) DeleteOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error {
	if _, err := s.db.Model(entity).Where(condition, params...).Delete(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s standalonePostgresqlDatabase) DeleteAllByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error {
	if _, err := s.db.Model(entity).Where(condition, params...).Delete(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s standalonePostgresqlDatabase) UpdateOneById(ctx context.Context, id string, entity interface{}) error {
	if _, err := s.db.Model(entity).Where("id = ? ", id).Update(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s standalonePostgresqlDatabase) UpdateOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error {
	if _, err := s.db.Model(entity).Where(condition, params...).Update(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s standalonePostgresqlDatabase) FindOneByFilter(ctx context.Context, entity interface{}, condition string, params ...interface{}) error {
	if err := s.db.Model(entity).Where(condition, params...).Select(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s standalonePostgresqlDatabase) FindOneById(ctx context.Context, id string, entity interface{}) error {
	if err := s.db.Model(entity).Where("id = ?", id).Select(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s standalonePostgresqlDatabase) BulkInsert(ctx context.Context, entities []interface{}) error {
	if _, err := s.db.Model(entities).Insert(); err != nil {
		return err
	}
	return nil
}

func (s standalonePostgresqlDatabase) Insert(ctx context.Context, entity interface{}) error {
	if _, err := s.db.Model(entity).Insert(); err != nil {
		return err
	}
	return nil
}

func New(options *options.PostgreSqlServerOptions) interfaces.PostgreSqlDatabaseProvider {
	db := pg.Connect(&pg.Options{
		User:            options.User,
		Password:        options.Password,
		Addr:            options.Addr,
		Database:        options.Database,
		ApplicationName: options.ApplicationName,
	})
	return &standalonePostgresqlDatabase{
		db: db,
	}
}
