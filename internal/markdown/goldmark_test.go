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
