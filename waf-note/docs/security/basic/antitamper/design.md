## WAF防篡改功能
旨在缓存静态资源`[]string{".htm", ".html", ".txt", ".shtml", ".js", ".css", ".jpg", ".png"}`，即使源站页面被篡改，waf仍然返回未修改的静态资源。

## CLOUDAPI
### 查询防篡改规则
`select id,name,domain,uri,protocol,status from tb_waf_custom_cache`

### 添加/修改/删除 防篡改规则

- Mongo ：
collection : domain-version-conf
document : appid + host
operation : antitamper(1) , antitamper_ts(timestamp)
- Mysql :
tb_waf_custom_cache
- redis-stream-master
1.xadd 追加消息
2.WafConfStream 消息信道
3.keyValues {
 	- appid : appid
	- host : domain
	- action:  delete ｜ add  ｜  edit
	- module :  antitamper }
	
### 刷新服务器防篡改页面到waf
- redis-stream-master
1.xadd 追加消息
2.WafConfStream 消息信道
3.keyValues {
 	- appid : appid
	- host : domain
	- action:  refresh
	- module :  antitamper }

## redis-agent-stream
A.Hset kv到各个地域的slave-redis
``` 
pipe = redis_conn.pipeline()
pipe.multi()
pipe.hset( redis_key,"antitamper_version",ver_str )
pipe.hset( redis_key,"antitamper",conf_str )  // mysql中的数据
pipe.execute()
```

B.通知redis-slave变更
1.xadd 追加消息
2.slave-waf-conf 消息信道
3.keyValues {
- appid : appid
- host : domain
- action:  delete ｜ add ｜ edit ｜ refresh
- module :  antitamper }
	
	
## Engine & NGINX
### init_worker（Engine）
load_policy --> load_config_from_redis5
``` 
local antitamper_ver, antitamper_policy = _load_rule(red, domain_key, "antitamper") //从slave-redis中获取antitamper的内容
domain_policy.antitamper_set = cjson.decode(antitamper_policy)
```

### rewrite （NGINX）
定义变量
```
set $fetch_skip 1;
set $store_skip 1;
set $cache_md5_key "";
```
### rewrite（Engine）
antitamper.check(req_obj, user_policy)
``` 
local antitamper_rule = user_policy.antitamper_set
遍历antitamper_rule，匹配当前host以及当前url的防篡改规则。

//没有匹配上则退出函数；
此时fetch_skip和store_skip都是1，跳过srcache_fetch和srcache_store

//匹配上，说明有防篡改的功能：
为当前请求的uri生成主和子md5，md5-key = main_md5 + ‘-’ + sub_md5
main_md5, sub_md5, md5_key  = gen_md5(path, fieldStr)

//设置cache_md5_key
ngx.var.cache_md5_key = md5_key

//如果当前请求的md5在共享内存或者redis中
fetch_skip= 0 store_skip = 1 在access阶段执行srcache_fetch，在filter阶段不执行srcache_store

//如果当前请求的md5不共享内存或者redis中
fetch_skip= 1 store_skip = 0 在access阶段不执行srcache_fetch，在filter阶段执行srcache_store
```

### ACCESS （NGINX）
srcache
``` 
srcache_fetch_skip $fetch_skip;
srcache_fetch GET /cloudwaf_handle_cache key=$cache_md5_key;
```
### ACCESS（Engine）
fetch_cache ： 子请求 从共享内存或者redis中获取cache，返回
``` 
local function fetch_cache(key)
    local cache_content = antitamper_shm:get(key)
    if not cache_content then
        local main_md5, sub_md5 = split_key_md5(key)
        if main_md5 ~= nil then
            -- 切换到redis5
            local redis = redis_access:redis5_new()   --get the content from redis
            local res,err = redis:hget(main_md5, sub_md5) --从redis里面获取cache内容，估计这里是防止共享内存淘汰机制，把cache淘汰了，所以用了redis来存储所有的数据
            if res then
                antitamper_shm:set(key,res,expire_time) -- 存入到共享内存中
                cache_content = res
            else
                if err then
                    waf_logger.err("Error in fetching cache from redis with err:"..err)
                end
            end
        end
    end
    if cache_content then
        ngx.print(cache_content) -- 返回200 + 内容
    else   --failed to fetch the cache content
        return ngx.exit(502) -- 502是子请求的返回，主请求检测到502则会往继续content阶段
    end
end
```

### Filter（NGINX）

``` 
srcache_store_skip $store_skip;
srcache_store PUT /cloudwaf_handle_cache key=$cache_md5_key;
```

### Filter（Engine）
store_cache 子请求，缓存当前请求响应到共享内存和redis中
``` 
local function store_cache(key)
    local cache_content = ngx_req_get_body_data()
    local main_md5, sub_md5 = split_key_md5(key)
    if main_md5 ~= nil then
        local redis = redis_access:redis5_new()
        --first,store the cache content into redis
        local ok,err = redis:hset(main_md5, sub_md5, cache_content)
        if not ok then
            local err_msg = "Failed to store the cache into redis"
            if err then
                err_msg = err_msg.." with err:"..err
            end
            waf_logger.err(err_msg)
        end
        --second,store the cache content into shared-dict
        antitamper_shm:set(key,cache_content,expire_time)
    else
        waf_logger.err("Failed to store the cache into redis for err cache key: "..key)
    end
end
```
