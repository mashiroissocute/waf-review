
## 最大回撤



## Expectancy Expectancy Ratio
在交易策略中，Expectancy 和 Expectancy Ratio 是衡量交易策略盈利能力的指标

1. **Expectancy**：期望值是指每笔交易的平均收益。计算公式如下：

```
Expectancy = (Win Rate * Average Win) - (Loss Rate * Average Loss)
```

其中：
- Win Rate 是胜率，即盈利交易的比例。
- Average Win 是平均盈利金额。
- Loss Rate 是亏损率，即亏损交易的比例。
- Average Loss 是平均亏损金额。

2. **Expectancy Ratio**：期望值比率是指每单位风险的平均收益。计算公式如下：

```
Expectancy Ratio = Expectancy / (Average Loss * sqrt(Trade Count))
```

其中：
- Expectancy 是期望值。
- Average Loss 是平均亏损金额。
- Trade Count 是交易次数。
- sqrt() 是平方根函数。

以下是一个使用 Pandas 计算 Expectancy 和 Expectancy Ratio 的示例：

```python
import pandas as pd
import numpy as np

# 创建一个示例 DataFrame，包含交易结果
data = {'profit': [100, -50, 200, -100, 300, -150]}
trades = pd.DataFrame(data)

# 计算胜率、亏损率、平均盈利和平均亏损
win_rate = (trades['profit'] > 0).mean()
loss_rate = 1 - win_rate
average_win = trades[trades['profit'] > 0]['profit'].mean()
average_loss = -trades[trades['profit'] < 0]['profit'].mean()

# 计算期望值
expectancy = (win_rate * average_win) - (loss_rate * average_loss)

# 计算期望值比率
trade_count = len(trades)
expectancy_ratio = expectancy / (average_loss * np.sqrt(trade_count))

# 打印结果
print(f"Expectancy: {expectancy}")
print(f"Expectancy Ratio: {expectancy_ratio}")
```

请注意，这个示例仅用于演示如何计算 Expectancy 和 Expectancy Ratio。在实际应用中，你需要根据你的交易策略和数据来计算这些指标。

## GAGR，即复合年均增长率（Compound Annual Growth Rate）

```python
(final_balance / starting_balance) ** (1 / (days_passed / 365)) - 1
```

## sharpe

```python

total_profit = trades['profit_abs'] / starting_balance
days_period = max(1, (max_date - min_date).days)

expected_returns_mean = total_profit.sum() / days_period
up_stdev = np.std(total_profit)

if up_stdev != 0:
    sharp_ratio = expected_returns_mean / up_stdev * np.sqrt(365)
else:
    # Define high (negative) sharpe ratio to be clear that this is NOT optimal.
    sharp_ratio = -100

# print(expected_returns_mean, up_stdev, sharp_ratio)
return sharp_ratio
```

## Sortino 

Sortino 是去掉盈利部分的sharpe，因此Sortino主要考虑下行风险。
```python
total_profit = trades['profit_abs'] / starting_balance
days_period = max(1, (max_date - min_date).days)

expected_returns_mean = total_profit.sum() / days_period

down_stdev = np.std(trades.loc[trades['profit_abs'] < 0, 'profit_abs'] / starting_balance)

if down_stdev != 0 and not np.isnan(down_stdev):
    sortino_ratio = expected_returns_mean / down_stdev * np.sqrt(365)
else:
    # Define high (negative) sortino ratio to be clear that this is NOT optimal.
    sortino_ratio = -100

# print(expected_returns_mean, down_stdev, sortino_ratio)
return sortino_ratio
```


## Calmar比率
Calmar比率 = (年化收益率 - 无风险利率) / 最大回撤

