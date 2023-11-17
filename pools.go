package godbutil

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	mysqlpool  *MysqlConnPool
	sqlitepool *SqliteConnPool
	redispool  *RedisConnPool
)

func GetMysqlInstance(db string) *gorm.DB {
	return mysqlpool.Handler(db)
}

func GetSqliteInstance(db string) *gorm.DB {
	return sqlitepool.Handler(db)
}

func GetRedisInstance(db string) *redis.Client {
	return redispool.Handler(db)
}

// ---------------- mysql -----------------------------

// InitMysqlHandle 初始化mysql句柄
func InitMysqlHandle(db string, db2 ...string) {
	mysqlpool = GetMysqlPool()
	mysqlpool.InitDataPool(db, db2...)
}

// AddMysqlInstance 添加mysql句柄
func AddMysqlInstance(db string) error {
	return mysqlpool.Add(db)
}

// RemoveMysqlInstance 移除mysql句柄
func RemoveMysqlInstance(db string) {
	mysqlpool.Remove(db)
}

// --------------- sqlite --------------------------

// InitSqliteHandle 初始化Sqlite句柄
func InitSqliteHandle(db string, db2 ...string) {
	sqlitepool = GetSqlitePool()
	sqlitepool.InitDataPool(db, db2...)
}

// AddSqliteInstance 添加Sqlite句柄
func AddSqliteInstance(db string) error {
	return sqlitepool.Add(db)
}

// RemoveSqliteInstance 移除Sqlite句柄
func RemoveSqliteInstance(db string) {
	sqlitepool.Remove(db)
}

// --------------- redis --------------------------

// InitRedisHandle 初始化redis句柄
func InitRedisHandle(db string, db2 ...string) {
	redispool = GetRedisPool()
	redispool.Add(db)
	for i := range db2 {
		redispool.Add(db2[i])
	}
}

// AddRedisInstance 添加redis句柄
func AddRedisInstance(db string) error {
	return redispool.Add(db)
}

// RemoveRedisInstance 移除redis句柄
func RemoveRedisInstance(db string) {
	redispool.Remove(db)
}
