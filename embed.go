package main

import "embed"

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
