package markdown

import (
	"strings"
	"testing"
	"time"
)

func TestParsePost(t *testing.T) {
	md := NewGoldmark()
	meta, body, err := md.Parse("---\ntitle: Hello\ndate: 2026-09-05\n---\n\nHello World.\n")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Title != "Hello" {
		t.Fatalf("title: got %q", meta.Title)
	}
	if meta.Type != "" {
		t.Fatalf("type: got %q", meta.Type)
	}
	if !meta.Date.Equal(time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("date: got %v", meta.Date)
	}
	if !strings.Contains(string(body), "Hello World.") {
		t.Fatalf("body: %s", body)
	}
}

func TestParseFeaturedImage(t *testing.T) {
	md := NewGoldmark()
	meta, _, err := md.Parse("---\ntitle: Hello\ndate: 2026-09-05\nfeatured_image: ./cover.jpg\n---\n\nHi.\n")
	if err != nil {
		t.Fatal(err)
	}
	if meta.FeaturedImage != "./cover.jpg" {
		t.Fatalf("featured_image: got %q", meta.FeaturedImage)
	}
	plain, _, err := md.Parse("---\ntitle: Hello\ndate: 2026-09-05\n---\n\nHi.\n")
	if err != nil {
		t.Fatal(err)
	}
	if plain.FeaturedImage != "" {
		t.Fatalf("missing featured_image: got %q", plain.FeaturedImage)
	}
}

func TestParsePage(t *testing.T) {
	md := NewGoldmark()
	meta, _, err := md.Parse("---\ntitle: About\ntype: page\n---\n\nAbout this site.\n")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Title != "About" || meta.Type != "page" {
		t.Fatalf("meta: %+v", meta)
	}
	if !meta.Date.IsZero() {
		t.Fatalf("page date should be zero, got %v", meta.Date)
	}
}

func TestParseInvalidDate(t *testing.T) {
	md := NewGoldmark()
	_, _, err := md.Parse("---\ntitle: Hello\ndate: not-a-date\n---\n\nHi.\n")
	if err == nil {
		t.Fatal("expected invalid date error")
	}
}

func TestHTMLRelativeLink(t *testing.T) {
	md := NewGoldmark()
	html, err := md.HTML("[photo](./photo.jpg)")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), `href="./photo.jpg"`) {
		t.Fatalf("relative link rewritten: %s", html)
	}
}

func TestHTMLRawHTML(t *testing.T) {
	md := NewGoldmark()
	html, err := md.HTML(`before <div class="x">raw</div> after`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), `<div class="x">raw</div>`) {
		t.Fatalf("raw HTML missing: %s", html)
	}
}

func TestHTMLImageStaysImage(t *testing.T) {
	md := NewGoldmark()
	for _, src := range []string{
		"![North ridge](./ridge.jpg)",
		"![North ridge](https://cdn.example.com/ridge.jpg)",
		"![](https://www.youtube.com/watch?v=dQw4w9WgXcQ)",
	} {
		html, err := md.HTML(src)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(html), "<video") {
			t.Fatalf("image rendered as video: %s\n%s", src, html)
		}
		if !strings.Contains(string(html), "<img") {
			t.Fatalf("missing img: %s\n%s", src, html)
		}
	}
	local, err := md.HTML("![North ridge](./ridge.jpg)")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(local), `src="./ridge.jpg"`) {
		t.Fatalf("relative image src: %s", local)
	}
	remote, err := md.HTML("![North ridge](https://cdn.example.com/ridge.jpg)")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(remote), `src="https://cdn.example.com/ridge.jpg"`) {
		t.Fatalf("remote image src: %s", remote)
	}
}

func TestHTMLVideoEmbed(t *testing.T) {
	md := NewGoldmark()
	local, err := md.HTML("![Walk to the creek](./walk.mp4)")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(local), `<video controls preload="metadata" src="./walk.mp4" aria-label="Walk to the creek"></video>`) {
		t.Fatalf("local video: %s", local)
	}
	remote, err := md.HTML("![Talk](https://cdn.example.com/talk.webm)")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(remote), `<video controls preload="metadata" src="https://cdn.example.com/talk.webm" aria-label="Talk"></video>`) {
		t.Fatalf("remote video: %s", remote)
	}
}

func TestHTMLVideoPoster(t *testing.T) {
	md := NewGoldmark()
	html, err := md.HTML("[![Walk](https://cdn.example.com/walk.jpg)](https://cdn.example.com/walk.mp4)")
	if err != nil {
		t.Fatal(err)
	}
	body := string(html)
	if strings.Contains(body, "<video") {
		t.Fatalf("poster rendered as video: %s", body)
	}
	if !strings.Contains(body, `href="https://cdn.example.com/walk.mp4"`) {
		t.Fatalf("video href: %s", body)
	}
	if !strings.Contains(body, `src="https://cdn.example.com/walk.jpg"`) {
		t.Fatalf("poster src: %s", body)
	}
	if !strings.Contains(body, "<a ") || !strings.Contains(body, "<img ") {
		t.Fatalf("poster markup: %s", body)
	}
	if !strings.Contains(body, `<span class="video-label">Walk</span>`) {
		t.Fatalf("poster label: %s", body)
	}

	empty, err := md.HTML("[![](./walk.jpg)](./walk.mp4)")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(empty), "video-label") {
		t.Fatalf("empty alt wrote a label: %s", empty)
	}
}

func TestHTMLGalleryList(t *testing.T) {
	md := NewGoldmark()
	html, err := md.HTML(strings.Join([]string{
		"- ![North ridge](./ridge.jpg)",
		"- [![Walk to the creek](./walk.jpg)](https://cdn.example.com/walk.mp4)",
	}, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), `class="gallery"`) {
		t.Fatalf("gallery class: %s", html)
	}
	if !strings.Contains(string(html), `<span class="video-label">Walk to the creek</span>`) {
		t.Fatalf("gallery poster label: %s", html)
	}

	plain, err := md.HTML("- ![North ridge](./ridge.jpg)\n- A note\n")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(plain), "gallery") {
		t.Fatalf("text list marked as gallery: %s", plain)
	}

	linked, err := md.HTML("- [![Other](./thumb.jpg)](./other-post/)\n")
	if err != nil {
		t.Fatal(err)
	}
	body := string(linked)
	if strings.Contains(body, "gallery") {
		t.Fatalf("post link marked as gallery: %s", body)
	}
	if !strings.Contains(body, `href="./other-post/"`) || !strings.Contains(body, "<img ") {
		t.Fatalf("post link markup: %s", body)
	}
	if strings.Contains(body, "video-label") {
		t.Fatalf("post link wrote a video label: %s", body)
	}
}

func TestHTMLCodeFenceUnchanged(t *testing.T) {
	md := NewGoldmark()
	html, err := md.HTML("```go\nfmt.Println(\"hi\")\n```\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), `<pre><code class="language-go">`) {
		t.Fatalf("code fence: %s", html)
	}
}

func TestHTMLWithoutFrontmatter(t *testing.T) {
	md := NewGoldmark()
	html, err := md.HTML("**Hi**")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), "<strong>Hi</strong>") {
		t.Fatalf("markdown not converted: %s", html)
	}
}
