[TOC]
## 核心思想
tRPC核心框架是采用基于接口编程的思想, 通过把框架功能抽象成一系列的插件组件, 注册到插件工厂, 并由插件工厂实例化插件. tRPC框架负责串联这些插件组件, 拼装出完整的框架功能. 我们可以把插件模型分为以下三层:
1. 框架设计层: 框架只定义标准接口，没有任何插件实现，与平台完全解耦；
2. 插件实现层: 将插件按框架标准接口封装起来即可实现框架插件；
3. 用户使用层: 业务开发只需要引入自己需要的插件即可，按需索取，拿来即用。


## 插件的标准接口
插件的标准接口是指标准的plugins包提供的插件注册和启动能力。
### Register 
trpc的plugins包提供了Register方法，该方法会构建全局的插件map：
``` 
plugin type => { plugin name => plugin factory }
插件类型 ： {插件名称 ： 插件工厂(接口)} 例如 database：{mysql ： mysqlFactory}
plugins = make(map[string]map[string]Factory) 
```
### Factory
plugins包还提供了一个插件工厂的接口Factory。
``` 
// Factory 插件工厂统一抽象，外部插件需要实现该接口，通过该工厂接口生成具体的插件并注册到具体的插件类型里面。
type Factory interface {
    // Type 插件的类型 如 selector log config tracing
    Type() string
    // Setup 根据配置项节点装载插件，需要用户自己先定义好具体插件的配置数据结构
	// Decode是一个解析配置文件的解析接口。框架实现了解析Yaml.Node的decode方法
    Setup(name string, dec Decoder) error
}
```
所有插件需要实现这两个方法，其中比较重要的是Setup方法。
每个插件需要自己去实现Setup方法，这个方法实现了插件启动的过程。
Setup方法一般包括了 
- 1.解析和校验配置文件 
- 2.初始化工作（创建所需资源）
- 3.注册插件到`trpc提供的标准功能插件 （config, log, filter, codec, selector, transport）`
### 注册插件
在程序的main包下，需要引用插件包，每个插件包引用了plugins包，并且有init方法和Factory接口实现方法，init方法会调用plugins包的Register方法，将自己的插件type和name注册到全局的插件map中。例如，007插件：
``` 
func init() {
    plugin.Register("m007", &m007Plugin{})
}
// m007Plugin 007监控插件工厂，实例化007插件，注册metrics，模调拦截器
type m007Plugin struct{}

// Type m007Plugin插件类型
func (p *m007Plugin) Type() string {
    return "metrics"
}

// Setup m007Plugin插件初始化
func (p *m007Plugin) Setup(name string, decoder plugin.Decoder) error {
    1.解析和校验配置文件 
	cfg := Config{}
    decoder.Decode(&cfg)
	2.初始化工作（创建所需资源pcgmonitor）
    err = pcgmonitor.Setup(&pcgmonitor.FrameSvrSetupInfo{
        FrameSvrInfo: pcgmonitor.FrameSvrInfo{
            App:          cfg.AppName,
            Server:       cfg.ServerName,
            IP:           cfg.IP,
            Container:    cfg.ContainerName,
            ConSetId:     cfg.ContainerSetId,
            Version:      cfg.Version,
            PhysicEnv:    cfg.PhysicEnv,
            UserEnv:      cfg.UserEnv,
            FrameCode:    cfg.FrameCode,
            DebugLogOpen: cfg.DebugLogOpen,
        },
        PolarisInfo: pcgmonitor.PolarisInfo{ // 拉007配置时使用, 不填，使用默认值
            Addrs: cfg.PolarisAddrs,
            Proto: cfg.PolarisProto,
        },
    })
	3.注册插件到`trpc提供的标准功能插件 （config, log, filter, codec, selector, transport）`
    // 3.1注册metrics
    metrics.RegisterMetricsSink(&M007Sink{})
    // 3.2注册客户端和服务端的拦截器
    filter.Register(name, PassiveModuleCallServerFilter, ActiveModuleCallClientFilter)
    return nil
}
```
因为很多插件都是注册到trpc提供的标准功能插件 （config, log, filter, codec, selector, transport）中，所以有必要对标准功能插件进行一些讲解。
## 标准功能插件
![enter image description here](/tencent/api/attachments/s3/url?attachmentid=4872830)
如图所示，是一个RPC框架完成一次RPC服务所要经历的流程。
该流程中需要使用到selector的服务发现、负载均衡、服务注册，发起请求和收到响应的filter，发起请求和收到请求的序列化和反序列化Codec，发送和接受网络数据的Transport，以及使用过程中的日志log和配置config。
这些组件是RPC框架所必须的，因此trpc如果要能够使用，必须要完成这些组件的一种默认实现方式。
此外，为了支持用户自定义上述组件的实现方式（例如自定义filter，自定义日志打印），每个组件还需要对外提供通用接口，来实现插件扩展的能力。
trpc针对RPC框架中必要的部分提供了扩展能力。

### filter插件实现和调用方法
拦截器原理关键点在于拦截器的 `功能`  以及  `顺序` 。
功能：拦截器可以拦截到接口的请求和响应，并对请求，响应，上下文进行处理（用通俗的语言阐述也就是 可以在 `请求接受前` 做一些事情， `请求处理后` 做一些事情），因此，拦截器从功能上说是分为两个部分的 前置（业务逻辑处理前） 和 后置（业务逻辑处理后）。
顺序：如下图所示，拦截器是有明确的顺序性，根据拦截器的注册顺序依次执行前置部分逻辑，并逆序执行拦截器的后置部分。
![enter image description here](/tencent/api/attachments/s3/url?attachmentid=4872831)
拦截器分为客户端拦截器和服务端拦截器
客户端拦截器是对rpc框架在`发起调用和收到响应`前后执行的动作
服务端拦截器是对rpc框架在`收到请求和返回响应`前后执行的动作
例如统计客户端调用的时间和服务端处理的时间。
https://iwiki.woa.com/pages/viewpage.action?pageId=274914183

#### Filter包：
- handlerFunc是业务逻辑处理，作为filter的最里层。
客户端的handlerFunc为callFunc，内部包括codec打解包，transport网络通信，以该函数前后为拦截器入口。
``` 
// HandleFunc 过滤器（拦截器）函数接口
type HandleFunc func(ctx context.Context, req interface{}, rsp interface{}) (err error)
```
- Filter是过滤的实现接口，用户可以自己实现Filter。
``` 
// Filter 过滤器（拦截器），根据dispatch处理流程进行上下文拦截处理
type Filter func(ctx context.Context, req interface{}, rsp interface{}, f HandleFunc) (err error)
ctx是请求的上下文
req是请求结构体
rsp是响应结构体
f是客户端或服务端的业务逻辑（根据REQ，生成RSP）的代码
```
- Register是通用注册接口
``` 
serverFilters = make(map[string]Filter)
clientFilters = make(map[string]Filter)

// Register 通过拦截器名字注册server client拦截器
func Register(name string, serverFilter Filter, clientFilter Filter) {
    lock.Lock()
    serverFilters[name] = serverFilter
    clientFilters[name] = clientFilter
    lock.Unlock()
}
```
#### Filter的调用入口：
- client ：
每个client需要有Invoke方法，Options结构体中存在filter.Chain（[]filter）成员，filter.Chain有handle方法。
``` 
// Invoke 启动后端调用，传入具体业务协议结构体 内部调用codec transport
func (c *client) Invoke(ctx context.Context, reqbody interface{}, rspbody interface{}, opt ...Option) error {
	// 读取配置参数，设置用户输入参数
    opts, err := c.getOptions(msg, opt...) 返回一个Options，Options有一个filter.Chain成员。该成员根据配置文件中filter的名字作为key，从clientFilters读取Filter。
	//filter的调用，c.callFunc主要是客户端的codec和transport流程
    return opts.Filters.Handle(ctx, reqbody, rspbody, c.callFunc(opts)) 
}
```
- handle方法递归的调用所有注册的Filter
``` 
// Chain 链式过滤器
type Chain []Filter

// Handle 链式过滤器递归处理流程
func (fc Chain) Handle(ctx context.Context, req, rsp interface{}, f HandleFunc) (err error) {
    if len(fc) == 0 {
        return f(ctx, req, rsp)
    }
    return fc[0](ctx, req, rsp, func(ctx context.Context, req interface{}, rsp interface{}) (err error) {
        return fc[1:].Handle(ctx, req, rsp, f)
    })
}
```
#### Filter例子：007上报
- 注册 
在007插件的Setup中，执行了
`filter.Register(name, PassiveModuleCallServerFilter, ActiveModuleCallClientFilter)`
- filter实现
PassiveModuleCallServerFilter是本服务作为server的时候，需要执行的逻辑
``` 
// PassiveModuleCallServerFilter  被调模调上报拦截器:上游调用自身, 处理函数结束时上报
func PassiveModuleCallServerFilter(ctx context.Context, req, rsp interface{}, handler filter.HandleFunc) error {
	-------- 处理请求前执行的逻辑 ----------
    begin := time.Now()
	
	-------- 处理请求 --------
    err := handler(ctx, req, rsp)

	-------- 处理请求后执行的逻辑 --------
    msg := trpc.Message(ctx)
    passiveMsg := new(pcgmonitor.PassiveMsg)
    // 统计处理请求的耗时ms
    passiveMsg.Time = float64(time.Since(begin)) / float64(time.Millisecond)
	// pcgmonitor在Setup中进行了初始化
    pcgmonitor.ReportPassive(passiveMsg)
    return err
}
```
ActiveModuleCallClientFilter是本服务作为客户端的时候，需要执行的逻辑
``` 
// ActiveModuleCallClientFilter  主调模调上报拦截器:自身调用下游，下游回包时上报
func ActiveModuleCallClientFilter(ctx context.Context, req, rsp interface{}, handler filter.HandleFunc) error {
	-------- 发送请求前执行的逻辑 ----------
    begin := time.Now()
	
	-------- 发送请求 --------
    err := handler(ctx, req, rsp)

	-------- 收到响应后执行的逻辑 --------
    activeMsg := new(pcgmonitor.ActiveMsg)
    // 发送请求到得到响应的耗时ms
    activeMsg.Time = float64(time.Since(begin)) / float64(time.Millisecond)
	pcgmonitor.ReportActive(activeMsg)
    return err
}
```

### Metric插件实现和调用方法
结合logging、metrics、tracing三大件有助于我们建立起一个更加全面的监控体系。这里是Metric
#### 常见的mertric
- 模调上报 （基于filter实现，详情见007filter）
	- 作为服务端，发响应到上游的时候，上报整个调用信息（主调和被调信息，回包的状态码，耗时等）
	- 作为客户端，收到下游回包的时候，上报整个调用信息（主调和被调信息，回包的状态码，耗时等）
- 属性上报 （基于metric实现）
	- 积累量（请求次数..）
	- 时刻量（QPS..）
	- 多维监控项
#### Metric包：
Metric包实现了一些常用的监控项，比如counter计数器，Gauge当前值，Timer时间花销，Histogram直方图。
counter有增加和减少的操作
Gauge有set的操作
Timer有计算时间段的操作
Histogram有增加sample的操作
操作以上数据的时候，都会有遍历所有的report插件来上报数据。因此插件只需要实现report方法就可以了（例如007，prometheus插件）。
#### runtime Metric插件
代码 ：`git.code.oa.com/trpc-go/trpc-metrics-runtime`
指标 ：
- golang协程数，线程数，CPU核心数
- 磁盘使用情况，内存使用情况，cpu使用情况
- 进程数，打开文件数
- tcp 打开socket数

### Tracing插件实现和调用方法
结合logging、metrics、tracing三大件有助于我们建立起一个更加全面的监控体系。这里是Tracing
#### pspanID，spanID和tracingID
https://cloud.tencent.com/developer/article/1832719
- tracingID会在所有服务中传递，并且不会改变
- spanID是当前服务生成的ID，pspanID是当前服务的上游服务生成的ID。


### 第三方Client
用户可以利用  `git.code.oa.com/trpc-go/trpc-go/client`  的Client接口操作网络调用
并实现自己的codec插件。例如实现和Taf互通的codec插件。
例如：https://iwiki.woa.com/pages/viewpage.action?pageId=188159612