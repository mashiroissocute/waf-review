## nginx http 请求处理模块的11个阶段

- 1. `post read`  realip
- 2. `server rewrite` rewrite, return, error page
- 3. `find config`
- 4. `rewrite` rewrite
- 5. `post rewrite`
- 6. `pre access` limit conn, limit req
- 7. `access` auth_basic, access,auth_request
- 8. `post access`
- 9. `pre content` try files, mirror
- 10. `content`  index, autoindex, concat
- 11. `log` access_log
- 注意
阶段一定分先后，前一阶段退出，后一阶段无法执行
同一个阶段中的不同模块，不一定全部执行

## post read 阶段：如何获取到客户端的真实ip 
``` 
        set $userrealip $remote_addr;  // 设置变量user real ip
        if ($http_x_forwarded_for != '') {set $userrealip $http_x_forwarded_for;} //取整个XFF
        if ($userrealip ~ ^([^,]+) ){set $userrealip $1;} //匹配第一个，前的XFF
        proxy_set_header X-Forwarded-For-Pound $userrealip; 
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; //设置反向代理的xFF
		// $proxy_add_x_forwarded_for变量包含客户端请求头中的"X-Forwarded-For"，与$remote_addr用逗号分开，如果没有"X-Forwarded-For" 请求头，则$proxy_add_x_forwarded_for等于$remote_addr。
        if ($http_x_real_ip != '') {set $userrealip $http_x_real_ip;} //取real IP
        proxy_set_header X-Real-IP $userrealip; //设置反向代理的realip
```

## server rewrite / location rewrite阶段： 常用指令 rewrite , return ，error_page ,  if

## find config阶段：
location匹配规则： 
- 常规 什么都不加
- := 
- ^~ 正则，且不再正则
- ~正则
- @内部location跳转

location匹配顺序：
- 1.精准匹配
- 2.^~
- 3.正则
- 4.最长前缀

## pre access 阶段： 1.如何限制用户并发连接数，2.如何限制用户一段时间内的处理请求数
``` 
limit_conn
limit_req // leaky bucket算法
```

## access阶段：
```
access 模块 使用allow或者deny
auth_basic 模块 使用密码登陆，nginx本地存储密码
auth_request 模块 使用上游服务验证客户身份
```

## pre content阶段：
```
try file 模块 使用本地文件回应
mirror 模块 发起子请求，以拷贝流量
```

## content阶段：
```
index 
autoindex
concat
```


# http响应过滤模块
```
copy_filter
psotpone_filter
header_filter
write_filter
```


# 变量
```
//定义变量
map
split_client
geo
```











