package web

import (
	"net/http"

	"encoding/json"
	"fmt"
	"html/template"
	"path"

	"github.com/goinon/httprouter"
	log "github.com/sirupsen/logrus"
)

// Context represents context of a Request.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	params         httprouter.Params
	Data           map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	ctx := new(Context)
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.params = params
	ctx.Data = make(map[string]interface{})
	return ctx
}

// Trigger handles and logs error by given status.
func (ctx *Context) Handle(status int, title string, err error) {
	switch status {
	case http.StatusNotFound:
		ctx.Data["Title"] = "Page Not Found"
	case http.StatusInternalServerError:
		ctx.Data["Title"] = "Internal Server Error"
		log.Error(2, "%s: %v", title, err)
	}
	ctx.HTML(status, fmt.Sprintf("status/%d", status))
}

// NotFound renders the 404 page.
func (ctx *Context) NotFound() {
	ctx.Handle(http.StatusNotFound, "404 Not Found", nil)
}

// ServerError renders the 500 page.
func (ctx *Context) ServerError(title string, err error) {
	ctx.Handle(http.StatusInternalServerError, title, err)
}

func (ctx *Context) HTML(status int, content string) {
	t, err := template.New("html").Parse(`{{define "T"}}<html><title>{{.}}</title><body>{{.}}</body></html>{{end}}`)
	err = t.ExecuteTemplate(ctx.ResponseWriter, "T", content)
	if err != nil {
		fmt.Println(err)
	}
}

func (ctx *Context) Render(templFile string, model interface{}) {
	t, _ := template.ParseFiles(path.Join(PAGES, templFile))
	ctx.Data["model"] = model
	log.Println(t.Execute(ctx.ResponseWriter, ctx.Data))
}

func (ctx *Context) RenderTempl(templFile, layoutFile string, model interface{}) {
	ctx.Data["model"] = model
	page := ParseTemplate(templFile, layoutFile, ctx.Data)
	ctx.ResponseWriter.Write(page)
}

func (ctx *Context) Layout(templFile string, model interface{}) {
	ctx.RenderTempl(templFile, LAYOUT, model)
}

func (ctx *Context) Admin(templFile string, model interface{}) {
	ctx.RenderTempl(templFile, ADMIN, model)
}

func (ctx *Context) RenderJson(model interface{}) {
	js, err := json.Marshal(model)
	if err != nil {
		panic(err)
	}
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	ctx.ResponseWriter.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, X-File-Type, Cache-Control, Origin, token")
	ctx.ResponseWriter.Header().Add("Access-Control-Expose-Headers", "Authorization")
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.Write(js)
}

func (ctx *Context) Redirect(urlStr string, code int) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, urlStr, code)
}
