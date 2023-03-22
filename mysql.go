package godbutil

import (
	"errors"
	"fmt"
	"github.com/kordar/goutil"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sync"
)

var instanceOfMysqlPool *MysqlConnPool
var onceOfMysql sync.Once
var mysqlConfig gorm.Config

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
	mysqlConfig = gorm.Config{}
	if dbLogLevel == "error" {
		mysqlConfig.Logger = logger.Default.LogMode(logger.Error)
	}
	if dbLogLevel == "warn" {
		mysqlConfig.Logger = logger.Default.LogMode(logger.Warn)
	}
	if dbLogLevel == "info" {
		mysqlConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	dbs := append(db2, db)
	for _, val := range dbs {
		if m.mysqlHandlers[val] != nil {
			continue
		}
		var err error
		err = m.Add(val)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	//关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

// Add 添加数据库实例
func (m *MysqlConnPool) Add(db string) error {
	if m.mysqlHandlers[db] != nil {
		return errors.New("MySQL实例已存在")
	}
	source := m.getDatabase(db)
	if obj, err := gorm.Open(mysql.Open(source), &mysqlConfig); err == nil {
		m.mysqlHandlers[db] = obj
		return nil
	} else {
		return err
	}
}

// Remove 移除句柄
func (m *MysqlConnPool) Remove(db string) {
	if m.mysqlHandlers[db] != nil {
		defer delete(m.mysqlHandlers, db)
		g := m.mysqlHandlers[db]
		if s, err := g.DB(); err == nil {
			s.Close()
		}
	}
}

// Handler /*
// 对外获取数据库连接对象db
func (m *MysqlConnPool) Handler(db string) (conn *gorm.DB) {
	return m.mysqlHandlers[db]
}
