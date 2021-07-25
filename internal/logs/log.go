package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"time"
)

var logs = &log.Logger{}

func Init(out ...io.Writer) {
	// 输出到文件
	logs.SetOutput(io.MultiWriter(append(out, os.Stdout)...))
}

var ErrDBNotConnect = errors.New("database not connected")

type Log struct{}

// TODO 可能需要考虑大量日志写入时的性能问题，可能可以改造成异步批量写入
// 	定时和定量写入
func (Log) Write(data []byte) (int, error) {
	if database.DB == nil {
		log.Println("database not connected")
		return 0, errors.WithStack(ErrDBNotConnect)
	}
	if err := database.NewTransaction(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
		return (&models.Log{Detail: string(data)}).Insert(ctx, tx)
	}); err != nil {
		return 0, err
	}
	return len(data), nil
}

func FileWriter() (*os.File, error) {
	path := viper.GetString("logValue.Filename")
	if path == "" {
		path = "logs/default.log"
		d, e := os.Open("logs")
		if e != nil {
			if err := os.MkdirAll("logs", 0666); err != nil {
				return nil, err
			}
		} else {
			d.Close()
		}
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
	LevelPanic = "panic"
)

func Info(v string, kv ...interface{}) {
	logs.Println(logValue(&LogValue{Level: LevelInfo, Value: v}, kv...))
}

func Debug(v string, kv ...interface{}) {
	logs.Println(logValue(&LogValue{Level: LevelDebug, Value: v}, kv...))
}

func Error(v string, kv ...interface{}) {
	logs.Println(logValue(&LogValue{Level: LevelError, Value: v}, kv...))
}

func Panic(v string, kv ...interface{}) {
	logs.Panicln(logValue(&LogValue{Level: LevelPanic, Value: v}, kv...))
}

func logValue(value *LogValue, kv ...interface{}) string {
	kvl := len(kv)
	if kvl%2 != 0 {
		return logKVV(value, kv)
	}
	value.Time = time.Now()
	value.Detail = map[string]interface{}{}
	for i := 0; i < kvl-1; i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			return logKVV(value, kv)
		}
		value.Detail[key] = fmt.Sprintf("%+v", kv[i+1])
	}
	return marshalLog(value)
}

func logKVV(value *LogValue, kv ...interface{}) string {
	value.Detail = map[string]interface{}{
		"kvv": kv,
	}
	return marshalLog(value)
}

func marshalLog(logValue *LogValue) string {
	value, err := json.Marshal(logValue)
	if err != nil {
		logs.Printf("marshalLog error %v\n", err)
		return err.Error()
	}
	return string(value)
}
