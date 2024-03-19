


## 背景
 waf控制台或者API接口在填写回源策略的时候，有两个选项，分别是ip回源和域名回源。
 由于waf的nginx转发机器在内网，并且由VIP提供了外网和内网打通的能力。如果不对回源ip或者域名进行限制的话，恶意用户填写内网ip或者域名地址，则可以通过waf从外网访问到内网，从而发起SSRF攻击。

 
## 知识点
服务端请求伪造（Server Side Request Forgery, SSRF）指的是攻击者在未能取得服务器所有权限时，利用服务器漏洞以服务器的身份发送一条构造好的请求给服务器所在内网。SSRF攻击通常针对外部网络无法直接访问的内部系统。
SSRF可服务器所在内网进行端口扫描，攻击运行在内网的应用。
常见的SSRF服务器漏洞：
当服务器用于转发用户请求的时候，用户配置了内网的ip或域名地址，就可能存在ssrf攻击。
此时需要过滤转发的ip或者域名。

## 防护策略 

- OSS添加 “回源域名和ip检测任务”
每个域名添加一个InnerSsrfCheckTask子任务，周期性扫描域名的回源信息，存在SSRF风险则将回源信息置空。（最好修改完之后，再加上一个告警，并提示原来的回源信息）
- OSS添加 “新增域名和编辑域名检测回源信息的逻辑”
在新增和编辑域名数据入库之后，开启协程检测回源信息，如果存在SSRF风险则将回源信息置空。（由于dns查询域名时，存在一定的时延，这会增加添加域名和编辑域名超时的风险。因此采用先入库，然后异步检测，有风险则修改的策略）



## 核心代码
- 1.如何判断某个ip是不是在某个网段内：将ip和网段都转二进制字符串，比较二者在掩码长度之内是否一致。

```golang
package main

import (
    "fmt"
    "strconv"
    "strings"
)

func ip2binary(ip string) string {
    str := strings.Split(ip, ".")
    var ipstr string
    for _, s := range str {
        i, err := strconv.ParseUint(s, 10, 8) //10进制 int8
        if err != nil {
            fmt.Println(err)
        }
        ipstr = ipstr + fmt.Sprintf("%08b", i)
    }
    return ipstr
}

func match(ip, iprange string) bool {
    ipb := ip2binary(ip)
    ipr := strings.Split(iprange, "/")
    masklen, err := strconv.ParseUint(ipr[1], 10, 32)
    if err != nil {
        fmt.Println(err)
        return false
    }
    iprb := ip2binary(ipr[0])
    return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
}

func isInnerIp(v string) bool {
    if match(v, "10.255.255.255/8") {
        return true
    }
    if match(v, "172.16.255.255/12") {
        return true
    }
    if match(v, "192.168.255.255/16") {
        return true
    }
    if match(v, "100.64.255.255/10") {
        return true
    }
    if match(v, "9.255.255.255/8") {
        return true
    }
    if match(v, "127.255.255.255/8") {
        return true
    }
    if match(v, "11.0.0.0/8") {
        return true
    }
    if match(v, "30.0.0.0/8") {
        return true
    }
    return false
}

func main() {
    if isInnerIp("1.1.1.1") {
        println("inner ip")
    } else {
        println("not innner ip")
    }
}
```

- 2.dns解析各种记录

```golang
func main() {
	var dnsClient = dns.Client{
		Timeout: 5 * time.Second, //dns一次查询的超时时间
	}
	var (
		rmsg   *dns.Msg
		domain = "lucas-sni-protection-domain.qcloudwaf.com."
	)

	nsAdderss := "10.123.119.98:53" //dns服务器地址
	msg := &dns.Msg{
		MsgHdr: dns.MsgHdr{},
	}
	msg.SetQuestion(domain, dns.TypeA)                 //查询dns的记录 例如 A AAAA CNAME
	rmsg, _, err := dnsClient.Exchange(msg, nsAdderss) //dns 查询有问题的时候 rmsg为nil
	if rmsg == nil || err != nil {
		println("dns query error")
		println("err : ", err.Error())
		println("is rmsg nil : ", rmsg == nil)
		return
	}

	for _, ans := range rmsg.Answer {
		aMsg := ans.(*dns.A)
		println("A", aMsg.A.String())
	}
}
```