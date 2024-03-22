## 实例 
strace查看nginx端口冲突导致reload没有成功的问题

- waf nginx tencentcloudwaf.conf中默认监听了127.0.0.1:9999端口。
用户加入了9999端口配置，然后reload。导致reload失败。

- 通过strace attach nginx master pid查看原因：
``` 
socket(PF_INET, SOCK_STREAM, IPPROTO_IP) = 106
setsockopt(106, SOL_SOCKET, SO_REUSEADDR, [1], 4) = 0
ioctl(106, FIONBIO, [1])                = 0
bind(106, {sa_family=AF_INET, sin_port=htons(9999), sin_addr=inet_addr("0.0.0.0")}, 16) = -1 EADDRINUSE (Address already in use)
write(4, "2022/07/11 15:29:38:624 [emerg] "..., 100) = 100
close(106)                              = 0
nanosleep({0, 500000000}, NULL)         = 0
```
`bind(106, {sa_family=AF_INET, sin_port=htons(9999), sin_addr=inet_addr("0.0.0.0")}, 16) = -1 EADDRINUSE (Address already in use)`
系统调用bind失败。

## strace 详解
strace常用来跟踪进程执行时的系统调用和所接收的信号。进程不能直接访问硬件设备，当进程需要访问硬件设备(比如读取磁盘文件，接收网络数据等等)时，必须由用户态模式切换至内核态模式，通过系统调用访问硬件设备。strace可以跟踪到一个进程产生的系统调用,包括参数，返回值，执行消耗的时间。

## 跟踪可执行程序
`strace ./nginx`
会直接运行nginx并输出系统调用信息。由于已经有一个nginx在运行了，所以大量的端口占用错误：
``` 
write(8, "2022/07/11 16:48:50 [emerg] 3861"..., 98) = 98
write(2, "nginx: [emerg] bind() to 0.0.0.0"..., 74nginx: [emerg] bind() to 0.0.0.0:8089 failed (98: Address already in use)
) = 74
close(11)                               = 0
socket(AF_INET, SOCK_STREAM, IPPROTO_IP) = 11
setsockopt(11, SOL_SOCKET, SO_REUSEADDR, [1], 4) = 0
ioctl(11, FIONBIO, [1])                 = 0
bind(11, {sa_family=AF_INET, sin_port=htons(8088), sin_addr=inet_addr("0.0.0.0")}, 16) = -1 EADDRINUSE (Address already in use)
write(8, "2022/07/11 16:48:50 [emerg] 3861"..., 98) = 98
write(2, "nginx: [emerg] bind() to 0.0.0.0"..., 74nginx: [emerg] bind() to 0.0.0.0:8088 failed (98: Address already in use)
) = 74
close(11)                               = 0
socket(AF_INET, SOCK_STREAM, IPPROTO_IP) = 11
setsockopt(11, SOL_SOCKET, SO_REUSEADDR, [1], 4) = 0
ioctl(11, FIONBIO, [1])                 = 0
bind(11, {sa_family=AF_INET, sin_port=htons(443), sin_addr=inet_addr("0.0.0.0")}, 16) = -1 EADDRINUSE (Address already in use)
write(8, "2022/07/11 16:48:50 [emerg] 3861"..., 97) = 97
write(2, "nginx: [emerg] bind() to 0.0.0.0"..., 73nginx: [emerg] bind() to 0.0.0.0:443 failed (98: Address already in use)
```

## 跟踪正在运行的程序
`strace -p pid`
然后发送reload到nginx，可以看到reload过程的系统调用


## 引用
https://linuxtools-rst.readthedocs.io/zh_CN/latest/tool/strace.html