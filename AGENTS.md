# Agent notes

Keep this blog generator small. Do not add templating languages, `<base>` tags, custom HTTP servers, or recursive post discovery unless asked.

## Layout

- `public/` is the site root and the rclone/R2 upload tree.
- Posts are one level deep: `public/<slug>/index.md` plus any media. Skip `public/assets/`.
- `scripts/config.js` holds `SITE_HEADER`, `SITE_FOOTER`, `LISTING_HEADING`, and `LISTING_INTRO`.
- `templates/header.html` and `templates/footer.html` wrap every page. Placeholders are `{{title}}`, `{{site_header}}`, and `{{site_footer}}`.
- `templates/listing.html` is the home page body. Placeholders are `{{listing_intro}}`, `{{listing_heading}}`, and `{{posts}}`. `LISTING_INTRO` is markdown; empty string omits the intro.
- `scripts/generate.js` is the only build step. It reads `scripts/config.js`.
- `.env` holds S3 provider, region, endpoint, access keys, bucket, and optional `CACHE_CONTROL`. Generate does not read it. Deploy does (`scripts/deploy.sh`).

## Generate

`npm run generate` (`node scripts/generate.js`):

1. Scan `public/*/` for `index.md`.
2. Require YAML frontmatter `title` and `date`.
3. Convert markdown with `marked` and write `public/<slug>/index.html` if that file is missing. Delete it to regenerate.
4. Write `public/index.html` (newest date first), including `LISTING_INTRO` above the post list when set.

Leave media files untouched. Markdown relative links (`./photo.jpg`) must stay relative in the HTML.

## URLs

- Listing links use a trailing slash: `hello-world/`, not `hello-world` and not `/hello-world/`.
- Site chrome uses root-absolute paths (`/`, `/assets/style.css`).
- Do not introduce a blog-base config. Posts live at the site root until someone moves them on purpose.

## Hosting (R2)

R2 does not serve `index.html` for directory paths. Production needs a Cloudflare **URL Rewrite** on the zone (proxied custom domain, not `r2.dev`): paths ending in `/` rewrite to `concat(http.request.uri.path, "index.html")`. That includes `/`.

Without the trailing slash in the browser URL, `./photo.jpg` resolves to `/photo.jpg`. Do not "fix" that with `<base href>`.

`npm run deploy` runs generate, then rclone sync of `public/` using R2 credentials from `.env` (flags, not `rclone.conf`). If `CACHE_CONTROL` is set, rclone passes `--header-upload Cache-Control: ...` so R2 stores it as object metadata. Set `CACHE_CONTROL_REUPLOAD=1` once after changing it so unchanged objects are rewritten. A Cloudflare Cache Rule with Eligible for cache is optional but recommended so HTML and extensionless URLs are not `DYNAMIC` (every request hitting R2). Do not wrap rclone in Node.

Setup steps for domain, rewrite, cache, R2 custom domain, API tokens, and `.env`: [SETUP.md](SETUP.md).

## Preview

`npm run dev` is `serve public`. Do not replace it with a hand-rolled server.
