package logger

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	//logger, err := NewLogger("test.log", DateFormat("2006-01-02 15:04:05"), TickerListen(time.NewTicker(time.Second)))
	logger, err := NewLogger("test.log", DateFormat("2006-01-02 15:04:05"), SignalListen())
	if err != nil {
		t.Error(err)
	}
	logger.ZapSugar.Infow("test", "name", "melon", "project", "testing")
	defer logger.ZapLogger.Sync()

	for {
		logger.ZapSugar.Infow("test", "name", "melon", "project", "testing")
		time.Sleep(time.Second)
	}
}
