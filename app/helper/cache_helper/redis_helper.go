package cache_helper

import (
	"context"
	"fmt"
	"gin_base/app/helper/exception_helper"
	"gin_base/app/helper/helper"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"
	"sync"
	"time"
)

var redisHepers = make(map[string]redisHelper)
var redisHepersRWMutex sync.RWMutex

type redisHelper struct {
	Client *redis.Client
}

// 获取redis助手单例
func RedisHelper(connectName ...string) redisHelper {
	currentConnectName := "default"
	if len(connectName) > 0 {
		currentConnectName = connectName[0]
	}
	redisHepersRWMutex.RLock() //读锁
	redis_helper, exists := redisHepers[currentConnectName]
	redisHepersRWMutex.RUnlock() //立即解开读锁
	if exists {
		return redis_helper
	}
	//不存在，需要初始化
	redisHepersRWMutex.Lock()         //写锁
	defer redisHepersRWMutex.Unlock() //初始化完毕自动解开写锁
	//防止在获取写锁已经被其他协程初始化
	if redis_helper, exists := redisHepers[currentConnectName]; exists {
		return redis_helper
	}

	redis_helper.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", helper.GetAppConfig().Redis[currentConnectName].Host, helper.GetAppConfig().Redis[currentConnectName].Port), // Redis 服务器地址
		Password: helper.GetAppConfig().Redis[currentConnectName].Password,                                                                         // 没有密码，默认值
		DB:       helper.GetAppConfig().Redis[currentConnectName].Select,                                                                           // 使用默认 DB
	})

	redisHepers[currentConnectName] = redis_helper

	return redis_helper
}

// 获取redis缓存
func (rh redisHelper) RedisGet(key string) (string, error) {
	res, err := rh.Client.Get(context.Background(), key).Result()
	return res, err
}

// 设置redis缓存
func (rh redisHelper) RedisSet(key string, value interface{}, expiration time.Duration) error {
	return rh.Client.Set(context.Background(), key, value, expiration).Err()
}

// 删除redis缓存
func (rh redisHelper) RedisDel(key string) error {
	return rh.Client.Del(context.Background(), key).Err()
}

// 获取redis分布式锁
func (rh redisHelper) RedisLock(key string, expirations ...time.Duration) string {
	key = "RedisLock:" + key
	ctx := context.Background()
	expiration := 60 * time.Second
	if len(expirations) > 0 {
		expiration = expirations[0]
	}
	id := ksuid.New().String()
	success, err := rh.Client.SetNX(ctx, key, id, expiration).Result()
	if err != nil || !success {
		id = ""
	}
	return id
}

// 循环等待获取锁，否则报错
func (rh redisHelper) RedisWaitLockOrException(key string, expiration time.Duration, wait time.Duration) string {
	start := time.Now()
	for {
		id := rh.RedisLock(key, expiration)
		if id != "" {
			return id
		}
		if time.Since(start) > wait {
			exception_helper.CommonException("上一个请求正在处理中，请稍后")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// 删除redis分布式锁
func (rh redisHelper) RedisUnLock(key string, lockId string) error {
	key = "RedisLock:" + key
	id, _ := rh.RedisGet(key)
	if id == lockId {
		return rh.RedisDel(key)
	}
	return nil
}

// 设置redis频率限制
func (rh redisHelper) RedisLimit(key string, limit int, expireTime int) bool {
	script := `
			local key = KEYS[1]
			local limit = tonumber(ARGV[1])
			local expire_time = tonumber(ARGV[2])
			
			local current = redis.call('incr', key)
			if current == 1 then
				redis.call('expire', key, expire_time)
			end
			
			if current > limit then
				return 0
			else
				return 1
			end
	`
	// 执行 Lua 脚本
	result, err := rh.Client.Eval(context.Background(), script, []string{key}, limit, expireTime).Result()
	if err != nil {
		return false
	}
	// 将结果转换为布尔值
	return result.(int64) == 1
}
