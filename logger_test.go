package logger

import "testing"

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger("test.log")
	if err != nil {
		t.Error(err)
	}
	logger.ZapSugar.Infow("test", "name", "melon", "project", "testing")
	defer logger.ZapLogger.Sync()
}
