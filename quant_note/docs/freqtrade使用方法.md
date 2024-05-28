# read note
## Introduction
- Develop Strategy
- trading bot (select exchange and market-pairs by config, WebUI and Telegram)
- data download
- backtesting
- parameters optimiztion (Find the best parameters for your strategy using hyperoptimization which employs machining learning methods)
- analyze (plotting)
- money management
- Edge 

## Basic
- All profit calculations include fees
- Pair naming : CCXT    
### Bot execution logic (bot loop  every few seconds)
- 1.fetch open trades from persistence
- 2.select tradable market-pairs 
- 3.download K line
- 4.exectute strategy(indicator、entry、exit)
- 5.check timeout for open orders
- 6.check existing postition
- 7.check exit orders
- 8.check enter orders
- 9.check trade count

### Backtesing / Hyperopt execution logic
- 1.load K line
- 2.exectute strategy(indicator、entry、exit)
- 3.run bot execution logic, **loops per candle**

## Configuration
###  Json File
#### Single File
指定配置文件`-c/--config` 
生成配置文件`freqtrade new-config --config config.json`

#### Multiple  Files
使用add_config_files 指定多个配置文件
``` json title="user_data/config.json"
"add_config_files": [
    "config1.json",
    "config-private.json"
]
```
冲突项将被合并，并且按照最后载入生效。
例如在config-private.json中定义了`exchange.key`，那么在config.json中定义的`exchange.key`将被忽略。
该特性可以用于分享config的同时，保留自己的隐私。

### Environment Variables
比Json配置文件的优先级更高，例如：
`export FREQTRADE__EXCHANGE__KEY=<yourExchangeKey>` 
将json file重的exchange.key进行替换。
所有的配置项都可以基于环境变量进行替换。

### Stratagy Parameters
有些参数可以通过stratagy代码中直接进行设置，例如minimal_roi、timeframe等等

### bash line Parameters
有些参数可以通过命令行直接设置，例如timeframe等参数

### prevalence
bash line > env var > json file > stratage code 

### config parameter explaination
refer to [the config explaination documentation](config-example.md).



## 关于仓位、余额、交易价格
max trades是最多同时存在多少交易中的交易对
balance是可用账户余额，通过tradable_balance_ratio来控制单bot能使用的余额百分比。
amount是单交易对资金上限，直接指定数量 或者 unlimited表示balance / max_trades
position是单交易对仓位，控制对amount的使用，例如50%。通过回调函数custom_stake_amount设置本次操作的仓位。
price是交易价格，分为entry price 和 exit price。
price可以配置为order book price(区分买入方向和卖出方向) 、 market price 、 limit price






## 策略
通过添加指标和交易规则来定制自己的策略

### 策略骨架
- 指标计算
- entry策略
- exit策略
- 止盈策略
- 止损策略
- 回调函数


### 指标计算
使用populate_indicators()方法，将指标计算结果存入dataframe。
计算指标会消耗CPU和内存，因此最好只计算需要的指标。
官方提供了ta-lib、pandas-ta、technical开箱即用
也可以自己import其他的包或者自己创造指标。


### DataFrame
有时候，需要用到以前的数据时，需要更多的历史数据，所以会用startup_candle_count变量配置。
不然会变成nan值，或者计算错误。


### entry策略
定义enter_long、enter_short、enter_tag
1 表示entry信号
0 表示不做动作

### exit策略
定义exit_long、exit_short、enter_tag
1 表示exit信号
0 表示不做动作


### 止盈
```python
minimal_roi = {
    "40": 0.0,
    "30": 0.01,
    "20": 0.02,
    "0": 0.04
}
```

### 止损
#### 静态止损
基于entry price计算止损价格
```python
stoploss = -0.10
```

#### 动态止损
```python
stoploss = -0.10
trailing_stop = True
```
其初始值为stoploss，但每次资产价格上涨时，该算法都会自动上调止损。止损将调整为始终为观察到的最高价格的-10%。

#### 自定义动态止损

#### 止损与杠杆
止损视作本金损失，在杠杆下也是如此。
假如-10%止损，那么10倍杠杆下，跌1%则触发止损。
保证金100，10倍杠杆下，持仓1000。
下跌1%，即损失10。







### 回调函数
#### 初始化时触发
- bot_start() 
#### 每隔一段时间触发 internals.process_throttle_secs
- bot_loop_start() 
- custom_stake_amount() 默认情况下，开始一笔trade时，头寸为固定的stake_amount配置。该函数运行调整初始头寸规模。
- adjust_trade_position() 可用于执行附加订单，减少或增加头寸规模。例如DCA定投，一般和custom_stake_amount搭配使用
- leverage() 调整杠杆大小
- custom_exit() 自定义exit信号
- custom_stoploss() 自定义止损
- custom_entry_price() and custom_exit_price() 默认情况下，会使用订单簿自动设置订单价格，也可以选择根据策略创建自定义订单价格(仅限价单可用)
- check_entry_timeout() and check_exit_timeout() 默认情况下，订单超时时间由unfilledtimeout配置，也可以选择根据策略创建自定义订单超时时间
- confirm_trade_entry() and confirm_trade_exit() 确认交易进入/退出。这是下订单之前调用的最后一个方法。可用于因为价格不满意等原因终止交易。
- adjust_entry_price() 当新k到来时，刷新/替换限价订单


## 插件
### pairlist插件
通过pairlist插件配置，bot可以交易的加密货币对。
可配置的Pairlist Handlers
- StaticPairList (default, if not configured differently)
- VolumePairList
- ProducerPairList
...


## 下载历史数据
### 下载
freqtrade download-data --config user_data/config.json --timerange 20220101-20231008 --timeframes 4h
- --config 制定jsonw文件，必须要包含exchange配置，使用whitePairs作为要下载的交易对
- --timerange 指定时间范围，格式为YYYYMMDD-YYYYMMDD
- --timeframes 指定时间粒度，1d表示1天，1h表示1小时

### 查看
freqtrade list-data --userdir user_data/


### 转换格式
freqtrade convert-data --format-from feather --format-to json --datadir user_data/data/binance --config user_data/config_backtest.json  
- 从feather格式转换为json格式
- --datadir 指定数据目录
- --config 指定json文件，使用whitePairs作为要转换的交易对


## 回测
freqtrade backtesting --strategy SampleStrategy --config user_data/config_backtest.json 

### notice 
- dry_run_wallet 账户余额
- --config 区分回测和实测
- --timeframe 指定时间区间

### 结果解释
```python

============================================================= BACKTESTING REPORT =============================================================
|       Pair |   Entries |   Avg Profit % |   Cum Profit % |   Tot Profit USDT |   Tot Profit % |     Avg Duration |   Win  Draw  Loss  Win% |
|------------+-----------+----------------+----------------+-------------------+----------------+------------------+-------------------------|
|   ETH/USDT |         4 |           2.50 |           9.99 |            99.985 |           0.50 |  1 day, 12:00:00 |     4     0     0   100 |
|   BTC/USDT |         9 |           0.76 |           6.80 |            68.103 |           0.34 | 2 days, 21:20:00 |     8     0     1  88.9 |
|   WLD/USDT |         2 |          -3.06 |          -6.13 |           -61.364 |          -0.31 |         12:00:00 |     1     0     1  50.0 |
|    OP/USDT |         3 |          -2.71 |          -8.12 |           -81.246 |          -0.41 |   1 day, 8:00:00 |     2     0     1  66.7 |
|   ARB/USDT |         3 |          -2.72 |          -8.17 |           -81.804 |          -0.41 | 3 days, 16:00:00 |     2     0     1  66.7 |
| MAGIC/USDT |         3 |          -2.73 |          -8.18 |           -81.830 |          -0.41 | 2 days, 16:00:00 |     2     0     1  66.7 |
|      TOTAL |        24 |          -0.58 |         -13.80 |          -138.157 |          -0.69 |  2 days, 8:00:00 |    19     0     5  79.2 |

```

- TOTAL是所有的交易对 
- Entries是交易次数


## 绘图
freqtrade plot-dataframe -p BTC/USDT --strategy SampleStrategy --config user_data/config_backtest.json 

## 调参数
freqtrade hyperopt --config config.json --hyperopt-loss ShortTradeDurHyperOptLoss --strategy <strategyname> -e 500 --spaces roi stoploss trailing


