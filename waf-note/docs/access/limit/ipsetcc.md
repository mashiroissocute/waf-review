## waf rs问题
Tgw的vip默认的流量保护是10Gb/s。
如果单个请求平均按照1kB计算，qps能达到1250000/s。
**这些流量到达到waf rs 集群，waf rs 集群会收到严重的大流量攻击，造成Dos，不能有效响应。**

Waf rs 集群分两种，独享集群，共享集群。
对于共享集群，当集群受到cc攻击时，大量499会影响到集群上所有的用户。
对于独享集群，waf应该阻止cc请求落到用户服务器，防止CC攻击到用户服务器。

## 用户服务器问题
**由于客户源站qps能力较弱，可能不是很高的cc攻击就会将客户打挂掉，这种情况不会触发waf侧的ipset防御机制。**
Cc对客户来讲属于异常，目前waf提供规则进行防御。
往往客户再被攻击的时候，需要临时去配置规则来防御。客户在被打挂的时候，配置规则会使得部分客户觉得压力很大，不知道怎么配。 
提供自动化的防御机制可以显著提升客户体验。

Cc 自动化防御: 自动化防御cc攻击，waf rs集群负载之内的cc攻击。

### 建立指标
按照30s进行窗口计算。统计30s内appid + host为key的各方面指标。包括：

- 超时指标： 
```
	- WAF：30s内有3000次访问中，且5%的超时比率，或者10次判定异常，超时比率达到10%判定异常。超时（499 + 504）
	- UPSTREAM ： 非499或504的请求，upstream_response_time平均超过1s
```
- QPS异常指标： 相比历史的QPS上涨太多
- IP指标：同个IP访问，30s内超过400次
- 负载指标：waf集群是否高负载

### 发现问题

- QPS无异常 & （UPSTREAM时间异常 ｜ WAF超时较多）时。说明用户源站异常，告警并触达用户。

- 集群负载过高 & QPS异常 & （UPSTREAM时间异常 ｜ WAF超时较多）时。说明WAF集群都扛不住的CC攻击量。需要对IP进行封堵，独立集群和cc隔离集群直接ipset封堵 IP指标和用户控制台配置的cc IP。共享集群先迁移到独立集群和cc隔离集群后，再ipset封堵 IP指标和用户控制台配置的cc IP

-  集群负载不高 & QPS异常 & （UPSTREAM时间异常 ｜ WAF超时较多）时。说明WAF集群扛得住，但是用户源站扛不住CC。WAF提供对CC攻击的源进行自动清洗的能力。方法是， 将IP指标和用户控制台配置的cc IP存放到mongo，并通过redis消息同步给bot-rpc-server。由bot-rps-server对指定ip源的cc进行拦截。bot-rpc-server里面的auto_cc。 redis ： SUB-cc-auto-reload
   




## IPSET
安装**iptables** 。
使用工具**ipset** 快速有效地设置和管理黑名单
### 创建黑名单
如何创建一个名为**blacklistv4** 的列表，其中可以包含1000000个IPv4地址（默认为65536）。
```
ipset create blacklistv4 hash:ip family inet maxelem 1000000
```
然后，将这个列表链接到**iptables** ，这样，添加到这个列表的IP地址就会被禁止。
```
iptables -I INPUT 1 -m set --match-set blacklistv4 src -j DROP
```
最后，可以将这个列表解除链接到**iptables**
``` 
iptables 删除 iptables -D INPUT -m set --match-set waf_black_ipset src -j DROP
```
### ipset操作
```
	
查看ipset  ipset list
创建set集合 ipset create blacklistv4 hash:ip family inet maxelem 1000000
销毁set集合 ipset destroy waf_black_ipset
向set集合加入元素 ipset add waf_black_ipset 1.1.1.1
向set集合删除元素 ipset del waf_black_ipset 1.1.1.1
查询元素是否在set中 ipset test waf_black_ipset 1.1.1.1
清除集合 ipset flush  waf_black_ipset
```
## IPSET API （waf-ipset-cc /manualServer / daemonserver） 
主动封堵、解除封堵ip。 发送请求到OSS，迁移流量到CC隔离集群。
ip相关的，接受请求并处理后，执行redis pubsub发送消息到ipset agent监听的redis。
迁移相关的，接受请求并处理后，发送请求到oss，执行VIP迁移。

## IPSET AGENT
监听redis的消息通道：
执行系统命令，操作ipset

