## 背景
某客户，在浏览器上的每一步操作，之前的选择全部丢失了。如何解决这个问题？

## 会话保持
https://learnku.com/docs/build-web-application-with-golang/060-session-and-data-storage/3188 
因为HTTP无状态，所以每次都需要登陆信息验证身份

- cookie：会话cookie（无过期时间，存储内存，进程消失则丢失。进程间不共享） ， 持久cookie（设置过期时间，存储在磁盘，过期时间内不会消失。进程间可以共享）。golang http.cookie 可以设置和读取cookie
- session通过cookie传递的时候，一般会讲cookie的过期时间设置为0，表示这是一个会话cookie。禁用cookie的时候，可以为页面返回的每个url加上sessionid的方法来解决。

### session的设计：
https://learnku.com/docs/build-web-application-with-golang/063-session-storage/3191
#### session的数据结构：
``` 
#list是一个双向链表
#*list.Element是一个值向链表node的指针
#string为sessionid

sessions map[string]*list.Element   

#list中的node为session
type SessionStore struct {
    sid          string                      // session id唯一标示
    timeAccessed time.Time                   // 最后访问时间
    value        map[interface{}]interface{} // session里面存储的值
}

#由于每个双向链表node都有一个指针对应。添加（插入队头）/删除（修改前一个节点的next指针和后一个节点的pre指针，并释放本节点的空间）/修改/查询都可以在O(1)时间内完成。
#每次访问了一个session元素（查询/修改/删除session中的某个value），会将改元素放在队头，并修改timeAccessed为访问时间。

#session的过期删除：
启动一个循环，从链表的末尾取session的timeAccessed发现 小于Now - MaxTime，则删除该条数据（list删除 + map元素删除）


#go list : https://studygolang.com/articles/11178
#go 数据结构 : https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter03/03.3.html
```

#### cookie和session的区别：

- 1.安全性
cookie将信息保存在客户端，如果不进行加密的话，无疑会暴露一些隐私信息，安全性很差，一般情况下敏感信息是经过加密后存储在cookie中，但很容易就会被窃取。cookie的风险：本地 cookie 中保存的用户名密码被破译，或 cookie 被其他网站收集（例如：1. appA 主动设置域 B cookie，让域 B cookie 获取；2. XSS，在 appA 上通过 javascript 获取 document.cookie，并传递给自己的 appB）。
而session只会将信息存储在服务端，如果存储在文件或数据库中，也有被窃取的可能，只是可能性比cookie小了太多。Session安全性方面比较突出的是存在会话劫持的问题，这是一种安全威胁。总体来讲，session的安全性要高于cookie；
- 2.性能
Cookie存储在客户端，消耗的是客户端的I/O和内存，而session存储在服务端，消耗的是服务端的资源。但是session对服务器造成的压力比较集中，而cookie很好地分散了资源消耗，就这点来说，cookie是要优于session的；
- 3.时效性
Cookie可以通过设置有效期使其较长时间内存在于客户端，而session一般只有比较短的有效期（用户主动销毁session或关闭浏览器后引发超时）；
- 4.其他
Cookie的处理在开发中没有session方便。而且cookie在客户端是有数量和大小的限制的，而session的大小却只以硬件为限制，能存储的数据无疑大了太多。
- 5.session的限制
session在集群模式下，需要维护数据中心


#### session会话劫持：
sessionid被攻击者截获，攻击者模仿客户发起正常请求。
**攻击者获取SessionID的方式有多种：**
```
1、 暴力破解：尝试各种Session ID，直到破解为止；
2、 预测：如果Session ID使用非随机的方式产生，那么就有可能计算出来；
3、 窃取：使用网络嗅探，XSS攻击等方法获得。
```
这个随机的Session ID往往是极其复杂的并且难于被预测出来，所以，对于第一、第二种攻击方式基本上是不太可能成功的。

对于第三种方式大多使用网络数据通讯层进行攻击获取，可以使用SSL进行防御。
目前有三种广泛使用的在Web环境中维护会话（传递Session ID）的方法：URL参数，隐藏域和Cookie。其中每一种都各有利弊，Cookie已经被证明是三种方法中最方便最安全的。从安全的观点，如果不是全部也是绝大多数针对基于Cookie的会话管理机制的攻击对于URL或是隐藏域机制同样适用，但是反过来却不一定，这就让Cookie成为从安全考虑的最佳选择。

**防御手段**
cookie httponly 和 token

- 1.sessionID 的值只允许 cookie 存储，而不是通过 URL 存储（cookie 不会像 URL 存储方式那么容易获取 sessionID），同时设置 cookie 的 httponly 为 true, 这个属性是设置是否可通过客户端脚本访问这个 cookie，可以防止这个 cookie 被 XSS 读取从而引起 session 劫持。
- 2.在每个请求里面加上 token，实现类似前面章节里面讲的防止 form 重复递交类似的功能，我们在每个请求里面加上一个隐藏的 token，然后每次验证这个 token，从而保证用户的请求都是唯一性。
- 3.给 session 额外设置一个创建时间的值，一旦过了一定的时间，我们销毁这个 sessionID，重新生成新的 session，这样可以一定程度上防止 session 劫持的问题。


## 身份认证
cookie和session都是可以记录用户与服务器的会话信息的，会话信息包括了身份验证信息，操作记录信息等等。但是某些时候，只是为了不用每次都输入密码即可得到授权，可以采用基于token的身份验证方法。Session 是一种**记录服务器和客户端会话状态的机制，使服务端有状态化，可以记录会话信息**。而 Token 是**令牌**，**访问资源接口（API）时所需要的资源凭证**。Token **使服务端无状态化，不会存储会话信息。**
### token
将用户信息加密后形成token返回给客户端，客户端再次请求的时候带上token，服务端解密token得到用户登陆信息。主要采用计算的方法解决了session的有状态的（分布式场景需要指定服务器，或者搭建session集群），且需要存储的缺点。
### JWT

- Token：简单的token组成为：uid(用户唯一的身份标识)、time(当前时间的时间戳)、sign（签名，token 的前几位以哈希算法压缩成的一定长度的十六进制字符串）。服务端验证客户端发送过来的 Token 时，还需要查询数据库获取用户信息，然后验证 Token 是否有效。

- JWT：一个完整的JWT是由三个部分组成，分别是Header头部、Payload数据部分、Signature签名三部分组成；将 Token 和 Payload（可以存储用户信息） 加密后存储于客户端，服务端只需要使用密钥解密进行校验（校验也是 JWT 自己实现的）即可，不需要查询或者减少查询数据库，因为 JWT 自包含了用户信息和加密的数据。
