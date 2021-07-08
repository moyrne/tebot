package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type CqHTTP struct{}

const (
	STFriend    = "friend"
	STGroup     = "group"
	STGroupSelf = "group_self"
	STOther     = "other"
)

// HTTP 文档 https://github.com/ishkong/go-cqhttp-docs/tree/main/docs/event
func (h CqHTTP) HTTP(c *gin.Context) {
	var params QMessage
	if err := c.BindJSON(&params); err != nil {
		// TODO log error
		return
	}

	fmt.Println("[params]", params)

	switch params.SubType {
	case STFriend:
		h.friend(c, params)
	case STGroup, STGroupSelf:
		h.group(c, params)
	case STOther:
		fallthrough
	default:
		// TODO log error
		return
	}
}

func (h CqHTTP) friend(c *gin.Context, params QMessage) {

}

func (h CqHTTP) group(c *gin.Context, params QMessage) {

}
