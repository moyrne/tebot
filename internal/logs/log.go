package logs

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"time"
)

var logs = &log.Logger{}

func Init(out ...io.Writer) {
	// 初始化
	writer := io.MultiWriter(os.Stdout)
	if len(out) != 0 {
		writer = io.MultiWriter(os.Stdout, out[0])
	}
	// 输出到文件
	logs.SetOutput(writer)
}

func FileWriter() (*os.File, error) {
	path := viper.GetString("Log.Filename")
	if path == "" {
		path = "logs/default.log"
	}
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
}

type LogValue struct {
	Level  string                 `json:"level"`
	Time   time.Time              `json:"time"`
	Value  string                 `json:"value"`
	Detail map[string]interface{} `json:"detail"`
}

const (
	LevelInfo  = "info"
	LevelDebug = "debug"
	LevelError = "error"
)

func Info(v string, kv ...interface{}) {
	Log(&LogValue{Level: LevelInfo, Value: v}, kv...)
}

func Debug(v string, kv ...interface{}) {
	Log(&LogValue{Level: LevelDebug, Value: v}, kv...)
}

func Error(v string, kv ...interface{}) {
	Log(&LogValue{Level: LevelError, Value: v}, kv...)
}

func Log(logValue *LogValue, kv ...interface{}) {
	kvl := len(kv)
	if kvl%2 != 0 {
		logKVV(logValue, kv)
		return
	}
	logValue.Time = time.Now()
	logValue.Detail = map[string]interface{}{}
	for i := 0; i < kvl-1; i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			logKVV(logValue, kv)
			return
		}
		logValue.Detail[key] = fmt.Sprintf("%v", kv[i+1])
	}
	marshalLog(logValue)
}

func logKVV(logValue *LogValue, kv ...interface{}) {
	logValue.Detail = map[string]interface{}{
		"kvv": kv,
	}
	marshalLog(logValue)
}

func marshalLog(logValue *LogValue) {
	value, err := json.Marshal(logValue)
	if err != nil {
		logs.Printf("marshalLog error %v\n", err)
		return
	}

	logs.Println(string(value))
}
