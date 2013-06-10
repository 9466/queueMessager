package client

import (
	"testing"
)

var redisAddr string = "192.168.24.80:6379"
var redisDb int = 0
var redisList string = "ltest2"
var redisTestVal string = "my remote test"

var client = NewClient(redisAddr, 0, 0)

func TestRedisPush(t *testing.T) {
	err := client.Push(redisList, []byte(redisTestVal))
	if err != nil {
		t.Error("Redis Push() faild, Err: ", err.Error())
	}
}

func TestRedisPop(t *testing.T) {
	val, err := client.Pop(redisList)
	if err != nil {
		t.Error("Redis Pop() faild, Err: ", err.Error())
	}
	if string(val) != redisTestVal {
		t.Error("Redis Pop() faild, val show be ", redisTestVal, ", but: ", string(val))
	}
}

func TestRedisClose(t *testing.T) {
	err := client.Close()
	if err != nil {
		t.Error("Redis Close() faild, Err: ", err.Error())
	}
}
