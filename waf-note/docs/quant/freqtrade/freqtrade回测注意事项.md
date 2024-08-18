## 规避回测中的trailing stoploss
1. 可以用采用informative pair的方法

主要的时间保持较短 例如1m 或者 5m，该时间用于交易 和 避免trailing 误差

informative的时间 例如4h 该时间用于判断大的趋势



2. 采用--time-datail的方法

测方法还要加上最好是把 trailing offset 调大 (tailing offset 最好是大于 stoploss ) 

可以把trailing可以适当调小，相当于到达offset之后，就止盈出来，但是该offset是大于stoploss的。
