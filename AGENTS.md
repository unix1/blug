# Agent notes

Keep this blog generator small. Do not add extra template engines, `<base>` tags, or recursive post discovery unless asked. `html/template` is the only template language. Markdown is [goldmark v2](https://github.com/yuin/goldmark) with [goldmark-meta/v2](https://github.com/yuin/goldmark-meta) for YAML frontmatter. Site YAML is `gopkg.in/yaml.v3`. Do not add other third-party packages unless asked.

`main.go` is the CLI entry (`go install github.com/unix1/blug@latest`). Packages under `internal/` are not importable outside this module: `config` (blug.yaml), `markdown` (`Markdown` interface; goldmark is the only implementation), `generate` (post/page/listing rules; takes `markdown.Markdown`), `site` (path constants), `subcmd` (generate, serve, deploy, init). Scaffold lives in `internal/subcmd/initialize.go`.

## Layout

- `public/` is the site root and the rclone/R2 upload tree.
- Posts and pages are one level deep: `public/<slug>/index.md` plus any media. Skip `public/assets/`.
- Pages use frontmatter `type: page` (no `date`). They share the same slug/URL layout as posts but are excluded from the home listing and render without a `<time>`.
- `blug.yaml` holds `site_header`, `site_title`, `site_subtitle`, `site_description`, `site_footer`, `listing_heading`, and `listing_intro`.
- `blug init` writes `templates/`. Generate reads only those site files; there is no embed fallback. Header/footer wrap every page (`{{.Title}}`, `{{.Description}}`, `{{.SiteHeader}}`, `{{.SiteFooter}}`). `listing.html` is the home body (`{{.ListingIntro}}`, `{{.ListingHeading}}`, `{{range .Posts}}` with `{{.Date}}`, `{{.Slug}}`, `{{.Title}}`). `article.html` is the post/page body (`{{.EntryTitle}}`, `{{.Date}}`, `{{.Body}}`); empty `{{.Date}}` omits `<time>` on pages. `listing_intro` is markdown; empty string omits the intro. Pass goldmark output and listing intro HTML as `template.HTML`.
- `blug generate` is the only build step. It reads `blug.yaml` from the current directory.
- `.env` holds S3 provider, region, endpoint, access keys, bucket, and optional `CACHE_CONTROL`. Generate does not read it. Deploy does (`scripts/deploy.sh`, run with the site as cwd).

## Generate

`blug generate`:

1. Scan `public/*/` for `index.md`.
2. If `type: page`, require YAML frontmatter `title` only. Otherwise require `title` and `date`.
3. Convert markdown with goldmark v2 and write `public/<slug>/index.html`. Posts include a `<time>`; pages do not. Document titles are `<entry title> - site_subtitle`. Skip the write when the file is already up to date.
4. Write `public/index.html` from posts only (newest date first), using `site_title` as the document title, including `listing_intro` above the post list when set.

Leave media files untouched. Markdown relative links (`./photo.jpg`) must stay relative in the HTML.

## URLs

- Listing links use a trailing slash: `hello-world/`, not `hello-world` and not `/hello-world/`.
- Site chrome uses root-absolute paths (`/`, `/assets/style.css`).
- Do not introduce a blog-base config. Posts live at the site root until someone moves them on purpose.

## Hosting (R2)

R2 does not serve `index.html` for directory paths. Production needs a Cloudflare **URL Rewrite** on the zone (proxied custom domain, not `r2.dev`): paths ending in `/` rewrite to `concat(http.request.uri.path, "index.html")`. That includes `/`.

Without the trailing slash in the browser URL, `./photo.jpg` resolves to `/photo.jpg`. Do not "fix" that with `<base href>`.

`blug deploy` runs generate, then rclone sync of `public/` using R2 credentials from `.env` (flags, not `rclone.conf`). If `CACHE_CONTROL` is set, rclone passes `--header-upload Cache-Control: ...` so R2 stores it as object metadata. Set `CACHE_CONTROL_REUPLOAD=1` once after changing it so unchanged objects are rewritten. A Cloudflare Cache Rule with Eligible for cache is optional but recommended so HTML and extensionless URLs are not `DYNAMIC` (every request hitting R2). Do not reimplement rclone in Go; exec `scripts/deploy.sh` with the site cwd.

Setup steps for domain, rewrite, cache, R2 custom domain, API tokens, and `.env`: [SETUP.md](SETUP.md).

## Preview

`blug serve` is stdlib `http.FileServer` on `public/` at `http://localhost:3000`.
