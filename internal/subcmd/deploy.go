package subcmd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/unix1/blug/internal/site"
)

func Deploy(siteRoot string, defaults embed.FS) error {
	if err := Generate(siteRoot); err != nil {
		return err
	}

	script, err := defaults.ReadFile(site.DeployScript)
	if err != nil {
		return err
	}
	cmd := exec.Command("sh", "-s")
	cmd.Dir = siteRoot
	cmd.Stdin = strings.NewReader(string(script))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	return nil
}
