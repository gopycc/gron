package web

import (
	"fmt"
	"strconv"

	"github.com/llangs/gron/mod"
)

func (ctx *Context) QString(name string) string {
	return ctx.params.ByName(name)
}

func (ctx *Context) QInt64(name string) int64 {
	value, _ := strconv.ParseInt(ctx.QString(name), 10, 64)
	return value
}

func (ctx *Context) QID(name string) mod.ID {
	return mod.ID(ctx.QInt64(name))
}

func (ctx *Context) QUint64(name string) uint64 {
	value, _ := strconv.ParseUint(ctx.params.ByName(name), 10, 64)
	return value
}

func (ctx *Context) QInt(name string) int {
	value, _ := strconv.Atoi(ctx.params.ByName(name))
	return value
}

func (ctx *Context) String(name string) string {
	return ctx.Request.FormValue(name)
}

func (ctx *Context) Int64(name string) int64 {
	value, _ := strconv.ParseInt(ctx.String(name), 10, 64)
	return value
}

func (ctx *Context) Uint64(name string) uint64 {
	value, _ := strconv.ParseUint(ctx.String(name), 10, 64)
	return value
}

func (ctx *Context) Int(name string) int {
	value, _ := strconv.Atoi(ctx.String(name))
	return value
}

func (ctx *Context) Mid(name string) mod.ID {
	return mod.ID(ctx.Int64(name))
}

func (ctx *Context) Bool(name string) bool {
	value, _ := strconv.ParseBool(ctx.String(name))
	return value
}

func (ctx *Context) MustInt64(name string) int64 {
	if ctx.String(name) == "" {
		panic(fmt.Errorf("%s is nil", name))
	}
	value, err := strconv.ParseInt(ctx.String(name), 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid %s:%v", name, err))
	}
	return value
}

func (ctx *Context) MustString(name string) string {
	if ctx.String(name) == "" {
		panic(fmt.Errorf("%s is nil", name))
	}
	return ctx.String(name)
}

func (ctx *Context) Pager() *mod.Pager {
	pager := &mod.Pager{Page: ctx.Int("page"), Limit: ctx.Int("limit")}
	return pager
}
