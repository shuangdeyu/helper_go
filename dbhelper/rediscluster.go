package dbhelper

import (
	"gopkg.in/redis.v5"
	"helper_go/comhelper"
	"log"
	"strings"
	"time"
)

// redis集群连接初始化
func NewRedisClusterInit() *redis.ClusterClient {
	if RedisClusterClient == nil {
		host := LoadIni("RedisCluster", "redis_host")
		hosts := strings.Split(host, ",")
		RedisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        hosts, // master主机
			Password:     LoadIni("RedisCluster", "redis_password"),
			DialTimeout:  time.Duration(comhelper.StringToInt(LoadIni("RedisCluster", "redis_timeout"))) * time.Second,
			PoolSize:     1000,             // 连接池大小(数量)
			PoolTimeout:  2 * time.Minute,  // 等待忙碌连接释放的时间
			IdleTimeout:  10 * time.Minute, // 空闲连接过多久后关闭
			ReadTimeout:  2 * time.Minute,  // 读取数据超时时间
			WriteTimeout: 1 * time.Minute,  // 写数据超时时间
		})
		_, err := RedisClusterClient.Ping().Result()
		if err != nil {
			log.Println("连接redis集群失败 ", err)
		}
	}
	return RedisClusterClient
}

// 获取缓存信息
func ClusterGet(key string) (string, error) {
	v, err := NewRedisClusterInit().Get(key).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}

// 设置缓存信息
func ClusterSave(key string, postfix string, val string) error {
	real_key := key + postfix
	err := NewRedisClusterInit().Set(real_key, val, 0).Err()
	return err
}

// 设置带缓存的超时时间
func ClusterSaveEx(key string, postfix string, val string, expire time.Duration) error {
	real_key := key + postfix
	err := NewRedisClusterInit().Set(real_key, val, expire).Err()
	return err
}

// 删除缓存信息
func ClusterDelete(key string) error {
	err := NewRedisClusterInit().Del(key).Err()
	return err
}

// 获取缓存（哈希）
func ClusterHget(key, index string) (string, error) {
	v, err := NewRedisClusterInit().HGet(key, index).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}

// 设置缓存（哈希）
func ClusterHset(key, index, data string) error {
	err := NewRedisClusterInit().HSet(key, index, data).Err()
	return err
}

// 删除缓存（哈希）
func ClusterHdelete(key, index string) error {
	err := NewRedisClusterInit().HDel(key, index).Err()
	return err
}

// 设置哈希缓存（带过期时间）
func ClusterHsetEx(key, postfix, index, data string, tim time.Duration) error {
	real_key := key + postfix
	err := NewRedisClusterInit().HSet(real_key, index, data).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	err = NewRedisClusterInit().Expire(real_key, tim).Err()
	return err
}

// 从左边插入数据
func ClusterLpush(key, postfix, value string) error {
	real_key := key + postfix
	err := NewRedisClusterInit().LPush(real_key, value).Err()
	return err
}

// 从左边插入数据(带过期时间)
func ClusterLpushEx(key, postfix, value string, tim time.Duration) error {
	real_key := key + postfix
	err := NewRedisClusterInit().LPush(real_key, value).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	err = NewRedisClusterInit().Expire(real_key, tim).Err()
	return err
}

// 从右边插入数据
func ClusterRpush(key, postfix, value string) error {
	real_key := key + postfix
	err := NewRedisClusterInit().RPush(real_key, value).Err()
	return err
}

// 从左边插入数据(带过期时间)
func ClusterRPushEx(key, postfix, value string, tim time.Duration) error {
	real_key := key + postfix
	err := NewRedisClusterInit().RPush(real_key, value).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	err = NewRedisClusterInit().Expire(real_key, tim).Err()
	return err
}

// 取出list数据(从左边取)
func ClusterLpop(key, postfix string) (string, error) {
	real_key := key + postfix
	ret, err := NewRedisClusterInit().LPop(real_key).Result()
	if err != nil {
		return "", err
	}
	return ret, err
}

// 取出list数据(从右边取)
func ClusterRpop(key, postfix string) (string, error) {
	real_key := key + postfix
	ret, err := NewRedisClusterInit().RPop(real_key).Result()
	if err != nil {
		return "", err
	}
	return ret, err
}
