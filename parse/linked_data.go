package parse

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	nurl "net/url"
	"time"

	"golang.org/x/net/html"
)

func fetchURL(url string) (doc *html.Node, err error) {
	var resp *http.Response
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Get(url)
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
	doc, err = html.Parse(resp.Body)
	return doc, err
}

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

func linkedDataCond(n *html.Node) (found bool) {
	a, found := getAttribute(n, "type")
	if !found {
		return false
	}
	return a.Val == "application/ld+json"
}

type LinkedData struct {
	Type          string      `json:"@type"`
	Context       string      `json:"@context"`
	Id            string      `json:"@id"`
	Image         []string    `json:"image"`
	Headline      string      `json:"headline"`
	DatePublished time.Time   `json:"datePublished"`
	Publisher     LDPublisher `json:"publisher"`
}
type LDPublisher struct {
	Name string `json:"name"`
	Logo LDLogo `json:"logo"`
}
type LDLogo struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func GetLinkedData(doc *html.Node) (nodes []LinkedData) {
	n, found := findNode(doc, linkedDataCond)
	data := make([]LinkedData, 0)
	if !found {
		return data
	}

	err := json.Unmarshal([]byte(n.FirstChild.Data), &data)
	if err != nil {
		// Try single entry
		singleLD := LinkedData{}
		err = json.Unmarshal([]byte(n.FirstChild.Data), &singleLD)
		if err != nil {
			return data
		}
		data = append(data, singleLD)
	}
	return data
}

func canonicalCond(n *html.Node) (found bool) {
	if n.Type == html.ElementNode && n.Data == "link" {
		a, found := getAttribute(n, "rel")
		if !found {
			return false
		}
		return a.Val == "canonical"
	}
	return false
}

func getCanonicalURL(doc *html.Node) (url string, err error) {
	n, found := findNode(doc, canonicalCond)
	if !found {
		return url, errors.New("not found")
	}
	var a html.Attribute
	a, found = getAttribute(n, "href")
	if !found {
		return url, errors.New("not found")
	}
	return a.Val, nil
}

func faviconCond(n *html.Node) (found bool) {
	if n.Type != html.ElementNode || n.Data != "link" {
		return found
	}
	a, found := getAttribute(n, "rel")
	if !found {
		return found
	}
	return a.Val == "icon"
}
func getFaviconURL(doc *html.Node, baseURL *nurl.URL) (url string, err error) {
	n, found := findNode(doc, faviconCond)
	if !found {
		return url, errors.New("not found")
	}
	a, found := getAttribute(n, "href")
	if !found {
		return url, errors.New("not found")
	}

	furl, err := nurl.Parse(a.Val)
	if err != nil {
		return url, err
	}
	log.Println(furl, baseURL)
	url = baseURL.ResolveReference(furl).String()
	return url, nil
}
