## encode/decode优化

CPU占用多优化
 
针对json解析部分（CPU占用耗时比例约为18%）进行优化

https://juejin.cn/post/7195491694052884517


使用snoic代替encoding，提升json解析速度。



## 锁竞争优化

通过pprof/mutex分析锁占用时间


用 atomic.Load/StoreXXX，atomic.Value, sync.Map 等代替 Mutex。



## http库优化

使用高效的第三方库，如用fasthttp替代 net/http

https://blog.csdn.net/RA681t58CJxsgCkJ31/article/details/130737645