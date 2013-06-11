#说明

Go语言练习项目。根据工作中的项目需求进行的实践。

**从redis中读取队列数据发送给后端的php fastcgi处理。**

该项目中使用到了redis,fastcgi客户端，使用了json,gorutine和signal信号处理，采用面对对象编程，基本把go中的基础知识都用了一遍。

##注意
其中使用的fcgiclient，需要自己先创建一个fcgiclient项目，项目中就一个文件，请点击下载：
fcgiclient.go https://gist.github.com/9466/5743027

##问题

1. 多并发时，系统提示 too many open files
2. fastcgi服务器上，出现大量 TIME_WAIT

第一个问题是因为自己的客户端机器限制了open files 

可以通过 

```
ulimit -a
```

查看自己的限制，MAC OS默认为256，修改：

```
ulimit -n 1024
```

第二个问题是TCP本身的机制问题，可以通过配置可以大大减少TIME_WAIT的数量：

```
vi /etc/sysctl.conf
```
增加如下几行：

```
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 1
net.ipv4.ip_local_port_range = 1024 65000
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.tcp_max_tw_buckets = 5000
```

具体解释参考：http://hi.baidu.com/dmkj2008/item/9aa9ea82c3947e5927ebd946

```
