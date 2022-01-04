package godbutil

import (
	"fmt"
	"github.com/kordar/goutil"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sync"
)

import (
	"gorm.io/driver/mysql"
)

var instanceOfMysqlPool *MysqlConnPool
var onceOfMysql sync.Once

// MysqlConnPool /*
// 数据库连接操作库
// 基于gorm封装开发
type MysqlConnPool struct {
	mysqlHandlers map[string]*gorm.DB
}

func GetMysqlPool() *MysqlConnPool {
	onceOfMysql.Do(func() {
		instanceOfMysqlPool = &MysqlConnPool{mysqlHandlers: make(map[string]*gorm.DB)}
	})
	return instanceOfMysqlPool
}

func (m *MysqlConnPool) getDatabase(db string) (source string) {
	user := goutil.GetSectionValue(db, "user")
	password := goutil.GetSectionValue(db, "password")
	host := goutil.GetSectionValue(db, "host")
	port := goutil.GetSectionValue(db, "port")
	database := goutil.GetSectionValue(db, "db")
	charset := goutil.GetSectionValue(db, "charset")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, database, "charset="+charset+"&parseTime=true")
}

// InitDataPool /*
// 初始化数据库连接(可在mail()适当位置调用)
func (m *MysqlConnPool) InitDataPool(db string, db2 ...string) (issucc bool) {
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
		if m.mysqlHandlers[val] != nil {
			continue
		}
		var err error
		source := m.getDatabase(val)
		m.mysqlHandlers[val], err = gorm.Open(mysql.Open(source), &config)
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
func (m *MysqlConnPool) Handler(db string) (conn *gorm.DB) {
	return m.mysqlHandlers[db]
}
