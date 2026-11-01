package markdown

import (
	"html/template"
	"time"
)

type Meta struct {
	Title string
	Date  time.Time
	Type  string
}

type Markdown interface {
	Parse(src string) (Meta, template.HTML, error)
	HTML(src string) (template.HTML, error)
}
