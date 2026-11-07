package subcmd

import (
	"os"

	"github.com/unix1/blug/internal/generate"
	"github.com/unix1/blug/internal/markdown"
)

func Generate(siteRoot string) error {
	stats, err := generate.Run(siteRoot, markdown.NewGoldmark())
	if err != nil {
		return err
	}
	stats.Print(os.Stdout)
	return nil
}
