package main

import (
	"embed"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/unix1/blug/internal/subcmd"
)

//go:embed templates/*.html
//go:embed scripts/deploy.sh
//go:embed blug.yaml
//go:embed .env.example
//go:embed .gitignore
//go:embed public/assets/style.css
//go:embed public/assets/favicon.svg
//go:embed public/index.md
//go:embed public/hello-world/index.md
var embedded embed.FS

// Set by GoReleaser: -X main.version={{.Tag}}
var version string

func versionString() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return ""
}

func usage() {
	fmt.Fprint(os.Stderr, `blug generate   write listing, post, and page HTML
blug serve      serve public/ at http://localhost:3000
blug deploy     generate, then rclone sync public/ to destination
blug init       write a starter site in the current directory
blug version    print version
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version", "--version":
		fmt.Println(versionString())
		return
	case "help", "-h", "--help":
		usage()
		return
	}

	siteRoot, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	switch os.Args[1] {
	case "generate":
		err = subcmd.Generate(siteRoot)
	case "serve":
		err = subcmd.Serve(siteRoot)
	case "deploy":
		err = subcmd.Deploy(siteRoot, embedded)
	case "init":
		err = subcmd.Init(siteRoot, embedded)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
