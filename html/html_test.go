package html

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestFindNode(t *testing.T) {
	h := "<div><div></div><div><span>hello</span></div></div>"
	doc, _ := html.Parse(strings.NewReader(h))
	f := func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "span"
	}
	n, found := findNode(doc, f)
	if !found {
		t.Error("nil")
	}
	if n.FirstChild.Data != "hello" {
		t.Error("wrong element", n.FirstChild.Data)
	}
}

func TestFindNodes(t *testing.T) {
	h := "<div><div></div><div><span>hello</span></div></div>"
	doc, _ := html.Parse(strings.NewReader(h))
	isMyNode := func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "div"
	}
	n := findNodes(doc, isMyNode)
	if n == nil {
		t.Error("nil")
	}
	if len(n) != 3 {
		t.Error("didn't find 3 elements")
	}
}
