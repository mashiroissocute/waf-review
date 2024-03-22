## 防盗链/保护资源
### referer
当前页面访问其他页面的资源时，在该请求重，浏览器会添加referer头部，表明该请求是从什么页面发起的。
服务器可以验证该请求重的referer来允许或者拒绝该请求。
例如 服务器配置：

``` 
referer [server name (特定url|正则表达式｜泛域名｜) blocked｜none]
if ($invalid_referer) {
	return 403
}
```

## secure_link （主要用户保护资源 比如付费下载）
当用户访问页面的时候，先根据用户信息生成加密的链接。用户使用该链接访问资源，链接通过验证则放行。否则拒绝访问。

```
 secure_link $arg_md5,$arg_expires;  #这里配置了2个参数一个是arg_md5，一个是arg_expires
 secure_link_md5 "$secure_link_expires$uri secret_key"; #secret_key为自定义的加密串（定期更换）   
 if ($secure_link = "") {
        return 403;       #资源不存在或哈希比对失败
 }
 if ($secure_link = "0") {
        return 403;      #时间戳过期 
 }
```
生成的url大致为：`http://img_server/test.jpg?md5=oa63dd5x_yi_E_eJLsAhHQ&expires=1523007651`
验证方法：nginx通过相同的md方法加密用户信息，与url的md5进行比对。

## cdn防盗链
通过比对referer实现
https://cloud.tencent.com/document/product/228/76113