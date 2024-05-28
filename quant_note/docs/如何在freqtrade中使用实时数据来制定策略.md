# 如何在freqtrade中使用实时ticket

我想使用不完整的蜡烛.
Freqtrade 不会为策略提供不完整的蜡烛图。使用不完整的蜡烛将导致重新绘制，从而导致“幽灵”购买策略，这些策略不可能在发生后进行回测和验证。
您可以通过使用dataprovider的订单簿或股票代码方法来使用“当前”市场数据 
但在回溯测试期间不能使用这些方法。

```
if self.dp.runmode.value in ('live', 'dry_run'):
    ticker = self.dp.ticker(metadata['pair'])
    dataframe['last_price'] = ticker['last']
    dataframe['volume24h'] = ticker['quoteVolume']
    dataframe['vwap'] = ticker['vwap']

```