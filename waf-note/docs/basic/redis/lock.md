## 单点redis分布式锁
https://xiaolincoding.com/redis/base/redis_interview.html#%E5%A6%82%E4%BD%95%E7%94%A8-redis-%E5%AE%9E%E7%8E%B0%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E7%9A%84

- 客户端唯一表示
- set nx
- 锁超时
- 续锁
- 释放锁： lua 保证原子性

## redis集群分布式锁 RedLock
单点redis锁的缺点： **Redis 主从复制模式中的数据是异步复制的，这样导致分布式锁的不可靠性**。如果在 Redis 主节点获取到锁后，在没有同步到其他节点时，Redis 主节点宕机了，此时新的 Redis 主节点依然可以获取锁，所以多个应用服务就可以同时获取到锁。

RedLock的思想是：
**是让客户端和多个独立的 Redis 节点依次请求申请加锁，如果客户端能够和半数以上的节点成功地完成加锁操作，那么我们就认为，客户端成功地获得分布式锁，否则加锁失败**。
这样一来，即使有某个 Redis 节点发生故障，因为锁的数据在其他节点上也有保存，所以客户端仍然可以正常地进行锁操作，锁的数据也不会丢失。

## 实战

``` 
//RedisLock Lock
type RedisLock struct {
    key      string
    token    string
    ticker   *time.Ticker
    duration time.Duration
}


//lock, err := utils.GetLock(job.Host.Domain, uuid.New().String(), 1*time.Minute)
func GetLock(key string, token string, duration time.Duration) (*RedisLock, error) {
    succ, err := models.TaskRedisCli.SetNX(key, token, duration).Result()
    if err != nil || !succ {
        return nil, err
    }
    l := &RedisLock{
        key:      key,
        token:    token,
        duration: duration,
        ticker:   nil,
    }
    return l, nil
}

//KeepAlive 向redis续命
func (l *RedisLock) KeepAlive() {
    go func() {
        if l.ticker != nil {
            tickerDuration := l.duration / 3
            if tickerDuration < time.Second*5 {
                tickerDuration = time.Second * 5
            }
            ticker := time.NewTicker(tickerDuration)
            l.ticker = ticker
            for range l.ticker.C {
                _, err := models.TaskRedisCli.Expire(l.key, l.duration).Result()
                if err != nil {
                    l.ticker = nil
                    ticker.Stop()
                    return
                }
            }
        }
    }()
}

//Stop stop
func (l *RedisLock) Stop() {
    if l.ticker != nil {
        ticker := l.ticker
        l.ticker = nil
        ticker.Stop()
    }
    l.UnLock()
}

//UnLock 解锁
func (l *RedisLock) UnLock() (bool, error) {
    luaScript := "if redis.call('get',KEYS[1]) == ARGV[1] then " +
        "return redis.call('del',KEYS[1]) else return 0 end"
    result, err := models.TaskRedisCli.Eval(luaScript, []string{l.key}, []string{l.token}).Result()
    if err != nil {
        return false, err
    }

    ret, ok := result.(int)
    if ok && ret != 0 {
        return true, nil
    }
    return false, nil
}
```

```
token需要时客户端的唯一标识 ： 例如 uuid.New().String()
// GetLock 基于redis实现的分布式锁，超时锁会自动释放。
 lock, err = GetLock(key, token, 30*time.Second)
 if err != nil {
        logger.Errorf("get lock by redis failed: %s", err.Error())
        return
    }
if lock == nil {
        logger.Errorf("get lock by redis failed, job maybe running")
        return
    }
 lock.KeepAlive() 		// 定时续锁
 defer lock.Stop()      // 停止拿锁
```

