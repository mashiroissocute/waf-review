## go的内置log库 log
https://darjun.github.io/2020/02/07/godailylib/log/
### Go Logger的优势和劣势

- 优势
它最大的优点是使用非常简单。我们可以设置任何 `io.Writer` 作为日志记录输出并向其发送要写入的日志。
- 劣势
```
- 仅限基本的日志级别
- 只有一个 `Print` 选项。不支持 `INFO` / `DEBUG` 等多个级别。
- 对于错误日志，它有 `Fatal` 和 `Panic` 
- Fatal日志通过调用 `os.Exit(1)` 来结束程序
- Panic日志在写入日志消息之后抛出一个panic
- 但是它缺少一个ERROR日志级别，这个级别可以在不抛出panic或退出程序的情况下记录错误
- 缺乏日志格式化的能力——例如记录调用者的函数名和行号，格式化日期和时间格式。等等。
- 不提供日志切割的能力。
```

## zap.log
- trpc默认log库 为 zap.loggoer
- supeross同样是使用 zap.sugaredlogger

https://www.liwenzhou.com/posts/Go/zap/
https://segmentfault.com/a/1190000022461706
### 为什么选择Uber-go zap
- 它同时提供了结构化日志记录和printf风格的日志记录
- 它非常的快 
### Zap提供了两种类型的日志记录器— `Sugared Logger` 和 `Logger` 。
在性能很好但不是很关键的上下文中，使用 `SugaredLogger` 。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。
在每一微秒和每一次内存分配都很重要的上下文中，使用 `Logger` 。它甚至比 `SugaredLogger` 更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
### 多writer和自定义协议传输log
http://docs.lvrui.io/2020/03/25/go%E8%AF%AD%E8%A8%80zap%E6%97%A5%E5%BF%97%E8%87%AA%E5%AE%9A%E4%B9%89%E8%BE%93%E5%87%BA/
### zaplog粗讲解
```
# 定义logger
logger := zap.New(zapcore.core, options...)
`zapcore.Core` 需要三个配置—— `Encoder` ， `WriteSyncer` ， `LogLevel`
Encoder : 决定日志格式
WriteSyncer ：解决日志写到哪里去 使用 `zapcore.AddSync()` 函数并且将打开的文件句柄传进去
level ：决定什么级别的日志要写
```




## trpc 日志配置
https://iwiki.woa.com/pages/viewpage.action?pageId=465532424

- trpc yaml文件中共有四种配置项，分别是global ， server ， client ， plugins
- config结构体为：
``` 
type Config struct {
    Global struct {
        Namespace      string `yaml:"namespace"`
        EnvName        string `yaml:"env_name"`
        ContainerName  string `yaml:"container_name"`
        LocalIP        string `yaml:"local_ip"`
        EnableSet      string `yaml:"enable_set"`                 // Y/N，是否启用Set分组，默认N
        FullSetName    string `yaml:"full_set_name"`              // set分组的名字，三段式：[set名].[set地区].[set组名]
        ReadBufferSize *int   `yaml:"read_buffer_size,omitempty"` // 网络收包缓冲区大小(单位B)：<=0表示禁用，不配使用默认值
    }
    Server struct {
        App      string
        Server   string
        BinPath  string `yaml:"bin_path"`
        DataPath string `yaml:"data_path"`
        ConfPath string `yaml:"conf_path"`
        Admin    struct {
            IP           string `yaml:"ip"` // 要绑定的网卡地址, 如127.0.0.1
            Nic          string
            Port         uint16 `yaml:"port"`          // 要绑定的端口号，如80，默认值9028
            ReadTimeout  int    `yaml:"read_timeout"`  // ms. 请求被接受到请求信息被完全读取的超时时间设置，防止慢客户端
            WriteTimeout int    `yaml:"write_timeout"` // ms. 处理的超时时间
            EnableTLS    bool   `yaml:"enable_tls"`    // 是否启用tls
        }
        Network       string           // 针对所有service的network 默认tcp
        Protocol      string           // 针对所有service的protocol 默认trpc
        Filter        []string         // 针对所有service的拦截器
        Service       []*ServiceConfig // 单个service服务的配置
        CloseWaitTime int              `yaml:"close_wait_time"` // ms. 关闭服务时,反注册后到真正停止服务之间的等待时间,来支持无损更新
    }
    Client  ClientConfig
    Plugins plugin.Config
}
```
- 日志属于插件服务，以生产环境为例
``` 
plugins:                                          #插件配置
  log:                                            #日志配置
    default:                                      #默认日志的配置，可支持多输出
      - writer: file                              #本地文件日志
        level: info                               #本地文件滚动日志的级别
        writer_config:
          filename: /usr/local/services/trpc/cloud_api/log/trpc.log                      #本地文件滚动日志存放的路径
          max_size: 10                              #本地文件滚动日志的大小 单位 MB
          max_backups: 10                           #最大日志文件数
          max_age: 7                                #最大日志保留天数
          compress:  false                          #日志文件是否压缩
      - writer: xingyun
        level: error          #大于等于这级别才发送星云告警
        remote_config:
          security_code: bfe99f7c00a611ec9a33525400dbc023  # 可在https://qcloud.oa.com/v3/platform/alarmSystem/alarmAccess申请
```
- 插件config为 ：
``` 
type Config map[string]map[string]yaml.Node
对应配置文件为
{"log":{"default":(file,xingyun) yaml.Node)}
一个日志logger对应多个输出通道。
```
- 注册多个通道的zapcore.core，并生成logger
``` 
func NewZapLogWithCallerSkip(c Config, callerSkip int) Logger {
    cores := make([]zapcore.Core, 0, len(c))
    levels := make([]zap.AtomicLevel, 0, len(c))
    for _, o := range c {
        writer, ok := writers[o.Writer]
        if !ok {
            fmt.Printf("log writer core:%s no registered!\n", o.Writer)
            return nil
        }

        decoder := &amp;Decoder{OutputConfig: &amp;o}
				# decoder，新建decoder，decoder为zapcore.core类型
        if err := writer.Setup(o.Writer, decoder); err != nil {
            fmt.Printf("log writer setup core:%s fail:%v!\n", o.Writer, err)
            return nil
        }

        cores = append(cores, decoder.Core)
        levels = append(levels, decoder.ZapLevel)
    }

    logger := zap.New(
				# 创建多个core的logger，同时往多个通道写日志
        zapcore.NewTee(cores...),
        zap.AddCallerSkip(callerSkip),
        zap.AddCaller(),
    )

    return &amp;zapLog{
        levels: levels,
        logger: logger,
    }
}
```

- trpc-log xingyun日志模块初始化
_ "git.code.oa.com/tsec_p3_tac_dev/plugin/trpc-log-xingyun"  引入这个包，不使用但是会调用init函数
``` 
将xingyun加入到日志插件中
const (
    pluginType = "log"
    pluginName = "xingyun"
)
func init() {
    log.RegisterWriter(pluginName, &amp;XYLoggerPlugin{})
}
```

``` 
// RegisterWriter 注册日志输出writer，支持同时多个日志实现
func RegisterWriter(name string, writer plugin.Factory) {
    writers[name] = writer
}
```

- 可以看出trpc log中调用writer.Setup来定义一个zapcore.core，因此xingyun库要实现一个Setup

``` 
// Setup 根据配置项节点装载插件, 用户自己先定义好具体插件的配置数据结构，通过decoder解析出来，开始实例化具体插件
func (p *XYLoggerPlugin) Setup(name string, configDec plugin.Decoder) error {
    ...代码略
    xyLogger := &amp;XYLogger{
        config:     xyconf,
        LogChannel: make(chan *buffer.Buffer, xyconf.ChannelCapacity),
    }

    encoderCfg := zapcore.EncoderConfig{
        TimeKey:        "Time",
        LevelKey:       "Level",
        NameKey:        "Name",
        CallerKey:      "Caller",
        MessageKey:     "Msg",
        StacktraceKey:  "StackTrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     log.NewTimeEncoder(conf.FormatConfig.TimeFmt),
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }
    encoder := zapcore.NewJSONEncoder(encoderCfg)
		
		#开启了写成去消费buffer的内容，即http post发送到xingyun
    go logConsume(&amp;xyconf, xyLogger.LogChannel)
		
		# 定义zapcore ， 规定了encoder格式和输出通道为xyLogger
    decoder.Core = zapcore.NewCore(
        encoder,
        zapcore.AddSync(xyLogger),
        zap.NewAtomicLevelAt(log.Levels[conf.Level]),
    )
}

#xyLogger输出通道为一个buffer，使用logger打印日志，会将内容填充到buffer中
type XYLogger struct {
    config     XYConfig
    LogChannel chan *buffer.Buffer
}
```





