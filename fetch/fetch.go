/* Package to act as interface for requesting URLs */
package fetch

import (
	"net/http"
	//nurl "net/url"
	"golang.org/x/net/html"
	nurl "net/url"
	"time"
)

type Crawler struct {
	client    *http.Client
	UserAgent string
}

func (c *Crawler) Fetch(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	// TODO: Timeout?
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err = c.client.Do(req)
	return
}

/* @Deprecated */
func Fetch(urlString string) (resp *http.Response, url *nurl.URL, err error) {
	// Check the URL is valid
	url, err = nurl.ParseRequestURI(urlString)
	if err != nil {
		return
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Get(urlString)
	return
}

/* @Deprecated */
func FetchDoc(urlString string) (doc *html.Node, url *nurl.URL, err error) {
	resp, url, err := Fetch(urlString)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err = html.Parse(resp.Body)
	return
}
