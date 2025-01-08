package mongo

import (
	"context"
	"time"

	"github.com/KalinduGandara/crm-system/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	mc *mongo.Cursor
}

type mongoSession struct {
	mongo.Session
}

func NewClient(connection string) (db.Client, error) {

	time.Local = time.UTC
	c, err := mongo.NewClient(options.Client().ApplyURI(connection))

	return &mongoClient{cl: c}, err

}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.cl.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) Database(dbName string) db.Database {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) UseSession(ctx context.Context, fn func(interface{}) error) error {
	return mc.cl.UseSession(ctx, func(sc mongo.SessionContext) error {
		return fn(sc)
	})
}

func (mc *mongoClient) StartSession() (interface{}, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.cl.Connect(ctx)
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}

func (md *mongoDatabase) Collection(colName string) db.Collection {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() db.Client {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) db.SingleResult {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (interface{}, error) {
	var updateOptions []*options.UpdateOptions
	for _, opt := range opts {
		if uo, ok := opt.(*options.UpdateOptions); ok {
			updateOptions = append(updateOptions, uo)
		}
	}
	return mc.coll.UpdateOne(ctx, filter, update, updateOptions...)
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) InsertMany(ctx context.Context, document []interface{}) ([]interface{}, error) {
	res, err := mc.coll.InsertMany(ctx, document)
	return res.InsertedIDs, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...interface{}) (db.Cursor, error) {
	var findOptions []*options.FindOptions
	for _, opt := range opts {
		if fo, ok := opt.(*options.FindOptions); ok {
			findOptions = append(findOptions, fo)
		}
	}
	findResult, err := mc.coll.Find(ctx, filter, findOptions...)
	return &mongoCursor{mc: findResult}, err
}

func (mc *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}) (db.Cursor, error) {
	aggregateResult, err := mc.coll.Aggregate(ctx, pipeline)
	return &mongoCursor{mc: aggregateResult}, err
}

func (mc *mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (interface{}, error) {
	var updateOptions []*options.UpdateOptions
	for _, opt := range opts {
		if uo, ok := opt.(*options.UpdateOptions); ok {
			updateOptions = append(updateOptions, uo)
		}
	}
	return mc.coll.UpdateMany(ctx, filter, update, updateOptions...)
}

func (mc *mongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...interface{}) (int64, error) {
	var countOptions []*options.CountOptions
	for _, opt := range opts {
		if co, ok := opt.(*options.CountOptions); ok {
			countOptions = append(countOptions, co)
		}
	}
	return mc.coll.CountDocuments(ctx, filter, countOptions...)
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (mr *mongoCursor) Close(ctx context.Context) error {
	return mr.mc.Close(ctx)
}

func (mr *mongoCursor) Next(ctx context.Context) bool {
	return mr.mc.Next(ctx)
}

func (mr *mongoCursor) Decode(v interface{}) error {
	return mr.mc.Decode(v)
}

func (mr *mongoCursor) All(ctx context.Context, result interface{}) error {
	return mr.mc.All(ctx, result)
}
