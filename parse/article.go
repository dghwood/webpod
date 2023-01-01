package parse

import (
	b64 "encoding/base64"
	readability "github.com/go-shiori/go-readability"
	"io"
	"log"
	"net/http"
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

func urlToDataURL(urlString string) (dataURL string, err error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return dataURL, err
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dataURL, err
	}
	base64Data := b64.StdEncoding.EncodeToString(body)

	dataURL = "data:" + contentType + ";base64," + base64Data
	return dataURL, nil
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

	iconURL, err := getFaviconURL(doc, url)
	if err == nil {
		article.Favicon = iconURL
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

	cURL, err := getCanonicalURL(doc)
	if err == nil {
		article.URL = cURL
	}

	/* download the article images */
	faviconDURL, err := urlToDataURL(article.Favicon)
	if err == nil {
		article.Favicon = faviconDURL
	}
	imageDURL, err := urlToDataURL(article.ImageURL)
	if err == nil {
		article.ImageURL = imageDURL
	}

	return article, err
}
