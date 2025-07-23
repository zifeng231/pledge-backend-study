package db

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"pledge-backend-study/config"
	"pledge-backend-study/log"
)

func InitRedis() *redis.Pool {
	log.Logger.Info("Init Redis")
	redisConf := config.Config.Redis
	//建立连接池
	RedisConn = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			dial, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisConf.Address, redisConf.Port))

			if err != nil {
				return nil, err
			}
			//验证密码
			_, err = dial.Do("AUTH", redisConf.Password)
			if err != nil {
				panic("redis auth err" + err.Error())
			}
			//选择数据库
			_, err = dial.Do("SELECT", redisConf.Db)
			if err != nil {
				panic("redis select db err" + err.Error())
			}
			return dial, nil

		},
		MaxIdle:   10,   // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxActive: 0,    // 最大的激活连接数，表示同时最多有N个连接   0 表示无穷大
		Wait:      true, // 如果连接数不足则阻
	}

	err := RedisConn.Get().Err()
	if err != nil {
		panic("redis connect err" + err.Error())
	}
	return RedisConn
}

// 手动写一个set方法
// value是任意类型
func RedisSet(key string, value interface{}, expire int) error {
	//获取连接
	conn := RedisConn.Get()
	//defer 是 Go 语言的关键字，表示延迟执行，即在当前函数返回之前执行括号里的语句
	defer func() {
		if err := conn.Close(); err != nil {
			log.Logger.Error("redis close err")
		}
	}()
	marshal, err := json.Marshal(value)
	if err != nil {
		log.Logger.Error("redis set err")
		return err
	}
	if expire > 0 {
		_, err = conn.Do("SETEX", key, expire, marshal)
	}
	if err != nil {
		return err
	}
	return nil
}

//获取value

func RedisGet(key string) (interface{}, error) {
	conn := RedisConn.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Logger.Error("redis close err")
		}
	}()
	bytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
