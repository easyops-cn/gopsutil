package net

import (
	"go.easyops.local/slog"
	"sync/atomic"
	"time"
)

var config SleepConfig

type SleepConfig struct {
	sleepMs         atomic.Int64
	loopBeforeSleep atomic.Int32
	logger          slog.Logger
}

func InitSleepConfig(ms int, loop int, logger slog.Logger) {
	config = SleepConfig{}
	config.sleepMs.Store(int64(ms))
	config.loopBeforeSleep.Store(int32(loop))
	config.logger = logger
}

func TimeSleep(i int) {
	j := int(config.loopBeforeSleep.Load())
	if j == 0 || (i+1)%j != 0 {
		return
	}
	ms := config.sleepMs.Load()
	if ms != 0 {
		config.logger.Infof("time sleep, sleep duration: %d, current loop count: %d", ms, i)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}
