package cqhttp

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func Test_randReply(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	replies := []string{"1", "2", "3"}
	gomonkey.ApplyMethodSeq(reflect.TypeOf(&rand.Rand{}), "Intn", []gomonkey.OutputCell{
		{Values: gomonkey.Params{0}},
		{Values: gomonkey.Params{1}},
		{Values: gomonkey.Params{2}},
	})
	assert.Equal(t, "1", randReply(replies))
	assert.Equal(t, "2", randReply(replies))
	assert.Equal(t, "3", randReply(replies))
}
