## OpenResty线程模型
OpenResty或者Nginx的每一个worker最常见的运行方式是单进程单线程多协程。

- 多协程可以看作是多个用户级线程，用户级线程的线程管理是由用户空间的线程库来完成的，不需要进行内核级别的线程切换（线程上下文切换），因此效率很高。
- 单线程则是操作系统内核级的线程，该线程是操作系统操作的最小单位。
- 因此Nginx单个worker的线程模型可以看作是多对一的线程模型：多个用户线程（协程）对应一个操作系统内核线程。
![alt text](image-2.png)

多对一的线程模型实际只有一个线程在绑定一个cpu核心执行。
因此如果用户线程中的一个调用了阻塞API，例如` os.execute("sleep " .. n)`，内核线程会被阻塞，所有其他的用户线程都将被阻塞，单线程模型下的进程同样被阻塞并让出cpu时间片。
避免阻塞的方法是调用非阻塞的API，例如`ngx.sleep()`，该API会让用户线程阻塞，切换到另外的用户线程执行，并不会导致内核线程阻塞。
https://moonbingbing.gitbooks.io/openresty-best-practices/content/ngx_lua/sleep.html

## 其他的线程模型补充
http://c.biancheng.net/view/1220.html