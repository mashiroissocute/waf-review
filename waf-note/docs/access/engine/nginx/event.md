https://segmentfault.com/a/1190000018359026

```
events { 
    accept_mutex on; #设置网路连接序列化，防止惊群现象发生，默认为on 
    
    multi_accept on; #设置一个进程是否同时接受多个网络连接，默认为off 
    
    use epoll; #事件驱动模型select|poll|kqueue|epoll|resig
    
    worker_connections 1024; #最大连接数，默认为512
}
```

If  `accept_mutex`  is enabled, worker processes will accept new connections by turn. Otherwise, all worker processes will be notified about new connections, and if volume of new connections is low, some of the worker processes may just waste system resources.