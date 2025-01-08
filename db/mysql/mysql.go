package mysql

import (
	"context"
	"fmt"

	"github.com/KalinduGandara/crm-system/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysqlClient struct {
	db *gorm.DB
}

type mysqlDatabase struct {
	db *gorm.DB
}

type mysqlCollection struct {
	db        *gorm.DB
	tableName string
}

type mysqlSingleResult struct {
	result *gorm.DB
}

type mysqlCursor struct {
	rows *gorm.DB
}

func NewClient(user, password, host string, port int, dbName string) (db.Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return &mysqlClient{db: db}, nil
}

func (mc *mysqlClient) Ping(ctx context.Context) error {
	sqlDB, err := mc.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (mc *mysqlClient) Database(dbName string) db.Database {
	return &mysqlDatabase{db: mc.db}
}

func (mc *mysqlClient) UseSession(ctx context.Context, fn func(interface{}) error) error {
	return fn(mc.db)
}

func (mc *mysqlClient) StartSession() (interface{}, error) {
	return mc.db, nil
}

func (mc *mysqlClient) Connect(ctx context.Context) error {
	// GORM handles connection pooling automatically
	return nil
}

func (mc *mysqlClient) Disconnect(ctx context.Context) error {
	sqlDB, err := mc.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (md *mysqlDatabase) Collection(tableName string) db.Collection {
	return &mysqlCollection{db: md.db, tableName: tableName}
}

func (md *mysqlDatabase) Client() db.Client {
	return &mysqlClient{db: md.db}
}

func (mc *mysqlCollection) FindOne(ctx context.Context, filter interface{}) db.SingleResult {
	result := mc.db.Table(mc.tableName).Where(filter).First(filter)
	return &mysqlSingleResult{result: result}
}

func (mc *mysqlCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	result := mc.db.Table(mc.tableName).Create(document)
	return result.RowsAffected, result.Error
}

func (mc *mysqlCollection) InsertMany(ctx context.Context, documents []interface{}) ([]interface{}, error) {
	result := mc.db.Table(mc.tableName).Create(documents)
	return nil, result.Error
}

func (mc *mysqlCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	result := mc.db.Table(mc.tableName).Where(filter).Delete(filter)
	return result.RowsAffected, result.Error
}

func (mc *mysqlCollection) Find(ctx context.Context, filter interface{}, opts ...interface{}) (db.Cursor, error) {
	rows := mc.db.Table(mc.tableName).Where(filter).Find(filter)
	return &mysqlCursor{rows: rows}, rows.Error
}

func (mc *mysqlCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...interface{}) (int64, error) {
	var count int64
	result := mc.db.Table(mc.tableName).Where(filter).Count(&count)
	return count, result.Error
}

func (mc *mysqlCollection) Aggregate(ctx context.Context, pipeline interface{}) (db.Cursor, error) {
	// MySQL does not support aggregation pipelines like MongoDB
	return nil, fmt.Errorf("aggregate not supported in MySQL")
}

func (mc *mysqlCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (interface{}, error) {
	result := mc.db.Table(mc.tableName).Where(filter).Updates(update)
	return result.RowsAffected, result.Error
}

func (mc *mysqlCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (interface{}, error) {
	result := mc.db.Table(mc.tableName).Where(filter).Updates(update)
	return result.RowsAffected, result.Error
}

func (sr *mysqlSingleResult) Decode(v interface{}) error {
	return sr.result.Scan(v).Error
}

func (mr *mysqlCursor) Close(ctx context.Context) error {
	// GORM handles closing automatically
	return nil
}

func (mr *mysqlCursor) Next(ctx context.Context) bool {
	// GORM handles iteration automatically
	return false
}

func (mr *mysqlCursor) Decode(v interface{}) error {
	return mr.rows.Scan(v).Error
}

func (mr *mysqlCursor) All(ctx context.Context, result interface{}) error {
	return mr.rows.Find(result).Error
}
