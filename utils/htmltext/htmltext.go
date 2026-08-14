package htmltext

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Text extracts plain text from HTML, skipping executable or styling content
// and inserting spaces at block boundaries so adjacent words do not merge.
func Text(raw string) string {
	nodes, err := html.ParseFragment(strings.NewReader(raw), &html.Node{Type: html.ElementNode, DataAtom: atom.Div, Data: "div"})
	if err != nil {
		return raw
	}
	var output strings.Builder
	for _, node := range nodes {
		collectText(&output, node)
	}
	return output.String()
}

// collectText appends text descendants while skipping script and style nodes.
func collectText(output *strings.Builder, node *html.Node) {
	if node.Type == html.TextNode {
		output.WriteString(node.Data)
		return
	}
	if node.Type == html.ElementNode && (node.Data == "script" || node.Data == "style") {
		return
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		collectText(output, child)
	}
	if node.Type == html.ElementNode && isTextBoundary(node.Data) {
		output.WriteByte(' ')
	}
}

// isTextBoundary reports whether a tag visually separates words in plain-text output.
func isTextBoundary(tag string) bool {
	switch strings.ToLower(tag) {
	case "blockquote", "br", "div", "h1", "h2", "h3", "h4", "h5", "h6", "hr", "li", "p", "pre", "table", "td", "th", "tr":
		return true
	default:
		return false
	}
}
