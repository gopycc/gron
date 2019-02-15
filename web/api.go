package web

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/llangs/gron/mod/errors"
)

type JsonResult struct {
	Code  errors.ErrorCode `json:"code"`
	Error string           `json:"error"`
	Data  interface{}      `json:"data"`
}

func (jr JsonResult) String() string {
	js, err := json.Marshal(jr)
	if err != nil {
		log.Fatalf("json marshal error: %v", err)
		return ""
	}
	return fmt.Sprintf("%s", js)
}

func Ok(data interface{}) *JsonResult {
	return &JsonResult{Code: errors.Success, Error: "", Data: data}
}

func (ctx *Context) JsonOk(data interface{}) {
	ctx.RenderJson(Ok(data))
}

func (ctx *Context) JsonError2(err string) {
	jr := &JsonResult{Code: -1, Error: err, Data: nil}
	ctx.RenderJson(jr)
}

func (ctx *Context) JsonError(err error) {
	jr := &JsonResult{Code: 0, Error: err.Error(), Data: nil}
	ctx.RenderJson(jr)
}
