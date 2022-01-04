package godbutil

import (
	"fmt"
	"github.com/kordar/goutil"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sync"
)

var instanceOfSqlitePool *SqliteConnPool
var onceOfSqlite sync.Once

// SqliteConnPool /*
// 数据库连接操作库
// 基于gorm封装开发
type SqliteConnPool struct {
	sqliteHandlers map[string]*gorm.DB
}

func GetSqlitePool() *SqliteConnPool {
	onceOfSqlite.Do(func() {
		instanceOfSqlitePool = &SqliteConnPool{sqliteHandlers: make(map[string]*gorm.DB)}
	})
	return instanceOfSqlitePool
}

func (m *SqliteConnPool) getDatabase(db string) (source string) {
	data := goutil.GetSectionValue("sqlite", "data")
	return fmt.Sprintf("%s/%s.db", data, db)
}

// InitDataPool /*
func (m *SqliteConnPool) InitDataPool(db string, db2 ...string) (issucc bool) {
	// 配置日志等级
	dbLogLevel := goutil.GetSystemValue("gorm_log_level")
	config := gorm.Config{}
	if dbLogLevel == "error" {
		config.Logger = logger.Default.LogMode(logger.Error)
	}
	if dbLogLevel == "warn" {
		config.Logger = logger.Default.LogMode(logger.Warn)
	}
	if dbLogLevel == "info" {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	dbs := append(db2, db)
	for _, val := range dbs {
		if m.sqliteHandlers[val] != nil {
			continue
		}
		var err error
		source := m.getDatabase(val)
		m.sqliteHandlers[val], err = gorm.Open(sqlite.Open(source), &config)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	//关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

// Handler /*
// 对外获取数据库连接对象db
func (m *SqliteConnPool) Handler(db string) (conn *gorm.DB) {
	return m.sqliteHandlers[db]
}
