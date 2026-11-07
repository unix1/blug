package generate

import (
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

func loadTemplates(siteRoot string) (templates, error) {
	header, err := parseTemplate(siteRoot, site.HeaderTemplate)
	if err != nil {
		return templates{}, err
	}
	footer, err := parseTemplate(siteRoot, site.FooterTemplate)
	if err != nil {
		return templates{}, err
	}
	listing, err := parseTemplate(siteRoot, site.ListingTemplate)
	if err != nil {
		return templates{}, err
	}
	article, err := parseTemplate(siteRoot, site.ArticleTemplate)
	if err != nil {
		return templates{}, err
	}
	return templates{header: header, footer: footer, listing: listing, article: article}, nil
}

func parseTemplate(siteRoot, name string) (*template.Template, error) {
	path := filepath.Join(siteRoot, site.TemplatesDir, name)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s/%s: %w", site.TemplatesDir, name, err)
	}
	tmpl, err := template.New(name).Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("%s/%s: %w", site.TemplatesDir, name, err)
	}
	return tmpl, nil
}
