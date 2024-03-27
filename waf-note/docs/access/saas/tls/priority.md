## 现象
用户http2.0协议，修改了加密套件之后，无法访问。

## 原因

- **加密套件按照从弱到强的顺序写在配置文件中。**
waf选择了比较弱的加密套件，该加密套件恰好在http2.0的黑名单中。
客户端http2.0协议，不认可该加密套件。浏览器报错：ERR_HTTP2_INADEQUATE_TRANSPORT_SECURITY

- **ERR_HTTP2_INADEQUATE_TRANSPORT_SECURITY解读**
https://blog.csdn.net/shiyong1949/article/details/109043720
https://stackoverflow.com/questions/64893062/how-to-fix-err-http2-inadequate-transport-security-in-iis-on-windows-server-2016
This problem is happening because of the HTTP/2. This basically means that the site started a HTTP/2 connection but there was a **blacklisted cypher negotiated**. SO the browser has prevented the access to the website. So, the usual solution for this is to reorder the cypher suites to meet the requirements of the HTTP/2.
Another solution is to disable HTTP/2 and only use HTTP/1.1. This can be addressed on the server-side by setting the following registry keys and the restarting the host Windows server

- 同时，myssl工具扫描域名显示：
`没有优先使用FS系列加密套件导致服务降级。`
问题原因也是因为加密套件按照从弱到强的顺序写在了配置文件中。FS系列的加密套件没有被优先使用。

# 修复方案
修改ssl_cipher加密套件的顺序，优先使用fs系列和强加密套件。
**加密套件优先级规则：**
按照waf以前的默认配置：ssl_ciphers  EECDH+CHACHA20:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;
制定了加密套件的优先级。

# 思考
出现该问题的原因，是因为对http2.0场景下的加密套件选择不够了解。