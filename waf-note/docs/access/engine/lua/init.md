### 前提
WAF-engine引入LUA是使用了lua-ngx-module而并非直接使用OpenResty
因此OpenResty里面的部分库，WAF并不支持直接使用。

### 配置
在配置文件中，定义了一些shared-dict和nginx不同阶段会执行的lua代码路径。只有init， init_worker, rewrite, access, header_filter, body_filter,log阶段。而没有balancer，ssl等阶段，因为这些在接入部分已经形成了静态的nginx配置。
```lua
lua_shared_dict antitamper_cache 512m;
lua_shared_dict cc_limit 300m;
lua_shared_dict udr_time 300m;
lua_shared_dict deny_time 300m;
lua_shared_dict captcha_state 100m;
# status_shared_dict 保存各个防护模块（api，bot，ai，menshen）的开关状态的信息
lua_shared_dict cloudwaf_status 32m;
# data_shared_dict 保存waf允许是的一些信息，比如本地ip地址，运行的worker数量等
lua_shared_dict cloudwaf_data 128m;
lua_shared_dict update_data 512m;
lua_shared_dict update_flag 128m;
lua_shared_dict update_domain_list 32m;
lua_shared_dict qps_data 256m;
lua_shared_dict attack_count 32m;
lua_shared_dict attack_detail 256m;
lua_shared_dict deny_count 32m;
lua_shared_dict bwlist_dict 1024m;
lua_shared_dict bwlist_version_dict 128m;

init_by_lua_file          
init_worker_by_lua_file   
rewrite_by_lua_file        
access_by_lua_file         
header_filter_by_lua_file  
body_filter_by_lua_file     
log_by_lua_file          
```

### init master
``` lua
#init全局变量：
ROOT_DIR ： 代码路径
LOCAL_IPADDR ： 本地ip地址
DEFAULT_HTML ： waf nginx的非200返回页面 502，504，40x
REGION ： 机器地域
GEOIP ： areaban模块的另个称呼

#init waf模块
config = {}， load_config到waf.config

#init waf logger 对ngx.log进一步封装消息格式

#init cloudwaf_status init各个防护模块的起始开关信息
	waf_status.set_menshen_mode(mode) 设置规则引擎走门神还是本地规则
	waf_status.set_waf_check(true)
    waf_status.set_owasp_check(true)
    waf_status.set_white_check(true)
    waf_status.set_areaban_check(true)
    waf_status.set_custom_check(true)
    waf_status.set_antitamper_check(true)
    waf_status.set_cc_check(true)
    waf_status.set_ai_check(true)
    waf_status.set_bot_check(true)
    waf_status.set_antileakage_check(true)
    waf_status.set_api_check(true)
    waf_status.set_bwlist_check(true)
    waf_status.set_innerdetect_check(true)
#init cloudwaf_data init部分waf运行时的信息 workernum ， ip
#init areaban.ipdb 读取ipdb lib来获得ip对应的地域信息，存在在areaban.ipdb{}中

#init captcha_l5 captcha是验证码的意思，调用第三方动态库l5Agent.api_init_route
```

### init worker
```lua
#---------init socket
sock_conf_access{127.0.0.1,28000,udp} 上传访问日志
sock_conf_attack{127.0.0.1,28001,udp} 上传攻击日志
sock_conf_report{127.0.0.1,28002,udp} 门神和规则结果不一致时候上报到report中
sock_conf_heart{127.0.0.1,28003,udp} 上传心跳日志
sock_conf_strace{127.0.0.1,28009,udp} 上报每个请求经历的检测模块
response socket只有webank才有
socket:init(conf)
将各个socket填入waf.socket{}


#---------worker_monitor
每个work调用是，增加wafdata的worknum，如果发现该变量和ngxwork不对应，会向monitor（第三方so）发送消息。

#---------update_domain_timer，循环执行任务，每隔一段时间就执行读取策略。该任务从woker启动开始以后就没有停止。
    -- 黑白名单策略加载, 从policy_shm中读取黑白名单策略到bwlist.policy
    -- local policy_shm = ngx.shared.bwlist_dict
    -- local version_shm = ngx.shared.bwlist_version_dict       
    bwlist:update_cache()
    
    -- 防护策略加载，从update_data_shm中读取防护策略到waf.policy
    -- local update_data_shm = ngx.shared.update_data
    -- local update_flag_shm = ngx.shared.update_flag 
    local status, err, ret = xpcall(update_domains,record_err,premature)

#--------- work 0 load_policy 开始从redis拉取策略

0.判断是否已经拉取过了，if waf_data.is_reload() then return 策略已经拉取完成了，不需要再次拉取。
ngx.shared.cloudwaf_data：get（ngx_relaod） ngx_reload存在shared dict中。
当nginx reload的时候，会有work init阶段执行。但是shared dict中的数据不会丢失，is_reload的状态不会丢失。也就是说拉取策略只需要master启动那一次就行了。

1.get_domain_list : hget customer-domain-data data  向某个地域的redis读取数据，所以这里只是单个地域的域名列表，
	{
		"saas" : [{"host": "bestpracticechina.com", "status": 1, "appid": 1257197664}, {"host": "qiye100.cn", "status": 1, "appid": 1251201367}, {"host": "packagexai.cn", "status": 0, "appid": 1256514921}]
		"clb" : [{"host": "cccc.com", "status": 1, "appid": 1257197664}, {"host": "qiye100.cn", "status": 1, "appid": 1251201367}, {"host": "packagexai.cn", "status": 0, "appid": 1256514921}]
		"saas_appid" : [{"appid": 1253752930, "status": 3}, {"appid": 1255461525, "status": 3} ... ]
		"clb_appid" : [{"appid": 1308627201, "status": 1}, {"appid": 1308875430, "status": 1} ... ]
	}
==> 排除status = 3
	valid_domain_list {"bestpracticechina.com" , "qiye100.cn" ...}
	
2.LoadPolicyByDomain(valid_domain_list,nil) 对于list中的每个域名，从redis5和redis中获取域名基础信息（key：clb【clb:dmain】saas【domain】）和防护信息(key: appid+domain) 并写到共享内存update_data_shm（ngx.shared.update_data）和update_flag_shm（ngx.shared.update_flag）中
-- update_data_shm:set(domain, domain_policy)
-- update_flag_shm:set(domain, domain_policy.version)
		 domain_policy{
                basic{
                    appid
                    uin
                    edition
                    proxy(is_cdn)
                    pan_domain(extensive_domain)
                    instanceid  
					botenable
                },
                enable (status = 1 --> true else --> false)
                rule_mode (mode = 1 --> 2  0 ---> 1 else ---> 0) 规则引擎
                ai_mode
                force_local_owasp 走本地规则引擎检测而非门神
                custom_set
                white_set
                areaban_set
                options
                bypass
                antitamper_set
                antileakage_set
                traffic_marking
                module_status{ //基础安全的策略开关，是否bypass
					web_security
					access_control
					cc_protection
					antitamper
					antileakage
					api_protection
				}
                bot_jsinject
                block_page
                webshell_mode
                version
				ulimit_args 参数数量限制
				acc_log
				rule_set{
					["21000000"] :  
                    ["22000000"] ...
                    ["23000000"]
                    ["24000000"]
				}
            }

            version{
                host_config_change(unix.time)
                customrule
                areaban
                options
                antitamper
                antileakage
                antileakage_set
                traffic_marking
                module_status
                bot_jsinject
                block_page
                webshell
            }
	}

3.waf_data.set_reload() 设置完成防护配置加载的标识
ngx.shared.cloudwaf_data：set("ngx_reload",true)


#--------- work 0 load_shm_config_timer 开始从redis拉取黑白名单策略

0.if waf_data.is_bwlist_reload() then return
黑白名单策略已经拉取完成了，不需要再次拉取。
ngx.shared.cloudwaf_data：get（ngx_bwlist_reload）
当nginx reload的时候，会有workinit阶段执行。但是shareddict中的数据不会丢失，is_reload的状态不会丢失。也就是说全量拉取黑白名单策略只需要master启动那一次就行了。接下来只需要更新。

1.get domain list 同load policy

2.更新shareddict  policy_shm，version_shm
local policy_shm = ngx.shared.bwlist_dict
local version_shm = ngx.shared.bwlist_version_dict

update_config 大致的内容是{"149.45.138.0/24": {"info": {"action": 40, "valid_ts": 1627449307}} .... }
version 内容 local ver = version..'_'..tostring(timestamp)
policy_shm:set(domain, update_config)
version_shm:set(domain, version)

3.waf_data.set_bwlist_reload()
ngx.shared.cloudwaf_data：set（ngx_bwlist_reload true）

#--------- hearbeat（放在日志章节详细解析）
每10s上报一次心跳日志，把heartinfo两个json通过heardsocket传出去。 heartinfo内容在log阶段填充。心跳日志的内容主要是按照域名维度统计：
收包长度（ub）/
回包长度（db）/
域名访问次数（ac）/
nginx 4xx，5xx次数（4x，5x）/
源站 4xx，5xx次数（u4，u5）/
攻击次数(at)
bot cc次数（b，cc）/
回源的QPS,这段时间内回源的次数（acu，回源的次数）

#--------- upload_attack_log（放在日志章节详细解析）
每10s上报一次攻击日志
日志链接https://iwiki.woa.com/pages/viewpage.action?pageId=1205944745

# waf_rpc_heart_recover 暂时不清楚用法
```