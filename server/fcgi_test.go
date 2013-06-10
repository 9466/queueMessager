package server

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	data := make(map[string]string)
	data["message"] = "test"
	client := NewServer("192.168.24.80", 9000)
	res, err := client.Request("/data/www/upload/test.php", data)
	if err != nil {
		t.Error(err.Error())
	}
	body := client.ParseBody(string(res))
	fmt.Println(body)
}

func TestRequest2(t *testing.T) {
	chs := make(chan int, 10)
	runtime.GOMAXPROCS(10)
	for i := 0; i < 10; i++ {
		go goTest(chs)
	}
	for i := 0; i < 10; i++ {
		<-chs
	}
}

func goTest(ch chan int) {
	data := make(map[string]string)
	data["message"] = "test"
	client := NewServer("192.168.24.80", 9000)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond)
		res, err := client.Request("/data/www/upload/test.php", data)
		if err != nil {
			fmt.Println(err.Error())
		}
		body := client.ParseBody(string(res))
		fmt.Println(body)
	}
	ch <- 1
}
