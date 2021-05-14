package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/openware/gocryptotrader/backtester/common"
	"github.com/openware/gocryptotrader/backtester/data/kline"
	"github.com/openware/pkg/currency"
	"github.com/openware/pkg/asset"
	gctkline "github.com/openware/pkg/kline"
	"github.com/openware/pkg/trade"
)

// LoadData retrieves data from an existing database using GoCryptoTrader's database handling implementation
func LoadData(startDate, endDate time.Time, interval time.Duration, exchangeName string, dataType int64, fPair currency.Pair, a asset.Item) (*kline.DataFromKline, error) {
	resp := &kline.DataFromKline{}
	switch dataType {
	case common.DataCandle:
		klineItem, err := getCandleDatabaseData(
			startDate,
			endDate,
			interval,
			exchangeName,
			fPair,
			a)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve database candle data for %v %v %v, %v", exchangeName, a, fPair, err)
		}
		resp.Item = klineItem
	case common.DataTrade:
		trades, err := trade.GetTradesInRange(
			exchangeName,
			a.String(),
			fPair.Base.String(),
			fPair.Quote.String(),
			startDate,
			endDate)
		if err != nil {
			return nil, err
		}
		klineItem, err := trade.ConvertTradesToCandles(
			gctkline.Interval(interval),
			trades...)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve database trade data for %v %v %v, %v", exchangeName, a, fPair, err)
		}
		resp.Item = klineItem
	default:
		return nil, fmt.Errorf("could not retrieve database data for %v %v %v, invalid data type received", exchangeName, a, fPair)
	}
	resp.Item.Exchange = strings.ToLower(resp.Item.Exchange)

	return resp, nil
}

func getCandleDatabaseData(startDate, endDate time.Time, interval time.Duration, exchangeName string, fPair currency.Pair, a asset.Item) (gctkline.Item, error) {
	return gctkline.LoadFromDatabase(
		exchangeName,
		fPair,
		a,
		gctkline.Interval(interval),
		startDate,
		endDate)
}
