package markdown

import (
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

func embedImageRenderer(next html.NodeRenderer) html.NodeRenderer {
	return html.NodeRendererFunc(func(w io.Writer, source []byte, n ast.Node, entering bool, rc renderer.Context) (ast.WalkStatus, error) {
		img, ok := n.(*ast.Image)
		if ok && entering && isVideoDest(img.Destination.Value(source)) {
			writeVideo(w, source, img, rc)
			return ast.WalkSkipChildren, nil
		}
		status, err := next.Render(w, source, n, entering, rc)
		if err != nil || !ok || !entering {
			return status, err
		}
		writePosterLabel(w, source, img, rc)
		return status, nil
	})
}

func writeVideo(w io.Writer, source []byte, n *ast.Image, rc renderer.Context) {
	bw := w.(util.BufWriter)
	_, _ = bw.WriteString(`<video controls preload="metadata" src="`)
	_, _ = html.ContextLinkURLWriter(rc).WriteString(n.Destination.Value(source))
	_ = bw.WriteByte('"')
	if alt := altText(n, source); alt != "" {
		_, _ = bw.WriteString(` aria-label="`)
		_, _ = html.ContextTextWriter(rc).WriteString(alt)
		_ = bw.WriteByte('"')
	}
	if !n.Title.IsEmpty() {
		_, _ = bw.WriteString(` title="`)
		_, _ = n.Title.WriteTo(html.ContextTextWriter(rc), source)
		_ = bw.WriteByte('"')
	}
	_, _ = bw.WriteString(`></video>`)
}

func writePosterLabel(w io.Writer, source []byte, img *ast.Image, rc renderer.Context) {
	link, ok := img.Parent().(*ast.Link)
	if !ok || !isVideoDest(link.Destination.Value(source)) {
		return
	}
	alt := altText(img, source)
	if alt == "" {
		return
	}
	bw := w.(util.BufWriter)
	_, _ = bw.WriteString(`<span class="video-label">`)
	_, _ = html.ContextTextWriter(rc).WriteString(alt)
	_, _ = bw.WriteString(`</span>`)
}

func markGalleries(doc ast.Node, source []byte) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		list, ok := n.(*ast.List)
		if !ok || !isGalleryList(list, source) {
			return ast.WalkContinue, nil
		}
		list.SetAttribute("class", text.NewMultiLineValueFromString("gallery", nil))
		return ast.WalkContinue, nil
	})
}

func isGalleryList(list *ast.List, source []byte) bool {
	items := significantChildren(list, source)
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if !isGalleryItem(item, source) {
			return false
		}
	}
	return true
}

func isGalleryItem(n ast.Node, source []byte) bool {
	if n.Kind() != ast.KindListItem {
		return false
	}
	kids := significantChildren(n, source)
	if len(kids) != 1 {
		return false
	}
	content := kids[0]
	if content.Kind() == ast.KindParagraph {
		kids = significantChildren(content, source)
		if len(kids) != 1 {
			return false
		}
		content = kids[0]
	}
	switch node := content.(type) {
	case *ast.Image:
		return true
	case *ast.Link:
		inner := significantChildren(node, source)
		if len(inner) != 1 {
			return false
		}
		if _, ok := inner[0].(*ast.Image); !ok {
			return false
		}
		return isVideoDest(node.Destination.Value(source))
	default:
		return false
	}
}

func significantChildren(n ast.Node, source []byte) []ast.Node {
	var out []ast.Node
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if isBlankText(c, source) {
			continue
		}
		out = append(out, c)
	}
	return out
}

func isBlankText(n ast.Node, source []byte) bool {
	t, ok := n.(*ast.Text)
	if !ok {
		return false
	}
	return strings.TrimSpace(t.Value.Value(source)) == ""
}

func isVideoDest(raw string) bool {
	raw = strings.TrimSpace(raw)
	if u, err := url.Parse(raw); err == nil && u.Path != "" {
		raw = u.Path
	}
	switch strings.ToLower(path.Ext(raw)) {
	case ".mp4", ".webm", ".ogv", ".mov":
		return true
	default:
		return false
	}
}

func altText(n ast.Node, source []byte) string {
	var b strings.Builder
	collectText(n, source, &b)
	return b.String()
}

func collectText(n ast.Node, source []byte, b *strings.Builder) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			b.WriteString(t.Value.Value(source))
			continue
		}
		collectText(c, source, b)
	}
}
