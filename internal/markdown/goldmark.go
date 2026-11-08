package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	meta "github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type Goldmark struct {
	parser   parser.Parser
	renderer html.Renderer
}

func NewGoldmark() *Goldmark {
	return &Goldmark{
		parser: parser.New(
			parser.WithExtensions(meta.Parser),
		),
		renderer: html.New(html.WithUnsafe()),
	}
}

func (g *Goldmark) Parse(src string) (Meta, template.HTML, error) {
	source := []byte(src)
	doc := g.parser.Parse(source)
	metaData, err := documentMeta(doc)
	if err != nil {
		return Meta{}, "", err
	}
	var buf bytes.Buffer
	if err := g.renderer.Render(&buf, source, doc); err != nil {
		return Meta{}, "", err
	}
	return metaData, template.HTML(buf.String()), nil
}

func (g *Goldmark) HTML(src string) (template.HTML, error) {
	_, html, err := g.Parse(src)
	return html, err
}

func documentMeta(doc ast.Node) (Meta, error) {
	d, ok := doc.(*ast.Document)
	if !ok {
		return Meta{}, fmt.Errorf("goldmark: root is not a document")
	}
	raw := d.Metadata()
	date, err := dateValue(raw["date"])
	if err != nil {
		return Meta{}, err
	}
	return Meta{
		Title: stringValue(raw["title"]),
		Date:  date,
		Type:  stringValue(raw["type"]),
	}, nil
}

func stringValue(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

func dateValue(v any) (time.Time, error) {
	if v == nil {
		return time.Time{}, nil
	}
	switch t := v.(type) {
	case time.Time:
		if t.IsZero() {
			return time.Time{}, fmt.Errorf("invalid date: %v", v)
		}
		return t, nil
	case string:
		if strings.TrimSpace(t) == "" {
			return time.Time{}, nil
		}
		for _, layout := range []string{"2006-01-02", time.RFC3339, time.RFC3339Nano} {
			if parsed, err := time.Parse(layout, t); err == nil {
				return parsed, nil
			}
		}
		return time.Time{}, fmt.Errorf("invalid date: %v", v)
	default:
		return time.Time{}, fmt.Errorf("invalid date: %v", v)
	}
}
