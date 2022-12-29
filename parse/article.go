package parse

import (
	readability "github.com/go-shiori/go-readability"
	"log"
	nurl "net/url"
)

type Article struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	ImageURL string `json:"image_url"`
	SiteName string `json:"site_name"`
	Favicon  string `json:"favicon"`
}

func ParseArticle(urlString string) (article Article, err error) {
	log.Println("parse article")
	url, err := nurl.ParseRequestURI(urlString)
	if err != nil {
		log.Println("not a URL")
		return article, err
	}
	log.Println(url)
	doc, err := fetchURL(urlString)
	if err != nil {
		return article, err
	}
	rArticle, err := readability.FromDocument(doc, url)
	if err != nil {
		return article, err
	}
	article = Article{
		URL:      urlString,
		Title:    rArticle.Title,
		Text:     rArticle.TextContent,
		ImageURL: rArticle.Image,
		Favicon:  rArticle.Favicon,
		SiteName: rArticle.SiteName,
	}

	linkedData := GetLinkedData(doc)
	for _, ld := range linkedData {
		if ld.Headline != "" {
			article.Title = ld.Headline
		}
		if len(ld.Image) > 0 {
			article.ImageURL = ld.Image[0]
		}
		if ld.Publisher.Logo.URL != "" {
			article.Favicon = ld.Publisher.Logo.URL
		}
		if ld.Publisher.Name != "" {
			article.SiteName = ld.Publisher.Name
		}
	}
	iconURL, err := getFaviconURL(doc)
	if err == nil {
		article.Favicon = iconURL
	}
	cURL, err := getCanonicalURL(doc)
	if err == nil {
		article.URL = cURL
	}
	return article, err
}
