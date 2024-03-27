## 问题背景
用户添加域名，选择TLS1.0-1.2，加密套件为通用型。
绑定host验证情况如下：
- 浏览器返回
``` 
err_ssl_version_or_cipher_mismatch
```
- NGINX ERROR LOG
``` 
2022/09/29 12:50:28:039 [error] 16620#0: *3933848794 SSL_do_handshake() failed (SSL: error:1417A0C1:SSL routines:tls_post_process_client_hello:no shared cipher) while SSL handshaking, client: 119.147.10.183, server: 0.0.0.0:443
```

## 问题原因
**用户证书的签名算法是ECDSAWithSHA384**。 其中SHA384是摘要算法，ECDSA是非对称加密算法。

- **加密：** CA使用该SHA384算法对证书信息进行HASH后，利用CA私钥和ECDSA算法加密HASH值，生成一个证书签名。证书签名append在证书信息后面。

- **解密：** 与加密过程相反，浏览器需要使用ECDSA和CA公钥解密证书签名。因此需要在密钥协商过程中，使用ECDSA身份认证算法。

**问题：** WAF服务器，TLS1.2版本没有配置ECDSA身份认证算法。因此Openssl HandShake报错，浏览器展示err_ssl_version_or_cipher_mismatch。

## 修复计划
- **临时解决方法**：1.让用户把TLS1.3协议也勾选上。或者 2.使用CLB-WAF接入
- **长期解决方法**：等封网结束后，WAF TLS1.2版本支持该身份认证算法。