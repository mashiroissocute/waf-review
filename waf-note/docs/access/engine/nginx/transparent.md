## 背景
Waf中的NGINX是对服务器有感知的，比如nginx 发向服务器的ip报文中，source ip是自己的ip而非客户端的ip。在服务端，客户端ip只能通过XFF或者XRI头来获取。但加入是四层的反向代理，没有了XFF等业务头，获取客户端ip就会比较困难。这样的反向代理并非透明代理。
## 透明代理
透明代理是十分普遍的需求，旨在让服务端感受不到前面有代理的存在。其中最为关键的点在于，服务端需要从nginx处收到source ip为客户端的ip报文，且发送dest ip为客户端ip的报文给nginx。

- nginx与上游服务的tcp链接建立和通讯时，从nginx发布的ip数据包，源地址由nginx服务所在机器的地址替换为客户端的ip地址;
- 从上游服务器返回的ip数据包，目的地址也是客户端的ip地址，通过特殊的路由方式，将返回的ip数据包路由到nginx进程；
- 最后nginx将收到的数据发送给对应的客户端。

## 实现方式
通过Linux 2.6.24中，socket新增的选项[IP_TRANSPARENT](https://links.jianshu.com/go?to=https%3A%2F%2Fman7.org%2Flinux%2Fman-pages%2Fman7%2Fip.7.html)，使socket可以接收目的地址没有配置的数据包，也`可以发送源地址不是本地地址的数据包`实现。使用这种方式，上游服务不需要任何修改就可以得到客户端的ip。
### NGINX配置
设置所有从nginx发送到上游服务器的ip数据包使用客户端的ip地址作为源地址。
1. nginx需要以 `root` 的用户运行，因为socket的 `IP_TRANSPARENT` 选项需要超级用户权限。
```
# in the "main"  context
use root;
```
2. nginx设置 `proxy_bind` ，使nginx在连接上游服务器时，启用socket的 `IP_TRANSPARENT` 选项。`remote_addr`是可以通过real_ip模块改变的。
```
server {
...
proxy_bind $remote_addr transparent;
}
```
### 服务端配置
配置服务器返回的ip数据包（其中的目的地址是客户端的ip）能返回到nginx所在的机器，并投递给nginx进程。
NGINX和服务同ip：
```
# 1. 将32上的上游服务回复的数据路由到nginx进程
# 1.1 对32上的上游服务返回的ip数据设置标记1
# NOTE: 如果32上增加access服务，需要修改这里--sport源端口，将新增的access服务端口包含进来
iptables -t mangle -A OUTPUT -p tcp --src 172.19.228.32 --sport 15010:15011 \
         -j MARK --set-xmark 0x1/0xffffffff
         
# 1.2 新增路由表，表id为100，将所有ip数据路由到lo网络接口         
# NOTE: 以下两条ip命令设置的route和rule都需要保存到/etc/rc.local中，以便在开机时启动
ip route add local 0.0.0.0/0 dev lo table 100
# 1.3 为标记为1的ip数据使用路由表100
ip rule add fwmark 1 lookup 100
```
NGINX和服务不同IP
```
# 1. 允许来自局域网网络接口bond0的ip数据
# 1.1 该设置主要避免来自32的ip数据被33屏蔽
iptables -I INPUT -i bond0 -j ACCEPT

# 1.2 禁用33主机上的reverse path filter。
# NOTE：因为nginx启动透明代理后，
#       33将在局域网网络接口bond0上收到源地址是客户端ip的数据(一般是公网ip), 
#       如果不禁用reverse path filter, 则因为无法在局域网网络接口回复公网的数据包，
#       导致33上收到的tcp sync包被33直接丢弃，导致tcp连接无法建立。
echo "net.ipv4.conf.all.rp_filter = 0" >> /etc/sysctl.conf
echo "net.ipv4.conf.bond0.rp_filter = 0" >> /etc/sysctl.conf

# 2. 将33上的上游服务回复的数据路由到32机器
# 2.1 对33上的上游服务发出的包设置标记1
# NOTE: 如果33上增加access服务，需要修改这里--sport源端口，将新增的access服务端口包含进来
iptables -t mangle -A OUTPUT -p tcp --src 172.19.228.33 --sport 15010:15011 \
         -j MARK --set-xmark 0x1/0xffffffff

# 2.2 新增路由表，表id为100，将所有ip数据路由到32机器
# NOTE: 以下两条ip命令设置的route和rule都需要保存到/etc/rc.local中，以便在开机时启动
ip route add default via 172.19.228.32 table 100
# 2.3 为标记为1的ip数据使用路由表100
ip rule add fwmark 1 lookup 100

# 3. 将33的上游服务返回的ip数据路由到32机器之后，再把这些数据路由到nginx进程
# 3.1 32上收到的来自33上游服务的数据设置标记1，使这些数据可以路由到nginx进程
# NOTE: 注意这一条命令在32主机上执行；
#       33上的增加access服务，需要修改这里的--sport源端口，将新增的access服务端口包含进来
# NOTE: 因为32上已经设置了对标记为1的ip数据使用策略路由表100，将其路由到lo网络接口
#       所以只需要对33返回的数据设置标记1即可
iptables -t mangle -A PREROUTING -p tcp --src 172.19.228.33 --sport 15010:15011 \
         -j MARK --set-xmark 0x1/0xffffffff
```
### 参考 
https://www.jianshu.com/p/0a975a997063

