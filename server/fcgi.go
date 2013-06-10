package server

import (
	"fcgiclient"
	"strings"
)

type myFcgi struct {
	client *fcgiclient.FCGIClient
	addr   string
	port   int
}

func NewServer(addr string, port int) (c *myFcgi) {
	c = new(myFcgi)
	c.addr = addr
	c.port = port
	return
}

// close fcgi pool
func (c *myFcgi) Close() (err error) {
	return nil
}

func (c *myFcgi) connect(retry bool) (err error) {
	if c.client == nil || retry {
		if c.client != nil {
			c.client = nil
		}
		c.client, err = fcgiclient.New(c.addr, c.port)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *myFcgi) Request(file string, data map[string]string) (retout []byte, err error) {
	err = c.connect(true)
	if err != nil {
		return nil, err
	}
	reqParams := ""
	for k, v := range data {
		reqParams += "&" + k + "=" + v
	}

	env := make(map[string]string)
	env["REQUEST_METHOD"] = "GET"
	env["SCRIPT_FILENAME"] = file
	env["REMOTE_ADDR"] = "127.0.0.1"
	env["SERVER_PROTOCOL"] = "HTTP/1.1"
	env["QUERY_STRING"] = reqParams

	retout, _, err = c.client.Request(env, reqParams)
	return
}

func (c *myFcgi) ParseBody(s string) string {
	body := ""
	data := strings.Split(s, "\r\n\r\n")
	if len(data) > 1 {
		for _, v := range data[1:] {
			body += v
		}
	}
	return body
}
