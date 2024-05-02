package main

import (
	"context"
	"golang-course-hw/hw3/domain"
	"golang-course-hw/hw3/generator"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

type Bookmark struct{}

func add(a, b int32) (int32, bool) { return a + b, true }

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	wg := sync.WaitGroup{}
	domain.WriteCandles(prices, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	cancel()
	wg.Wait()
}
