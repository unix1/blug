package subcmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/unix1/blug/internal/site"
)

func Init(siteRoot string, defaults embed.FS) error {
	starterMD := filepath.ToSlash(filepath.Join(site.PublicDir, site.StarterPost, site.IndexMarkdown))
	assets := filepath.Join(site.PublicDir, site.AssetsDir)
	files := []struct {
		src  string
		dest string
	}{
		{site.ConfigFile, site.ConfigFile},
		{site.EnvExample, site.EnvExample},
		{site.EnvExample, site.EnvFile},
		{site.Gitignore, site.Gitignore},
		{filepath.ToSlash(filepath.Join(assets, site.StyleSheet)), filepath.ToSlash(filepath.Join(assets, site.StyleSheet))},
		{filepath.ToSlash(filepath.Join(assets, site.Favicon)), filepath.ToSlash(filepath.Join(assets, site.Favicon))},
		{starterMD, starterMD},
		{filepath.ToSlash(filepath.Join(site.TemplatesDir, site.HeaderTemplate)), filepath.ToSlash(filepath.Join(site.TemplatesDir, site.HeaderTemplate))},
		{filepath.ToSlash(filepath.Join(site.TemplatesDir, site.FooterTemplate)), filepath.ToSlash(filepath.Join(site.TemplatesDir, site.FooterTemplate))},
		{filepath.ToSlash(filepath.Join(site.TemplatesDir, site.ListingTemplate)), filepath.ToSlash(filepath.Join(site.TemplatesDir, site.ListingTemplate))},
		{filepath.ToSlash(filepath.Join(site.TemplatesDir, site.ArticleTemplate)), filepath.ToSlash(filepath.Join(site.TemplatesDir, site.ArticleTemplate))},
	}

	for _, file := range files {
		dest := filepath.Join(siteRoot, file.dest)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if _, err := os.Stat(dest); err == nil {
			fmt.Printf("skip %s (exists)\n", file.dest)
			continue
		} else if !os.IsNotExist(err) {
			return err
		}
		raw, err := defaults.ReadFile(file.src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dest, raw, 0o644); err != nil {
			return err
		}
		fmt.Printf("wrote %s\n", file.dest)
	}
	return nil
}
