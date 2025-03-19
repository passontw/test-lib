package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sl.framework.com/resource/conf"
	"sl.framework.com/trace"
	"strconv"
	"strings"
	"time"
)

var (
	RedisKeyPrefix = "g32"
	flag           = conf.Section("nacos", "namespaceId")
	LoginKey       = fmt.Sprintf("%s-%s-%s:%s",
		RedisKeyPrefix, conf.Section("beego", "appName"), flag, "%s_dealer_session")
	TaskKey = fmt.Sprintf("%s-%s-%s:%s",
		RedisKeyPrefix, conf.Section("beego", "appName"), flag, "deal_tasks_schedule_time")
)

var ctx = context.Background()

type RCache struct {
	cache redis.UniversalClient
}

var c *RCache

func init() {
	addr := conf.Section("redis", "host")
	port := conf.Section("redis", "port")
	endpoint := fmt.Sprintf("%s:%s", addr, port)
	host, err := parseHost(endpoint)
	if err != nil {
		trace.Error("the redis %s resolution is incorrect, please check error: %s", endpoint, err.Error())
	}
	username := conf.Section("redis", "username")
	pwd := conf.Section("redis", "pwd")
	_, err = initRedis(username, pwd, host...)
	if err == nil {
		trace.Notice("redis initialization successful")
	} else {
		trace.Error("redis initialization failed, error: %s", err.Error())
	}
}

// initRedis 获取redis对象
func initRedis(username, password string, addr ...string) (*RCache, error) {
	trace.Notice("Redis: user: [%s], pwd: [%s], addr: %v", username, password, addr)
	c = new(RCache)

	params := &redis.UniversalOptions{
		Addrs:    addr,
		PoolSize: 200,
		DB:       0,
	}

	if username != "" {
		params.Username = username
	}
	if password != "" {
		params.Password = password
	}

	mode := conf.Section("redis", "mode")
	dbIdx := conf.Section("redis", "db")
	if mode != "cluster" {
		if len(addr) > 0 {
			params.Addrs = []string{addr[0]}
			params.DB, _ = strconv.Atoi(dbIdx)
			trace.Notice("Redis is [single] mode, addr: %v db: %d", params.Addrs, params.DB)
		}
		c.cache = redis.NewUniversalClient(params) // 单节点客户端
	} else {
		trace.Notice("Redis is [cluster] mode")
		c.cache = redis.NewClusterClient(params.Cluster()) // 集群客户端
	}

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := c.cache.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	trace.Notice("Redis: user: [%s], pwd: [%s], addr: %v connection successful", username, password, params.Addrs)
	return c, nil
}

// parseHost 解析 host 地址，如果包含 "-"，返回地址组合，否则返回原地址
func parseHost(host string) ([]string, error) {
	parts := strings.Split(host, ":")
	if len(parts) != 2 {
		return nil, errors.New("the host address is incorrect") // 地址格式不正确
	}

	// 提取主机地址和端口部分
	hostAddress := parts[0]
	portRange := parts[1]

	// 如果端口范围包含 "-"，生成所有端口的地址组合
	if strings.Contains(portRange, "-") {
		var addresses []string
		// 分割端口范围
		rangeParts := strings.Split(portRange, "-")
		if len(rangeParts) != 2 {
			return nil, errors.New("the port range format is incorrect")
		}
		// 将端口字符串转换为整数
		s, err1 := strconv.Atoi(rangeParts[0])
		e, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil {
			return nil, errors.New("error converting port string to integer")
		}

		for i := s; i <= e; i++ {
			addresses = append(addresses, fmt.Sprintf("%s:%d", hostAddress, i))
		}
		return addresses, nil
	}
	return []string{host}, nil
}

func Get() *RCache {
	return c
}

func (r *RCache) IsExist(key string) bool {
	_, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *RCache) Delete(key string) error {
	_, err := r.cache.Del(ctx, key).Result()
	return err
}

// GetSequence 自增序列
func (r *RCache) GetSequence(key string) (int64, error) {
	seq, err := r.cache.Incr(ctx, key).Result()
	return seq, err
}

func (r *RCache) GetString(key string) (string, error) {
	return r.cache.Get(ctx, key).Result()
}

func (r *RCache) GetInt(key string) (int, error) {
	return r.cache.Get(ctx, key).Int()
}

func (r *RCache) GetInt64(key string) (int64, error) {
	return r.cache.Get(ctx, key).Int64()
}

func (r *RCache) GetFloat64(key string) (float64, error) {
	return r.cache.Get(ctx, key).Float64()
}

// SetString 设置KEY值 永不过期
func (r *RCache) SetString(key string, val string) error {
	err := r.cache.Set(ctx, key, val, 0).Err()
	return err
}

// SetStringExpire 设置KEY值
func (r *RCache) SetStringExpire(key string, val string, expiration time.Duration) error {
	err := r.cache.Set(ctx, key, val, expiration).Err()
	return err
}

// SetNX 分布式锁
func (r *RCache) SetNX(key string, val string, expiration time.Duration) (bool, error) {
	return r.cache.SetNX(ctx, key, val, expiration).Result()
}

// PushList 存入队列(左进)
func (r *RCache) PushList(key string, val string) error {
	err := r.cache.LPush(ctx, key, val).Err()
	return err
}

// PopList 移除并获取列表最后一个元素(右出)
func (r *RCache) PopList(key string) string {
	return r.cache.RPop(ctx, key).Val()
}

// BRPopList 阻塞式的移除并获取列表最后一个元素(右出)
func (r *RCache) BRPopList(key string) string {
	rec := r.cache.BRPop(ctx, time.Duration(0), key).Val()
	if len(rec) == 2 {
		return rec[1]
	}
	return ""
}

// GetIndexList 通过索引获取列表中的元素
func (r *RCache) GetIndexList(key string, dx int64) string {
	return r.cache.LIndex(ctx, key, dx).Val()
}

func (r *RCache) HSet(key string, fields string, val string) error {
	return r.cache.HSet(ctx, key, fields, val).Err()
}

func (r *RCache) HGet(key string, fields string) string {
	val := r.cache.HGet(ctx, key, fields).Val()
	return val
}

func (r *RCache) HGetAll(key string) map[string]string {
	val := r.cache.HGetAll(ctx, key).Val()
	return val
}

func (r *RCache) HDel(key string, fields string) error {
	return r.cache.HDel(ctx, key, fields).Err()
}

// BRPopLPush 队列转移
func (r *RCache) BRPopLPush(source, destination string) string {
	return r.cache.BRPopLPush(ctx, source, destination, 0).Val()
}

// RPopLPush 队列转移
func (r *RCache) RPopLPush(source, destination string) string {
	return r.cache.RPopLPush(ctx, source, destination).Val()
}
