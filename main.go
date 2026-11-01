package main

import (
	"fmt"
	"os"

	"github.com/unix1/blug/internal/subcmd"
)

func usage() {
	fmt.Fprint(os.Stderr, `blug generate   write listing, post, and page HTML
blug serve      serve public/ at http://localhost:3000
blug deploy     generate, then rclone sync public/ to destination
blug init       write a starter site in the current directory
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	siteRoot, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	switch os.Args[1] {
	case "generate":
		err = subcmd.Generate(siteRoot, embedded)
	case "serve":
		err = subcmd.Serve(siteRoot)
	case "deploy":
		err = subcmd.Deploy(siteRoot, embedded)
	case "init":
		err = subcmd.Init(siteRoot, embedded)
	case "help", "-h", "--help":
		usage()
		return
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
