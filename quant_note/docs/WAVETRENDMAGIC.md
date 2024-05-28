--------------------------------------------------------------------+--------+-----------+----------+--------------------------+--------------+-------------------------------+-----------------+----------------+-------------------------------+
|   Best |     Epoch |   Trades |    Win  Draw  Loss  Win% |   Avg profit |                        Profit |    Avg duration |      Objective |           Max Drawdown (Acct) |
|--------+-----------+----------+--------------------------+--------------+-------------------------------+-----------------+----------------+-------------------------------|
| * Best |    1/1000 |      245 |    171     0    74  69.8 |        2.21% |   197280.064 USDT (9,864.00%) | 0 days 03:01:00 | -197,280.06394 |     14559.430 USDT   (13.57%) |
| * Best |   20/1000 |      245 |    172     0    73  70.2 |        2.21% |   198461.474 USDT (9,923.07%) | 0 days 03:03:00 | -198,461.47415 |     14559.430 USDT   (13.26%) |
|   Best |   35/1000 |      264 |    184     0    80  69.7 |        2.09% |  209177.108 USDT (10,458.86%) | 0 days 03:00:00 | -209,177.10797 |     14559.430 USDT   (13.29%) |
|   Best |  218/1000 |      244 |    171     0    73  70.1 |        2.25% |  220046.808 USDT (11,002.34%) | 0 days 03:07:00 | -220,046.80774 |     14559.430 USDT   (11.70%) |
|   Best |  234/1000 |      245 |    172     0    73  70.2 |        2.24% |  220449.806 USDT (11,022.49%) | 0 days 03:09:00 | -220,449.80628 |     14559.430 USDT   (11.69%) |
|   Best |  241/1000 |      245 |    172     0    73  70.2 |        2.25% |  220674.349 USDT (11,033.72%) | 0 days 03:10:00 | -220,674.34871 |     14559.430 USDT   (11.68%) |
|   Best |  283/1000 |      247 |    173     0    74  70.0 |        2.24% |  224604.158 USDT (11,230.21%) | 0 days 03:06:00 | -224,604.15796 |     14559.430 USDT   (11.50%) |
|   Best |  337/1000 |      245 |    172     0    73  70.2 |        2.26% |  226120.012 USDT (11,306.00%) | 0 days 03:06:00 | -226,120.01216 |     14559.430 USDT   (11.43%) |
Epochs ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╸                                                  696/1000  70% • 0:37:45 • 0:09:56
User interrupted..

Best result:

   337/1000:    245 trades. 172/0/73 Wins/Draws/Losses. Avg profit   2.26%. Median profit   4.02%. Total profit 226120.01216133 USDT (11306.00%). Avg duration 3:06:00 min. Objective: -226120.01216


    # Buy hyperspace params:
    buy_params = {
        "entry_long_adx_threshold": 27,
        "entry_long_adx_threshold_trend": 61,
        "entry_long_derta": -8,
        "entry_long_derta_trend": -10,
        "exit_long_adx_threshold": 27,
        "exit_long_derta": 23,
    }

    # Sell hyperspace params:
    sell_params = {
        "entry_short_adx_threshold": 21,  # value loaded from strategy
        "entry_short_adx_threshold_trend": 31,  # value loaded from strategy
        "entry_short_derta": -30,  # value loaded from strategy
        "entry_short_derta_trend": -28,  # value loaded from strategy
        "exit_short_adx_threshold": 36,  # value loaded from strategy
        "exit_short_derta": 12,  # value loaded from strategy
    }

    # ROI table:  # value loaded from strategy
    minimal_roi = {
        "0": 0.087,
        "87": 0.072,
        "207": 0.04,
        "547": 0
    }

    # Stoploss:
    stoploss = -0.234  # value loaded from strategy

    # Trailing stop:
    trailing_stop = False  # value loaded from strategy
    trailing_stop_positive = None  # value loaded from strategy
    trailing_stop_positive_offset = 0.0  # value loaded from strategy
    trailing_only_offset_is_reached = False  # value loaded from strategy
    

    # Max Open Trades:
    max_open_trades = 1  # value loaded from strategy
Result for strategy WTLBStrategyV2
=============================================================== BACKTESTING REPORT ==============================================================
|            Pair |   Entries |   Avg Profit % |   Cum Profit % |   Tot Profit USDT |   Tot Profit % |   Avg Duration |   Win  Draw  Loss  Win% |
|-----------------+-----------+----------------+----------------+-------------------+----------------+----------------+-------------------------|
| MAGIC/USDT:USDT |       245 |           2.26 |         552.56 |        226120.012 |       11306.00 |        3:06:00 |   172     0    73  70.2 |
|           TOTAL |       245 |           2.26 |         552.56 |        226120.012 |       11306.00 |        3:06:00 |   172     0    73  70.2 |
======================================================= LEFT OPEN TRADES REPORT ========================================================
|   Pair |   Entries |   Avg Profit % |   Cum Profit % |   Tot Profit USDT |   Tot Profit % |   Avg Duration |   Win  Draw  Loss  Win% |
|--------+-----------+----------------+----------------+-------------------+----------------+----------------+-------------------------|
|  TOTAL |         0 |           0.00 |           0.00 |             0.000 |           0.00 |           0:00 |     0     0     0     0 |
=========================================================== ENTER TAG STATS ===========================================================
|   TAG |   Entries |   Avg Profit % |   Cum Profit % |   Tot Profit USDT |   Tot Profit % |   Avg Duration |   Win  Draw  Loss  Win% |
|-------+-----------+----------------+----------------+-------------------+----------------+----------------+-------------------------|
| TOTAL |       245 |           2.26 |         552.56 |        226120.012 |       11306.00 |        3:06:00 |   172     0    73  70.2 |
===================================================== EXIT REASON STATS =====================================================
|   Exit Reason |   Exits |   Win  Draws  Loss  Win% |   Avg Profit % |   Cum Profit % |   Tot Profit USDT |   Tot Profit % |
|---------------+---------+--------------------------+----------------+----------------+-------------------+----------------|
|           roi |     167 |    167     0     0   100 |           6.14 |        1024.89 |         357168    |        1024.89 |
|   exit_signal |      72 |      5     0    67   6.9 |          -4.6  |        -330.84 |        -124391    |        -330.84 |
|     stop_loss |       6 |      0     0     6     0 |         -23.58 |        -141.49 |          -6657.45 |        -141.49 |
==================== SUMMARY METRICS ====================
| Metric                      | Value                   |
|-----------------------------+-------------------------|
| Backtesting from            | 2023-01-25 03:00:00     |
| Backtesting to              | 2023-11-03 07:30:00     |
| Max open trades             | 1                       |
|                             |                         |
| Total/Daily Avg Trades      | 245 / 0.87              |
| Starting balance            | 2000 USDT               |
| Final balance               | 228120.012 USDT         |
| Absolute profit             | 226120.012 USDT         |
| Total profit %              | 11306.00%               |
| CAGR %                      | 45883.47%               |
| Sortino                     | 8.14                    |
| Sharpe                      | 4.79                    |
| Calmar                      | 6700.04                 |
| Profit factor               | 2.67                    |
| Expectancy (Ratio)          | 922.94 (0.50)           |
| Trades per day              | 0.87                    |
| Avg. daily profit %         | 40.09%                  |
| Avg. stake amount           | 51827.064 USDT          |
| Total trade volume          | 12697630.739 USDT       |
|                             |                         |
| Long / Short                | 4 / 241                 |
| Total profit Long %         | 708.27%                 |
| Total profit Short %        | 10597.73%               |
| Absolute profit Long        | 14165.4 USDT            |
| Absolute profit Short       | 211954.612 USDT         |
|                             |                         |
| Best Pair                   | MAGIC/USDT:USDT 552.56% |
| Worst Pair                  | MAGIC/USDT:USDT 552.56% |
| Best trade                  | MAGIC/USDT:USDT 8.84%   |
| Worst trade                 | MAGIC/USDT:USDT -23.61% |
| Best day                    | 33798.588 USDT          |
| Worst day                   | -8543.353 USDT          |
| Days win/draw/lose          | 70 / 167 / 35           |
| Avg. Duration Winners       | 3:03:00                 |
| Avg. Duration Loser         | 3:15:00                 |
| Max Consecutive Wins / Loss | 12 / 4                  |
| Rejected Entry signals      | 0                       |
| Entry/Exit Timeouts         | 0 / 0                   |
|                             |                         |
| Min balance                 | 2105.886 USDT           |
| Max balance                 | 235210.693 USDT         |
| Max % of account underwater | 47.80%                  |
| Absolute Drawdown (Account) | 11.43%                  |
| Absolute Drawdown           | 14559.43 USDT           |
| Drawdown high               | 125354.879 USDT         |
| Drawdown low                | 110795.449 USDT         |
| Drawdown Start              | 2023-09-12 15:00:00     |
| Drawdown End                | 2023-09-29 07:30:00     |
| Market change               | -53.62%                 |
=========================================================

Backtested 2023-01-25 03:00:00 -> 2023-11-03 07:30:00 | Max open trades : 1
=========================================================================== STRATEGY SUMMARY ===========================================================================
|       Strategy |   Entries |   Avg Profit % |   Cum Profit % |   Tot Profit USDT |   Tot Profit % |   Avg Duration |   Win  Draw  Loss  Win% |              Drawdown |
|----------------+-----------+----------------+----------------+-------------------+----------------+----------------+-------------------------+-----------------------|
| WTLBStrategyV2 |       245 |           2.26 |         552.56 |        226120.012 |       11306.00 |        3:06:00 |   172     0    73  70.2 | 14559.43 USDT  11.43% |
========================================================================================================================================================================

For more details, please look at the detail tables above
