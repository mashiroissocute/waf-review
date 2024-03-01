
# 1. 购买海外云服务器
腾讯云服务器 HK 2c2g 20g 3m 
86rmb/month

# 2. 配置和安装shadowsocket server
https://github.com/zhouaini528/scientific_internet_access/blob/master/ALIYUN.md
```
// 下载ssr配置脚本
wget --no-check-certificate https://raw.githubusercontent.com/teddysun/shadowsocks_install/master/shadowsocksR.sh

chmod +x shadowsocksR.sh

./shadowsocksR.sh 2>&1 | tee shadowsocksR.log
```

脚本运行分如下阶段：
- 设置个人配置 ： 加密算法、端口、协议、密码
生产配置文件 /etc/shadowsocks.json
```
{
    "server":"0.0.0.0",
    "server_ipv6":"[::]",
    "server_port":19933,
    "local_address":"127.0.0.1",
    "local_port":1080,
    "password":"WOshizhu1996ca",
    "timeout":120,
    "method":"aes-256-cfb",
    "protocol":"origin",
    "protocol_param":"",
    "obfs":"plain",
    "obfs_param":"",
    "redirect":"",
    "dns_ipv6":false,
    "fast_open":false,
    "workers":1
}
```
- 安装运行依赖
yum install -y python3 openssl openssl-devel curl wget unzip gcc automake autoconf make libtool
- 下载安装加密库libsodium
```
libsodium_file="libsodium-1.0.18"
libsodium_url="https://github.com/jedisct1/libsodium/releases/download/1.0.18-RELEASE/libsodium-1.0.18.tar.gz"

解压后 编译
./configure --prefix=/usr && make && make install
```
- 下载ssr
```
shadowsocks_r_file="shadowsocksr-3.2.2"
shadowsocks_r_url="https://github.com/shadowsocksrr/shadowsocksr/archive/3.2.2.tar.gz"
```
- 下载ssr启动脚本 放到etc/init.d目录，进行系统管理
```
https://raw.githubusercontent.com/teddysun/shadowsocks_install/master/shadowsocksR
```

- 启动
```
/etc/init.d/shadowsocks start

其实是运行
usr/local/shadowsocks/server.py -c `/etc/shadowsocks.json`
但是server.py中指定的运行环境为python, 所以我手动执行了
python3 /usr/local/shadowsocks/server.py -c /etc/shadowsocks.json
```
- 检查
![enter image description here](/tencent/api/attachments/s3/url?attachmentid=3715154)

# 3. 打开安全组
放通tcp:19933端口访问

# 4. 配置和安装shadowsocket client
![enter image description here](/tencent/api/attachments/s3/url?attachmentid=3715126)



# 5. 资源
链接: https://pan.baidu.com/s/1VlSGBjA3FsArKO0FCGSUgA  密码: jd68
