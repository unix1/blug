package generate

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/unix1/blug/internal/site"
)

type templates struct {
	header  *template.Template
	footer  *template.Template
	listing *template.Template
	article *template.Template
}

const articleSrc = `<article>
<h1>{{.EntryTitle}}</h1>
{{if .Date}}<time datetime="{{.Date}}">{{.Date}}</time>
{{end}}{{.Body}}
</article>
`

func loadTemplates(siteRoot string, defaults embed.FS) (templates, error) {
	header, err := parseTemplate(siteRoot, defaults, site.HeaderTemplate)
	if err != nil {
		return templates{}, err
	}
	footer, err := parseTemplate(siteRoot, defaults, site.FooterTemplate)
	if err != nil {
		return templates{}, err
	}
	listing, err := parseTemplate(siteRoot, defaults, site.ListingTemplate)
	if err != nil {
		return templates{}, err
	}
	article, err := template.New("article").Parse(articleSrc)
	if err != nil {
		return templates{}, err
	}
	return templates{header: header, footer: footer, listing: listing, article: article}, nil
}

func parseTemplate(siteRoot string, defaults embed.FS, name string) (*template.Template, error) {
	path := filepath.Join(siteRoot, site.TemplatesDir, name)
	raw, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		raw, err = defaults.ReadFile(filepath.ToSlash(filepath.Join(site.TemplatesDir, name)))
		if err != nil {
			return nil, fmt.Errorf("template %s: %w", name, err)
		}
	}
	tmpl, err := template.New(name).Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("template %s: %w", name, err)
	}
	return tmpl, nil
}
