package logs

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func JSONMarshalIgnoreErr(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

type FileHook struct {
	*os.File
}

func (f *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FileHook) Fire(e *logrus.Entry) error {
	serialized, err := e.Logger.Formatter.Format(e)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = f.Write(serialized)
	return errors.WithStack(err)
}

func NewFileHook(path string) (file *FileHook, close func(), err error) {
	if path == "" {
		path = "logs/default.log"
		d, e := os.Open("logs")
		if e != nil {
			if err := os.MkdirAll("logs", 0666); err != nil {
				return nil, func() {}, err
			}
		} else {
			d.Close()
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	return &FileHook{File: f}, func() { f.Close() }, errors.WithStack(err)
}
