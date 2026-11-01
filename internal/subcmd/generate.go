package subcmd

import (
	"embed"
	"os"

	"github.com/unix1/blug/internal/generate"
	"github.com/unix1/blug/internal/markdown"
)

func Generate(siteRoot string, defaults embed.FS) error {
	stats, err := generate.Run(siteRoot, defaults, markdown.NewGoldmark())
	if err != nil {
		return err
	}
	stats.Print(os.Stdout)
	return nil
}
