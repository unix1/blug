package generate

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/unix1/blug/internal/config"
	"github.com/unix1/blug/internal/markdown"
	"go.yaml.in/yaml/v4"
)

type stubMarkdown struct{}

func (stubMarkdown) Parse(src string) (markdown.Meta, template.HTML, error) {
	meta, body, err := stubFrontmatter(src)
	return meta, template.HTML(body), err
}

func (stubMarkdown) HTML(src string) (template.HTML, error) {
	return template.HTML("<p>" + src + "</p>"), nil
}

func stubFrontmatter(src string) (markdown.Meta, string, error) {
	src = strings.TrimPrefix(src, "\n")
	if !strings.HasPrefix(src, "---") {
		return markdown.Meta{}, src, nil
	}
	parts := strings.SplitN(src, "---", 3)
	if len(parts) < 3 {
		return markdown.Meta{}, src, nil
	}
	var raw struct {
		Title string `yaml:"title"`
		Date  string `yaml:"date"`
		Type  string `yaml:"type"`
	}
	if err := yaml.Unmarshal([]byte(parts[1]), &raw); err != nil {
		return markdown.Meta{}, "", err
	}
	var date time.Time
	if strings.TrimSpace(raw.Date) != "" {
		parsed, err := time.Parse("2006-01-02", raw.Date)
		if err != nil {
			return markdown.Meta{}, "", err
		}
		date = parsed
	}
	return markdown.Meta{Title: raw.Title, Date: date, Type: raw.Type}, strings.TrimPrefix(parts[2], "\n"), nil
}

func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func testTemplates(t *testing.T) templates {
	t.Helper()
	must := func(name, src string) *template.Template {
		t.Helper()
		tmpl, err := template.New(name).Parse(src)
		if err != nil {
			t.Fatal(err)
		}
		return tmpl
	}
	return templates{
		header:  must("header", "<title>{{.Title}}</title>\n"),
		footer:  must("footer", "<footer>{{.SiteFooter}}</footer>\n"),
		listing: must("listing", "{{if .Body}}<section class=\"listing-intro\">{{.Body}}</section>{{end}}<h1>{{.ListingHeading}}</h1>\n{{range .Posts}}<a href=\"{{.Slug}}/\">{{.Title}}</a>\n{{end}}"),
		article: must("article", "<article>\n<h1>{{.EntryTitle}}</h1>\n{{if .Date}}<time datetime=\"{{.Date}}\">{{.Date}}</time>\n{{end}}{{.Body}}\n</article>\n"),
	}
}

func TestRun(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"blug.yaml": `site_header: Header
site_title: Site Title
site_subtitle: Sub
site_description: Desc
site_footer: Footer
listing_heading: Posts
`,
		"templates/header.html":  "<title>{{.Title}}</title>\n",
		"templates/footer.html":  "<footer>{{.SiteFooter}}</footer>\n",
		"templates/listing.html": "<h1>{{.ListingHeading}}</h1>\n{{range .Posts}}<a href=\"{{.Slug}}/\">{{.Title}}</a>\n{{end}}",
		"templates/article.html": "<article>{{.EntryTitle}}{{.Body}}</article>\n",
		"public/index.md":        "\n",
		"public/hello/index.md":  "---\ntitle: Hello\ndate: 2026-01-01\n---\n\nHi.\n",
	})

	stats, err := Run(root, stubMarkdown{})
	if err != nil {
		t.Fatal(err)
	}
	if stats.GeneratedPosts != 1 || stats.GeneratedPages != 0 || !stats.WroteListing {
		t.Fatalf("stats: %+v", stats)
	}

	post := filepath.Join(root, "public/hello/index.html")
	listing := filepath.Join(root, "public/index.html")
	if _, err := os.Stat(post); err != nil {
		t.Fatalf("post: %v", err)
	}
	if _, err := os.Stat(listing); err != nil {
		t.Fatalf("listing: %v", err)
	}
}

func TestLoadTemplatesMissing(t *testing.T) {
	_, err := loadTemplates(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "templates/header.html") {
		t.Fatalf("got %v", err)
	}
}

func TestFindEntries(t *testing.T) {
	publicDir := t.TempDir()
	writeTree(t, publicDir, map[string]string{
		"older/index.md":  "---\ntitle: Older\ndate: 2026-01-01\n---\n\nolder-body\n",
		"newer/index.md":  "---\ntitle: Newer\ndate: 2026-02-01\n---\n\nnewer-body\n",
		"about/index.md":  "---\ntitle: About\ntype: page\n---\n\nabout-body\n",
		"assets/index.md": "---\ntitle: Style notes\ndate: 2026-03-01\n---\n\nshould-skip\n",
		"empty/keep.txt":  "no index.md\n",
		"readme.md":       "not a directory\n",
		"index.md":        "home intro should not be a post\n",
	})

	posts, pages, err := findEntries(publicDir, stubMarkdown{})
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 || posts[0].Slug != "newer" || posts[1].Slug != "older" {
		t.Fatalf("posts: %+v", posts)
	}
	if posts[0].Title != "Newer" || !strings.Contains(string(posts[0].Body), "newer-body") {
		t.Fatalf("newer: %+v", posts[0])
	}
	if len(pages) != 1 || pages[0].Slug != "about" || !pages[0].IsPage {
		t.Fatalf("pages: %+v", pages)
	}
	if !pages[0].Date.IsZero() {
		t.Fatalf("page date should stay zero: %+v", pages[0])
	}
}

func TestFindEntriesMissingTitle(t *testing.T) {
	publicDir := t.TempDir()
	writeTree(t, publicDir, map[string]string{
		"hello/index.md": "---\ndate: 2026-01-01\n---\n\nHi.\n",
	})
	_, _, err := findEntries(publicDir, stubMarkdown{})
	if err == nil || !strings.Contains(err.Error(), "missing a title") {
		t.Fatalf("got %v", err)
	}
}

func TestFindEntriesMissingDate(t *testing.T) {
	publicDir := t.TempDir()
	writeTree(t, publicDir, map[string]string{
		"hello/index.md": "---\ntitle: Hello\n---\n\nHi.\n",
	})
	_, _, err := findEntries(publicDir, stubMarkdown{})
	if err == nil || !strings.Contains(err.Error(), "missing a date") {
		t.Fatalf("got %v", err)
	}
}

func TestListingIntro(t *testing.T) {
	emptyDir := t.TempDir()
	writeTree(t, emptyDir, map[string]string{"index.md": "  \n"})
	empty, err := listingIntro(emptyDir, stubMarkdown{})
	if err != nil || empty != "" {
		t.Fatalf("empty intro: %q %v", empty, err)
	}

	introDir := t.TempDir()
	writeTree(t, introDir, map[string]string{"index.md": "welcome"})
	intro, err := listingIntro(introDir, stubMarkdown{})
	if err != nil {
		t.Fatal(err)
	}
	if intro != "<p>welcome</p>" {
		t.Fatalf("intro: %q", intro)
	}
}

func TestListingIntroMissing(t *testing.T) {
	_, err := listingIntro(t.TempDir(), stubMarkdown{})
	if err == nil || !strings.Contains(err.Error(), "public/index.md") {
		t.Fatalf("got %v", err)
	}
}

func TestRenderPost(t *testing.T) {
	html, err := renderPost(config.Config{SiteSubtitle: "Sub", SiteFooter: "Footer"}, testTemplates(t), entry{
		Title: "Hello",
		Date:  time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		Body:  "hello-body",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<title>Hello - Sub</title>") {
		t.Fatalf("title: %s", html)
	}
	if !strings.Contains(html, `<time datetime="2026-02-01">2026-02-01</time>`) {
		t.Fatalf("time: %s", html)
	}
	if !strings.Contains(html, "hello-body") {
		t.Fatalf("body: %s", html)
	}
}

func TestRenderStaticPage(t *testing.T) {
	html, err := renderStaticPage(config.Config{SiteSubtitle: "Sub", SiteFooter: "Footer"}, testTemplates(t), entry{
		Title:  "About",
		Body:   "about-body",
		IsPage: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<title>About - Sub</title>") {
		t.Fatalf("title: %s", html)
	}
	if strings.Contains(html, "<time") {
		t.Fatalf("page should omit time: %s", html)
	}
	if !strings.Contains(html, "about-body") {
		t.Fatalf("body: %s", html)
	}
}

func TestRenderListing(t *testing.T) {
	publicDir := t.TempDir()
	writeTree(t, publicDir, map[string]string{"index.md": "welcome"})
	html, err := renderListing(config.Config{
		SiteTitle:      "Site Title",
		ListingHeading: "Posts",
		SiteFooter:     "Footer",
	}, testTemplates(t), publicDir, []entry{
		{Title: "Newer", Slug: "newer", Date: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
		{Title: "Older", Slug: "older", Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}, stubMarkdown{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<title>Site Title</title>") {
		t.Fatalf("title: %s", html)
	}
	if !strings.Contains(html, `<section class="listing-intro">`) || !strings.Contains(html, "<p>welcome</p>") {
		t.Fatalf("intro: %s", html)
	}
	if !strings.Contains(html, `<a href="newer/">Newer</a>`) || !strings.Contains(html, `<a href="older/">Older</a>`) {
		t.Fatalf("links: %s", html)
	}
}

func TestWriteFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.html")

	wrote, err := writeFile(path, "one")
	if err != nil || !wrote {
		t.Fatalf("first write: wrote=%v err=%v", wrote, err)
	}

	wrote, err = writeFile(path, "one")
	if err != nil || wrote {
		t.Fatalf("unchanged: wrote=%v err=%v", wrote, err)
	}

	wrote, err = writeFile(path, "two")
	if err != nil || !wrote {
		t.Fatalf("changed: wrote=%v err=%v", wrote, err)
	}
	raw, err := os.ReadFile(path)
	if err != nil || string(raw) != "two" {
		t.Fatalf("content: %q %v", raw, err)
	}
}

func TestDocumentTitle(t *testing.T) {
	if got := documentTitle("Hello", "Blog"); got != "Hello - Blog" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatDate(t *testing.T) {
	date := time.Date(2026, 1, 2, 1, 0, 0, 0, time.FixedZone("ahead", 5*3600))
	if got := formatDate(date); got != "2026-01-01" {
		t.Fatalf("got %q", got)
	}
}

func TestCountLabel(t *testing.T) {
	if got := countLabel(1, "post"); got != "1 post" {
		t.Fatalf("got %q", got)
	}
	if got := countLabel(2, "post"); got != "2 posts" {
		t.Fatalf("got %q", got)
	}
}
