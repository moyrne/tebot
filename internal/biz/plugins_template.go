package biz

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

// TODO 转移至 pkg

type Template struct {
	*template.Template
}

func NewTemplate(name, parse string) (*Template, error) {
	temp, err := template.New(name).Parse(parse)
	return &Template{temp}, errors.WithStack(err)
}

func (t *Template) Execute(data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, data); err != nil {
		return "", errors.WithStack(err)
	}
	return buf.String(), nil
}
