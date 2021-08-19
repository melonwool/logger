# logger
基于zap仅实现接收信号日志重新写入

# 使用方式
可以使用zap Logger 和 ZapSugar 所有方法...
```go
func main(){
	
	fileName := "test.log"
	logger, err := NewLogger(fileName)
	if err != nil {
		panic(err)
	}
	defer logger.ZapLogger.Sync()
	logger.ZapSugar.Infow("access.log", "a", 1, "b", 2, "c", "ok")
	//  {"level":"info","ts":1629346460.2836978,"caller":"runtime/asm_amd64.s:1371","msg":"access.log","a":1,"b":2,"c":"ok"}
}
```