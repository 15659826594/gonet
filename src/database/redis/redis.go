package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	once sync.Once
	ctx  context.Context
	rdbs map[int]*redis.Client
)

func init() {
	ctx = context.Background()
}

func NewClient(index int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "",
		DB:          index,
		PoolSize:    10,
		PoolTimeout: 5 * time.Second,
	})
	if rdbs == nil {
		rdbs = map[int]*redis.Client{}
	}
	rdbs[index] = rdb
	return rdb
}

func Instance(dbIndex ...int) *redis.Client {
	index := 0
	if len(dbIndex) > 0 {
		index = dbIndex[0]
	}
	once.Do(func() {
		rdb := NewClient(index)

		configResult, err := rdb.ConfigGet(ctx, "notify-keyspace-events").Result()
		if err != nil {
			fmt.Printf("获取 Redis 配置失败：%v\n", err)
		} else {
			fmt.Printf("当前 notify-keyspace-events 配置：%v\n", configResult)
		}
		// 订阅过期事件频道
		pubsub := rdb.PSubscribe(ctx, "__keyevent@0__:expired")

		// 处理消息
		go func() {
			for msg := range pubsub.Channel() {
				fmt.Printf("Key %s expired at %s\n", msg.Payload, time.Now())
			}
		}()

		// 测试：创建一个 5 秒后过期的 key
		testKey := "test:expire:" + time.Now().Format("15:04:05")
		err = rdb.Set(ctx, testKey, "this will expire in 5 seconds", 5*time.Second).Err()
		if err != nil {
			fmt.Printf("设置测试 key 失败：%v\n", err)
		} else {
			fmt.Printf("已创建测试 key: %s (5 秒后过期)\n", testKey)
		}

		// 验证连接
		pingResult, err := rdb.Ping(ctx).Result()
		if err != nil {
			fmt.Printf("Redis 连接失败：%v\n", err)
		} else {
			fmt.Printf("Redis 连接成功：%s\n", pingResult)
		}
	})
	return rdbs[index]
}

func Close() {
	for _, client := range rdbs {
		client.Close()
	}
}

func UseLua(userid, prodid string) bool {
	//编写脚本 - 检查数值，是否够用，够用再减，否则返回减掉后的结果
	var luaScript = redis.NewScript(`
		local userid=KEYS[1];
		local prodid=KEYS[2];
		local qtKey="sk:"..prodid..":qt";
		local userKey="sk:"..prodid..":user";
		local userExists=redis.call("sismember",userKey,userid);
		if tonumber(userExists)==1 then
		 return 2;
		end
		local num=redis.call("get",qtKey);
		if tonumber(num)<=0 then
		 return 0;
		else
		 redis.call("decr",qtKey);
		 redis.call("SAdd",userKey,userid);
		end
		return 1;
	`)
	//执行脚本
	n, err := luaScript.Run(ctx, Instance(), []string{userid, prodid}).Result()
	if err != nil {
		return false
	}
	switch n {
	case int64(0):
		fmt.Println("抢购结束")
		return false
	case int64(1):
		fmt.Println(userid, "：抢购成功")
		return true
	case int64(2):
		fmt.Println(userid, "：已经抢购了")
		return false
	default:
		fmt.Println("发生未知错误！")
		return false
	}
}
