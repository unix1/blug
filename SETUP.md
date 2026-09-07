# Setup

Run a blug blog locally and host it on Cloudflare R2.

## Cloudflare

Do this before `npm run deploy`. The rewrite and cache rules belong on the **zone** (your proxied custom domain), not on `r2.dev`.

### Domain

1. Create a [Cloudflare account](https://dash.cloudflare.com/sign-up).
2. Add your domain to Cloudflare and follow the DNS setup instructions.
3. Wait until DNS has propagated and the zone is active.

### URL rewrite

R2 serves exact object keys. It does not map `/` or `/hello-world/` to `index.html`. A rewrite tells Cloudflare to fetch `index.html` while the browser URL keeps the trailing slash, so relative media links like `./photo.jpg` resolve.

1. Open the domain → **Rules** → **Overview**.
2. **Create rule** → **URL Rewrite Rule**.
3. Name: `Rewrite requests */ to */index.html`.
4. **If incoming requests match**: Custom filter expression.
   - Field: URI Path
   - Operator: ends with
   - Value: `/`
5. **Then**:
   - Path: **Rewrite to** → **Dynamic** → `concat(http.request.uri.path, "index.html")`
   - Query: **Preserve**
6. Deploy the rule.

That also covers `/` → `/index.html`. The rewrite is internal: the address bar stays `/hello-world/`.

### Cache

Cloudflare does not cache HTML by default. A cache-everything rule makes all files eligible, including the generated pages. Without it, HTML stays `DYNAMIC` and every request hits R2.

1. Open the domain → **Caching** → **Cache Rules** (or **Rules** → **Templates**).
2. Use the **Cache everything** template, or create a rule named `Cache everything`.
3. When incoming requests match: **All incoming requests**.
4. Then: **Eligible for cache**.
5. Leave the rest at defaults and deploy.

Deploy also uploads `CACHE_CONTROL` from `.env` as object metadata so Cloudflare Cache has an edge TTL to use.

### R2 bucket

1. Create an R2 bucket. Any name is fine; you will put it in `.env` later.
2. In the bucket **Settings**, add a **custom domain**: the domain from the previous step.
   - The zone must already be active in the same Cloudflare account. If it is not ready, this step fails.
3. Create an **account** R2 API token with Object Read & Write access to that bucket:
   - R2 overview → **Manage** next to **API Tokens** → **Create Account API token**
   - Copy **Access Key ID**, **Secret Access Key**, and the S3 API endpoint (`https://<ACCOUNT_ID>.r2.cloudflarestorage.com`). Keep these private.

## Local development

1. On [github.com/unix1/blug](https://github.com/unix1/blug), click **Use this template** and create a repository of your own.
   - Without GitHub, copy the files from the `main` branch wherever you like.
2. Clone that repository:

   ```bash
   git clone <your-repo-url>
   cd <your-repo>
   ```

3. Copy `.env.example` to `.env` and fill in the R2 values:

   ```
   S3_PROVIDER=Cloudflare
   S3_REGION=auto
   S3_ENDPOINT=https://<ACCOUNT_ID>.r2.cloudflarestorage.com
   S3_ACCESS_KEY_ID=...
   S3_SECRET_ACCESS_KEY=...
   S3_BUCKET=<your-bucket-name>
   CACHE_CONTROL="public, max-age=3600"
   ```

   Quote `CACHE_CONTROL`. Deploy sources `.env` as a shell file, so an unquoted value with commas is parsed as a command.

4. Install and generate:

   ```bash
   npm install && npm run generate
   ```

5. Preview:

   ```bash
   npm run dev
   ```

   Open http://127.0.0.1:3000 (or the URL printed in the logs).

6. To iterate, add a post under `public/<slug>/index.md` and run `npm run generate` again, then refresh. Generate writes post HTML only when `index.html` is missing; delete that file to regenerate an existing post.

## Production deploy

Install [rclone](https://rclone.org/) if you do not already have it. After Cloudflare hosting and `.env` are set:

```bash
npm run deploy
```

That runs generate, then `rclone sync` of `public/` to the R2 bucket.
