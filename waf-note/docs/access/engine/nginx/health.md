## 配置 ：

- waf `check interval=3000 rise=2 fall=5 timeout=5000` 
- 标准 ：https://cloud.tencent.com/developer/article/1397237
- 
``` 
Syntax: check interval=milliseconds [fall=count] [rise=count] [timeout=milliseconds] [default_down=true|false] [type=tcp|http|ssl_hello|mysql|ajp] [port=check_port]
Default: 如果没有配置参数，默认值是：interval=30000 fall=5 rise=2 timeout=1000 default_down=true type=tcp
Context: upstream

该指令可以打开后端服务器的健康检查功能。
指令后面的参数意义是：
interval：向后端发送的健康检查包的间隔。
fall(fall_count): 如果连续失败次数达到fall_count，服务器就被认为是down。
rise(rise_count): 如果连续成功次数达到rise_count，服务器就被认为是up。
timeout: 后端健康请求的超时时间。
【default_down】: 设定初始时服务器的状态，如果是true，就说明默认是down的，如果是false，就是up的。默认值是true，也就是一开始服务器认为是不可用，要等健康检查包达到一定成功次数以后才会被认为是健康的。
type：健康检查包的类型，现在支持以下多种类型
tcp：简单的tcp连接，如果连接成功，就说明后端正常。
ssl_hello：发送一个初始的SSL hello包并接受服务器的SSL hello包。
http：发送HTTP请求，通过后端的回复包的状态来判断后端是否存活。
mysql: 向mysql服务器连接，通过接收服务器的greeting包来判断后端是否存活。
ajp：向后端发送AJP协议的Cping包，通过接收Cpong包来判断后端是否存活。
port: 指定后端服务器的检查端口。你可以指定不同于真实服务的后端服务器的端口，比如后端提供的是443端口的应用，你可以去检查80端口的状态来判断后端健康状况。默认是0，表示跟后端server提供真实服务的端口一样。该选项出现于Tengine-1.4.0。
```