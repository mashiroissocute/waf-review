## OpenResty是什么？

- 兼具开发效率和性能的服务端开发平台
- 基于NGINX实现，但是适用范围远远的超过了NGINX的基于本地文件系统的静态文件服务器功能和基于缓存及负载均衡的反向代理功能。
- 核心是lua-ngx-module （c语言编写的nginx扩展模块），该模块将LuaJIT嵌入NGINX。使用lua是满足服务器的开发需求。

## OpenResty和Nginx和lua-nginx-module和lua-resty-*之间的关系是什么？
- Nginx作为底层网络库，提供nginx最主要的两个功能（静态文件服务器和反向代理）。
- lua-nginx-module和nginx的其他c模块一样，将luajit嵌入nginx，提供api控制nginx的运行。
- `lua-resty-*`库包含了openresty官方的和第三方的resty库，这些库使用lua编写，提供例如kafka，mysql等功能。`lua-resty-*`库最主要的是基于cosocket实现了同步非阻塞编程，提供的api一般来说比lua或者lua-nginx-module提供的api要更加高效

## OpenResty干什么？
- 使用OpenResty 来构建业务系统的比例并不高，因为OpenResty没有像 Python（Django）、Go（Gin）那样有成熟的 Web 框架和生态。
- 使用者大都用OpenResty 来处理入口流量，例如在API网关、CDN和WAF中应用。

## OpenResty的特性：
- 详尽的文档和完备的测试框架
- `同步非阻塞的编程模式 `
	- 同步与异步 ： 代码执行顺序，是否上一个代码块执行完毕之后才是开始执行下一个代码块。
	- 阻塞与非阻塞：进程或线程在执行的过程中是否因为I/O等因素而挂起。 非阻塞是C10K等高并发的关键。
- `动态` ： OpenResty使用脚本语言Lua来控制逻辑。因此可以动态的控制路由，上游，证书，请求，响应等，主要是因为这些模块的数据可以动态变化。NGINX在配置文件拟定后，边无法控制运行时逻辑。改变配置文件后需要reload才会生效。

## 为什么推荐利用系统包工具安装？
- OpenResty在OpenSSL上打了自己的补丁并维护了自己的LuaJit分支，自己源码安装需要的步骤和出错的概率会增大。


## 为什么推荐安装前切换系统包的仓库？
- 官方库一般并不维护第三方提供的包，而OpenResty中自己维护了luajit等组件，这种第三方包并不会被官方的库保存。因此最好是将Openresty加入到仓库后安装

## OpenResty的学习重点：
- 同步非阻塞的编程模式；
- 不同阶段的作用；
- LuaJIT 和 Lua 的不同之处；
- OpenResty API 和周边库；
- 协程和 cosocket；
- 单元测试框架和性能测试工具；
- 火焰图和周边工具链；
- 性能优化。

## OpenResty启动和重载
- 启动openresty -p /usr/local/openresty/  -c conf/nginx.conf  （-p 指定前缀目录 -c 指定配置文件）
- 重载openresty -p /usr/local/openresty/  -s reload
- openresty其实是nginx bin的软连接

## Lua 代码内容的变更，需要重启 OpenResty 服务才会生效，这样显然不方便调试，那么有没有什么即时生效的方法呢？
- 因为lua是脚本语言，nginx.conf中使用`* _by_lua_file`时，使用lua_code_cache off可以避免修改了lua后需要reload nginx的操作。但是每个ngx_lua处理的请求将运行在一个独立的Lua VM实例，非常影响性能。使用lua_code_cache on， Lua 代码在第一个请求时会被加载，并默认缓存起来。所以在你每次修改 Lua 源文件后，都必须重新加载 OpenResty 才会生效。会将lua代码缓存到内存，性能高但是修改了lua代码也需要reload nginx。
- 语法: lua_code_cache on | off
- 使用的上下文：http, server, location, location if
- 作用：lua_code_cache是nginx_lua模块的一条指令。它为 `*_by_lua_file`(如 set_by_lua_file , content_by_lua_file) 这些指令以及Lua模块, 开启或关闭Lua代码缓存。如果关闭，每个ngx_lua处理的请求将运行在一个独立的Lua VM实例里，0.9.3版本后有效. 所以 set_by_lua_file, content_by_lua_file,access_by_lua_file, 等等指令引用的Lua文件将不再缓存到内存， 并且所有Lua模块每次都会从头重新加载. 这样开发者就可以避免改代码然后reload nginx的操作。但是, 那些直接写在 nginx.conf 里的代码比如由 set_by_lua, content_by_lua, access_by_lua, and rewrite_by_lua 指定的代码不会在你编辑他们时实时更新，因为只有发送HUP信号通知Nginx才会正确重新加载Nginx的config文件。
- -生产环境下千万别关闭Lua代码缓存，只能用在开发模式下，因为对性能有十分大的影响（每次IO读取和编译Lua代码消耗很大， 简单的hello world都会慢一个数量级）

## nginx是如何找到lua文件的目录的？
- lua-package-path配置

## `*-nginx-module` 和 `lua-resty-*` 的区别
- 前者是c语言开发的nginx模块，后者是lua语言开发的lua库
- 前者需要编译进nginx，其中最重要的是lua-nginx-module，可以让lua嵌入nginx中，是所有lua-resty-*库运行起来的基础
- 二者可以实现相同的功能，比如访问mysql和redis等。但是基本是后者更具有优势
- lua-resty-* 运行环境是lua，因此可移植性比较强
- 在不使用OpenResty的时候，只基于nginx + lua-nginx-module仍然可以自己去下载并运行lua-resty-*库

## OpenResty包管理工具
- OpenResty有官方的`*-nginx-module`（20+） 和 `lua-resty-*`（18）包。
- 同时，也存在大量的第三方包，比如发起http请求，与KAFKA通信等。不应该使用lua里面的包来实现这些功能，而应该使用lua-resty-*来满足，因为lua里的包可能是阻塞的，而openresty尽量使用非阻塞的代码。
- OPM包管理工具，这是OpenResty自带的包管理工具，在bin目录下自动安装。能够下载`lua-resty-*`库，并整理库的安装目录。
- LUAROCKS需自行安装，可以制作自己的项目包并上传。
- AWESOME-RESTY最全的包管理。


## OPM项目
- templet.lua是根据tt2文件自动生成的
- lua默认是全局变量，容易冲突，不好查问题
- 局部变量效率更高
- lua主要是面向过程的编程，而非面向对象的编程
- https://github.com/openresty/opm