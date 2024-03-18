### 知识扩充
ngx.var 和 ngx.ctx的区别：
- ngx.var需要在nginx.conf中提前set，作为一个变量，可以用ngx.var.xx访问
- ngx.ctx不需要提前set，作为一个lua table访问

ngx.ctx 和 var是针对请求的， waf_status waf_data是针对整个waf的

### rewrite阶段 主要负责添加防护模块，填充ngx.ctx内容

- ngx.ctx.request 
- ngx.ctx.policy
- ngx.ctx.bypass_geoip , ngx.ctx.bypass_custom ...

``` 
-- 1.default server 的location中 set status = 400
   2.commonloaction中定义了, error-page为502 504时重定向到502 504location
-- set $waf_inner_error 1, set $error_status 
1.if ngx.var.waf_inner_error == "1" 
	error_response(error_status) 返回对应的错误页面并ngx.exit()
	
2.初始化trace，跟踪经过的防护模块
ngx.ctx.reqstrace = {
            owasp=0,
            custom_rule=0,
            ip_bw=0,
            api_sec=0,
            jsinject=0,
            bot=0,
            cc=0,
            ai=0,
            antitamper=0,
            antileakage=0,
            areaban=0,
            business_risk=0,
            reputation=0,
            webshell=0,
            autodeny=0,
            threaten=0
        }
		
3.cdnwaf需要鉴权header，失败返回ngx.exit(ngx.HTTP_UNAUTHORIZED) 401


4. 判断bypasswaf和waf开关，关闭则直接退出rewrite阶段，后续操作不继续。

4.load policy，并写入ctx（ctx的生命周期跟随请求）
user_policy = get_host_policy() ==>  user_policy = waf.policy[host]
clb-waf场景下，如果host为ip则取不到policy，会从约定的header stgw-orgservername中取host来匹配policy
ngx.ctx.user_policy = user_policy

5.初始化req_obj,并写入从ctx
req_obj {
    uuid = nil,
    uri = nil,
    header_raw = nil,
    request_uri = nil,
    method = nil,
    scheme = nil,
    client_ip = nil,
    client_ip_bot = nil,
    server_host = nil,
    server_port = nil,
    headers = nil,
    raw_headers = nil,
    cookies = nil,
    cookies_raw = nil,
    params = nil,
    multipart_headers = nil,
    body = nil,
    file_body = nil,
    upload_flag = false,
    query_string = nil,
    post_body = nil,
    headers_names = nil,
    params_names = nil,
    post_names = nil,
    cookies_names = nil,
    loaded = false,
    waf_flag = nil,
    real_ip
}
ngx.ctx.request = req_obj

6.根据user_policy中的module_status数据，按需加载基础安全模块的插件（web安全规则，访问控制规则，cc防护规则，网页防篡改规则，信息防泄漏规则，API防护规则）
对应规则为：
#web_security
ngx.ctx.bypass_owasp
ngx.ctx.bypass_ai
#access_control
ngx.ctx.bypass_geoip
ngx.ctx.bypass_custom
#cc_protection
ngx.ctx.bypass_cc
#antitamper
ngx.ctx.bypass_antitamper
#antileakage
ngx.ctx.bypass_antileakage
#api_protection
ngx.ctx.bypass_api	

7.防篡改检测
saaswaf防篡改检测
ngx.ctx.antileakage_set = antileakage_rule
```