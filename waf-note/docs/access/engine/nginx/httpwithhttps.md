## 如何同个端口提供http和https服
### 配置文件
无法同一个端口listen ssl和非ssl共存。无法同时提供http和https服务
### error page
配置文件listen ssl
访问使用http时，nginx返回错误页面：
``` 
<html>
<head><title>400 The plain HTTP request was sent to HTTPS port</title></head>
<body>
<center><h1>400 Bad Request</h1></center>
<center>The plain HTTP request was sent to HTTPS port</center>
<hr><center>nginx/1.21.5</center>
</body>
</html>
```
拿着400 The plain HTTP request was sent to HTTPS port去nginx源码查看：
``` 
static char ngx_http_error_497_page[] =
"<html>" CRLF
"<head><title>400 The plain HTTP request was sent to HTTPS port</title></head>"
CRLF
"<body>" CRLF
"<center><h1>400 Bad Request</h1></center>" CRLF
"<center>The plain HTTP request was sent to HTTPS port</center>" CRLF
;
```
其实是返回的497 code。
那么使用error page 497 就可以实现重定向到https协议。例如:
``` 
listen 8000 ssl;
error_page 497 https://$host$uri;
```





