package dbhelper

import (
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	"gopkg.in/redis.v5"
	"helper_go/comhelper"
	"log"
	"time"
)

// redis连接初始化
func NewRedisInit() *redis.Client {
	if RedisClient == nil {
		RedisClient = redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         LoadIni("Redis", "redis_host"),
			Password:     LoadIni("Redis", "redis_password"),
			DB:           comhelper.StringToInt(LoadIni("Redis", "redis_database")),
			DialTimeout:  time.Duration(comhelper.StringToInt(LoadIni("Redis", "redis_timeout"))) * time.Second,
			PoolSize:     1000,             // 连接池大小(数量)
			PoolTimeout:  2 * time.Minute,  // 等待忙碌连接释放的时间
			IdleTimeout:  10 * time.Minute, // 空闲连接过多久后关闭
			ReadTimeout:  2 * time.Minute,  // 读取数据超时时间
			WriteTimeout: 1 * time.Minute,  // 写数据超时时间
		})
		_, err := RedisClient.Ping().Result()
		if err != nil {
			log.Println("连接redis失败 ", err)
		}
	}
	return RedisClient
}

// 获取缓存信息
func Get(key string) (string, error) {
	v, err := NewRedisInit().Get(key).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}
func GetByMap(key string) map[string]interface{} {
	val, _ := Get(key)
	data, _ := comhelper.UnSerialize(val)
	if data != nil {
		var ret_data map[string]interface{}
		ret_data = make(map[string]interface{})
		for k, v := range data.(php_serialize.PhpArray) {
			d, ok := v.(php_serialize.PhpArray)
			if ok {
				var tmp_data map[string]interface{}
				tmp_data = make(map[string]interface{})
				for k1, v1 := range d {
					tmp_data[k1.(string)] = v1
				}
				ret_data[k.(string)] = tmp_data
			} else {
				if v == nil {
					v = ""
				}
				ret_data[k.(string)] = v
			}
		}
		return ret_data
	}
	return nil
}

// 设置缓存信息
func Save(key string, postfix string, val string) error {
	real_key := key + postfix
	err := NewRedisInit().Set(real_key, val, 0).Err()
	return err
}

// 设置带缓存的超时时间
func Save_ex(key string, postfix string, val string, expire time.Duration) error {
	real_key := key + postfix
	err := NewRedisInit().Set(real_key, val, expire).Err()
	return err
}

// 追加缓存
func Append(key string, postfix string, value string) error {
	realKey := key + postfix
	err := NewRedisInit().Append(realKey, value).Err()
	return err
}

// 删除缓存信息
func Delete(key string) error {
	err := NewRedisInit().Del(key).Err()
	return err
}

// 获取缓存（哈希）
func Hget(key, index string) (string, error) {
	v, err := NewRedisInit().HGet(key, index).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}

// 设置缓存（哈希）
func HGetAll(key string) (map[string]string, error) {
	return NewRedisInit().HGetAll(key).Result()
}

// 获取缓存（哈希）- 转换成map类型
func HgetByMap(key, index string) map[string]interface{} {
	val, _ := Hget(key, index)
	data, _ := comhelper.UnSerialize(val)
	if data != nil {
		var ret_data map[string]interface{}
		ret_data = make(map[string]interface{})
		for k, v := range data.(php_serialize.PhpArray) {
			d, ok := v.(php_serialize.PhpArray)
			if ok {
				var tmp_data map[string]interface{}
				tmp_data = make(map[string]interface{})
				for k1, v1 := range d {
					tmp_data[k1.(string)] = v1
				}
				ret_data[k.(string)] = tmp_data
			} else {
				ret_data[k.(string)] = v
			}
		}
		return ret_data
	}
	return nil
}

// 设置缓存（哈希）
func Hset(key, index, data string) error {
	err := NewRedisInit().HSet(key, index, data).Err()
	return err
}

// 删除缓存（哈希）
func Hdelete(key, index string) error {
	err := NewRedisInit().HDel(key, index).Err()
	return err
}

// 设置哈希缓存（带过期时间）
func Hset_ex(key, postfix, index, data string, tim time.Duration) error {
	real_key := key + postfix
	err := NewRedisInit().HSet(real_key, index, data).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	err = NewRedisInit().Expire(real_key, tim).Err()
	return err
}

// 从左边插入数据
func Lpush(key, postfix, value string) error {
	real_key := key + postfix
	err := NewRedisInit().LPush(real_key, value).Err()
	return err
}

// 从左边插入数据(带过期时间)
func Lpush_ex(key, postfix, value string, tim time.Duration) error {
	real_key := key + postfix
	err := NewRedisInit().LPush(real_key, value).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	err = NewRedisInit().Expire(real_key, tim).Err()
	return err
}

// 从右边插入数据
func Rpush(key, postfix, value string) error {
	real_key := key + postfix
	err := NewRedisInit().RPush(real_key, value).Err()
	return err
}

// 从左边插入数据(带过期时间)
func RPush_ex(key, postfix, value string, tim time.Duration) error {
	real_key := key + postfix
	err := NewRedisInit().RPush(real_key, value).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	err = NewRedisInit().Expire(real_key, tim).Err()
	return err
}

// 取出list数据(从左边取)
func Lpop(key, postfix string) (string, error) {
	real_key := key + postfix
	ret, err := NewRedisInit().LPop(real_key).Result()
	if err != nil {
		return "", err
	}
	return ret, err
}

// 取出list数据(从右边取)
func Rpop(key, postfix string) (string, error) {
	real_key := key + postfix
	ret, err := NewRedisInit().RPop(real_key).Result()
	if err != nil {
		return "", err
	}
	return ret, err
}

// php序列化数据
func PhpSerialize(data map[string]interface{}) string {
	encoder := php_serialize.NewSerializer()
	var new_info php_serialize.PhpArray
	new_info = make(php_serialize.PhpArray)
	for k, v := range data {
		new_info[k] = v
	}
	val, _ := encoder.Encode(new_info)
	return val
}

// 判断Key是否存在
func Exists(key string) (bool, error) {
	ret, err := NewRedisInit().Exists(key).Result()
	return ret, err
}

// 增加值
func HIncrBy(key, field string, v int) (int64, error) {
	ret, err := NewRedisInit().HIncrBy(key, field, int64(v)).Result()
	return ret, err
}

// 增加值
func HIncr_Ex(key, field string, v int, tim time.Duration) (int64, error) {
	ret, err := NewRedisInit().HIncrBy(key, field, int64(v)).Result()
	// 设置过期时间
	err = NewRedisInit().Expire(key, tim).Err()
	return ret, err
}
