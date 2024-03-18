## SRCACHE
这是OpenResty提供的一个nginx第三方模块，默认不编译进入nginx。通过--add-module来编译进NGINX。

## srcache_fetch
**context:** http, server, location, location if

**phase:** post-access
This directive registers an access phase handler that will issue an Nginx subrequest to lookup the cache.
When the subrequest returns status code other than  `200` , than a cache miss is signaled and the control flow will continue to the later phases including the content phase configured by [ngx_http_proxy_module](http://nginx.org/en/docs/http/ngx_http_proxy_module.html), [ngx_http_fastcgi_module](http://nginx.org/en/docs/http/ngx_http_fastcgi_module.html), and others. If the subrequest returns  `200 OK` , then a cache hit is signaled and this module will send the subrequest's response as the current main request's response to the client directly.
This directive will always run at the end of the access phase, such that [ngx_http_access_module](http://nginx.org/en/docs/http/ngx_http_access_module.html)'s [allow](http://nginx.org/en/docs/http/ngx_http_access_module.html#allow) and [deny](http://nginx.org/en/docs/http/ngx_http_access_module.html#deny) will always run before this.
You can use the [srcache_fetch_skip](https://github.com/openresty/srcache-nginx-module#srcache_fetch_skip) directive to disable cache look-up selectively.

**meaning:**

- 在access阶段发起子请求，主请求hold住。如果子请求响应200，则使用子请求响应作为主请求响应。如果子请求响应大于200，则主请求继续往下执行。一般该请求用户取缓存
- 指令开关：
```
	srcache_fetch_skip 1 srcache_fetch指令生效
	srcache_fetch_skip 0 跳过srcache_fetch指令
```

## srcache_store 
**context:** http, server, location, location if

**phase:** output-filter
This directive registers an output filter handler that will issue an Nginx subrequest to save the response of the current main request into a cache backend. The status code of the subrequest will be ignored.
You can use the [srcache_store_skip](https://github.com/openresty/srcache-nginx-module#srcache_store_skip) and [srcache_store_max_size](https://github.com/openresty/srcache-nginx-module#srcache_store_max_size) directives to disable caching for certain requests in case of a cache miss.
Since the  `v0.12rc7`  release, both the response status line, response headers, and response bodies will be put into the cache. By default, the following special response headers will not be cached:

- Connection
- Keep-Alive
- Proxy-Authenticate
- Proxy-Authorization
- TE
- Trailers
- Transfer-Encoding
- Upgrade
- Set-Cookie

You can use the [srcache_store_pass_header](https://github.com/openresty/srcache-nginx-module#srcache_store_pass_header) and/or [srcache_store_hide_header](https://github.com/openresty/srcache-nginx-module#srcache_store_hide_header) directives to control what headers to cache and what not.

**meaning:**

- 在filter阶段发起子请求，且不阻塞主请求流程。子请求的响应会被忽略。子请求通过`ngx.req.get_body_data()`当前请求的响应。一般该子请求用于更新缓存。
- 指令开关：
```
	srcache_store_skip 1 srcache_store指令生效
	srcache_store_skip 0 跳过srcache_store指令
```