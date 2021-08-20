# logger
基于zap仅实现接收信号日志重新写入

# 使用方式
可以使用zap Logger 和 ZapSugar 所有方法...
```go
import "github.com/melonwool/logger"

func main(){
	
	fileName := "test.log"
	logger, err := logger.NewLogger(fileName)
	// 定义ts 时间格式, 定时检测文件是否存在，不存在的话重新创建
	// logger, err := logger.NewLogger(fileName, logger.DateFormat("2006-01-02 15:04:05"),logger.TickerListen(time.NewTicker(time.Second)))
	// 通过信号SIGUSR1 来进行文件重新创建
	// logger, err := logger.NewLogger(fileName, logger.SignalListen())
	if err != nil {
		panic(err)
	}
	defer logger.ZapLogger.Sync()
	logger.ZapSugar.Infow("access.log", "a", 1, "b", 2, "c", "ok")
	//  {"level":"info","ts":1629346460.2836978,"caller":"runtime/asm_amd64.s:1371","msg":"access.log","a":1,"b":2,"c":"ok"}
}
```