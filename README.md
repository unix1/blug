# blug

**blug** is a simple static blog and site generator.

1. Write blog posts or static pages in markdown in `public/<my-custom-slug>/index.md`.
2. Run `blug generate` to generate HTML from markdown.
3. Run `blug serve` to preview your site locally.
4. Run `blug deploy` to deploy your site to Cloudflare or another S3 compatible storage.

`rclone` is required for the deploy step. You can install it via your preferred method (Linux distro, homebrew on MacOS, etc.).

## Install

Option 1: download the prebuilt binary from GitHub releases and place it in your path.

Linux and macOS binaries are on the [GitHub Releases](https://github.com/unix1/blug/releases) page.

Option 2: if you use Go, you can run

```bash
go install github.com/unix1/blug@latest
```

See [SETUP.md](SETUP.md) for local preview and Cloudflare R2 hosting.

## Create a new blog

In a directory of your choosing where you'd like to keep your blog files run

```
blug init
```

This will create the following directories and files

```
blug.yaml            # site configuration: header, footer, title, description
.env                 # S3 credentials for publishing your site
templates/           # header.html, footer.html, listing.html, article.html you can modify
public/
  index.md           # home page intro you write
  index.html         # generated home page, including list of posts (don't modify this)
  assets/style.css   # make the site your own by updating CSS
  hello-world/
    index.md         # you write this page or post
    index.html       # generated post (don't modify this)
    photo.jpg        # optional media you provide
```

You can and should modify any of these to your liking, except the generated `index.html` files.

## Write a post

Create a new directory under `public` with any name you like. Then create `public/<slug>/index.md` with markdown contents such as:

```markdown
---
title: Hello World
date: 2026-08-30
---

This is sample text. Text can reference arbitrary media such as a [photo](./photo.jpg).
```

`title` and `date` are required. Keep media in the same directory and use relative links.

## Write a page

Pages work just like posts, but with `type: page` and no date. They are not automatically listed on the home page, but you can link to them from anywhere you like.

Create `public/<slug>/index.md`:

```markdown
---
title: Privacy Policy
type: page
---

Your content here.
```

`title` and `type: page` are required. Link to pages the same way as posts - if your slug was `privacy-policy`, your link will be to `/privacy-policy/`.

## Site config

Edit `blug.yaml`:
- `site_header`: header text that's displayed at the top of every page
- `site_title`: value of the title tag on the home page
- `site_subtitle`: subtitle appended to the page title on all pages except home
- `site_description`: value of the meta description tag
- `site_footer`: footer text that's displayed at the bottom of every page
- `listing_heading`: header text to call your list of posts - e.g. Posts, News, Updates or whatever you like

## Generate your site and preview locally

At any time you can generate your site HTML and preview locally on your computer.

To generate site HTML run

```
blug generate
```

To preview your site locally run

```
blug serve
```

Note that `blug` will try to start a web server on `http://localhost:3000`. If that port is busy it will try another random port that's available. In either case, it will print the full address in its output.

## Commands reference

```bash
blug init        # write a starter site in the current directory
blug generate    # write listing, post, and page HTML
blug serve       # serve public/ at http://localhost:3000 (or another available port)
blug deploy      # generate, then rclone sync public/ to destination
blug version     # print version
blug help        # print this help
```

## Hosting

You can take the contents of `public/` directory and copy them to the web server of your choosing.

However, `blug` comes with support for S3 compatible storage that can be fronted by a web service like a CDN. `blug` has been tested on Cloudflare R2.

To set up Cloudflare with `blug` see [SETUP.md](SETUP.md).

After you've done so and configured your `.env` variables you can run

```
blug deploy
```

to publish your site. Exciting stuff!

## License

[Apache License 2.0](LICENSE)
