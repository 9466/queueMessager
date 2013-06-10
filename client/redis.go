package client

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

type myRedis struct {
	pool *redis.Pool
}

func NewClient(addr string, poolSize int, timeOut time.Duration) (c *myRedis) {
	if poolSize == 0 {
		poolSize = 3
	}
	if timeOut == 0 {
		timeOut = 3
	}
	newConn := func() (redis.Conn, error) {
		conn, err := redis.DialTimeout("tcp4", addr, time.Second*timeOut, time.Second*timeOut, time.Second*timeOut)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	c = new(myRedis)
	c.pool = redis.NewPool(newConn, poolSize)
	return
}

// use redis list LPOP to get a key
func (c *myRedis) Pop(list string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	data, err := conn.Do("Lpop", list)
	if err != nil {
		return nil, err
	}
	val, ok := data.([]byte)
	if !ok {
		return nil, errors.New("data format error")
	}
	return val, nil
}

// use redis list RPUSH to add a key
func (c *myRedis) Push(list string, val []byte) (err error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err = conn.Do("Rpush", list, val)
	return
}

// close redis pool
func (c *myRedis) Close() (err error) {
	return c.pool.Close()
}
