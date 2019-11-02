package router

import (
	"bytes"
	"html/template"
	"regexp"
	"time"

	"github.com/YanxinTang/blog/config"
	"github.com/YanxinTang/blog/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func Date(t time.Time) string {
	return t.Format("Jan _2, 2006")
}

func Markdown(t string) string {
	input := bytes.Replace([]byte(t), []byte("\r"), nil, -1)
	cr := utils.NewChromaRenderer()
	output := blackfriday.Run(input, blackfriday.WithRenderer(cr))
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	p.AllowAttrs("class").Globally()
	p.AllowAttrs("style").OnElements("span", "p")
	return string(p.SanitizeBytes(output))
}

func Safe(t string) template.HTML {
	return template.HTML(t)
}

func Summary(t string) string {
	input := bytes.Replace([]byte(t), []byte("\r"), nil, -1)
	output := blackfriday.Run(input, blackfriday.WithRenderer(utils.NewSummaryRenderer()))
	return string(bluemonday.UGCPolicy().SanitizeBytes(output))
}

func Config() config.ConfigStruct {
	return config.Config
}
