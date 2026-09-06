import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import matter from "gray-matter";
import { marked } from "marked";
import { LISTING_HEADING, SITE_FOOTER, SITE_HEADER } from "./config.js";

const ROOT = path.join(path.dirname(fileURLToPath(import.meta.url)), "..");
const PUBLIC = path.join(ROOT, "public");
const TEMPLATES = path.join(ROOT, "templates");
const SKIP_DIRS = new Set(["assets"]);

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

function loadTemplate(name) {
  return fs.readFileSync(path.join(TEMPLATES, name), "utf8");
}

function renderPage(title, content) {
  const header = loadTemplate("header.html")
    .replaceAll("{{title}}", escapeHtml(title))
    .replaceAll("{{site_header}}", escapeHtml(SITE_HEADER));
  const footer = loadTemplate("footer.html").replaceAll("{{site_footer}}", escapeHtml(SITE_FOOTER));
  return header + content + footer;
}

function toDate(value) {
  if (value instanceof Date && !Number.isNaN(value.getTime())) {
    return value;
  }
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    throw new Error(`Invalid date: ${value}`);
  }
  return parsed;
}

function formatDate(date) {
  return date.toISOString().slice(0, 10);
}

function findPosts() {
  const posts = [];

  for (const entry of fs.readdirSync(PUBLIC, { withFileTypes: true })) {
    if (!entry.isDirectory() || SKIP_DIRS.has(entry.name)) {
      continue;
    }

    const mdPath = path.join(PUBLIC, entry.name, "index.md");
    if (!fs.existsSync(mdPath)) {
      continue;
    }

    const raw = fs.readFileSync(mdPath, "utf8");
    const { data, content } = matter(raw);

    if (!data.title) {
      throw new Error(`${entry.name}/index.md is missing a title`);
    }
    if (!data.date) {
      throw new Error(`${entry.name}/index.md is missing a date`);
    }

    posts.push({
      slug: entry.name,
      title: String(data.title),
      date: toDate(data.date),
      content,
    });
  }

  posts.sort((a, b) => b.date - a.date);
  return posts;
}

function writePost(post) {
  const body = marked.parse(post.content);
  const html = renderPage(
    post.title,
    `<article>
<h1>${escapeHtml(post.title)}</h1>
<time datetime="${formatDate(post.date)}">${formatDate(post.date)}</time>
${body}
</article>
`,
  );
  return writeFile(path.join(PUBLIC, post.slug, "index.html"), html);
}

function writeListing(posts) {
  const items = posts
    .map(
      (post) =>
        `  <li><time datetime="${formatDate(post.date)}">${formatDate(post.date)}</time> <a href="${post.slug}/">${escapeHtml(post.title)}</a></li>`,
    )
    .join("\n");

  const html = renderPage(
    SITE_HEADER,
    `<h1>${escapeHtml(LISTING_HEADING)}</h1>
<ul class="post-list">
${items}
</ul>
`,
  );
  return writeFile(path.join(PUBLIC, "index.html"), html);
}

function writeFile(path, content) {
  if (fs.existsSync(path) && fs.readFileSync(path, "utf8") === content) {
    return false;
  }
  fs.writeFileSync(path, content);
  return true;
}

const posts = findPosts();
let skipped = 0;
let generated = 0;
for (const post of posts) {
  const wrote = writePost(post);
  if (!wrote) {
    skipped++;
    continue;
  }
  generated++;
}
console.log(`Generated: ${generated} post${generated === 1 ? "" : "s"}`);
console.log(`Skipped generation: ${skipped} post${skipped === 1 ? "" : "s"}`);
if (generated > 0 || !fs.existsSync(path.join(PUBLIC, "index.html"))) {
  const wrote = writeListing(posts);
  if (wrote) {
    console.log("Wrote index listing");
  } else {
    console.log("Skipping index listing write, already up to date");
  }
} else {
  console.log("Skipping index listing write, no posts were generated");
}
console.log(`Total: ${posts.length} post${posts.length === 1 ? "" : "s"}`);
