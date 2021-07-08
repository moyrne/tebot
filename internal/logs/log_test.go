package logs

import (
	"bytes"
	"errors"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	gomonkey.ApplyFunc(time.Now, func() time.Time {
		ti, _ := time.Parse("2006-01-02T15:04:05.999999999-07:00", "2021-07-08T13:35:25.618779118+08:00")
		return ti
	})
	out := bytes.NewBuffer(nil)
	Init(out)
	Info("info log", "err", errors.New("info log error"))
	Debug("debug log", "err", errors.New("debug log error"))
	Error("error log", "err", errors.New("error log error"))
	outs := strings.Split(out.String(), "\n")
	assert.JSONEq(t, outs[0], `{"level":"info","time":"2021-07-08T13:35:25.618779118+08:00","value":"info log","detail":{"err":"info log error"}}`)
	assert.JSONEq(t, outs[1], `{"level":"debug","time":"2021-07-08T13:35:25.618779118+08:00","value":"debug log","detail":{"err":"debug log error"}}`)
	assert.JSONEq(t, outs[2], `{"level":"error","time":"2021-07-08T13:35:25.618779118+08:00","value":"error log","detail":{"err":"error log error"}}`)
}
