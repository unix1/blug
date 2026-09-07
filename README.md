# blug

A static blog generator. Posts are markdown folders under `public/`. `scripts/generate.js` turns them into HTML.

See [SETUP.md](SETUP.md) for local preview and Cloudflare R2 hosting.

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

This supports S3-compatible hosting, such as Cloudflare R2. Put `S3_*` values in `.env` (see `.env.example`). `npm run deploy` passes those to rclone; no `rclone.conf` is needed. Set `CACHE_CONTROL` in `.env` to control Cache-Control on uploaded objects.

R2 serves exact object keys. It does not map `/hello-world/` to `hello-world/index.html`. Production needs a URL rewrite so paths ending in `/` fetch `index.html`, plus a cache-everything rule so HTML is eligible for cache. Post listing links include a trailing slash (`hello-world/`) so relative media links resolve.

Full dashboard steps: [SETUP.md](SETUP.md).

## License

[MIT](LICENSE)
