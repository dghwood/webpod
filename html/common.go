/* Convience functions for finding common HTML Nodes */
package html

import (
	"errors"
	"golang.org/x/net/html"
	nurl "net/url"
)

/* Canonical URLs */
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

func GetCanonicalURL(doc *html.Node) (url string, err error) {
	n, found := findNode(doc, canonicalCond)
	if !found {
		return url, errors.New("not found")
	}
	var a html.Attribute
	a, found = getAttribute(n, "href")
	if !found {
		return url, errors.New("not found")
	}
	// Check if the URL is valid
	_, err = nurl.ParseRequestURI(a.Val)
	if err != nil {
		return url, err
	}
	return a.Val, nil
}

/* Favicon URLs */
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
func GetFaviconURL(doc *html.Node, baseURL *nurl.URL) (url string, err error) {
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
	url = baseURL.ResolveReference(furl).String()
	return url, nil
}
