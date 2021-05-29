package couchbase

import (
	"context"
	"github.com/couchbase/gocb/v2"
	"github.com/ereb-or-od/kenobi/pkg/couchbase/interfaces"
	"time"
)

type couchbaseDatabase struct {
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

func (c *couchbaseDatabase) Insert(ctx context.Context, id string, entity interface{}) error {
	_, err := c.collection.Upsert(id, entity, &gocb.UpsertOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *couchbaseDatabase) BulkInsert(ctx context.Context, entities map[string]interface{}) error {
	var items []gocb.BulkOp

	for id, entity := range entities {
		items = append(items, &gocb.InsertOp{ID: id, Value: entity})
		items = append(items, &gocb.InsertOp{ID: id, Value: entity})
	}

	return c.collection.Do(items, &gocb.BulkOpOptions{})
}

func (c *couchbaseDatabase) FindOneById(ctx context.Context, id string, entity interface{}) error {
	getResult, err := c.collection.Get(id, &gocb.GetOptions{})
	if err != nil {
		return err
	}
	if err := getResult.Content(entity); err != nil {
		return err
	}
	return nil
}

func (c *couchbaseDatabase) FindByOneQuery(ctx context.Context, entity interface{}, query string, params ...interface{}) error {
	rows, err := c.cluster.Query(query,
		&gocb.QueryOptions{
			Adhoc:                true,
			PositionalParameters: []interface{}{params},
		},
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return rows.One(entity)
}

func (c *couchbaseDatabase) FindByQuery(ctx context.Context, mapToEntity func(data []byte) error, query string, params ...interface{}) error {
	rows, err := c.cluster.Query(query,
		&gocb.QueryOptions{
			Adhoc:                true,
			PositionalParameters: []interface{}{params},
		},
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rowData := make([]byte, 0)
		err := rows.Row(rowData)
		if err != nil {
			return err
		}
		if err := mapToEntity(rowData); err != nil {
			return err
		}
	}
	return nil
}

func (c *couchbaseDatabase) DeleteOneById(ctx context.Context, id string) error {
	_, err := c.collection.Remove(id, &gocb.RemoveOptions{})
	return err
}

func (c *couchbaseDatabase) MutateOneById(ctx context.Context, id string, pathAndValues map[string]interface{}) error {
	mops := make([]gocb.MutateInSpec, 0)
	for path, value := range pathAndValues {
		mops = append(mops, gocb.UpsertSpec(path, value, &gocb.UpsertSpecOptions{}))
	}
	_, err := c.collection.MutateIn(id, mops, &gocb.MutateInOptions{})
	return err
}

func (c *couchbaseDatabase) ReplaceOneById(ctx context.Context, id string, entity interface{}) error {
	_, err := c.collection.Replace(id, entity, &gocb.ReplaceOptions{})
	return err
}

func (c *couchbaseDatabase) ExecuteNoneQuery(ctx context.Context, query string, params ...interface{}) error {
	_, err := c.cluster.Query(query,
		&gocb.QueryOptions{
			Adhoc:                true,
			PositionalParameters: []interface{}{params},
		},
	)
	return err
}

func New(connectionString string, username string, password string, bucketName string) (interfaces.CouchbaseDatabase, error) {
	cluster, err := gocb.Connect(connectionString,
		gocb.ClusterOptions{
			Username: username,
			Password: password,
		})
	if err != nil {
		return nil, err
	}

	bucket := cluster.Bucket(bucketName)
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return nil, err
	}
	return &couchbaseDatabase{collection: bucket.DefaultCollection(), cluster: cluster}, nil
}
