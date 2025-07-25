package sleepconfig

import (
	"go.easyops.local/slog"
	"sync/atomic"
	"time"
)

var config SleepConfig

type SleepConfig struct {
	enable          bool
	sleepMs         atomic.Int64
	loopBeforeSleep atomic.Int32
	logger          slog.Logger
}

func InitSleepConfig(ms int, loop int, logger slog.Logger) {
	config = SleepConfig{}
	config.enable = true
	config.sleepMs.Store(int64(ms))
	config.loopBeforeSleep.Store(int32(loop))
	config.logger = logger
}

func TimeSleep(i int, total int) {
	if !config.enable {
		return
	}
	j := int(config.loopBeforeSleep.Load())
	if j == 0 || (i+1)%j != 0 {
		return
	}
	ms := config.sleepMs.Load()
	if ms != 0 {
		if (i+1) == j || (i+1) == (total/j)*j || (i+1)%10000 == 0 {
			config.logger.Infof("time sleep, sleep duration: %d, current loop count: %d", ms, i)
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}
