/* Basic HTML parsing functionality to find html nodes */
package html

import (
	"golang.org/x/net/html"
)

func findNodes(doc *html.Node, nodeSelector func(*html.Node) bool) (nodes []*html.Node) {
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if nodeSelector(c) {
			nodes = append(nodes, c)
		}
		n := findNodes(c, nodeSelector)
		if n != nil {
			nodes = append(nodes, n...)
		}
	}
	return nodes
}

func findNode(doc *html.Node, nodeSelector func(*html.Node) bool) (node *html.Node, found bool) {
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if nodeSelector(c) {
			return c, true
		}
		n, found := findNode(c, nodeSelector)
		if found {
			return n, true
		}

	}
	return node, false
}

func getAttribute(n *html.Node, name string) (attr html.Attribute, found bool) {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr, true
		}
	}
	return attr, false
}
