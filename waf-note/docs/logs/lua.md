
## init_worker 定义用于日志通信的socket
``` 
waf.socket.heart = socket_heart {ip(127.0.0.1),type(udp/tcp),port(2003), }
...
waf.socket.access = socket_access
waf.socket.attack = socket_attack
```

## 心跳日志：

``` 
init_worker : {

	init_socket : {
		sock_conf_heart.host = sock_ip
		sock_conf_heart.port = port_heart
		sock_conf_heart.sock_type = sock_type
		-- 心跳要保证每次都要上报
		sock_conf_heart.flush_limit = 32
		
		local socket_heart,err = socket:init(sock_conf_heart)
		waf.socket.heart = socket_heart
	}
	
	#	(10s 循环调用一次，10s上报一次统计数据，10s后就清除)
	heartbeat : {
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
	}
}

```

``` 
log : {
	main() {
		waf_heart.statistic(req_obj,action_flag,user_policy,cc_status) {
			# 写入 info / info2 / domain_num
			# req_obj: ngx.ctx.request
			# action_flag: ngx.ctx.action_flag
			# user_policy: ngx.ctx.user_policy
			# cc_status: string ("c" / "b")
			
			会进行统计10内的:
			收包长度（ub）/
			回包长度（db）/
			域名访问次数（ac）/
			nginx 4xx，5xx次数（4x，5x）/
			原站 4xx，5xx次数（u4，u5）/
			攻击次数(at)
			bot cc次数（b，cc）/
			回源的QPS,这段时间内回源的次数（acu，回源的次数）
			
			基本信息有:
			edition（e）/
			ip(ip) /
			time(t) /
			workid(wid) /
		}
	}
	
}

```


## 攻击日志

- 以规则引擎拦截为例：

``` 
access {
	main {
		if mode == 0 or user_policy.force_local_owasp == 1 then
              -- 走原来的规则引擎
              rule_result, action = owasp_detect.check(req_obj, user_policy) {
			  			--和规则进行正则匹配，返回匹配结果，匹配模式
						local result,payload,location, action = rule_check(rule,req_obj){
							result,payload = regex_match(target,cmp_pattern)
						}
						--记录日志
						attack_log_ctx(rule,detect_type,payload,location,req_obj,action){
							temp_log.rule = rule
							temp_log.detect_type = detect_type
							temp_log.payload = payload
							temp_log.location = location
							temp_log.request = req_obj
							temp_log.action = action
							local attack_log = cloudwaf_ctx.get_attacklog() or {}
							attack_log[#attack_log + 1] = temp_log
							cloudwaf_ctx.set_attacklog(attack_log)
						}
					}
			  }
	}

}
```


```
# 每次请求都会操作
log {
	main {
		local attack_log = waf_ctx.get_attacklog(){
			#按照逻辑填充attack_log
			#并传送到socket
			socket_attack:log(cjson.encode(attack_log_table).."\n")
		}
	}

}
```

## 访问日志：

``` 
log {
	if need_report_access_log(user_policy, req_obj) == true then --关闭waf不记录访问日志
            request_log = handle_requestlog(req_obj) {
						       if socket_access then
        socket_access:log(cjson.encode(request_log_table).."\n")		
			}


}

```




