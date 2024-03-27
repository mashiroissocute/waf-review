

## 第三方模块 ngx_http_stgw_guard_module.c

### 使用方法
```
ngx_http_stgw_guard_module:
用于在cpu，内存，网卡高负载情况下的自身过载保护。

stgw_guard的默认配置为

stgw_guard interval=1 factor=3 trans_time=30 shm_size=10m net_dev=eth1;
stgw_guard_load 95;
stgw_guard_mem 80;
stgw_guard_net 90;
stgw_guard_limit 3000r/s;

其中stgw_guard指令若没有配置，则表示不开启过载保护，若配置了则开启过载保护
stgw_guard_load stgw_guard_mem stgw_guard_net stgw_guard_limit 指令不配置时为上述默认值
其中，对于https，limit为3000r/s，对于http，limit默认设置为https的2倍。
stgw_guard为http级配置，为全局配置
stgw_guard_load stgw_guard_mem stgw_guard_net stgw_guard_limit为server级配置，可以针对不同的server设置不通的阈值

stgw_guard指令的各个参数：
interval=1 表示处于正常状态时，每1s采集一次cpu，mem，net系统信息
factor=3 表示的是，当处于限流状态时,每3倍interval时间采集一次cpu，mem，net系统信息
trans_time=30 主要为了防止突然的陡升陡降，去毛刺。其表示
处于正常状态时，连续30s处于高负载情况下，会切换到限流状态
处于限流状态时，连续30s处于低负载情况下，会切换到正常状态
shm_size=10m 表示默认分配10m的空间用于存储限流节点
对于http，限流节点以server_ip:server_port:domain区分
对于https，限流节点以server_ip:server_port区分
由于每个节点平均字节数一般不超过64字节，因此10m足以支撑10w以上服务
net_dev表示的是对网卡限流时指定的网卡设备
```

### 疑惑
#### 1.对于https， 为什么 限流节点的key， server_ip:server_port 里边不包含domain? 是有什么特殊的考虑吗？ 还是说处理所处的阶段 取不到 domain?
对，有可能取不到。这个是负载很高的时候，针对HTTPS请求，握手的时候进行限制。我们相当于是过载保护，SSL握手很多的时候，非常耗CPU，就提前在握手的时候把请求限制了，来防止RSA计算机型攻击。


### 源码
```c
#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>

#include <sys/ioctl.h>
#include <linux/sockios.h>
#include <net/if.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/socket.h>

#if (NGX_HTTP_DUMP)
#include <ngx_http_dump_buff_helper.h>
#endif

#ifndef SIOCETHTOOL
#define SIOCETHTOOL     0x8946
#endif

#define NGX_MEMINFO_FILE            "/proc/meminfo"
#define NGX_CPUINFO_FILE            "/proc/stat"
#define NGX_NETINFO_FILE            "/proc/net/dev"
#define NGX_MAX_NETDEV_COUNT        5

/* CMDs currently supported */
#define ETHTOOL_GSET        0x00000001 /* Get settings. */

typedef __uint32_t __u32;       /* ditto */
typedef __uint16_t __u16;       /* ditto */
typedef __uint8_t __u8;         /* ditto */

/* This should work for both 32 and 64 bit userland. */
struct ethtool_cmd {
        __u32   cmd;
        __u32   supported;      /* Features this interface supports */
        __u32   advertising;    /* Features this interface advertises */
        __u16   speed;          /* The forced speed, 10Mb, 100Mb, gigabit */
        __u8    duplex;         /* Duplex, half or full */
        __u8    port;           /* Which connector port */
        __u8    phy_address;
        __u8    transceiver;    /* Which transceiver to use */
        __u8    autoneg;        /* Enable or disable autonegotiation */
        __u32   maxtxpkt;       /* Tx pkts before generating tx int */
        __u32   maxrxpkt;       /* Rx pkts before generating rx int */
        __u32   reserved[4];
};


typedef struct {
    u_char                       color;
    u_char                       dummy;
    u_short                      len;
    int64_t                      last;
    int64_t                      remain;
    int64_t                      drop;
    time_t                       flag_time;
    u_char                       data[1];
} ngx_http_stgw_guard_node_t;

typedef struct {
    ngx_rbtree_t                  rbtree;
    ngx_rbtree_node_t             sentinel;
} ngx_http_stgw_guard_rbtree_t;

typedef struct {
    ngx_int_t net_rate;
    ngx_str_t net_dev;
} ngx_http_stgw_guard_netdev_t;

typedef struct {
    ngx_flag_t                    enable;
    time_t                        interval;
    time_t                        factor;
    time_t                        trans_time;
    ngx_http_stgw_guard_netdev_t  net_dev[NGX_MAX_NETDEV_COUNT];
    ngx_int_t                     netdev_count;
    size_t                        shm_size;
    ngx_shm_zone_t                *shm_zone;
    ngx_slab_pool_t               *shpool;
    ngx_http_stgw_guard_rbtree_t  *rbtree;

} ngx_http_stgw_guard_main_conf_t;

typedef struct {
    ngx_int_t                     load;
    ngx_int_t                     mem;
    ngx_int_t                     net;
    int64_t                       limit_rate;
} ngx_http_stgw_guard_conf_t;

typedef struct {
    size_t total;
    size_t free;
    size_t buffers;
    size_t cached;
} ngx_http_stgw_guard_meminfo_t;

typedef struct {
    size_t user;
    size_t nice;
    size_t sys;
    size_t idle;
    size_t iowait;
    size_t hardirq;
    size_t softirq;
} ngx_http_stgw_guard_cpuinfo_t;

typedef struct {
    size_t send;
    size_t recv;
    ngx_msec_t msec;
} ngx_http_stgw_guard_netinfo_t;


static void *ngx_http_stgw_guard_create_conf(ngx_conf_t *cf);
static char *ngx_http_stgw_guard_merge_conf(ngx_conf_t *cf, void *parent,
    void *child);
static ngx_int_t ngx_http_stgw_guard_init(ngx_conf_t *cf);
static ngx_int_t
ngx_http_stgw_guard_update_load(ngx_log_t *log);
#if 0
static ngx_int_t
ngx_http_stgw_guard_update_mem(ngx_log_t *log);
static ngx_int_t
ngx_http_stgw_guard_update_net(ngx_log_t *log, ngx_http_stgw_guard_main_conf_t *glcf);
#endif
static ngx_int_t
ngx_http_stgw_guard_init_shm(ngx_shm_zone_t *shm_zone, void *data);
static void
ngx_http_stgw_guard_rbtree_insert_value(ngx_rbtree_node_t *temp,
    ngx_rbtree_node_t *node, ngx_rbtree_node_t *sentinel);
static ngx_int_t
ngx_http_stgw_guard_limit_req(ngx_http_stgw_guard_main_conf_t *glcf, ngx_http_stgw_guard_conf_t *srv_conf, ngx_log_t *log, ngx_str_t *key);
static ngx_int_t
ngx_http_stgw_guard_lookup(ngx_http_stgw_guard_main_conf_t *glcf, ngx_http_stgw_guard_conf_t *srv_conf, ngx_log_t *log, ngx_uint_t hash, ngx_str_t *key);
static ngx_int_t
ngx_http_stgw_guard_handler(ngx_http_stgw_guard_main_conf_t *glcf, ngx_http_stgw_guard_conf_t *srv_conf, ngx_log_t *log, ngx_str_t *key);
static char *
ngx_http_stgw_guard_limit(ngx_conf_t *cf, ngx_command_t *cmd, void *conf);
#if 0
static ngx_int_t
ngx_http_stgw_guard_getmeminfo(ngx_http_stgw_guard_meminfo_t *meminfo, ngx_log_t *log);
#endif
static ngx_int_t
ngx_http_stgw_guard_getcpuinfo(ngx_http_stgw_guard_cpuinfo_t *cpuinfo, ngx_log_t *log);
#if 0
static ngx_int_t
ngx_http_stgw_guard_getnetinfo(ngx_http_stgw_guard_netinfo_t *netinfo, ngx_log_t *log, ngx_str_t *net_dev);
#endif
static char *
ngx_http_stgw_guard(ngx_conf_t *cf, ngx_command_t *cmd, void *conf);
static void *
ngx_http_stgw_guard_create_main_conf(ngx_conf_t *cf);
static ngx_int_t
ngx_http_stgw_guard_http_handler(ngx_http_request_t *r);
ngx_int_t
ngx_http_stgw_guard_https_handler(ngx_http_connection_t *hc, ngx_connection_t *c);
static ngx_int_t
ngx_http_stgw_guard_get_net_dev_rate(char *net_dev);

static ngx_command_t  ngx_http_stgw_guard_commands[] = {

    { ngx_string("stgw_guard"),
      NGX_HTTP_MAIN_CONF|NGX_CONF_TAKE5,
      ngx_http_stgw_guard,
      0,
      0,
      NULL },

    { ngx_string("stgw_guard_load"),
      NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_CONF_TAKE1,
      ngx_conf_set_num_slot,
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_stgw_guard_conf_t, load),
      NULL },

    { ngx_string("stgw_guard_mem"),
      NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_CONF_TAKE1,
      ngx_conf_set_num_slot,
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_stgw_guard_conf_t, mem),
      NULL },

    { ngx_string("stgw_guard_net"),
      NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_CONF_TAKE1,
      ngx_conf_set_num_slot,
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_stgw_guard_conf_t, net),
      NULL },

    { ngx_string("stgw_guard_limit"),
      NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_CONF_TAKE1,
      ngx_http_stgw_guard_limit,
      NGX_HTTP_SRV_CONF_OFFSET,
      0,
      NULL },

      ngx_null_command
};


static ngx_http_module_t  ngx_http_stgw_guard_module_ctx = {
    NULL,                                   /* preconfiguration */
    ngx_http_stgw_guard_init,               /* postconfiguration */

    ngx_http_stgw_guard_create_main_conf,   /* create main configuration */
    NULL,                                   /* init main configuration */

    ngx_http_stgw_guard_create_conf,        /* create server configuration */
    ngx_http_stgw_guard_merge_conf,         /* merge server configuration */

    NULL,          /* create location configuration */
    NULL            /* merge location configuration */
};


ngx_module_t  ngx_http_stgw_guard_module = {
    NGX_MODULE_V1,
    &ngx_http_stgw_guard_module_ctx,          /* module context */
    ngx_http_stgw_guard_commands,             /* module directives */
    NGX_HTTP_MODULE,                        /* module type */
    NULL,                                   /* init master */
    NULL,                                   /* init module */
    NULL,                                   /* init process */
    NULL,                                   /* init thread */
    NULL,                                   /* exit thread */
    NULL,                                   /* exit process */
    NULL,                                   /* exit master */
    NGX_MODULE_V1_PADDING
};


static time_t    ngx_http_stgw_guard_cached_load_exptime = 0;
#if 0
static time_t    ngx_http_stgw_guard_cached_mem_exptime = 0;
static time_t    ngx_http_stgw_guard_cached_net_exptime = 0;
#endif
static time_t    ngx_http_stgw_guard_cached_first_abnormal_time = 0;
static time_t    ngx_http_stgw_guard_cached_first_normal_time = 0;
static ngx_int_t ngx_http_stgw_guard_cached_load = 0;
static ngx_int_t ngx_http_stgw_guard_cached_mem = 0;
static ngx_int_t ngx_http_stgw_guard_cached_net = 0;
static ngx_http_stgw_guard_cpuinfo_t ngx_http_stgw_guard_cached_cpuinfo;
static ngx_http_stgw_guard_netinfo_t ngx_http_stgw_guard_cached_netinfo[NGX_MAX_NETDEV_COUNT];


static ngx_int_t
ngx_http_stgw_guard_update_load(ngx_log_t *log)
{
    ngx_int_t rc;
    ngx_http_stgw_guard_cpuinfo_t cpuinfo;
    size_t t1_used, t1_all, t2_used, t2_all;


    rc = ngx_http_stgw_guard_getcpuinfo(&cpuinfo, log);
    if (rc == NGX_ERROR) {
        ngx_http_stgw_guard_cached_load = NGX_CONF_UNSET;
        return NGX_ERROR;
    }
    //first get the cpuinfo
    if(ngx_http_stgw_guard_cached_cpuinfo.user == 0) {
        ngx_http_stgw_guard_cached_load = NGX_CONF_UNSET;
        ngx_http_stgw_guard_cached_cpuinfo = cpuinfo;
        return NGX_OK;
    }

    t1_used = ngx_http_stgw_guard_cached_cpuinfo.user + ngx_http_stgw_guard_cached_cpuinfo.nice +
            ngx_http_stgw_guard_cached_cpuinfo.sys + ngx_http_stgw_guard_cached_cpuinfo.iowait +
            ngx_http_stgw_guard_cached_cpuinfo.hardirq + ngx_http_stgw_guard_cached_cpuinfo.softirq;
    t1_all = t1_used + ngx_http_stgw_guard_cached_cpuinfo.idle;

    t2_used = cpuinfo.user + cpuinfo.nice + cpuinfo.sys + cpuinfo.iowait + cpuinfo.hardirq + cpuinfo.softirq;
    t2_all = t2_used + cpuinfo.idle;

    if(t2_all == t1_all) {
        ngx_http_stgw_guard_cached_load = NGX_CONF_UNSET;
        return NGX_ERROR;
    }

    ngx_http_stgw_guard_cached_load = (t2_used - t1_used) * 100/ (t2_all - t1_all);
    ngx_http_stgw_guard_cached_cpuinfo = cpuinfo;

    return NGX_OK;
}

#if 0
static ngx_int_t
ngx_http_stgw_guard_update_mem(ngx_log_t *log)
{
    ngx_int_t      rc;
    ngx_http_stgw_guard_meminfo_t  m;


    rc = ngx_http_stgw_guard_getmeminfo(&m, log);
    if (rc == NGX_ERROR) {
        ngx_http_stgw_guard_cached_mem = NGX_CONF_UNSET;
        return NGX_ERROR;
    }

    if(m.total > 0) {
        ngx_http_stgw_guard_cached_mem = (m.total - m.free - m.cached - m.buffers) * 100 / m.total;
        if(ngx_http_stgw_guard_cached_mem < 0) {
            ngx_http_stgw_guard_cached_mem = NGX_CONF_UNSET;
            return NGX_ERROR;
        }
        return NGX_OK;
    }
    else {
        ngx_http_stgw_guard_cached_mem = NGX_CONF_UNSET;
        return NGX_ERROR;
    }
}

static ngx_int_t
ngx_http_stgw_guard_update_net(ngx_log_t *log, ngx_http_stgw_guard_main_conf_t *glcf)
{
    ngx_int_t rc;
    ngx_http_stgw_guard_netinfo_t netinfo[NGX_MAX_NETDEV_COUNT];
    size_t t1_send, t1_recv, t2_send, t2_recv;
    ngx_msec_t ms;
    int i;
    int flag = 0;
    ngx_int_t max_cached_net;
    ngx_int_t cached_net;

    cached_net = max_cached_net = 0;

    for(i = 0; i< glcf->netdev_count; i++) {
        rc = ngx_http_stgw_guard_getnetinfo(&netinfo[i], log, &glcf->net_dev[i].net_dev);
        if (rc == NGX_ERROR) {
            ngx_http_stgw_guard_cached_net = NGX_CONF_UNSET;
            return NGX_ERROR;
        }
    }

    for(i = 0; i< glcf->netdev_count; i++) {
        if(ngx_http_stgw_guard_cached_netinfo[i].recv == 0) {
            ngx_http_stgw_guard_cached_netinfo[i] = netinfo[i];
            flag = 1;
        }
    }
    if(flag == 1) {
        ngx_http_stgw_guard_cached_net = NGX_CONF_UNSET;
        return NGX_OK;
    }

    for(i = 0; i< glcf->netdev_count; i++) {
        t1_send = ngx_http_stgw_guard_cached_netinfo[i].send;
        t1_recv = ngx_http_stgw_guard_cached_netinfo[i].recv;

        t2_send = netinfo[i].send;
        t2_recv = netinfo[i].recv;
        ms = (netinfo[i].msec - ngx_http_stgw_guard_cached_netinfo[i].msec);

        if ((ngx_int_t) (t2_send - t1_send) < 0 || (ngx_int_t) (t2_recv - t1_recv) < 0 || ms <= 0)
            continue;
        cached_net = (t2_send - t1_send) > (t2_recv - t1_recv) ? (t2_send - t1_send) * 100 * 1000/(glcf->net_dev[i].net_rate * ms * 1000000 / 8) : (t2_recv - t1_recv) * 100 * 1000/ (glcf->net_dev[i].net_rate * ms * 1000000 / 8);
        if(cached_net > max_cached_net)
            max_cached_net = cached_net;

        ngx_http_stgw_guard_cached_netinfo[i] = netinfo[i];
    }
    ngx_http_stgw_guard_cached_net = max_cached_net;

    return NGX_OK;
}
#endif

static ngx_int_t
ngx_http_stgw_guard_get_net_dev_rate(char *net_dev)
{
    struct ifreq ifr;
    int fd;

    /* Setup our control structures. */
    memset(&ifr, 0, sizeof(ifr));
    strcpy(ifr.ifr_name, net_dev);

    /* Open control socket. */
    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd < 0) {
        return NGX_CONF_UNSET;
    }

    int err;
    struct ethtool_cmd ep;

    ep.cmd = ETHTOOL_GSET;
    ifr.ifr_data = (caddr_t)&ep;
    err = ioctl(fd, SIOCETHTOOL, &ifr);
    if(err != 0) {
        close(fd);
        return NGX_CONF_UNSET;
    }

    //only support >=1000Mb/s <= 40000Mb/s
    if(ep.speed < 1000 || ep.speed > 40000) {
        close(fd);
        return NGX_CONF_UNSET;
    }

    close(fd);
    return ep.speed;
}

static ngx_int_t
ngx_http_stgw_guard_init_shm(ngx_shm_zone_t *shm_zone, void *data)
{

    size_t                     len;
    ngx_http_stgw_guard_main_conf_t *glcf;
    ngx_http_stgw_guard_main_conf_t *oglcf = data;

    glcf = shm_zone->data;
    glcf->shm_zone = shm_zone;
    glcf->shpool = (ngx_slab_pool_t *) shm_zone->shm.addr;

    if(data && oglcf->rbtree) {
        glcf->rbtree = oglcf->rbtree;
        ngx_log_error(NGX_LOG_DEBUG, shm_zone->shm.log, 0,
            "stgw_guard oglcf exist, use oglcf");
        return NGX_OK;
    }

    if (shm_zone->shm.exists) {
        glcf->rbtree = glcf->shpool->data;

        return NGX_OK;
    }

    glcf->rbtree = ngx_slab_alloc(glcf->shpool, sizeof(ngx_http_stgw_guard_rbtree_t));
    if (glcf->rbtree == NULL) {
        ngx_log_error(NGX_LOG_ERR, shm_zone->shm.log, 0,
            "ngx_slab_alloc return null");
        return NGX_ERROR;
    }

    glcf->shpool->data = glcf->rbtree;

    ngx_rbtree_init(&glcf->rbtree->rbtree, &glcf->rbtree->sentinel,
                    ngx_http_stgw_guard_rbtree_insert_value);

    len = sizeof(" in stgw_guard zone \"\"") + shm_zone->shm.name.len;

    glcf->shpool->log_ctx = ngx_slab_alloc(glcf->shpool, len);
    if (glcf->shpool->log_ctx == NULL) {
        return NGX_ERROR;
    }

    ngx_sprintf(glcf->shpool->log_ctx, " in stgw_guard zone \"%V\"",
                &shm_zone->shm.name);

    glcf->shpool->log_nomem = 0;

    return NGX_OK;
}

static void
ngx_http_stgw_guard_rbtree_insert_value(ngx_rbtree_node_t *temp,
    ngx_rbtree_node_t *node, ngx_rbtree_node_t *sentinel)
{
    ngx_rbtree_node_t          **p;
    ngx_http_stgw_guard_node_t   *lrn, *lrnt;

    for ( ;; ) {

        if (node->key < temp->key) {

            p = &temp->left;

        } else if (node->key > temp->key) {

            p = &temp->right;

        } else { /* node->key == temp->key */

            lrn = (ngx_http_stgw_guard_node_t *) &node->color;
            lrnt = (ngx_http_stgw_guard_node_t *) &temp->color;

            p = (ngx_memn2cmp(lrn->data, lrnt->data, lrn->len, lrnt->len) < 0)
                ? &temp->left : &temp->right;
        }

        if (*p == sentinel) {
            break;
        }

        temp = *p;
    }

    *p = node;
    node->parent = temp;
    node->left = sentinel;
    node->right = sentinel;
    ngx_rbt_red(node);
}


static ngx_int_t
ngx_http_stgw_guard_limit_req(ngx_http_stgw_guard_main_conf_t *glcf, ngx_http_stgw_guard_conf_t *srv_conf, ngx_log_t *log, ngx_str_t *key)
{
    uint32_t                     hash;
    ngx_int_t                    rc;
    ngx_int_t                   ratio;

    ratio = 1;
    //http ratio double
    if(key->data[0] == '1')
        ratio = 2;
    hash = ngx_crc32_short(key->data, key->len);

    ngx_shmtx_lock(&glcf->shpool->mutex);

    rc = ngx_http_stgw_guard_lookup(glcf, srv_conf, log, hash, key);

    ngx_shmtx_unlock(&glcf->shpool->mutex);

    if (rc == NGX_BUSY) {
        return NGX_HTTP_SERVICE_UNAVAILABLE;
    }

    ngx_log_error(NGX_LOG_DEBUG, log, 0,
        "not limit request by key \"%V\"  cur cpu:%d, mem:%d, net:%d, conf cpu:%d, mem:%d, net:%d limit_rate:%lr/s",
        key, ngx_http_stgw_guard_cached_load, ngx_http_stgw_guard_cached_mem, ngx_http_stgw_guard_cached_net,
        srv_conf->load, srv_conf->mem, srv_conf->net, srv_conf->limit_rate * ratio);

    return NGX_DECLINED;
}

static ngx_int_t
ngx_http_stgw_guard_lookup(ngx_http_stgw_guard_main_conf_t *glcf, ngx_http_stgw_guard_conf_t *srv_conf, ngx_log_t *log, ngx_uint_t hash, ngx_str_t *key)
{
    size_t                      size;
    ngx_int_t                   rc;
    int64_t                     now;
    ngx_rbtree_node_t          *node, *sentinel;
    ngx_http_stgw_guard_node_t  *lr;
    ngx_int_t                   ratio;
    struct timeval tv;
    gettimeofday(&tv, NULL);

    now = tv.tv_sec;

    node = glcf->rbtree->rbtree.root;
    sentinel = glcf->rbtree->rbtree.sentinel;

    //for http, the limit rate is double
    ratio = 1;
    if(key->data[0] == '1')
        ratio = 2;
    while (node != sentinel) {

        if (hash < node->key) {
            node = node->left;
            continue;
        }

        if (hash > node->key) {
            node = node->right;
            continue;
        }

        /* hash == node->key */

        lr = (ngx_http_stgw_guard_node_t *) &node->color;

        rc = ngx_memn2cmp(key->data, lr->data, key->len, (size_t) lr->len);

        if (rc == 0) {

            if (now > lr->last) {
                lr->remain = srv_conf->limit_rate * ratio;
                lr->last = now;
            }
            lr->remain--;

            ngx_log_error(NGX_LOG_DEBUG, log, 0,
                   "remain = %l, key=\"%V\", rate=%l r/s", lr->remain, key, srv_conf->limit_rate * ratio);

            if (lr->remain < 0) {
                lr->drop++;
                if((lr->flag_time + 1 <= ngx_time()) && lr->drop > 0) {
                    lr->flag_time = ngx_time();
                    ngx_log_error(NGX_LOG_ERR, log, 0,
                        "limiting %d requests by key \"%V\" cur cpu:%d, mem:%d, net:%d, conf cpu:%d, mem:%d, net:%d limit_rate:%lr/s",
                        lr->drop, key, ngx_http_stgw_guard_cached_load, ngx_http_stgw_guard_cached_mem, ngx_http_stgw_guard_cached_net,
                        srv_conf->load, srv_conf->mem, srv_conf->net, srv_conf->limit_rate * ratio);
                    lr->drop = 0;
                }
#if (NGX_HTTP_NEW_STAT)
				ngx_add_global_stat(ngx_stat_worker_shm, "limit_req_num", 1, STAT_SUM);
#endif
                return NGX_BUSY;
            }

            return NGX_OK;
        }

        node = (rc < 0) ? node->left : node->right;
    }

    size = offsetof(ngx_rbtree_node_t, color)
           + offsetof(ngx_http_stgw_guard_node_t, data)
           + key->len;


    ngx_log_error(NGX_LOG_DEBUG, log, 0,
                   "ngx_slab_alloc_locked a new node");

    node = ngx_slab_alloc_locked(glcf->shpool, size);

    if (node == NULL) {

        ngx_log_error(NGX_LOG_WARN, log, 0,
                      "could not allocate guard node%s", glcf->shpool->log_ctx);
        return NGX_OK;
    }

    node->key = hash;

    lr = (ngx_http_stgw_guard_node_t *) &node->color;

    lr->len = (u_short) key->len;
    lr->remain = srv_conf->limit_rate * ratio;

    ngx_memcpy(lr->data, key->data, key->len);

    ngx_rbtree_insert(&glcf->rbtree->rbtree, node);

    lr->last = now;
    lr->drop = 0;
    lr->flag_time = ngx_time();

    return NGX_OK;
}

static ngx_int_t
ngx_http_stgw_guard_handler(ngx_http_stgw_guard_main_conf_t *glcf, ngx_http_stgw_guard_conf_t *srv_conf, ngx_log_t *log, ngx_str_t *key)
{

    if (key->len == 0 || key->len >= 65535) {
        ngx_log_error(NGX_LOG_WARN, log, 0,
                      "the value of the \"%V\" key "
                      "is more than 65535 bytes or null",
                      key);
        return NGX_DECLINED;
    }

#if (NGX_HTTP_DUMP)
    ngx_heavy_loading = 0;
#endif

    /* load */
    if (srv_conf->load != NGX_CONF_UNSET) {

        if (ngx_http_stgw_guard_cached_load_exptime < ngx_time()) {

            ngx_http_stgw_guard_update_load(log);
            ngx_http_stgw_guard_cached_load_exptime = ngx_time() + glcf->interval;
            ngx_log_error(NGX_LOG_DEBUG, log, 0,
                       "http stgw_guard handler load: current:%d srv_conf:%d",
                       ngx_http_stgw_guard_cached_load,
                       srv_conf->load);
            if(ngx_http_stgw_guard_cached_load != NGX_CONF_UNSET && ngx_http_stgw_guard_cached_load > srv_conf->load)
                ngx_http_stgw_guard_cached_load_exptime = ngx_http_stgw_guard_cached_load_exptime + (glcf->factor -1) * glcf->interval;
        }

        if(ngx_http_stgw_guard_cached_load != NGX_CONF_UNSET) {

            if (ngx_http_stgw_guard_cached_load > srv_conf->load) {

#if (NGX_HTTP_DUMP)
                ngx_heavy_loading = 1;
#endif

                ngx_http_stgw_guard_cached_first_normal_time = 0;
                if(ngx_http_stgw_guard_cached_first_abnormal_time == 0)
                    ngx_http_stgw_guard_cached_first_abnormal_time = ngx_time();
                ngx_log_error(NGX_LOG_WARN, log, 0,
                              "stgw_guard load limited, current:%d srv_conf:%d",
                              ngx_http_stgw_guard_cached_load,
                              srv_conf->load);

                if(ngx_http_stgw_guard_cached_first_abnormal_time + glcf->trans_time < ngx_time())
                    return ngx_http_stgw_guard_limit_req(glcf, srv_conf, log, key);
                else
                    return NGX_DECLINED;
            }
        }
    }

    /* mem */
#if 0
    if (srv_conf->mem != NGX_CONF_UNSET) {

        if (ngx_http_stgw_guard_cached_mem_exptime < ngx_time()) {

            ngx_http_stgw_guard_update_mem(log);
            ngx_http_stgw_guard_cached_mem_exptime = ngx_time() + glcf->interval;
            ngx_log_error(NGX_LOG_DEBUG, log, 0,
                           "http stgw_guard handler mem: current:%d srv_conf:%d",
                           ngx_http_stgw_guard_cached_mem,
                           srv_conf->mem);
            if(ngx_http_stgw_guard_cached_mem != NGX_CONF_UNSET && ngx_http_stgw_guard_cached_mem > srv_conf->mem)
                ngx_http_stgw_guard_cached_mem_exptime = ngx_http_stgw_guard_cached_mem_exptime + (glcf->factor -1) * glcf->interval;
        }

        if (ngx_http_stgw_guard_cached_mem != NGX_CONF_UNSET) {

            if (ngx_http_stgw_guard_cached_mem > srv_conf->mem) {

#if (NGX_HTTP_DUMP)
                ngx_heavy_loading = 1;
#endif

                ngx_http_stgw_guard_cached_first_normal_time = 0;
                if(ngx_http_stgw_guard_cached_first_abnormal_time == 0)
                    ngx_http_stgw_guard_cached_first_abnormal_time = ngx_time();
                ngx_log_error(NGX_LOG_WARN, log, 0,
                              "stgw_guard mem limited, "
                              "current:%d srv_conf:%d",
                              ngx_http_stgw_guard_cached_mem,
                              srv_conf->mem);

                if(ngx_http_stgw_guard_cached_first_abnormal_time + glcf->trans_time < ngx_time())
                    return ngx_http_stgw_guard_limit_req(glcf, srv_conf, log, key);
                else
                    return NGX_DECLINED;
            }
        }
    }
#endif

    /* net */
#if 0
    if (srv_conf->net != NGX_CONF_UNSET && glcf->netdev_count > 0) {

        if (ngx_http_stgw_guard_cached_net_exptime < ngx_time()) {

            ngx_http_stgw_guard_update_net(log, glcf);
            ngx_http_stgw_guard_cached_net_exptime = ngx_time() + glcf->interval;
            ngx_log_error(NGX_LOG_DEBUG, log, 0,
                           "http stgw_guard handler net: current:%d srv_conf:%d",
                           ngx_http_stgw_guard_cached_net,
                           srv_conf->net);
            if(ngx_http_stgw_guard_cached_net != NGX_CONF_UNSET && ngx_http_stgw_guard_cached_net > srv_conf->net)
                ngx_http_stgw_guard_cached_net_exptime = ngx_http_stgw_guard_cached_net_exptime + (glcf->factor -1) * glcf->interval;
        }

        if (ngx_http_stgw_guard_cached_net != NGX_CONF_UNSET) {

            if (ngx_http_stgw_guard_cached_net > srv_conf->net) {

#if (NGX_HTTP_DUMP)
                ngx_heavy_loading = 1;
#endif

                ngx_http_stgw_guard_cached_first_normal_time = 0;
                if(ngx_http_stgw_guard_cached_first_abnormal_time == 0)
                    ngx_http_stgw_guard_cached_first_abnormal_time = ngx_time();
                ngx_log_error(NGX_LOG_WARN, log, 0,
                              "stgw_guard net limited, "
                              "current:%d srv_conf:%d",
                              ngx_http_stgw_guard_cached_net,
                              srv_conf->net);

                if(ngx_http_stgw_guard_cached_first_abnormal_time + glcf->trans_time < ngx_time())
                    return ngx_http_stgw_guard_limit_req(glcf, srv_conf, log, key);
                else
                    return NGX_DECLINED;
            }
        }
    }
#endif
    //limit begin after glcf->trans_time of ngx_http_stgw_guard_cached_first_abnormal_time

    if(ngx_http_stgw_guard_cached_first_normal_time == 0)
        ngx_http_stgw_guard_cached_first_normal_time = ngx_time();

    if(ngx_http_stgw_guard_cached_first_normal_time + glcf->trans_time < ngx_time())
        ngx_http_stgw_guard_cached_first_abnormal_time = 0;

    return NGX_DECLINED;
}

ngx_int_t
ngx_http_stgw_guard_https_handler(ngx_http_connection_t *hc, ngx_connection_t *c)
{

    ngx_http_stgw_guard_main_conf_t   *glcf;
    ngx_http_stgw_guard_conf_t   *srv_conf;
    ngx_int_t                    ret;
    char *ip;
    unsigned short port;
    ngx_str_t                    key;
    u_char                       buf[1024];
    ngx_memzero(buf,sizeof(buf));

    glcf = ngx_http_get_module_main_conf(hc->conf_ctx, ngx_http_stgw_guard_module);
    srv_conf = ngx_http_get_module_srv_conf(hc->conf_ctx, ngx_http_stgw_guard_module);

    if (!glcf->enable) {
        ngx_log_error(NGX_LOG_DEBUG, c->log, 0,
            "stgw_guard !glcf->enable return NGX_DECLINED");
        return NGX_DECLINED;
    }

    //http key=ip:port:domain https key=ip:port
    port = ntohs(((struct sockaddr_in *)(c->local_sockaddr))->sin_port);
    ip = inet_ntoa(((struct sockaddr_in *)(c->local_sockaddr))->sin_addr);
    if(ip == NULL)
    {
        ngx_log_error(NGX_LOG_ERR, c->log, 0,
            "https: invalid ip of local_sockaddr");
        return NGX_DECLINED;
    }
    ngx_snprintf(buf, sizeof(buf), "0:%s:%d", ip, port);

    key.len = ngx_strlen(buf);
    key.data = buf;
    ret = ngx_http_stgw_guard_handler(glcf, srv_conf, c->log, &key);
    c->stgw_guard = 1;
    ngx_log_error(NGX_LOG_DEBUG, c->log, 0,
            "https: ngx_http_stgw_guard_handler ret:%d key \"%V\"", ret, &key);
    return ret;
}

static ngx_int_t
ngx_http_stgw_guard_http_handler(ngx_http_request_t *r)
{

    ngx_http_stgw_guard_main_conf_t   *glcf;
    ngx_http_stgw_guard_conf_t   *srv_conf;
    ngx_int_t                    ret;
    ngx_str_t                    key;
    char *ip;
    unsigned short port;
    ngx_str_t                    server_name;
    u_char                       buf[256];
    ngx_http_core_srv_conf_t  *cscf;

    glcf = ngx_http_get_module_main_conf(r, ngx_http_stgw_guard_module);
    srv_conf = ngx_http_get_module_srv_conf(r, ngx_http_stgw_guard_module);

    ngx_memzero(buf,sizeof(buf));

    if (!glcf->enable) {
        ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0,
            "stgw_guard !glcf->enable return NGX_DECLINED");
        return NGX_DECLINED;
    }

    if (r->main->limit_req_set) {
        ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0,
            "r->main->limit_req_set is true");
        return NGX_DECLINED;
    }
    if (r->connection->stgw_guard == 1) {
        ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0,
            "r->connection->stgw_guard is set");
        r->connection->stgw_guard = 0;
        return NGX_DECLINED;
    }

    if (r->headers_in.server.len) {
        server_name.len = r->headers_in.server.len;
        server_name.data = r->headers_in.server.data;

    } else {
        cscf = ngx_http_get_module_srv_conf(r, ngx_http_core_module);
        server_name.len = cscf->server_name.len;
        server_name.data = cscf->server_name.data;
    }

    //http key=ip:port:domain https key=ip:port
    port = ntohs(((struct sockaddr_in *)(r->connection->local_sockaddr))->sin_port);
    ip = inet_ntoa(((struct sockaddr_in *)(r->connection->local_sockaddr))->sin_addr);
    if(ip == NULL || server_name.len <= 0)
    {
        ngx_log_error(NGX_LOG_ERR, r->connection->log, 0,
            "https: invalid ip of local_sockaddr or server_name");
        return NGX_DECLINED;
    }
    
    key.len = ngx_snprintf(buf, sizeof(buf), "1:%s:%d:%V", ip, port, &server_name) - buf;
    key.data = buf;
    
    ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0,
            "http: ip:%s, port:%d, server_name:%s, key \"%V\"", ip, port, server_name.data, &key);
    ret = ngx_http_stgw_guard_handler(glcf, srv_conf, r->connection->log, &key);
    ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0,
            "http: ngx_http_stgw_guard_handler ret:%d key \"%V\"", ret, &key);
    r->main->limit_req_set = 1;
    return ret;
}
static void *
ngx_http_stgw_guard_create_main_conf(ngx_conf_t *cf)
{
    ngx_http_stgw_guard_main_conf_t  *glcf;

    glcf = ngx_pcalloc(cf->pool, sizeof(ngx_http_stgw_guard_main_conf_t));
    if (glcf == NULL) {
        ngx_conf_log_error(NGX_LOG_ERR, cf, 0,
            "ngx_pcalloc ngx_http_stgw_guard_main_conf_t fail");
        return NGX_CONF_ERROR;
    }

    glcf->enable = 0;
    glcf->interval = 1;
    glcf->factor = 3;
    glcf->trans_time = 30;
    glcf->shm_size = 10 * 1024 * 1024;
    glcf->shm_zone = NULL;
    glcf->shpool = NULL;
    glcf->rbtree = NULL;

    return glcf;
}

static void *
ngx_http_stgw_guard_create_conf(ngx_conf_t *cf)
{
    ngx_http_stgw_guard_conf_t  *srv_conf;

    srv_conf = ngx_pcalloc(cf->pool, sizeof(ngx_http_stgw_guard_conf_t));
    if (srv_conf == NULL) {
        ngx_conf_log_error(NGX_LOG_ERR, cf, 0,
            "ngx_pcalloc ngx_http_stgw_guard_conf_t fail");
        return NGX_CONF_ERROR;
    }

    srv_conf->load = NGX_CONF_UNSET;
    srv_conf->mem = NGX_CONF_UNSET;
    srv_conf->net = NGX_CONF_UNSET;
    srv_conf->limit_rate = NGX_CONF_UNSET_SIZE;

    return srv_conf;
}

static char *
ngx_http_stgw_guard_merge_conf(ngx_conf_t *cf, void *parent, void *child)
{

    ngx_http_stgw_guard_conf_t  *prev = parent;
    ngx_http_stgw_guard_conf_t  *srv_conf = child;

    ngx_conf_merge_value(srv_conf->load, prev->load, 95);
    ngx_conf_merge_value(srv_conf->mem, prev->mem, 90);
    ngx_conf_merge_value(srv_conf->net, prev->net, 90);
    ngx_conf_merge_value(srv_conf->limit_rate, prev->limit_rate, 3000);

    return NGX_CONF_OK;
}

static char *
ngx_http_stgw_guard(ngx_conf_t *cf, ngx_command_t *cmd, void *conf)
{
    size_t                             size;
    char                               *s, *p;
    char                               *delim = ",";
    ngx_str_t                          *value, name;
    ngx_uint_t                         i;
    ngx_str_t shm_name = ngx_string("stgw_guard_shm");
    ngx_int_t   net_rate, interval, factor, trans_time;
    ngx_http_stgw_guard_main_conf_t *glcf = conf;
    value = cf->args->elts;
    glcf->netdev_count = 0;

    interval = 1;
    factor = 3;
    trans_time = 30;
    size = 10 * 1024 * 1024;
    for (i = 1; i < cf->args->nelts; ++i) {
        if (ngx_strncmp(value[i].data, "shm_size=", 9) == 0) {
            name.data = value[i].data + 9;
            name.len = value[i].len - 9;
            size = ngx_parse_size(&name);

            if (size == (size_t) NGX_ERROR) {
                ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                    "invalid shm size \"%V\"", &value[i]);
                return NGX_CONF_ERROR;
            }

            if (size < 8 * ngx_pagesize) {
                ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                    "shm \"%V\" is too small", &value[i]);
                return NGX_CONF_ERROR;
            }
            continue;
        }
        if (ngx_strncmp(value[i].data, "interval=", 9) == 0) {
            name.len = value[i].len - 9;
            name.data = value[i].data + 9;
            interval = ngx_atoi(name.data, name.len);
            if (interval == NGX_ERROR || interval <= 0) {
                ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                    "invalid interval \"%V\" ", &value[i]);
                return NGX_CONF_ERROR;
            }
            continue;
        }
        if (ngx_strncmp(value[i].data, "factor=", 7) == 0) {
            name.len = value[i].len - 7;
            name.data = value[i].data + 7;
            factor = ngx_atoi(name.data, name.len);
            if (factor == NGX_ERROR || factor <= 0) {
                ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                    "invalid factor \"%V\" ", &value[i]);
                return NGX_CONF_ERROR;
            }
            continue;
        }
        if (ngx_strncmp(value[i].data, "trans_time=", 11) == 0) {
            name.len = value[i].len - 11;
            name.data = value[i].data + 11;
            trans_time = ngx_atoi(name.data, name.len);
            if (trans_time == NGX_ERROR || trans_time <= 0) {
                ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                    "invalid trans_time \"%V\" ", &value[i]);
                return NGX_CONF_ERROR;
            }
            continue;
        }
        if (ngx_strncmp(value[i].data, "net_dev=", 8) == 0) {
            name.len = value[i].len - 8;
            name.data = value[i].data + 8;
            if(name.len <= 0)
                continue;
            s = (char *) name.data;
            p = strtok(s, delim);
            while(p != NULL) {
                net_rate = ngx_http_stgw_guard_get_net_dev_rate(p);
                if(net_rate != NGX_CONF_UNSET) {
                    glcf->net_dev[glcf->netdev_count].net_rate = net_rate;
                    glcf->net_dev[glcf->netdev_count].net_dev.len = strlen(p);
                    glcf->net_dev[glcf->netdev_count].net_dev.data = (u_char *) p;
                    glcf->netdev_count++;
                    if(glcf->netdev_count >= NGX_MAX_NETDEV_COUNT)
                        break;
                }
                p = strtok(NULL, delim);
            }
            continue;
        }
        ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
            "invalid parameter \"%V\"", &value[i]);
        return NGX_CONF_ERROR;
    }

    glcf->interval = interval;
    glcf->factor = factor;
    glcf->trans_time = trans_time;
    glcf->shm_size = size;

    if(glcf->netdev_count == 0) {
        ngx_conf_log_error(NGX_LOG_ERR, cf, 0,
            "can't get a net device");
    }
    glcf->shm_zone = ngx_shared_memory_add(cf, &shm_name, glcf->shm_size,
                                     &ngx_http_stgw_guard_module);
    if (glcf->shm_zone == NULL) {
        ngx_conf_log_error(NGX_LOG_ERR, cf, 0,
            "ngx_shared_memory_add ngx_http_stgw_guard_module fail");
        return NGX_CONF_ERROR;
    }

    glcf->shm_zone->init = ngx_http_stgw_guard_init_shm;
    glcf->shm_zone->data = glcf;
    glcf->enable = 1;

    for(i = 0; i< NGX_MAX_NETDEV_COUNT; i++)
        ngx_http_stgw_guard_cached_netinfo[i].recv = 0;

    return NGX_CONF_OK;
}

static char *
ngx_http_stgw_guard_limit(ngx_conf_t *cf, ngx_command_t *cmd, void *conf)
{
    u_char                            *p;
    size_t                             len;
    ngx_str_t                         *value;
    ngx_int_t                          rate, scale;

    ngx_http_stgw_guard_conf_t *srv_conf = conf;
    value = cf->args->elts;

    len = value[1].len;
    p = value[1].data + len - 3;

    if (ngx_strncmp(p, "r/s", 3) == 0) {
        scale = 1;
        len -= 3;

    } else if (ngx_strncmp(p, "r/m", 3) == 0) {
        scale = 60;
        len -= 3;
    }
    else {
        ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                           "invalid rate \"%V\"", &value[1]);
        return NGX_CONF_ERROR;
    }

    rate = ngx_atoi(value[1].data, len);
    if (rate <= 0) {
        ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                           "invalid rate \"%V\"", &value[1]);
        return NGX_CONF_ERROR;
    }

    srv_conf->limit_rate = rate / scale;

    return NGX_CONF_OK;
}


static ngx_int_t
ngx_http_stgw_guard_init(ngx_conf_t *cf)
{
    ngx_http_handler_pt        *h;
    ngx_http_core_main_conf_t  *cmcf;

    cmcf = ngx_http_conf_get_module_main_conf(cf, ngx_http_core_module);

    h = ngx_array_push(&cmcf->phases[NGX_HTTP_PREACCESS_PHASE].handlers);
    if (h == NULL) {
        ngx_conf_log_error(NGX_LOG_ERR, cf, 0, "push ngx_http_stgw_guard_http_handler return NULL");
        return NGX_ERROR;
    }

    *h = ngx_http_stgw_guard_http_handler;

    return NGX_OK;
}


#if 0
static ngx_file_t                   ngx_meminfo_file;
static ngx_file_t                   ngx_netinfo_file;
#endif
static ngx_file_t                   ngx_cpuinfo_file;

#if 0
static u_char g_net_buf[131070];

static ngx_int_t
ngx_http_stgw_guard_getmeminfo(ngx_http_stgw_guard_meminfo_t *meminfo, ngx_log_t *log)
{
    u_char              content[1024];
    ngx_fd_t            fd;
    ssize_t             n;
    size_t tmp;
    ngx_memzero(meminfo, sizeof(ngx_http_stgw_guard_meminfo_t));
    ngx_memzero(content,sizeof(content));
    if (ngx_meminfo_file.fd == 0) {

        fd = ngx_open_file(NGX_MEMINFO_FILE, NGX_FILE_RDONLY,
                           NGX_FILE_OPEN,
                           NGX_FILE_DEFAULT_ACCESS);

        if (fd == NGX_INVALID_FILE) {
            ngx_log_error(NGX_LOG_ERR, log, ngx_errno,
                          ngx_open_file_n " \"%s\" failed",
                          NGX_MEMINFO_FILE);

            return NGX_ERROR;
        }

        ngx_meminfo_file.name.data = (u_char *) NGX_MEMINFO_FILE;
        ngx_meminfo_file.name.len = ngx_strlen(NGX_MEMINFO_FILE);

        ngx_meminfo_file.fd = fd;
    }

    ngx_meminfo_file.log = log;
    n = ngx_read_file(&ngx_meminfo_file, content, sizeof(content) - 1, 0);
    if (n == NGX_ERROR) {
        ngx_log_error(NGX_LOG_ERR, log, ngx_errno,
                      ngx_read_file_n " \"%s\" failed",
                      NGX_MEMINFO_FILE);

        return NGX_ERROR;
    }

    if (ngx_strstr( content, "MemAvailable" ) ) {
        n = sscanf((char *) content, "MemTotal: %ld %*s\nMemFree: %ld %*s\nMemAvailable: %ld %*s\nBuffers: %ld %*s\nCached: %ld",
                    &meminfo->total, &meminfo->free, &tmp, &meminfo->buffers, &meminfo->cached );

        if( n != 5 ) {
            ngx_log_error( NGX_LOG_ERR, log, ngx_errno,
                           ngx_read_file_n " \"%s\" parse failed",
                           NGX_MEMINFO_FILE );
            return NGX_ERROR;
        }

        return NGX_OK;
    }

    n = sscanf((char *) content, "MemTotal: %ld %*s\nMemFree: %ld %*s\nBuffers: %ld %*s\nCached: %ld",
                &meminfo->total, &meminfo->free, &meminfo->buffers, &meminfo->cached );

    if( n != 4 ) {
        ngx_log_error( NGX_LOG_ERR, log, ngx_errno,
                       ngx_read_file_n " \"%s\" parse failed",
                       NGX_MEMINFO_FILE );
        return NGX_ERROR;
    }

    return NGX_OK;
}
#endif

static ngx_int_t
ngx_http_stgw_guard_getcpuinfo(ngx_http_stgw_guard_cpuinfo_t *cpuinfo, ngx_log_t *log)
{
    u_char              content[1024];
    ngx_fd_t            fd;
    ssize_t             n;
    ngx_memzero(cpuinfo, sizeof(ngx_http_stgw_guard_cpuinfo_t));
    ngx_memzero(content,sizeof(content));

    if (ngx_cpuinfo_file.fd == 0) {

        fd = ngx_open_file(NGX_CPUINFO_FILE, NGX_FILE_RDONLY,
                           NGX_FILE_OPEN,
                           NGX_FILE_DEFAULT_ACCESS);

        if (fd == NGX_INVALID_FILE) {
            ngx_log_error(NGX_LOG_ERR, log, ngx_errno,
                          ngx_open_file_n " \"%s\" failed",
                          NGX_CPUINFO_FILE);

            return NGX_ERROR;
        }

        ngx_cpuinfo_file.name.data = (u_char *) NGX_CPUINFO_FILE;
        ngx_cpuinfo_file.name.len = ngx_strlen(NGX_CPUINFO_FILE);

        ngx_cpuinfo_file.fd = fd;
    }

    ngx_cpuinfo_file.log = log;
    n = ngx_read_file(&ngx_cpuinfo_file, content, sizeof(content) - 1, 0);
    if (n == NGX_ERROR) {
        ngx_log_error(NGX_LOG_ERR, log, ngx_errno,
                      ngx_read_file_n " \"%s\" failed",
                      NGX_CPUINFO_FILE);
        ngx_cpuinfo_file.fd = 0;
        return NGX_ERROR;
    }
    n = sscanf((char *) content, "cpu %ld %ld %ld %ld %ld %ld %ld",
        &cpuinfo->user, &cpuinfo->nice, &cpuinfo->sys, &cpuinfo->idle,
        &cpuinfo->iowait, &cpuinfo->hardirq, &cpuinfo->softirq);
    if(n != 7) {
        ngx_log_error(NGX_LOG_ERR, log, ngx_errno,
                      ngx_read_file_n " \"%s\" parse failed",
                      NGX_CPUINFO_FILE);

        return NGX_ERROR;
    }
    return NGX_OK;
}

#if 0
static ngx_int_t
ngx_http_stgw_guard_getnetinfo(ngx_http_stgw_guard_netinfo_t *netinfo, ngx_log_t *log, ngx_str_t *net_dev)
{

    ngx_fd_t            fd;
    ssize_t             n;
    char *p;
    p = NULL;
    size_t              packets,errs,drop,fifo,frame,compressed,multicast;
    ngx_time_t          *tp;
    int num;
    if (net_dev->len == 0) {

        ngx_log_error(NGX_LOG_ERR, log, 0,
            "net net_dev is null");
        return NGX_ERROR;
    }
    ngx_memzero(netinfo, sizeof(ngx_http_stgw_guard_netinfo_t));
    if (ngx_netinfo_file.fd == 0) {

        fd = ngx_open_file(NGX_NETINFO_FILE, NGX_FILE_RDONLY,
                           NGX_FILE_OPEN,
                           NGX_FILE_DEFAULT_ACCESS);

        if (fd == NGX_INVALID_FILE) {
            ngx_log_error(NGX_LOG_ERR, log, ngx_errno,
                          ngx_open_file_n " \"%s\" failed",
                          NGX_NETINFO_FILE);

            return NGX_ERROR;
        }

        ngx_netinfo_file.name.data = (u_char *) NGX_NETINFO_FILE;
        ngx_netinfo_file.name.len = ngx_strlen(NGX_NETINFO_FILE);

        ngx_netinfo_file.fd = fd;
    }

    ngx_netinfo_file.log = log;
    ngx_netinfo_file.offset = 0;
    num = 0;
    while(1) {
        n = ngx_read_file(&ngx_netinfo_file, g_net_buf + ngx_netinfo_file.offset,
            4096, ngx_netinfo_file.offset);
        if (n == NGX_ERROR) {
            ngx_log_error(NGX_LOG_ERR, log, ngx_errno,
                          ngx_read_file_n " \"%s\" failed",
                          NGX_NETINFO_FILE);

            ngx_netinfo_file.fd = 0;
            return NGX_ERROR;
        }
        num ++;
        ngx_log_error(NGX_LOG_DEBUG, log, ngx_errno,
                       "num:%d, n:%d, offset:%d", num, n, ngx_netinfo_file.offset);
        if (n == 0 || ngx_netinfo_file.offset >= (off_t) sizeof(g_net_buf) - 4096 || num > 20)
            break;
    }
    p = ngx_strstr(g_net_buf, net_dev->data);
    if(p == NULL) {
        ngx_log_error(NGX_LOG_WARN, log, ngx_errno,
                       "find net_dev \"%s\" fail, offset %d", net_dev->data, ngx_netinfo_file.offset);
        return NGX_ERROR;
    }
    p = p + net_dev->len;
    n = sscanf(p, ":%ld %ld %ld %ld %ld %ld %ld %ld %ld",
        &netinfo->recv, &packets, &errs, &drop, &fifo, &frame, &compressed, &multicast, &netinfo->send);
    if(n != 9) {
        ngx_log_error(NGX_LOG_WARN, log, ngx_errno,
                      ngx_read_file_n " \"%s\" parse failed",
                      NGX_NETINFO_FILE);

        return NGX_ERROR;
    }
    tp = ngx_timeofday();
    netinfo->msec = (ngx_msec_t) (tp->sec * 1000 + tp->msec);
    return NGX_OK;
}
#endif

```