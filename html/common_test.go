package html

import (
	"golang.org/x/net/html"
	nurl "net/url"
	"strings"
	"testing"
)

func TestFavicon(t *testing.T) {
	h := `
	<div>
		<head>
			<link rel="icon" href="/favicon-48x48.cbbd161b.png"/>
		</head>
		<div>
			<span>hello world</span>
		</div>
	</div>`
	doc, err := html.Parse(strings.NewReader(h))

	url := "https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch"
	//doc, err := fetchURL(url)
	if err != nil {
		t.Error(err)
	}
	u, err := nurl.ParseRequestURI(url)
	if err != nil {
		t.Error(err)
	}
	furl, err := GetFaviconURL(doc, u)
	if err != nil {
		t.Error(err)
	}
	if furl != "https://developer.mozilla.org/favicon-48x48.cbbd161b.png" {
		t.Error("furl is wrong")
	}
}

func TestCanonical(t *testing.T) {
	h := `
	<div>
		<head>
			<link rel="canonical" href="http://www.google.com"/>
		</head>
		<div>
			<span>hello world</span>
		</div>
	</div>`
	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		t.Error(err)
	}
	url, err := GetCanonicalURL(doc)
	if err != nil {
		t.Error(err)
	}
	if url != "http://www.google.com" {
		t.Error("url is wrong", url)
	}
}
