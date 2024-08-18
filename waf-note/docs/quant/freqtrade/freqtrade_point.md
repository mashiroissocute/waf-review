## eps

eps允许在回测的时候，为每个pair允许多个trade。即使当前pair在trade中，当信号来临时，仍然为该pair添加trade。

eps仅仅支持在回测中进行实验。并不支持实盘。


## adjust_trade_position()
在实盘中，进允许通过adjust_trade_position回调来调整仓位。
adjust_trade_position进调整一个trade的仓位，并不像eps那样开启多个trade。


在实盘中，应该是使用adjust_trade_position来追加金额或者减少金额。
同时adjust_trade_position支持回测。 我们应该放弃掉使用eps


## custom_stoploss

在dca订单中，stoploss不会跟着dca订单走，需要使用custom_stoploss去调整dca订单后的stoploss.


注意custom_stoploss只能相较于stoploss向上移动，否则无效。

并且custom_stoploss和tailing stop loss不能一起使用。

```python
use_custom_stoploss = True

def custom_stoploss(self, pair: str, trade: 'Trade', current_time: datetime,
                    current_rate: float, current_profit: float, after_fill: bool, 
                    **kwargs) -> Optional[float]:

    if after_fill: 
        # After an additional order, start with a stoploss of 10% below the new open rate
        return stoploss_from_open(0.10, current_profit, is_short=trade.is_short, leverage=trade.leverage)
    # Make sure you have the longest interval first - these conditions are evaluated from top to bottom.
    if current_time - timedelta(minutes=120) > trade.open_date_utc:
        return -0.05
    elif current_time - timedelta(minutes=60) > trade.open_date_utc:
        return -0.10
    return None
```

