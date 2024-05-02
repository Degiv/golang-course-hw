package domain

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

func WriteCandles(prices <-chan Price, wg *sync.WaitGroup) {
	pricesAsCandle := make(chan Candle)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(pricesAsCandle)

		for price := range prices {
			pricesAsCandle <- Candle{
				Ticker: price.Ticker,
				Period: zeroValue,
				Open:   price.Value,
				High:   price.Value,
				Low:    price.Value,
				Close:  price.Value,
				TS:     price.TS,
			}
		}
	}()

	candles1m := writeCandles(pricesAsCandle, CandlePeriod1m, "hw3/candles_1m.csv", wg)
	candles2m := writeCandles(candles1m, CandlePeriod2m, "hw3/candles_2m.csv", wg)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range writeCandles(candles2m, CandlePeriod10m, "hw3/candles_10m.csv", wg) {
		}
	}()
}

func writeCandles(in <-chan Candle, period CandlePeriod, filename string, wg *sync.WaitGroup) <-chan Candle {
	out := make(chan Candle)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)

		aggregatedMap := map[string]Candle{} // from ticker to candle

		file, _ := os.Create(filename)
		defer file.Close()

		writer := csv.NewWriter(file)

		for candle := range in {
			candlePTS, _ := PeriodTS(period, candle.TS)
			aggregated, ok := aggregatedMap[candle.Ticker]
			if ok {
				aggregatedPTS, _ := PeriodTS(period, aggregated.TS)
				isNewPeriod := candlePTS.After(aggregatedPTS)
				if isNewPeriod {
					writeCandle(writer, aggregated)
					out <- aggregated
					aggregatedMap[candle.Ticker] = NewCandleByCandle(candle, period)
				} else {
					aggregatedMap[candle.Ticker] = updated(aggregated, candle)
				}
			} else {
				aggregatedMap[candle.Ticker] = NewCandleByCandle(candle, period)
			}
		}
	}()

	return out
}

func updated(toUpdate Candle, from Candle) Candle {
	if toUpdate.High < from.High {
		toUpdate.High = from.High
	}

	if toUpdate.Low > from.Low {
		toUpdate.Low = from.Low
	}

	toUpdate.Close = from.Close
	return toUpdate
}

func NewCandleByCandle(candle Candle, period CandlePeriod) Candle {
	ts, _ := PeriodTS(period, candle.TS)
	return Candle{
		Ticker: candle.Ticker,
		Period: period,
		Open:   candle.Open,
		High:   candle.High,
		Low:    candle.Low,
		Close:  candle.Close,
		TS:     ts,
	}
}

func writeCandle(writer *csv.Writer, candle Candle) {
	writer.Write([]string{
		candle.Ticker,
		candle.TS.String(),
		fmt.Sprintf("%f", candle.Open),
		fmt.Sprintf("%f", candle.High),
		fmt.Sprintf("%f", candle.Low),
		fmt.Sprintf("%f", candle.Close)})
	writer.Flush()
}
