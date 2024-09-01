## 响应码


### 301 Moved permanent
nginx作为代理服务器，请求的host匹配到的server中包含永久重定向命令。

出现场景1:
``` 
server {
    listen       80;
    server_name  xxx.com;
    rewrite ^/(.*) https://$host/$1 permanent;
}
```

### 302 Moved Temporarily 
nginx作为代理服务器，请求的host匹配到的server中包含临时重定向命令。

出现场景1:https强制跳转
``` 
server {
    listen       80;
    server_name  xxx.com;
    rewrite ^/(.*) https://$host/$1 redirect;
}
```
Waf返回：302 Moved Temporarily
请求重定向：curl + -L / 浏览器自动支持


### 400 Bad Request 
客户端错误

nginx作为代理服务器，无法理解当前请求，或者认为请求的参数有误。

出现场景1：没有与sever请求的host匹配，走默认localhost.conf。设置状态码为400。 

Waf返回：TencetWafEngine中，在rewrite阶段发现返回状态码为400时，返回默认4x.html:很抱歉，您提交的请求存在异常，请向网站管理员确认并获取正确的访问方式
	
出现场景2: request header过大（cookie）超出nginx读取header的缓冲区。


### 502 Bad GateWay

nginx作为网关或者代理工作的服务器，尝试向上游服务器发起请求时，收到了无效的响应（收到响应，但无效）。

出现场景：关闭nginx代理的上游服务器，关闭上游服务防火墙。（建立连接，但是收到无效响应）

nginx error log ：connect() failed (111: Connection refused) while connecting to upstream。

Waf返回：“很抱歉，你提交的请求无法正常响应，请联系网站管理员处理”

### 504 GateWay Time-Out 
nginx作为网关或者代理工作的服务器向上游服务器发起请求时，未能在一定的时间内从上游服务器收到响应（超时，没有收到响应），与proxy_read_time/proxy_connect_timeout有联系。

出现场景1：开启nginx代理的上游服务，但是开启上游服务防火墙，阻止端口。 （无法建立连接，将在proxy_connect_timeout时间后超时返回504）

nginx error log : upstream timed out (110: Connection timed out) while connecting to upstream。

Waf返回：“很抱歉，你提交的请求无法正常响应，请联系网站管理员处理”

出现场景2: 开启nginx代理的上游服务，关闭上游服务器防火墙，但是上游服务处理时间超过proxy_read_timeout。（成功建立连接，但是nginx在proxy_read_timeout时间内没有收到上游服务器的响应，超时返回）

nginx error log : upstream timed out (110: Connection timed out) while reading response header from upstream。

Waf返回：“很抱歉，你提交的请求无法正常响应，请联系网站管理员处理”


## 首部
分为4 种类型的首部字段：通用首部字段、请求首部字段、响应首部字段和实体内容首部字段。
### 1.1通用首部字段
|首部字段名 |说明 |
|:--|
|Cache-Control |控制缓存的行为 |
|Connection |控制不再转发给代理的首部字段、管理持久连接 |
|Date |创建报文的日期时间 |
|Pragma |报文指令 |
|Trailer |报文末端的首部一览 |
|Transfer-Encoding |指定报文主体的传输编码方式 |
|Upgrade |升级为其他协议 |
|Via |代理服务器的相关信息 |
|Warning |错误通知 |

### 1.2请求首部字段
|首部字段名 |说明 |
|:--|
|Accept |用户代理可处理的媒体类型 |
|Accept-Charset |优先的字符集 |
|Accept-Encoding |优先的内容编码 |
|Accept-Language |优先的语言（自然语言） |
|Authorization |Web 认证信息 |
|Expect |期待服务器的特定行为 |
|From |用户的电子邮箱地址 |
|Host |请求资源所在服务器 |
|If-Match |比较实体标记（ETag） |
|If-Modified-Since |比较资源的更新时间 |
|If-None-Match |比较实体标记（与 If-Match 相反） |
|If-Range |资源未更新时发送实体 Byte 的范围请求 |
|If-Unmodified-Since |比较资源的更新时间（与 If-Modified-Since 相反） |
|Max-Forwards |最大传输逐跳数 |
|Proxy-Authorization |代理服务器要求客户端的认证信息 |
|Range |实体的字节范围请求 |
|Referer |对请求中 URI 的原始获取方 |
|TE |传输编码的优先级 |
|User-Agent |HTTP 客户端程序的信息 |

### 1.3响应首部字段
|首部字段名 |说明 |
|:--|
|Accept-Ranges |是否接受字节范围请求 |
|Age |推算资源创建经过时间 |
|ETag |资源的匹配信息 |
|Location |令客户端重定向至指定 URI |
|Proxy-Authenticate |代理服务器对客户端的认证信息 |
|Retry-After |对再次发起请求的时机要求 |
|Server |HTTP 服务器的安装信息 |
|Vary |代理服务器缓存的管理信息 |
|WWW-Authenticate |服务器对客户端的认证信息 |

### 1.4实体首部字段 
|首部字段名 |说明 |
|:--|
|Allow |资源可支持的 HTTP 方法 |
|Content-Encoding |实体主体适用的编码方式 |
|Content-Language |实体主体的自然语言 |
|Content-Length |实体主体的大小 |
|Content-Location |替代对应资源的 URI |
|Content-MD5 |实体主体的报文摘要 |
|Content-Range |实体主体的位置范围 |
|Content-Type |实体主体的媒体类型 |
|Expires |实体主体过期的日期时间 |
|Last-Modified |资源的最后修改日期时间 |


## 什么时候http连接断开

- 长时间未通信，断开
- 重试超过次数，断开
- close，主动断开

## 长连接

HTTP长连接依赖TCP连接不关闭。
所以不要调用close，并且定期发送心跳，保持tcp连接不断开。



## HTTP报文解析
解析请求行：
当读取到足够的数据后，Nginx会尝试解析HTTP请求的请求行。请求行包含请求方法、请求URI和HTTP版本信息。Nginx会根据这些信息判断请求的类型，并进行相应的处理。

解析请求头：
解析完请求行后，Nginx会继续解析请求头。请求头包含多个键值对，描述请求的元数据。Nginx会将解析到的请求头存储在一个内部表示的结构体中，以便后续处理。

解析请求体：
如果请求方法需要请求体（如POST、PUT等），Nginx会在解析完请求头后继续读取并解析请求体。请求体的大小可以通过请求头中的Content-Length字段或Transfer-Encoding字段来确定。


## HTTPS

### http的安全性问题
- 明文传输数据，内容可能被窃听
- 报文完整性没有验证，内容可能被篡改
- 通信方身份没有验证，身份可能遭遇伪装
### TLS
TLS四次握手

[tls](./tls.md)

## CURL

curl https://xx.qcloudwaf.com:443 --resolve xx.qcloudwaf.com -Ikv 

(i v展示头部信息等)

(L 连续跳转)

(使用 -k 或者 –insecure 选项，来忽略签名认证的警告。 这样就可以让curl命令执行不安全的SSL连接，进而去获取数据。)