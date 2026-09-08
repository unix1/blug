# Setup

Agents can also do all of this work. These steps are here for agent consumption and manual reference.

## Cloudflare

Host your blog on Cloudflare, one of the top CDNs in the world for free.

### Domain

1. Create a [Cloudflare account](https://dash.cloudflare.com/sign-up).
2. Add your domain to Cloudflare and follow the DNS setup instructions.
3. Wait until DNS has propagated and the zone is active.

### CDN setup

We will set up two CDN rules.

#### Rule 1: redirect all */ requests to */index.html

> ℹ️ Why? R2 serves exact object keys. It does not map a directory path to `index.html`. The rewrite sends any request ending in `/` to `index.html` in that folder. That is common web server behavior that R2 does not offer on its own.

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

#### Rule 2: make all objects eligible for cache

> ℹ️ Why? A cache-everything rule is optional but recommended so the CDN can cache HTML. Cloudflare does not cache HTML by default.

1. Open the domain → **Caching** → **Cache Rules** (or **Rules** → **Templates**).
2. Use the **Cache everything** template, and name the rule `Cache everything`.
3. When incoming requests match: **All incoming requests**.
4. Then: **Eligible for cache**.
5. Leave the rest at defaults and deploy.

### R2 bucket

> ℹ️ R2 has a generous free tier. Cloudflare may still ask for payment information. Check [Cloudflare pricing](https://developers.cloudflare.com/r2/pricing/) so that is not a surprise.

1. Create an R2 bucket. Any name is fine; you will use it later when configuring deploy.
2. In the bucket **Settings**, add a **custom domain**: the domain from the previous step.
   - The zone must already be active in the same Cloudflare account. If it is not ready, this step fails.
3. Create an account-level API token with read and write access to that R2 bucket. Copy all values including Access Key ID, Secret Access Key, and the S3 API endpoint.

## Local development

1. On [github.com/unix1/blug](https://github.com/unix1/blug), click **Use this template** and create a repository of your own.
   - Note: `git` is not required. If you'd rather not use it, copy the files from the `main` branch wherever you like. You can download the [latest blug zip archive](https://github.com/unix1/blug/archive/refs/heads/main.zip) from GitHub and unzip in your preferred location.
2. If you are using `git` for convenience, clone your newly created repository (if not, you can skip this step):

   ```bash
   git clone <your-repo-url>
   cd <your-repo>
   ```

3. Copy `.env.example` to `.env` and fill in the R2 values:

   ```
   S3_PROVIDER=Cloudflare
   S3_REGION=auto
   S3_ENDPOINT=
   S3_ACCESS_KEY_ID=
   S3_SECRET_ACCESS_KEY=
   S3_BUCKET=
   CACHE_CONTROL="public, max-age=5"
   CACHE_CONTROL_REUPLOAD=0
   ```

   > ℹ️ Note: `max-age=5` above will only cache your objects for 5 seconds in the CDN. This value should be OK to get started while you are iterating on your site. Once you are set up, you may want to change this to a larger value. See the [CDN Cache Control](https://blug.blog/cdn-cache-control/) blog post for more details.

4. Install and generate:

   ```bash
   npm install && npm run generate
   ```

5. Preview:

   ```bash
   npm run dev
   ```

   Open http://localhost:3000 (or the URL printed in the logs).

6. To iterate, add a post under `public/<slug>/index.md` and run `npm run generate` again, then refresh.

## Production deploy

Install [rclone](https://rclone.org/) if you do not already have it. After Cloudflare hosting and `.env` are set:

```bash
npm run deploy
```

That runs generate, then `rclone sync` of `public/` to the R2 bucket.

You are done. Open your custom domain: the site is live. Congratulations — great job! You are hosting on one of the top CDNs in the world for free. Go write blog posts. Do not worry about hosting.
