package main

import (
	"encoding/json"
	"fmt"
	"github.com/9466/goconfig"
	"log"
	"os"
	"os/signal"
	"queue/client"
	"queue/server"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type QueueRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func queue(i int64, c *goconfig.ConfigFile, ch chan int64) {
	// start rutine
	logger.Println("rutine ", i+1, " start successed.")
	log.Println("rutine ", i+1, " start successed.")
	// client: redis config
	queueName, err := c.GetString("common", "queueName")
	errHandle(err)
	redisHost, err := c.GetString("redis", "host")
	errHandle(err)
	redisPort, err := c.GetString("redis", "port")
	errHandle(err)
	//redisDb, err := c.GetInt64("redis", "db")
	//errHandle(err)
	// server: fcgi config
	fcgiHost, err := c.GetString("fcgi", "host")
	errHandle(err)
	fcgiPort, err := c.GetInt64("fcgi", "port")
	errHandle(err)
	fcgiPath, err := c.GetString("fcgi", "path")
	errHandle(err)
	// init connect
	client := client.NewClient(redisHost+":"+redisPort, 0, 0)
	server := server.NewServer(fcgiHost, int(fcgiPort))

	for shutdown == false {
		time.Sleep(100 * time.Millisecond)
		// get Content
		val, err := client.Pop(queueName)
		if err != nil || len(val) < 1 {
			// redis red no key, continue, not log
			continue
		}
		fmt.Println("rutine ", i+1, " got value: ", string(val))

		// send to Server
		data := make(map[string]string)
		data["message"] = string(val)
		res, err := server.Request(fcgiPath, data)
		if err != nil || len(res) < 1 {
			//logger.Println("queue ", i+1, " error: ", err.Error())
			log.Println("queue ", i+1, " error: ", err.Error())
			// if failed, retry, back into queue
			client.Push(queueName, val)
			continue
		}

		// parse body
		body := server.ParseBody(strings.TrimSpace(string(res)))
		fmt.Println(body)
		// parse json
		var bodyJson QueueRes
		json.Unmarshal([]byte(body), &bodyJson)
		if bodyJson.Code != 0 {
			//logger.Println("queue ", i+1, " error: ", bodyJson.message)
			log.Println("queue ", i+1, " error: ", bodyJson.Message)
			// if failed, retry, back into queue
			client.Push(queueName, val)
			continue
		}
		fmt.Println(bodyJson.Code, bodyJson.Message)
	}
	// shutdown
	client.Close()
	server.Close()
	ch <- 1
	logger.Println("rutine ", i+1, " shutdown.")
	log.Println("rutine ", i+1, " shutdown.")
}

// custom log
var logger *log.Logger

// receive os signal for shutdown
var shutdown = false

func main() {
	var conFile string = "queue.conf"
	var i int64

	// init config
	c, err := goconfig.ReadConfigFile(conFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// init log
	logFile, err := c.GetString("log", "logFile")
	if err != nil {
		log.Fatalln(err.Error())
	}
	w, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer w.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	logger = log.New(w, "", log.Ldate|log.Ltime|log.Llongfile)

	gnum, err := c.GetInt64("common", "rutineNum")
	errHandle(err)

	ch := make(chan int64, gnum)

	runtime.GOMAXPROCS(int(gnum))
	for i = 0; i < gnum; i++ {
		go queue(i, c, ch)
	}

	// trap signal
	sch := make(chan os.Signal, 10)
	signal.Notify(sch, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT,
		syscall.SIGHUP, syscall.SIGSTOP, syscall.SIGQUIT)
	go func(ch <-chan os.Signal) {
		sig := <-ch
		log.Print("signal recieved " + sig.String())
		if sig == syscall.SIGHUP {
			log.Println("queue restart now...")
		}
		shutdown = true
	}(sch)

	// wait channel
	for i = 0; i < gnum; i++ {
		<-ch
	}
}

func errHandle(e error) {
	if e != nil {
		logger.Fatalln(e.Error())
	}
}
