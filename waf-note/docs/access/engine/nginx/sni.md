
## 背景
sni背景
TLS握手发生在HTTP报文传输之前，此时对于服务端而言，只有客户端提供的Port和IP可以拿到，拿不到Host。因此只有在TLS握手的时候，提供sni，服务端才能够拿到host相关的信息来选择证书。
https://blog.csdn.net/wzj_110/article/details/110149984


## 关于证书验证
- **waf侧的防护域名证书一定需要正确**（**有效 且 证书和防护域名对应**）,否则会提示不安全的连接。
客户端将首先和waf建立https连接。在TLS握手时，sni为防护域名，waf的server_name也为防护域名。
因此WAF会提供用户配置的证书给客户端，客户端拿到证书将校验有效性（是否有效 且 是否和请求域名对应）。
有效性验证失败将提示不安全的连接。

- **允许waf侧的证书正确，源站证书不正确**
客户端浏览器与waf侧建立连接时，waf侧的证书是正确的，客户端与waf建立安全的https连接。
waf和源站建立连接时，默认不会对证书进行强校验，类似curl -k。因此源站的证书可以是无效的（但必须符合证书格式，否则源站nginx -t是会ermerg的）。
这种情况下，绕过waf，直接访问源站，将提示不安全的连接。

## 关于host ，proxy_ssl_name
- `proxy_host` : proxy_pass后面的路径（lucas-sni-protection-domain.qcloudwzgj.com_443.1045123.upstream）
- `host` and `http_host`: 如果客户端发过来的请求的header中有host字段时，http_host和host都是原始请求的host。但是waf作为反向代理的时候，可以发送给后端的报文是可以修改host的。例如，我们目前的修改为：
``` 
 set $userheader $host;
 if ($http_host != '') {set $userheader $http_host;}
 proxy_set_header Host $userheader;
```


## 实验链接
https://blog.dianduidian.com/post/nginx%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86%E5%BD%93%E5%90%8E%E7%AB%AF%E4%B8%BAhttps%E6%97%B6%E7%9A%84%E4%B8%80%E4%BA%9B%E7%BB%86%E8%8A%82%E5%92%8C%E5%8E%9F%E7%90%86/


## 最终结论

waf-nginx配置中，默认proxy_ssl_verify off，不校验源站证书有效性。
并且默认不发送sni信息：proxy_ssl_server_name off;

但如果源站对sni没有进行特殊处理（例如，部分用户自己实现服务器并对sni进行校验，在单个或多个域名回源的场景下，waf向源站传递的host必须要填写正确，目前我们填写的是防护域名，如果源站指定了特定的host访问，可能会502。该情况只能先提供单联系waf进行处理）。

解决方法是: 需要配置proxy_ssl_server_name on;
此时，sni默认发送防护域名。可以通过proxy_ssl_name进行修改。

最终结论，域名回源场景，nginx向后端传递的sni和host其实依赖用户源站的设计方式。如果用户存在访问不通的问题，需要提工单，我们分析其服务器处理方式，来定制sni。