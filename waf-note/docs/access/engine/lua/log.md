## 日志类型

- socket_heart 按照域名维度统计心跳日志，在initworker阶段启动，每10秒上报一次统计的内容。
- socket_attack 攻击日志，在initworker阶段启动，每10秒上报一次统计的内容。
- socket_access 记录访问日志，每次请求到来都发送访问日志。
- socket_response 记录响应日志，每次请求到来都发送响应日志。

## socket定义 
``` 
init阶段定义了用于发送日志的udp socket
waf.socket.heart = socket_heart {ip(127.0.0.1),type(udp/tcp),port(2003), }
...
waf.socket.access = socket_access
waf.socket.attack = socket_attack
```

## 心跳日志：

``` 
# ----1 init work
heartbeat : #(10s 循环调用一次，10s上报一次统计数据，10s后就清除)
{
    waf_heart.new() : {
			初始化 info
			local info = {}
			-- 心跳包过大，分为两个解决
			local info2 = {}
			local domain_num = 0
		}
    
	local heart_info,heart_info2,domain_num = waf_heart.get_heart_info() : {
			return info,info2,domain_num
		}
	waf_heart.empty_info()
	socket_heart:log(cjson.encode(heart_info).."\n")
 }
}

# ----2 log statistic

waf_heart.statistic(req_obj,action_flag,user_policy,cc_status) {
			# 写入 info / info2 / domain_num
			# req_obj: ngx.ctx.request
			# action_flag: ngx.ctx.action_flag
			# user_policy: ngx.ctx.user_policy
			# cc_status: string ("cc" / "bot" / “acl”)
			
			会统计当个域名10内的:
			收包长度（ub）/
			回包长度（db）/
			域名访问次数（ac）/
			nginx 4xx，5xx次数（4x，5x）/
			原站 4xx，5xx次数（u4，u5）/
			攻击次数(at)
			bot cc acl次数（b，cc）/
			回源的QPS,这段时间内回源的次数（acu，回源的次数）
			...	
			
			写入到info中，info写不下则写到info2中
			info在worker间是不共享的。
			所有肯定还需要有服务再聚合一下。
}

```

## 攻击日志
``` 
# ------0 关于攻击日志的shared dict
local attack_count = ngx.shared.attack_count
local deny_count = ngx.shared.deny_count
local attack_detail = ngx.shared.attack_detail

# ------1 init worker每10s调用一次upload_attack_log
upload_attack_log 发送attack_detail和attack_count
attack_detail里面的内容按照不同的攻击和维度聚合
-- ip黑名单 采用 custom::ip-103.13.247.119 的方式聚合
-- 网页防篡改 key = ip + full_uri
-- api安全 采用 域名/ip/rule_id 的方式聚合
-- 威胁情报 采用 custom::ip-103.13.247.119 的方式聚合
-- 自动封禁 采用 custom::ip-103.13.247.119 的方式聚合
-- 地域封禁 采用 域名/城市 的key
-- 其他带有ruleid的攻击按照 域名/ip/rule_id 方式聚合


# ------2 log 阶段调用handle_attacklog函数上传本次请求的攻击日志
填充攻击信息 在access阶段，各个类型的检测，发现攻击则忘ngx.ctx.attack_log里面加入攻击信息
{
rule
detect_type : 检测模块 areaban ， webshell ..
payload 
location
req_obj ：req_obj
action ： deny ， log ..
}
大部分的检测都是串行的，一般attack_log里面不会有多条日志

- 如果是自定义策略，CC，地域封禁，bot拦截，IP黑名单，网页防篡改，api安全等模块检测的攻击日志，聚合记录在ngx.ctx.attack_log中，周期发送。
- 其他的模块检测的攻击日志，直接通过socket_attack发送出去。
```

## 访问（请求）日志：
``` 
handle_requestlog(req_obj) {
# 主要是讲req_obj和ngx.var和reqstrace（经过的检测模块）里面的信息发送出去。
# 单次请求就会发送，没有周期概念
socket_access:log(cjson.encode(request_log_table).."\n")		
}
```

## 响应日志：
``` 
handle_requestlog(req_obj) {
# ngx.ctx.response.var里面的信息发送出去。
# 单次请求就会发送，没有周期概念
socket_response:log(cjson.encode(resp).."\n")		
}
```