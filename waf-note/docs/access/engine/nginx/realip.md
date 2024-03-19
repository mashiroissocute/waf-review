## 常见的clientIp获取方式

### TCP连接四元组，可以拿到直接上游的地址
remote_addr

### X-Forworaded-For / X-Real-IP
假设链路为：
client -->  cdn(proxy1) --> (clb)proxy2 --> (nginx)proxy3 --> server
#### 数据来源 
一般客户端不会设置XFF XRI header，只有代理才会设置XFF XRI header。
代理设置XFF有两种方式：

- proxy_set_header X-Forwarded-For $remote_addr;
- proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

proxy_add_x_forward_for =  xff + remote_addr https://blog.csdn.net/bigtree_3721/article/details/72820594

代理设置X-Real-IP
一般将该ip设置为自己认为的clientIP，见waf的配置

#### XFF伪造的问题
https://www.cnblogs.com/skychx/p/X-Forwarded-For-get-real-IP.html
- 只有当前代理可控，为当前代理设置可信任的IP，找xff中该ip前一个
- 链路代理可控，入口代理这是xff为remoteaddr，后续代理追加
- 代理个数清晰，从右往左计算ip

### nginx Real ip模块 可以修改remote addr
（WAF没有使用，因为没有获取已知可以信任的上游IP）
在HTTP的POST Read阶段，可以直接修改Remoteaddr

- set_real_ip_from 设置信任的ip，只有这个ip发过来内容时，才会修改addr
- real_ip_header 从哪里去取addr ， 默认x-real-ip，可选XFF的最后一位（信任ip的上游）等。
_ real_ip_recursive 当XFF最后一个ip和remote addr地址相同的时候，忽略这个ip，再往左取ip

## WAF
如何获取ClientIp 以及如何设置XFF和XRI
``` 
		set $userrealip $remote_addr; //1.设置为remote_addr （这里的remote_addr可以通过real ip模块修改）
        if ($http_x_forwarded_for != '') {set $userrealip $http_x_forwarded_for;}
        if ($userrealip ~ ^([^,]+) ){set $userrealip $1;} //2.如果有XFF，取XFF第一位
        if ($http_x_real_ip != '') {set $userrealip $http_x_real_ip;} //3.如果有XRI，取XRI
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; //追加XFF，（XFF + remoteaddr）
		proxy_set_header X-Real-IP $userrealip; //设置XRI为，自己认为的clientIP
```