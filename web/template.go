package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

const (
	PAGES  = "pages"
	LAYOUT = "layout.html"
	ADMIN  = "admin/layout.html"
)

var (
	TimeZoneOffset int64
)

var funcMaps = template.FuncMap{
	"html": func(text string) template.HTML {
		return template.HTML(text)
	},
	"loadtimes": func(startTime time.Time) string {
		return fmt.Sprintf("%dms", time.Now().Sub(startTime)/1000000)
	},
	"url": func(url string) string {
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			return url
		}

		return "http://" + url
	},
	"add": func(a, b int) int {
		return a + b
	},
	"formatdate": func(t time.Time) string {
		return t.Format(time.RFC822)
	},
	"formattime": func(t time.Time) string {
		now := time.Now()
		duration := now.Sub(t)
		if duration.Seconds() < 60 {
			return fmt.Sprintf("刚刚")
		} else if duration.Minutes() < 60 {
			return fmt.Sprintf("%.0f 分钟前", duration.Minutes())
		} else if duration.Hours() < 24 {
			return fmt.Sprintf("%.0f 小时前", duration.Hours())
		}
		t = t.Add(time.Hour * time.Duration(TimeZoneOffset))
		return t.Format("2006-01-02 15:04")
	},
	"formatdatetime": func(t time.Time) string {
		return t.Add(time.Hour * time.Duration(TimeZoneOffset)).Format("2006-01-02 15:04:05")
	},
	"nl2br": func(text string) template.HTML {
		return template.HTML(strings.Replace(text, "\n", "<br>", -1))
	},
	"include": func(filename string, data map[string]interface{}) template.HTML {
		var buf bytes.Buffer
		t, err := template.ParseFiles(path.Join("pages", filename))
		if err != nil {
			panic(err)
		}
		err = t.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		return template.HTML(buf.Bytes())
	},
}

func ParseTemplate(templFile, layoutFile string, data map[string]interface{}) []byte {
	var buf bytes.Buffer
	t := template.New(templFile).Funcs(funcMaps)
	baseBytes, err := ioutil.ReadFile(path.Join(PAGES, layoutFile))
	if err != nil {
		panic(err)
	}
	t, err = t.Parse(string(baseBytes))
	if err != nil {
		panic(err)
	}
	t, err = t.ParseFiles(path.Join(PAGES, templFile))
	if err != nil {
		panic(err)
	}
	err = t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
