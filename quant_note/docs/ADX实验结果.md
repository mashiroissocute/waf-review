# ADX 实验


## 大而全
一个bot跑所有的200+合约。
重点关注max_open_orders设置。
交易量从大到小 
#### 20230101-20231023
mor: 一起优化
优化器: OnlyProfitHyperOptLoss
结果 ：
|   Best |  217/2000 |      372 |    287     0    85  77.2 |        3.91% | 1265562.069 USDT (6,327,810.34%) | 0 days 11:34:00 | -1,265,562.06899 |    667892.954 USDT   (44.03%) | <br>

参数：
```
    "buy": {
      "entry_adx_threshold": 28
    },
    "sell": {
      "exit_adx_threshold": 53
    },
    "protection": {},
    "roi": {
      "0": 0.672,
      "1739": 0.14,
      "4267": 0.068,
      "6885": 0
    },
    "stoploss": {
      "stoploss": -0.307
    },
    "trailing": {
      "trailing_stop": true,
      "trailing_stop_positive": 0.016,
      "trailing_stop_positive_offset": 0.041999999999999996,
      "trailing_only_offset_is_reached": true
    },
    "max_open_trades": {
      "max_open_trades": 1
    }
```

#### 优化sharp




### 10 pair 市值从大到小 
优势在于币种熟悉 + 波动较小，可以控制回测

```
"BTC/USDT:USDT",
"ETH/USDT:USDT",
"BNB/USDT:USDT",
"XRP/USDT:USDT",
"SOL/USDT:USDT",
"ADA/USDT:USDT",
"DOGE/USDT:USDT",
"TRX/USDT:USDT",
"LINK/USDT:USDT",
"MATIC/USDT:USDT"
```

#### 20230101-20231023
mor: 一起优化
优化器: OnlyProfitHyperOptLoss
结果:
76/2000:    813 trades. 600/0/213 Wins/Draws/Losses. Avg profit   1.70%. Median profit   2.99%. Total profit 340.09343960 USDT (1700.47%). Avg duration 9:27:00 min. Objective: -4.66135
参数
```python
{
  "strategy_name": "ADXDIStrategy",
  "params": {
    "buy": {
      "entry_adx_threshold": 38
    },
    "sell": {
      "exit_adx_threshold": 32
    },
    "protection": {},
    "roi": {
      "0": 0.181,
      "546": 0.086,
      "2411": 0.056,
      "5585": 0
    },
    "stoploss": {
      "stoploss": -0.117
    },
    "trailing": {
      "trailing_stop": true,
      "trailing_stop_positive": 0.011,
      "trailing_stop_positive_offset": 0.028999999999999998,
      "trailing_only_offset_is_reached": true
    },
    "max_open_trades": {
      "max_open_trades": 4
    }
  },
  "ft_stratparam_v": 1,
  "export_time": "2023-10-26 06:38:44.508203+00:00"
}
```

#### 20230601-20231023
#### 优化profit  (OnlyProfitHyperOptLoss)
mor: = 一起优化
优化器: OnlyProfitHyperOptLoss
结果: 415/2000:    215 trades. 185/0/30 Wins/Draws/Losses. Avg profit   1.55%. Median profit   2.45%. Total profit 35970.27598762 USDT (1798.51%). Avg duration 7:41:00 min. Objective: -35970.27599

```python
{
    # Buy hyperspace params:
    buy_params = {
        "entry_adx_threshold": 23,
    }

    # Sell hyperspace params:
    sell_params = {
        "exit_adx_threshold": 35,
    }

    # ROI table:
    minimal_roi = {
        "0": 0.245,
        "646": 0.17,
        "1498": 0.068,
        "5538": 0
    }

    # Stoploss:
    stoploss = -0.237

    # Trailing stop:
    trailing_stop = True
    trailing_stop_positive = 0.01
    trailing_stop_positive_offset = 0.012
    trailing_only_offset_is_reached = True
    

    # Max Open Trades:
    max_open_trades = 1
}
```


#### 优化sharp daily (sharp daily)
mor: = 一起优化
优化器: SharpeHyperOptLossDaily
结果：  228/2000:    266 trades. 216/0/50 Wins/Draws/Losses. Avg profit   1.84%. Median profit   2.79%. Total profit 15108.23432504 USDT ( 755.41%). Avg duration 9:28:00 min. Objective: -6.02408
```python
{
    # Buy hyperspace params:
    buy_params = {
        "entry_adx_threshold": 38,
    }

    # Sell hyperspace params:
    sell_params = {
        "exit_adx_threshold": 29,
    }

    # ROI table:
    minimal_roi = {
        "0": 0.508,
        "1575": 0.298,
        "2410": 0.102,
        "7965": 0
    }

    # Stoploss:
    stoploss = -0.172

    # Trailing stop:
    trailing_stop = True
    trailing_stop_positive = 0.01
    trailing_stop_positive_offset = 0.02
    trailing_only_offset_is_reached = True
    

    # Max Open Trades:
    max_open_trades = 2
}
```




### 3
mor: = pair num
时间: 20230601-20231023




## 小而美
一个bot跑 2-3倍max_open_orders合约。