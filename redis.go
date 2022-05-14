package godbutil

import (
	"github.com/go-redis/redis"
	"github.com/kordar/goutil"
	"log"
	"sync"
	"time"
)

var instanceOfRedisPool *RedisConnPool
var onceOfRedis sync.Once

// RedisConnPool /*
// Redis连接操作库
// 基于go-redis封装开发
type RedisConnPool struct {
	redisHandlers map[string]*redis.Client
}

func GetRedisPool() *RedisConnPool {
	onceOfRedis.Do(func() {
		instanceOfRedisPool = &RedisConnPool{redisHandlers: make(map[string]*redis.Client)}
	})
	return instanceOfRedisPool
}

// GetOptions 获取配置
func (p *RedisConnPool) GetOptions(section string) redis.Options {

	addr := goutil.GetSectionValue(section, "addr")
	password := goutil.GetSectionValue(section, "password")
	db := goutil.GetSectionValueInt(section, "db")
	poolSize := goutil.GetSectionValueInt(section, "poolSize")
	maxRetries := goutil.GetSectionValueInt(section, "maxRetries")
	idleTimeout := goutil.GetSectionValueInt(section, "idleTimeout")
	minIdleConns := goutil.GetSectionValueInt(section, "minIdleConns")

	return redis.Options{
		Addr:         addr,                                     // Redis地址
		Password:     password,                                 // Redis账号
		DB:           db,                                       // Redis库
		PoolSize:     poolSize,                                 // Redis连接池大小
		MaxRetries:   maxRetries,                               // 最大重试次数
		IdleTimeout:  time.Duration(idleTimeout) * time.Second, // 空闲链接超时时间
		MinIdleConns: minIdleConns,                             // 空闲连接数量
	}
}

// InitDataPool 初始化redis连接池
func (p *RedisConnPool) InitDataPool(db string) bool {
	options := p.GetOptions(db)
	client := redis.NewClient(&options)
	if ok := p.Ping(client); ok {
		p.redisHandlers[db] = client
		return true
	} else {
		return false
	}
}

// Ping 测试连接
func (p *RedisConnPool) Ping(client *redis.Client) bool {
	pong, err := client.Ping().Result()
	if err == redis.Nil {
		log.Println("Redis异常")
		return false
	} else if err != nil {
		log.Println("失败:", err)
		return false
	} else {
		log.Println(pong)
		return true
	}
}

// Handler 对外获取Redis连接对象client
func (p *RedisConnPool) Handler(db string) (conn *redis.Client) {
	return p.redisHandlers[db]
}
