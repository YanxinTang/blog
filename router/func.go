package router

import (
	"bytes"
	"html/template"
	"regexp"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func Date(t time.Time) string {
	return t.Format("Jan _2, 2006")
}

func Markdown(t string) string {
	input := []byte(t)
	input = bytes.Replace(input, []byte("\r"), nil, -1)
	output := blackfriday.Run(input)
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return string(p.SanitizeBytes(output))
}

func Safe(t string) template.HTML {
	return template.HTML(t)
}
