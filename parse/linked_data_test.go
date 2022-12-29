package parse

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestFetchURL(t *testing.T) {
	url := "https://www.theguardian.com/commentisfree/2022/dec/21/politicians-candidates-tax-returns-mandatory-congress-stocks"
	_, err := fetchURL(url)
	if err != nil {
		t.Error(err)
	}
}

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

func TestLinkedData(t *testing.T) {
	doc, err := fetchURL("https://www.economist.com/leaders/2022/12/15/the-french-exception")
	if err != nil {
		t.Error(err)
	}
	n, found := findNode(doc, linkedDataCond)
	if !found {
		t.Error(err)
	}
	if n == nil {
		t.Error("node is nil")
	}
	//t.Error(GetLinkedData(doc))
}
