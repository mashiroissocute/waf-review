## gdb调试Core文件
1. gdb查看core文件内容：gdb sbin/nginx_bin core.629 
2. 打印函数栈：bt 
3. 进入某个函数栈：f $n
4. 查看变量  p ((ngx_shm_zone_t*)(cycle->shared_memory.part->elts))[0]
5. 还可以通过地址直接转结构体 ：
``` 
(gdb) p *(ngx_http_conf_ctx_t*)0x7f5a53304ca8 
$9 = {main_conf = 0x7f5a6b608c50, srv_conf = 0x7f5a53304cc0, loc_conf = 0x7f5a53304f48}


(gdb) p *(ngx_http_core_srv_conf_t*)((*(ngx_http_conf_ctx_t*)0x7f5a53304ca8)->srv_conf[0]) 

$15 = {server_names = {elts = 0x7f5a533052c0, nelts = 1, size = 32, nalloc = 4, pool = 0x7f5a6b607010}, ctx = 0x7f5a53304ca8, 
  file_name = 0x7f5a53301bd6 "/usr/local/services/nginx_install-1.0//conf/auto_vhost/work.yintongcard.com_8001.conf", line = 10, server_name = {len = 20, 
    data = 0x7f5a50eb3e48 "work.yintongcard.com/usr/local/services/nginx_install-1.0//html"}, connection_pool_size = 512, request_pool_size = 4096, client_header_buffer_size = 4096, 
  large_client_header_buffers = {num = 4, size = 1048576}, client_header_timeout = 60000, ignore_invalid_headers = 1, merge_slashes = 1, underscores_in_headers = 1, spdy_default_off = 1, 
  dysvr_clean_time = 180, vsvc_id = 162835, vip_vpcid = -1, vport = 0, vip = {len = 0, data = 0x0}, listen = 1, captures = 0, invalid = 0, delete_time = 0, http2 = 0, spdy = 0, 
  named_locations = 0x7f5a521970b0}
```



## gdb调试运行中进程
1. gdb进入运行中进程：gdb -p pid
2. 设置断点 b  src/core/ngx_slab.c:173 或者 b ngx_ssl_session_cache_init
3. 查看断点 info break 
4. 删除断点 delete 2（info break展示的序号）
5. 继续执行到断点 c
6. 单步执行 n

## 查看变量的技巧
core文件内容为该进程实际使用的物理内存的“快照”。分析core dump文件可以获取应用程序崩溃时的现场信息，如程序运行时的CPU寄存器值、堆栈指针、栈数据、函数调用栈等信息。
因此core文件是包含了程序运行时变量的值，不论变量的内存实际是在栈还是堆上，都可以打印这个变量。
打印void 指针变量的时候，需要知道这个指针变量的具体类型，可以通过查看源码获得。然后进行一个类型转换。
p ((ngx_shm_zone_t*)(cycle->shared_memory.part->elts))[0]
如果强行打印void 指针变量，会得到Attempt to dereference a generic pointer. 错误

