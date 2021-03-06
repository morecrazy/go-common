package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"third/redigo/redis"
	"time"
)

type Cache struct {
	redisPool *redis.Pool
	Redis     RedisConfig
}

const (
	Success     int = 1
	KeyNotFound int = 2
	RedisError  int = 3
)

func CheckRedisReturnValue(err error) int {
	if err != nil && strings.Contains(err.Error(), "nil returned") {
		return KeyNotFound
	} else if err == nil {
		return Success
	} else {
		return RedisError
	}
}

func InitRedisPool(my_redis *RedisConfig) (*Cache, error) {
	cache := new(Cache)
	cache.Redis = *my_redis
	cache.RedisPool()
	/*
		err := pool.TestOnBorrow(pool.Get(), time.Now())
		if err != nil {
			fmt.Println("init cache error :", my_redis, err)
			return nil, err
		}*/
	return cache, nil
}

func (cache *Cache) RedisPool() *redis.Pool {
	if cache.redisPool == nil {
		cache.NewRedisPool(&cache.Redis)
	}
	return cache.redisPool
}

func (cache *Cache) NewRedisPool(my_redis *RedisConfig) {
	cache.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			//fmt.Println(*my_redis)
			var connect_timeout time.Duration = time.Duration(my_redis.ConnectTimeout) * time.Second
			var read_timeout = time.Duration(my_redis.ReadTimeout) * time.Second
			var write_timeout = time.Duration(my_redis.WriteTimeout) * time.Second

			//c, err := redis.DialTimeout(config.Redis.Network, config.Redis.Address, connect_timeout, read_timeout, write_timeout)
			c, err := redis.DialTimeout("tcp", my_redis.RedisConn, connect_timeout, read_timeout, write_timeout)

			if err != nil {
				fmt.Println("DialTimeout", my_redis.RedisConn)
				return nil, err
			}
			if len(my_redis.RedisPasswd) > 0 {
				if _, err := c.Do("AUTH", my_redis.RedisPasswd); err != nil {
					c.Close()
					return nil, err
				}
			}

			if my_redis.RedisDb != "" {
				if _, err := c.Do("SELECT", my_redis.RedisDb); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		}, /*
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				fmt.Println("PING")
				return err
			},*/
		MaxIdle:     my_redis.MaxIdle,
		MaxActive:   my_redis.MaxActive,
		IdleTimeout: time.Duration(my_redis.IdleTimeout) * time.Second,
		Wait:        true,
	}
}

func (cache *Cache) Get(key string) ([]byte, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()

	res, err := redis.Bytes(conn.Do("GET", key))
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Incr(key string) (int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("INCR", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) IncrInt64(key string) (int64, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int64(conn.Do("INCR", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Incrby(key string, value int) (int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("INCRBY", key, value))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) IncrbyInt64(key string, value int) (int64, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int64(conn.Do("INCRBY", key, value))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Decr(key string) (int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("DECR", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Decrby(key string, value int) (int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("DECRBY", key, value))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) MGet(key []interface{}) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := conn.Do("MGET", key...)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) MGetValue(keys []interface{}) ([]interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Values(conn.Do("MGET", keys...))
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}
func (cache *Cache) HSet(key, field string, value interface{}) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := conn.Do("HSET", key, field, value)

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}
func (cache *Cache) HMset(value []interface{}) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := conn.Do("HMSET", value...)

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) HGet(key, field string) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := conn.Do("HGET", key, field)

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) HIncrby(key, field string, value interface{}) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := conn.Do("HINCRBY", key, field, value)

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Hmget(key string, fields []string) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()

	var args []interface{}
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}

	res, err := conn.Do("HMGET", args...)

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) GetString(key string) (string, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.String(conn.Do("GET", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) GetStringMap(key string) (map[string]string, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.StringMap(conn.Do("HGETALL", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) HGetAll(key string) ([]byte, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Bytes(conn.Do("HGETALL", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) GetInt(key string) (int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("GET", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) GetInt64(key string) (int64, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int64(conn.Do("GET", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) GetInts(key string) ([]int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Ints(conn.Do("GET", key))

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Expire(key string, timeout int) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", key, timeout)

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}

	return err
}

func (cache *Cache) Set(key string, bytes interface{}, timeout int) error {
	var err error
	conn := cache.RedisPool().Get()
	defer conn.Close()
	if timeout == -1 {
		_, err = conn.Do("SET", key, bytes)
	} else {
		_, err = conn.Do("SET", key, bytes, "EX", timeout)
	}

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return err
}

func (cache *Cache) Del(key string) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return err
}

func (cache *Cache) Exists(key string) (bool, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	var flag bool
	exists, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
		Logger.Error(err.Error())
		return flag, err
	}
	if exists == 1 {
		flag = true
	}
	return flag, nil
}

func (cache *Cache) Zrange(key string, start, end int, withscores bool) ([]string, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	var res []string
	var err error
	if withscores {
		res, err = redis.Strings(conn.Do("ZRANGE", key, start, end, "withscores"))
	} else {
		res, err = redis.Strings(conn.Do("ZRANGE", key, start, end))
	}
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) ZrangeInts(key string, start, end int, withscores bool) ([]int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	var res []int
	var err error
	if withscores {
		res, err = redis.Ints(conn.Do("ZRANGE", key, start, end, "withscores"))
	} else {
		res, err = redis.Ints(conn.Do("ZRANGE", key, start, end))
	}
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Zrevrange(key string, start, end int, withscores bool) ([]int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	var res []int
	var err error
	if withscores {
		res, err = redis.Ints(conn.Do("ZREVRANGE", key, start, end, "withscores"))
	} else {
		res, err = redis.Ints(conn.Do("ZREVRANGE", key, start, end))
	}

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

//add by yuan xiang
func (cache *Cache) ZrevrangeString(key string, start, end int, withscores bool) ([]string, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	var res []string
	var err error
	if withscores {
		res, err = redis.Strings(conn.Do("ZREVRANGE", key, start, end, "withscores"))
	} else {
		res, err = redis.Strings(conn.Do("ZREVRANGE", key, start, end))
	}

	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) ZrevrangeByScore(key string, max_num, min_num int, withscores bool, offset, count int) ([]int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	var res []int
	var err error
	if !withscores {
		res, err = redis.Ints(conn.Do("ZREVRANGEBYSCORE", key, max_num, min_num, "limit", offset, count))
	} else {
		res, err = redis.Ints(conn.Do("ZREVRANGEBYSCORE", key, max_num, min_num, "withscores", "limit", offset, count))
	}
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}
func (cache *Cache) ZrangeByScore(key string, min_num, max_num int64, withscores bool, offset, count int) ([]int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	var res []int
	var err error
	if withscores {
		res, err = redis.Ints(conn.Do("ZREVRANGEBYSCORE", key, max_num, min_num, "limit", offset, count))
	} else {
		res, err = redis.Ints(conn.Do("ZREVRANGEBYSCORE", key, max_num, min_num, "withscores", "limit", offset, count))
	}
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Zscore(key, member string) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()

	res, err := conn.Do("ZSCORE", key, member)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Zrevrank(key, member string) (int64, error) {
	/*
		获取排名
		Returns a 0-based value indicating the descending rank of
		``value`` in sorted set ``key``
		for SimpleSortedSet

		** rank start from 0 **

		update by wuql 2016-7-25
	*/
	conn := cache.RedisPool().Get()
	defer conn.Close()
	rank, err := redis.Int64(conn.Do("ZREVRANK", key, member))
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		return 0, err
	} else {
		return rank + 1, nil
	}
}

// 判断某个对象是否存在 存在返回true 不存在返回false
// for SimpleSortedSet
// update by wuql 2016-6-28
func (cache *Cache) IsMember(key, member string) bool {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := conn.Do("ZSCORE", key, member)
	// res为score(value)
	if res != nil && err == nil {
		return true
	} else {
		return false
	}
}

// 批量插入 for SortedSet（有序集合）
// score_member_list: score1, member2, score2, member2...
// update by wuql 2016-6-24
func (cache *Cache) ZaddBatch(key string, score_member_list []interface{}) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	// 必须是:score(value) member(key)  数据对
	// len(score_member_list) >= 2 && len(score_member_list)%2 == 0
	res, err := conn.Do("ZADD", score_member_list...)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

// ZREM命令从有序集合存储键删除指定成员
// update by wuql 2016-6-24
func (cache *Cache) Zrem(key, member string) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()

	res, err := conn.Do("ZREM", key, member)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

// Set the value of key ``name`` to ``value`` if key doesn't exist"
// "Set the value of key ``key`` to ``value`` if key exist"
// update by wuql 2016-6-27
func (cache *Cache) ENSet(key string, value interface{}) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	is_exist, err := cache.Exists(key)
	if err == nil {
		if is_exist == true {
			res, err := conn.Do("SETEX", key, value)
			if nil != err && !strings.Contains(err.Error(), "nil returned") {
				err = NewInternalError(CacheErrCode, err)
			}
			return res, err
		} else {
			res, err := conn.Do("SETNX", key, value)
			if nil != err && !strings.Contains(err.Error(), "nil returned") {
				err = NewInternalError(CacheErrCode, err)
			}
			return res, err
		}
	} else {
		return nil, err
	}
}

func (cache *Cache) Zadd(key string, value, member interface{}) (interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()

	res, err := conn.Do("ZADD", key, value, member)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

//add by yuan xiang
func (cache *Cache) Zcard(key string) (int, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()

	res, err := redis.Int(conn.Do("ZCARD", key))
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Sadd(key string, items string) (int, error) {
	//var err error
	conn := cache.RedisPool().Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("SADD", key, items))
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return res, err
}

func (cache *Cache) Rpush(key string, value interface{}) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", key, value)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return err
}

func (cache *Cache) Rpop(key string) (value interface{}, err error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	value, err = conn.Do("RPOP", key)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return
}

func (cache *Cache) RpushBatch(keys []interface{}) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", keys...)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return err
}

func (cache *Cache) Lrange(key string, start, end int) ([]interface{}, error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	result, err := redis.Values(conn.Do("LRANGE", key, start, end))
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return result, err
}

func (cache *Cache) Lrem(key string, value interface{}) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("LREM", key, 0, value)
	if nil != err && !strings.Contains(err.Error(), "nil returned") {
		err = NewInternalError(CacheErrCode, err)
	}
	return err
}

func (cache *Cache) Push(key string, bydata []byte) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("LPUSH", key, bydata)
	return err
}

func (cache *Cache) Publish(channel, msg string) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("PUBLISH", channel, msg)
	return err
}

func (cache *Cache) Llen(key string) (key_len int, err error) {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	key_len, err = redis.Int(conn.Do("LLEN", key))
	return
}

//Set the value of key ``name`` to ``value`` that expires in ``time`` seconds
func (cache *Cache) Setex(name string, value, time int64) error {
	conn := cache.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", name, time, value)
	return err
}

func (cache *Cache) SetTimeLock(id string, time_out int64) (flag bool, err error) {
	key := fmt.Sprintf("tlock:%s", id)
	is_exist, err := cache.Exists(key)
	if err != nil {
		return
	}
	if is_exist {
		return
	}
	if err = cache.Setex(key, 0, time_out); err != nil {
		return
	}
	flag = true
	return
}

// by liudan 2016.07.29
func (cache *Cache) GetJsonObj(key string, obj interface{}) error {
	data, err := cache.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

func (cache *Cache) SaveJsonObj(key string, obj interface{}, timeout int) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return cache.Set(key, data, timeout)
}
