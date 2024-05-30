- 基于ADX判断突破趋势
- 基于价格和均线判断回踩
- DCA订单


## 策略
```python
# --- Do not remove these libs ---
from freqtrade.strategy import IStrategy,merge_informative_pair
from typing import Dict, List
from functools import reduce
from pandas import DataFrame
from datetime import datetime
from typing import Optional
from freqtrade.strategy import CategoricalParameter, DecimalParameter, IntParameter
from freqtrade.persistence import PairLocks
# --------------------------------

import talib.abstract as ta
import freqtrade.vendor.qtpylib.indicators as qtpylib
from freqtrade.persistence import Trade
from typing import Optional, Tuple, Union


class FutureDMIPRICEStrategy(IStrategy):

    INTERFACE_VERSION: int = 3
    can_short = True
    position_adjustment_enable = True
    use_custom_stoploss = True

    minimal_roi = {
        "0": 0.1
    }
    
    stoploss =  -0.1
    
    # trailing_stop = False
    # trailing_stop_positive = 0.2
    # trailing_stop_positive_offset = 0.6
    # trailing_only_offset_is_reached = False
    

    order_types = {
        'entry': 'market',
        'exit': 'market',
        'stoploss': 'market',
        'stoploss_on_exchange': True
    }

    # Optional order time in force.
    order_time_in_force = {
        'entry': 'GTC',
        'exit': 'GTC'
    }

    
    adxWindow = IntParameter(7, 21, default=14, space="buy")
    adxThr = IntParameter(15, 35, default=25, space="buy")
    emaThr = IntParameter(5, 55, default=24, space="buy")
    


    # Optimal timeframe for the strategy
    timeframe = '1m'
    inf_tf = '4h'
    
    # DCA Order
    max_entry_position_adjustment = 3
    # DCA factor
    max_dca_multiplier = 5.5
    
        
    def informative_pairs(self):
        # get access to all pairs available in whitelist.
        pairs = self.dp.current_whitelist()
        # Assign tf to each pair so they can be downloaded and cached for strategy.
        informative_pairs = [(pair, self.inf_tf) for pair in pairs]
        return informative_pairs
    

    def populate_indicators(self, dataframe: DataFrame, metadata: dict) -> DataFrame:
        adxWindow = self.adxWindow.value
        

        informative = self.dp.get_pair_dataframe(pair=metadata['pair'], timeframe=self.inf_tf)
        informative['plus_di'] = ta.PLUS_DI(informative,adxWindow)
        informative['minus_di'] = ta.MINUS_DI(informative,adxWindow)
        
        informative['ema'] = ta.EMA(informative, timeperiod=self.emaThr.value)
        dataframe = merge_informative_pair(dataframe, informative, self.timeframe, self.inf_tf, ffill=True)    

        return dataframe

    def populate_entry_trend(self, dataframe: DataFrame, metadata: dict) -> DataFrame:
        """
        Based on TA indicators, populates the buy signal for the given dataframe
        :param dataframe: DataFrame
        :return: DataFrame with buy column
        """
        dataframe.loc[
            (
                (dataframe['close'] < dataframe[f'ema_{self.inf_tf}'])
                & 
                (dataframe[f'plus_di_{self.inf_tf}'] > dataframe[f'minus_di_{self.inf_tf}']) & (dataframe[f'plus_di_{self.inf_tf}']>self.adxThr.value)
            ),
            'enter_long'] = 1
        
        dataframe.loc[
            (
                (dataframe['close'] > dataframe[f'ema_{self.inf_tf}'])
                & 
                (dataframe[f'plus_di_{self.inf_tf}'] < dataframe[f'minus_di_{self.inf_tf}']) & (dataframe[f'minus_di_{self.inf_tf}']>self.adxThr.value)
            ),
            'enter_short'] = 1

        return dataframe

    def populate_exit_trend(self, dataframe: DataFrame, metadata: dict) -> DataFrame:
        """
        Based on TA indicators, populates the sell signal for the given dataframe
        :param dataframe: DataFrame
        :return: DataFrame with buy column
        """
        dataframe.loc[
            (
            ),
            'exit_long'] = 0

        dataframe.loc[
            (
            ),
            'exit_short'] = 0

        return dataframe

    def leverage(self, pair: str, current_time: datetime, current_rate: float,
                 proposed_leverage: float, max_leverage: float, entry_tag: Optional[str],
                 side: str, **kwargs) -> float:
        # Return 3.0 in all cases.
        # Bot-logic must make sure it's an allowed leverage and eventually adjust accordingly.
        return 1
    
    # DCA the initial order (opening trade)
    def custom_stake_amount(self, pair: str, current_time: datetime, current_rate: float,
                            proposed_stake: float, min_stake: Optional[float], max_stake: float,
                            leverage: float, entry_tag: Optional[str], side: str,
                            **kwargs) -> float:

        # We need to leave most of the funds for possible further DCA orders
        # This also applies to fixed stakes
        return proposed_stake / self.max_dca_multiplier
    
    # DCA left order (append trade)
    def adjust_trade_position(self, trade: Trade, current_time: datetime,
                              current_rate: float, current_profit: float,
                              min_stake: Optional[float], max_stake: float,
                              current_entry_rate: float, current_exit_rate: float,
                              current_entry_profit: float, current_exit_profit: float,
                              **kwargs
                              ) -> Union[Optional[float], Tuple[Optional[float], Optional[str]]]:

        if current_profit > 0.05 and trade.nr_of_successful_exits == 0:
            # Take half of the profit at +5%
            return -(trade.stake_amount / 2)
        
        
        if current_profit > -0.05:
            return None

        # Obtain pair dataframe (just to show how to access it)
        dataframe, _ = self.dp.get_analyzed_dataframe(trade.pair, self.timeframe)
        # Only buy when not actively falling price.
        last_candle = dataframe.iloc[-1].squeeze()
        previous_candle = dataframe.iloc[-2].squeeze()
        if last_candle['close'] < previous_candle['close']:
            return None

        filled_entries = trade.select_filled_orders(trade.entry_side)
        count_of_entries = trade.nr_of_successful_entries
        # Allow up to 3 additional increasingly larger buys (4 in total)
        # Initial buy is 1x
        # If that falls to -5% profit, we buy 1.25x more, average profit should increase to roughly -2.2%
        # If that falls down to -5% again, we buy 1.5x more
        # If that falls once again down to -5%, we buy 1.75x more
        # Total stake for this trade would be 1 + 1.25 + 1.5 + 1.75 = 5.5x of the initial allowed stake.
        # That is why max_dca_multiplier is 5.5
        try:
            # This returns first order stake size
            stake_amount = filled_entries[0].stake_amount
            # This then calculates current safety order size
            stake_amount = stake_amount * (1 + (count_of_entries * 0.25))
            return stake_amount
        except Exception as exception:
            return None

        return None
    
    # DCA stop loss
    def custom_stoploss(self, pair: str, trade: 'Trade', current_time: datetime,
                        current_rate: float, current_profit: float, after_fill: bool, 
                        **kwargs) -> Optional[float]:

        if after_fill: 
            # After an additional order, start with a stoploss of 10% below the new open rate
            return stoploss_from_open(0.10, current_profit, is_short=trade.is_short, leverage=trade.leverage)
        return None
```


## 配置
```json
{
    "max_open_trades": 50,
    "stake_currency": "USDT",
    "stake_amount": "unlimited",
    "dry_run": true,
    "dry_run_wallet": 5000,
    "db_url": "sqlite:///dmi_price_future.dry_run.sqlite",
    "cancel_open_orders_on_exit": false,
    "trading_mode": "futures",
    "margin_mode": "isolated",
    "unfilledtimeout": {
        "entry": 10,
        "exit": 10,
        "exit_timeout_count": 0,
        "unit": "minutes"
    },
    "entry_pricing": {
        "price_side": "other",
        "use_order_book": true,
        "order_book_top": 1,
        "price_last_balance": 0.0,
        "check_depth_of_market": {
            "enabled": false,
            "bids_to_ask_delta": 1
        }
    },
    "exit_pricing":{
        "price_side": "other",
        "use_order_book": true,
        "order_book_top": 1
    },
    "exchange": {
        "name": "binance",
        "key": "",
        "secret": "",
        "ccxt_config": {},
        "ccxt_async_config": {},
        "pair_whitelist": [
            ".*/USDT:USDT"
        ],
        "pair_blacklist": [
            "PAXG/USDT:USDT",
            "FDUSD/USDT:USDT",
            "USDC/USDT:USDT",
            "BUSD/USDT:USDT",
            "TUSD/USDT:USDT",
            "EUR/USDT:USDT",
            "USDP/USDT:USDT",
            "USTC/USDT:USDT",
            ".*DOWN/USDT:USDT"
        ]
    },
    "pairlists": [
        {
            "method": "StaticPairList"
        },
        {
            "method": "ShuffleFilter"   
        }
    ],
    "telegram": {
        "enabled": true,
        "token": "6575611686:AAH42t2vEWHxJUmKv2l5X9tiPOvDsROQECo",
        "chat_id": "-4260008784"
    },
    "api_server": {
        "enabled": false,
        "listen_ip_address": "0.0.0.0",
        "listen_port": 8082,
        "verbosity": "error",
        "enable_openapi": false,
        "jwt_secret_key": "deb1cd340cd8b8244ec9df7b361c0e078f134ed9881eb6f598b83ae651c708be",
        "ws_token": "7SUlvTqmPXXBwGIBHIastfhRTdW-Vmsg8A",
        "CORS_origins": [],
        "username": "freqtrader",
        "password": "li19960205"
    },
    "bot_name": "freqtrade",
    "initial_state": "running",
    "force_entry_enable": true,
    "internals": {
        "process_throttle_secs": 60
    }
}
```