- ADX突破
- 均线回踩
- supertrend atr止盈止损

效果差

```python
# --- Do not remove these libs ---
from freqtrade.strategy import IStrategy
from typing import Dict, List
from functools import reduce
from pandas import DataFrame
import time
import numpy as np
import pandas as pd
from datetime import datetime
from typing import Optional
from freqtrade.strategy import CategoricalParameter, DecimalParameter, IntParameter
from freqtrade.persistence import PairLocks
# --------------------------------

import talib.abstract as ta
import freqtrade.vendor.qtpylib.indicators as qtpylib


class ATRADXPRICEStrategy24(IStrategy):

    INTERFACE_VERSION: int = 3
    can_short = False

    minimal_roi = {
        "0": 1000
    }
    
    stoploss =  -0.99
    
    # trailing_stop = True
    # trailing_stop_positive = 0.02
    # trailing_stop_positive_offset = 0.1
    # trailing_only_offset_is_reached = True
    

    order_types = {
        'entry': 'market',
        'exit': 'market',
        'stoploss': 'market',
        'stoploss_on_exchange': False
    }

    # Optional order time in force.
    order_time_in_force = {
        'entry': 'GTC',
        'exit': 'GTC'
    }

    plot_config = {
        'main_plot': {
            "ema": {},
            "up":{},
            "down":{},
        },
        'subplots': {
            "dmi": {
                "minus_di": {},
                "plus_di": {}
            }
        },
    }
    
    startup_candle_count = 24
    
    n1 = IntParameter(5, 15, default=10, space="buy")
    n2 = IntParameter(15, 25, default=21, space="buy")
    rollWindow = IntParameter(2, 6, default=4, space="buy")
    adxWindow = IntParameter(7, 21, default=14, space="buy")
    adxThr = IntParameter(15, 35, default=25, space="buy")
    emaThr = IntParameter(5, 55, default=24, space="buy")
    superTrendMultiper = IntParameter(1, 7, default=5, space="sell")
    superTrendWindow = IntParameter(7, 21, default=14, space="sell")



    # Optimal timeframe for the strategy
    timeframe = '4h'

    def populate_indicators(self, dataframe: DataFrame, metadata: dict) -> DataFrame:
        n1 = self.n1.value
        n2 = self.n2.value
        window = self.rollWindow.value
        adxWindow = self.adxWindow.value
        

        # ap = (dataframe['high'] + dataframe['low'] + dataframe['close']) / 3
        # esa = ap.ewm(span=n1, min_periods=n1).mean()
        # d = ap.sub(esa).abs().ewm(span=n1, min_periods=n1).mean()
        # ci = (ap - esa) / (0.015 * d)
        # tci = ci.ewm(span=n2, min_periods=n2).mean()

        # dataframe['wt1'] = tci
        # dataframe['wt2'] = dataframe['wt1'].rolling(window=window).mean()

        dataframe['plus_di'] = ta.PLUS_DI(dataframe,adxWindow)
        dataframe['minus_di'] = ta.MINUS_DI(dataframe,adxWindow)
        
        dataframe['ema'] = ta.EMA(dataframe, timeperiod=self.emaThr.value)
        superT = self.supertrend(dataframe, self.superTrendMultiper.value, self.superTrendWindow.value)
        dataframe['up'] = superT['up']
        dataframe['down'] = superT['down']
    
        return dataframe

    def populate_entry_trend(self, dataframe: DataFrame, metadata: dict) -> DataFrame:
        """
        Based on TA indicators, populates the buy signal for the given dataframe
        :param dataframe: DataFrame
        :return: DataFrame with buy column
        """
        dataframe.loc[
            (
                ####################################
                # long condition:
                # 1. wt1 cross above wt2
                # 2. wt1 - wt2 > entry_long_derta
                # 3. adx > entry_long_adx_threshold
                ####################################
                # qtpylib.crossed_above(
                #     dataframe['wt1'],
                #     dataframe['wt2']
                # )
                (dataframe['close'] < dataframe['ema'])
                & 
                (dataframe['plus_di'] > dataframe['minus_di']) & (dataframe['plus_di']>self.adxThr.value)
                # & (dataframe['wt1'] < 0)
            ),
            'enter_long'] = 1

        # dataframe.loc[
        #     (
        #         ####################################
        #         # long condition:
        #         # 1. wt1 cross above wt2
        #         # 2. wt1 - wt2 > entry_long_derta
        #         # 3. adx > entry_long_adx_threshold
        #         ####################################
        #         (dataframe['close'] > dataframe['ema'])
        #         & 
        #         (dataframe['plus_di'] < dataframe['minus_di']) & (dataframe['minus_di']>self.adxThr.value)
        #         # & (dataframe['wt1'] < 0)
        #     ),
        #     'enter_short'] = 0

        return dataframe

    def populate_exit_trend(self, dataframe: DataFrame, metadata: dict) -> DataFrame:
        """
        Based on TA indicators, populates the sell signal for the given dataframe
        :param dataframe: DataFrame
        :return: DataFrame with buy column
        """
        dataframe.loc[
            (
                (dataframe['close'] < dataframe['down'])
            ),
            'exit_long'] = 1

        # dataframe.loc[
        #     (
        #         ####################################
        #         # exit long condition:
        #         # 1. wt2 - wt1 > exit_long_derta
        #         # 2. adx < exit_long_adx_threshold
        #         ####################################
        #         # qtpylib.crossed_above(
        #         #     dataframe['wt1'],
        #         #     dataframe['wt2']
        #         # )
        #         # &
        #         (dataframe['plus_di'] > dataframe['minus_di']) & (dataframe['plus_di']>25)
        #     ),
        #     'exit_short'] = 0

        return dataframe

    # def leverage(self, pair: str, current_time: datetime, current_rate: float,
    #              proposed_leverage: float, max_leverage: float, entry_tag: Optional[str],
    #              side: str, **kwargs) -> float:
    #     # Return 3.0 in all cases.
    #     # Bot-logic must make sure it's an allowed leverage and eventually adjust accordingly.
    #     return 1
    
    
    def supertrend(self, dataframe: DataFrame, multiplier, period):
        # start_time = time.time()

        df = dataframe.copy()
        last_row = dataframe.tail(1).index.item()

        df['TR'] = ta.TRANGE(df)
        df['ATR'] = ta.SMA(df['TR'], period)

        # st = 'ST_' + str(period) + '_' + str(multiplier)
        # stx = 'STX_' + str(period) + '_' + str(multiplier)

        # Compute basic upper and lower bands
        BASIC_UB = ((df['high'] + df['low']) / 2 + multiplier * df['ATR']).values
        BASIC_LB = ((df['high'] + df['low']) / 2 - multiplier * df['ATR']).values

        FINAL_UB = np.zeros(last_row + 1)
        FINAL_LB = np.zeros(last_row + 1)
        ST = np.zeros(last_row + 1)
        CLOSE = df['close'].values

        # Compute final upper and lower bands
        for i in range(period, last_row):
            FINAL_UB[i] = BASIC_UB[i] if BASIC_UB[i] < FINAL_UB[i - 1] or CLOSE[i - 1] > FINAL_UB[i - 1] else FINAL_UB[i - 1]
            FINAL_LB[i] = BASIC_LB[i] if BASIC_LB[i] > FINAL_LB[i - 1] or CLOSE[i - 1] < FINAL_LB[i - 1] else FINAL_LB[i - 1]

        df_Up = pd.DataFrame(FINAL_UB, columns=['up'])
        df = pd.concat([df, df_Up],axis=1)
        df_Down = pd.DataFrame(FINAL_LB, columns=['down'])
        df = pd.concat([df, df_Down],axis=1)
        
        
        # # Set the Supertrend value
        # for i in range(period, last_row):
        #     ST[i] = FINAL_UB[i] if ST[i - 1] == FINAL_UB[i - 1] and CLOSE[i] <= FINAL_UB[i] else \
        #             FINAL_LB[i] if ST[i - 1] == FINAL_UB[i - 1] and CLOSE[i] >  FINAL_UB[i] else \
        #             FINAL_LB[i] if ST[i - 1] == FINAL_LB[i - 1] and CLOSE[i] >= FINAL_LB[i] else \
        #             FINAL_UB[i] if ST[i - 1] == FINAL_LB[i - 1] and CLOSE[i] <  FINAL_LB[i] else 0.00
        # df_ST = pd.DataFrame(ST, columns=[st])
        # df = pd.concat([df, df_ST],axis=1)

        # # Mark the trend direction up/down
        # df[stx] = np.where((df[st] > 0.00), np.where((df['close'] < df[st]), 'down',  'up'), np.NaN)

        df.fillna(0, inplace=True)

        # end_time = time.time()
        # print("total time taken this loop: ", end_time - start_time)

        return DataFrame(index=df.index, data={
            'up' : df['up'],
            'down' : df['down']
        })
```


