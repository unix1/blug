package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	root := t.TempDir()
	raw := `site_header: Header
site_title: Title
site_subtitle: Sub
site_description: Desc
site_footer: Footer
listing_heading: Posts
`
	if err := os.WriteFile(filepath.Join(root, "blug.yaml"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SiteHeader != "Header" || cfg.SiteTitle != "Title" || cfg.SiteSubtitle != "Sub" {
		t.Fatalf("site fields: %+v", cfg)
	}
	if cfg.SiteDescription != "Desc" || cfg.SiteFooter != "Footer" {
		t.Fatalf("site chrome: %+v", cfg)
	}
	if cfg.ListingHeading != "Posts" {
		t.Fatalf("listing: %+v", cfg)
	}
}

func TestLoadMissing(t *testing.T) {
	_, err := Load(t.TempDir())
	if err == nil {
		t.Fatal("expected missing config error")
	}
	if !strings.Contains(err.Error(), "blug.yaml") {
		t.Fatalf("error should mention blug.yaml: %v", err)
	}
}
