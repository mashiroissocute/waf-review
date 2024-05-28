# eps

eps允许在回测的时候，为每个pair允许多个trade。即使当前pair在trade中，当信号来临时，仍然为该pair添加trade。

eps仅仅支持在回测中进行实验。并不支持实盘。


# adjust_trade_position()
在实盘中，进允许通过adjust_trade_position回调来调整仓位。
adjust_trade_position进调整一个trade的仓位，并不像eps那样开启多个trade。


在实盘中，应该是使用adjust_trade_position来追加金额或者减少金额。
同时adjust_trade_position支持回测。 我们应该放弃掉使用eps
