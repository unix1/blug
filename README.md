# blug

A static blog generator. Posts are markdown folders under `public/`. `scripts/generate.js` turns them into HTML.

## Layout

```
.env                 # S3 credentials for rclone (not committed)
templates/           # header.html, footer.html, listing.html
scripts/generate.js
scripts/config.js    # header, footer, listing heading, optional home intro
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

Edit `scripts/config.js` for the site title, footer, listing heading, and optional home-page intro (`LISTING_INTRO`, markdown shown above the post list). Leave `LISTING_INTRO` empty to omit it.

## Commands

```bash
npm install
npm run generate   # write listing and post HTML
npm run dev        # serve public/ at http://127.0.0.1:3000
npm run deploy     # generate, then rclone sync public/ to destination
```

## Hosting

This supports S3-compatible hosting, such as Cloudflare R2. Put `S3_*` values in `.env` (see `.env.example`). `npm run deploy` passes those to rclone; no `rclone.conf` is needed.

R2 serves exact object keys. It does not map `/hello-world/` to `hello-world/index.html`. Add one **URL Rewrite** on the domain (Rules → Overview → URL Rewrite Rule):

- If URI Path ends with `/`
- Rewrite path dynamically to `concat(http.request.uri.path, "index.html")`

That also covers `/` → `/index.html`.

Post links on the listing include a trailing slash (`hello-world/`) so the browser treats the page as a directory and relative media links resolve. The rewrite is internal: the address bar stays `/hello-world/`.
