package router

import (
	"html/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func Date(t time.Time) string {
	return t.Format("Jan _2, 2006")
}

func Markdown(t string) string {
	return string(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.Run([]byte(t))))
}

func Safe(t string) template.HTML {
	return template.HTML(t)
}
