## Nginx和LuaVM的配合方式

- OpenResty 的 master 和 worker 进程中，都包含一个 LuaJIT VM。在同一个nginx进程内的所有协程，都会共享这个 VM，并在这个 VM 中运行 Lua 代码。因此lua协程竞争的是LuaVM，Nginx协程竞争的是cpu。
- Nginx使用协程来处理多个请求，每个请求同样会对应一个lua协程。lua协程的生命周期由nginx的时间机制控制。
- lua协程有yelid和resume接口供nginx来调用。
- lua协程遇到IO后自己yield并注册nginx回调，时间到来时由nginx唤醒。Lua 的协程会与 NGINX 的事件机制相互配合。如果 Lua 代码中出现类似查询 MySQL 数据库这样的 I/O 操作，就会先调用 Lua 协程的 yield 把自己挂起，然后在 NGINX 中注册回调；在 I/O 操作完成（也可能是超时或者出错）后，再由 NGINX 回调 resume 来唤醒 Lua 协程。这样就完成了 Lua 协程和 NGINX 事件驱动的配合，避免在Lua 代码中写回调。
- 如果 Lua 代码中没有 I/O 或者 sleep 操作，比如全是密集的加解密运算，那么 Lua 协程就会一直占用 LuaJIT VM，直到处理完整个请求。

## OpenResty API
- OpenResty的API都在lua-nginx-module中和lua-resty-core中。lua-nginx-module中的api可以直接使用，不需要require，比如ngx.sleep、ngx.config、ngx.timer等。但是lua-resty-core中的api需要require使用，比如 ngx.ssl、ngx.base64、ngx.errlog、ngx.process、ngx.re.split、ngx.resp.add_header、ngx.balancer、ngx.semaphore、ngx.ocsp 这些 API 。
- OpenResty 的 API 主要分为下面几个大类：https://github.com/openresty/lua-nginx-module/#nginx-api-for-lua
```
	- 处理请求和响应；
	- SSL 相关；
	- shared dict；
	- cosocket；
	- 处理四层流量；
	- process 和 worker；
	- 获取 NGINX 变量和配置；
	- 字符串、时间、编解码等通用功能。
```
- 每个API都有可以执行的lua阶段，超出这个阶段执行就会报错。在你使用 API 之前，一定记得要先查阅文档，确定其能否在代码的上下文中使用。
- 所有的API都是非阻塞的。例如A协程调用ngx.sleep，会挂起A协程，此时cpu处理B协程。 A协程的上下文会被保存起来，并由 NGINX 的事件机制来唤醒。

## OpenResty变量

- 变量有三种：全局变量、局部变量、模块变量（是全局变量的一种更清晰表达方式）、跨阶段变量。
- 不推荐定义全局变量
- 模块只会被每个worker加载一次，并且worker间隔离。同个worker的所有请求都共享模块的变量。
- openresty没有做模块读写加锁，所以最好是不要写模块变量，否则出现race导致难以定位的bug。访问模块变量的时候，你最好保持只读，而不要尝试去修改，不然在高并发的情况下会出现 race。这种 bug 依靠单元测试是无法发现的，它在线上偶尔会出现，并且很难定位。因为在对模块变量进行写操作的时候，OpenResty 并不会加锁，这时就会产生竞争，模块变量的值就会被多个请求同时更新。
- 跨阶段变量可以使用ngx.ctx来定义，ngx.ctx是一个lua table。生命周期同该请求保持一致。子请求或者内部重定向可能会导致重新创建ctx或者ctx内容清空。

## 变量的生命周期和作用范围

- shared dict，全部worker且全部请求共享，请求全程有效。
- 全局变量，单worker内全部请求共享，定义之后的阶段有效。一般在init阶段定义全局变量。
- 模块变量，单worker内全部请求共享，定义之后的阶段有效。在第一个请求到来的时候，加载模块，此时为模块变量分配了空间。模块代码并非lua req中的阶段代码，不是每个请求到来都会初始化模块，而只是对模块数据进行读或修改。模块变量存在race的问题。
- 局部变量，单worker内单个请求，无共享，只在当前文件中有效。每个请求到来，都会执行lua req的阶段代码，例如conntent_by_lua，因此在该文件中定义的local变量则是局部变量，仅在当前请求当前文件中是有效的。
- ngx.ctx，单worker内单个请求，请求全程有效 ，数据结构为lua table。
- ngx.var，单worker内单个请求，支持在nginx的c模块和lua代码间共享，请求全程有效 nginx变量，config中需要定义。

## API文档和测试案例

- lua-nginx-module
- 多阅读/t目录的实例
- 多阅读文档


## OpenResty如何处理请求和响应的？

- nginx由静态的配置文件驱动，但是OpenResty由Lua Api驱动，可以提供更多的灵活性和可编程性。
- nginx的内置变量http://nginx.org/en/docs/http/ngx_http_core_module.html#variables , 可以使用lua-nginx-module的api ngx.var.xxx来获取。绝大部分的 ngx.var 是只读的，只有很少数的变量是可写的
- **请求行**：ngx.req.get_method, ngx.req.set_method , ngx.req.set_uri, ngx.req.set_uri_args
- **请求头**：ngx.req.get_header，ngx.req.set_header, ngx.req.clear_header。OpenResty 并没有提供获取某一个
指定请求头的 API，也就是没有 ngx.req.header['host'] 这种形式。如果你有这样的需求，那就需要借助NGINX 的变量 $http_xxx 来实现了，那么在 OpenResty中，就是 ngx.var.http_xxx 这样的获取方式。例如，ngx.var.http_x_real_ip。
- **请求体**：ngx.req.get_body_data() ， OpenResty不会主动读取请求体的内容。对于比较大(超过client_body_buffer_size)的请求体，OpenResty会把内容保存在磁盘的临时文件中。ngx.req.set_body_data和ngx.req.set_body_file可以接受字符串和本地磁盘文件作为输入参数，来改写请求体。
- **响应行**：ngx.exit(ngx.HTTP_BAD_REQUEST)。OpenResty 的 HTTP 状态码中，有一个特别的常量：ngx.OK。当 ngx.exit(ngx.OK) 时，请求会退出当前处理阶段，进入下一个阶段，而不是直接返回给客户端。
- **响应头**：ngx.header["X-My-Header"] = 'blah blah'   ， ngx_resp.add_header("","")
- **响应体**：ngx.say (最后多个换行符), ngx.print


## Shared Dict

- lua-nginx-module提供的api
- 可以用于hash 和 queue两种数据结构
- 基于红黑树实现，性能很好，但也有自己的局限：必须事先在 Nginx 的配置文件中，声明共享内存的大小，并且不能在运行期更改。因为提前存储大小，所以采用lru的方式在淘汰数据。
- shared dict 只能缓存字符串类型的数据，不支持复杂的 Lua 数据类型。这也就意味着，当我需要存放 table 等复杂的数据类型时，我将不得不使用 json 或者其他的方法，来序列化和反序列化，这自然会带来不小的性能损耗。
- 共享字典本身，它对外提供了 20 多个 Lua API，不过所有的这些 API 都是原子操作，你不用担心多个 worker 和高并发的情况下的竞争问题。
- API分类：字典读写类，队列操作类，管理类（查看keys ， 查看空间等）


## cosocket

- lua-nginx-module提供的api
- 可以用于非阻塞的进行网络通信，原理是lua协程遇到网络 I/O 时，它会交出控制权（yield），把网络事件注册到 Nginx 监听列表中，并把权限交给 Nginx；当有Nginx 事件达到触发条件时，便唤醒对应的协程继续处理（resume）。
- cosocket 是各种 `lua-resty-*` 非阻塞库的基础
- 早期的 OpenResty 版本中，如果你想要去与 Redis、memcached 这些服务交互的话，需要使用 redis2-nginx-module、redis-nginx-module 和 memc-nginx-module这些 C 模块. 这些模块至今仍然在 OpenResty 的发行
包中。不过，cosocket 功能加入以后，它们都已经被 lua-resty-redis 和 lua-resty-memcached 替代，基本上没人再去使用 C 模块连接外部服务了。
- cosocket 支持 TCP、UDP 和Unix Domain Socket。
- API： connect ， settimeout ， send ，  receive ， receiveutil ， receiveany ， close
- nginx指令 ： lua_socket_connect_timeout， lua_socket_pool_size， lua_socket_buffer_size等
- 重复的时候，API的优先级高于nginx指令
- 连接池： 在close之前，调用setkeepalive将不用的socket放入连接池，下次调用connect的时候会优先从连接池中取得socket
- 使用范围：归咎于 Nginx 内核的各种限制，cosocket API 在 set_by_lua， log_by_lua， header_filter_by_lua和 body_filter_by_lua中是无法使用的。而在 init_by_lua和 init_worker_by_lua中暂时也不能用，不过 Nginx 内核对这两个阶段并没有限制，后面可以增加对这它们的支持。可以采用timer的机制绕过该限制。

- 了解一下为什么init阶段，不能使用cosocket或者sleep

## timer

- timer接口替代了系统的crontab调用，更加便于维护。
- API: ngx.timer.at , ngx.timer.every , ngx.timer.pendingcount , ngx.timer.runingcount 
- nginx指令:  lua_max_pending_timers 和 lua_max_running_timers这两个指令，来对定时任务数量进行限制。前者代表等待执行的定时任务的最大值，后者代表当前正在运行的定时任务的最大值。因为一旦开始了定时任务就无法停止，这两个命令避免定时任务过多导致资源耗尽。

## 特权进程

- 特权进程具有master一样的权限。
- lua-resty-core中提供了process的API, ngx.process。
- 特权进程在init master中开启，在init worker中工作。例如，定时执行系统命令来发送reload命令给master或者清理日志等。


## ngx.pipe

- 在lua中执行os.exit是阻塞的命令，使用ngx.pipe提供的api可以变成非阻塞的操作。

## 正则

- ngx.re.* lua-resty-core提供的api
- lua_regex_match_limit 指令，限制正则匹配回溯的次数，避免cpu满载