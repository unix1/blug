# blug

**blug** is a simple static blog and site generator.

1. Write blog posts or static pages in markdown in `public/<my-page-path>/index.md`.
2. Run `blug generate` to generate HTML from markdown.
3. Run `blug serve` to preview your site locally.
4. Run `blug deploy` to deploy your site to Cloudflare or another S3 compatible storage.

`rclone` is required for the deploy step. You can install it via your preferred method (Linux distro, homebrew on MacOS, etc.).

Install:

```bash
go install github.com/unix1/blug@latest
```

See [SETUP.md](SETUP.md) for local preview and Cloudflare R2 hosting.

## Layout

```
blug.yaml            # header, titles, footer, listing heading, optional home intro
.env                 # S3 credentials for rclone (not committed)
templates/           # header.html, footer.html, listing.html, article.html
public/
  index.html         # generated listing
  assets/style.css
  hello-world/
    index.md         # you write this
    index.html       # generated post
    photo.jpg        # optional media
```

## Write a post

Create `public/<slug>/index.md`:

```markdown
---
title: Hello World
date: 2026-08-30
---

Text. Media next to this file: [photo](./photo.jpg).
```

`title` and `date` are required. Keep media in the same directory and use relative links.

## Write a page

Pages use the same folder layout as posts, but with `type: page` and no date. They are not listed on the home page.

Create `public/<slug>/index.md`:

```markdown
---
title: Privacy Policy
type: page
---

Your content here.
```

`title` and `type: page` are required. Link to pages the same way as posts (`/privacy-policy/`).

Edit `blug.yaml` for the site header, listing page title (`site_title`), post/page title suffix (`site_subtitle`), meta description (`site_description`), footer, listing heading, and optional home-page intro (`listing_intro`, markdown shown above the post list). Posts and pages use `<title> - site_subtitle` as the document title. Leave `listing_intro` empty to omit it.

`blug init` writes `templates/`. Edit those files to change the HTML; generate does not fall back to copies inside the binary. They use Go `html/template` fields such as `{{.Title}}` and `{{range .Posts}}`.

## Commands

```bash
blug init        # starter site in the current directory
blug generate    # write listing, post, and page HTML
blug serve       # serve public/ at http://localhost:3000
blug deploy      # generate, then rclone sync public/ to destination
```

## Hosting

This supports S3-compatible hosting, such as Cloudflare R2. Put `S3_*` values in `.env` (see `.env.example`). `blug deploy` passes those to rclone; no `rclone.conf` is needed. Set `CACHE_CONTROL` in `.env` to control Cache-Control on uploaded objects.

R2 serves exact object keys. It does not map `/hello-world/` to `hello-world/index.html`. Production needs a URL rewrite so paths ending in `/` fetch `index.html`. A cache-everything rule is optional but recommended so the CDN can cache HTML. Post listing links include a trailing slash (`hello-world/`) so relative media links resolve.

Full steps: [SETUP.md](SETUP.md).

## License

[MIT](LICENSE)
