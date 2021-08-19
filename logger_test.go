package logger

import "testing"

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger("test.log", DateFormat("2006-01-02 15:04:05"))
	if err != nil {
		t.Error(err)
	}
	logger.ZapSugar.Infow("test", "name", "melon", "project", "testing")
	defer logger.ZapLogger.Sync()
}
