package generate

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/unix1/blug/internal/config"
	"github.com/unix1/blug/internal/markdown"
	"github.com/unix1/blug/internal/site"
)

var skipDirs = map[string]bool{site.AssetsDir: true}

type entry struct {
	Slug   string
	Title  string
	Date   time.Time
	Body   template.HTML
	IsPage bool
}

type listPost struct {
	Title string
	Date  string
	Slug  string
}

type pageData struct {
	Title          string
	Description    string
	SiteHeader     string
	SiteFooter     string
	EntryTitle     string
	Date           string
	Body           template.HTML
	ListingIntro   template.HTML
	ListingHeading string
	Posts          []listPost
}

type Stats struct {
	GeneratedPosts int
	SkippedPosts   int
	GeneratedPages int
	SkippedPages   int
	WroteListing   bool
	Posts          int
	Pages          int
}

func Run(siteRoot string, defaults embed.FS, md markdown.Markdown) (Stats, error) {
	cfg, err := config.Load(siteRoot)
	if err != nil {
		return Stats{}, err
	}
	tmpls, err := loadTemplates(siteRoot, defaults)
	if err != nil {
		return Stats{}, err
	}

	publicDir := filepath.Join(siteRoot, site.PublicDir)
	posts, pages, err := findEntries(publicDir, md)
	if err != nil {
		return Stats{}, err
	}

	var stats Stats
	stats.Posts = len(posts)
	stats.Pages = len(pages)

	for _, post := range posts {
		html, err := renderPost(cfg, tmpls, post)
		if err != nil {
			return stats, err
		}
		wrote, err := writeFile(filepath.Join(publicDir, post.Slug, site.IndexHTML), html)
		if err != nil {
			return stats, err
		}
		if wrote {
			stats.GeneratedPosts++
		} else {
			stats.SkippedPosts++
		}
	}
	for _, page := range pages {
		html, err := renderStaticPage(cfg, tmpls, page)
		if err != nil {
			return stats, err
		}
		wrote, err := writeFile(filepath.Join(publicDir, page.Slug, site.IndexHTML), html)
		if err != nil {
			return stats, err
		}
		if wrote {
			stats.GeneratedPages++
		} else {
			stats.SkippedPages++
		}
	}

	html, err := renderListing(cfg, tmpls, posts, md)
	if err != nil {
		return stats, err
	}
	wroteListing, err := writeFile(filepath.Join(publicDir, site.IndexHTML), html)
	if err != nil {
		return stats, err
	}
	stats.WroteListing = wroteListing
	return stats, nil
}

func (s Stats) Print(w io.Writer) {
	fmt.Fprintf(w, "Generated: %s, %s\n", countLabel(s.GeneratedPosts, "post"), countLabel(s.GeneratedPages, "page"))
	fmt.Fprintf(w, "Skipped generation: %s, %s\n", countLabel(s.SkippedPosts, "post"), countLabel(s.SkippedPages, "page"))
	if s.WroteListing {
		fmt.Fprintln(w, "Wrote index listing")
	} else {
		fmt.Fprintln(w, "Skipping index listing write, already up to date")
	}
	fmt.Fprintf(w, "Total: %s, %s\n", countLabel(s.Posts, "post"), countLabel(s.Pages, "page"))
}

func findEntries(publicDir string, md markdown.Markdown) ([]entry, []entry, error) {
	entries, err := os.ReadDir(publicDir)
	if err != nil {
		return nil, nil, fmt.Errorf("%s/: %w", site.PublicDir, err)
	}

	var posts []entry
	var pages []entry
	for _, item := range entries {
		if !item.IsDir() || skipDirs[item.Name()] {
			continue
		}
		mdPath := filepath.Join(publicDir, item.Name(), site.IndexMarkdown)
		raw, err := os.ReadFile(mdPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, nil, err
		}

		meta, body, err := md.Parse(string(raw))
		if err != nil {
			return nil, nil, fmt.Errorf("%s/%s: %w", item.Name(), site.IndexMarkdown, err)
		}
		if strings.TrimSpace(meta.Title) == "" {
			return nil, nil, fmt.Errorf("%s/%s is missing a title", item.Name(), site.IndexMarkdown)
		}

		ent := entry{
			Slug:   item.Name(),
			Title:  meta.Title,
			Body:   body,
			IsPage: meta.Type == "page",
		}
		if ent.IsPage {
			pages = append(pages, ent)
			continue
		}
		if meta.Date.IsZero() {
			return nil, nil, fmt.Errorf("%s/%s is missing a date", item.Name(), site.IndexMarkdown)
		}
		ent.Date = meta.Date
		posts = append(posts, ent)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, pages, nil
}

func documentTitle(entryTitle, subtitle string) string {
	return entryTitle + " - " + subtitle
}

func renderPost(cfg config.Config, tmpls templates, post entry) (string, error) {
	article, err := execTemplate(tmpls.article, pageData{
		EntryTitle: post.Title,
		Date:       formatDate(post.Date),
		Body:       post.Body,
	})
	if err != nil {
		return "", err
	}
	return renderPage(cfg, tmpls, documentTitle(post.Title, cfg.SiteSubtitle), article)
}

func renderStaticPage(cfg config.Config, tmpls templates, page entry) (string, error) {
	article, err := execTemplate(tmpls.article, pageData{
		EntryTitle: page.Title,
		Body:       page.Body,
	})
	if err != nil {
		return "", err
	}
	return renderPage(cfg, tmpls, documentTitle(page.Title, cfg.SiteSubtitle), article)
}

func listingIntro(cfg config.Config, md markdown.Markdown) (template.HTML, error) {
	src := strings.TrimSpace(cfg.ListingIntro)
	if src == "" {
		return "", nil
	}
	body, err := md.HTML(src)
	if err != nil {
		return "", err
	}
	return template.HTML("<section class=\"listing-intro\">\n" + string(body) + "</section>\n"), nil
}

func renderListing(cfg config.Config, tmpls templates, posts []entry, md markdown.Markdown) (string, error) {
	intro, err := listingIntro(cfg, md)
	if err != nil {
		return "", err
	}
	items := make([]listPost, len(posts))
	for i, post := range posts {
		items[i] = listPost{
			Title: post.Title,
			Date:  formatDate(post.Date),
			Slug:  post.Slug,
		}
	}
	listing, err := execTemplate(tmpls.listing, pageData{
		ListingIntro:   intro,
		ListingHeading: cfg.ListingHeading,
		Posts:          items,
	})
	if err != nil {
		return "", err
	}
	return renderPage(cfg, tmpls, cfg.SiteTitle, listing)
}

func renderPage(cfg config.Config, tmpls templates, title, content string) (string, error) {
	data := pageData{
		Title:       title,
		Description: cfg.SiteDescription,
		SiteHeader:  cfg.SiteHeader,
		SiteFooter:  cfg.SiteFooter,
	}
	header, err := execTemplate(tmpls.header, data)
	if err != nil {
		return "", err
	}
	footer, err := execTemplate(tmpls.footer, data)
	if err != nil {
		return "", err
	}
	return header + content + footer, nil
}

func execTemplate(tmpl *template.Template, data pageData) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func writeFile(path, content string) (bool, error) {
	if existing, err := os.ReadFile(path); err == nil && string(existing) == content {
		return false, nil
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func countLabel(n int, singular string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %ss", n, singular)
}

func formatDate(date time.Time) string {
	return date.UTC().Format("2006-01-02")
}
