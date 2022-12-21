package parse

import (
	readability "github.com/go-shiori/go-readability"
	//"log"
	"time"
)

type Article struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	ImageURL string `json:"image_url"`
	SiteName string `json:"site_name"`
	Favicon  string `json:"favicon"`
}

func ParseArticle(url string) (Article, bool) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		//log.Info("failed to parse %s, %v\n", url, err)
		return Article{}, true
	}
	return Article{
		Title:    article.Title,
		Text:     article.TextContent,
		ImageURL: article.Image,
		Favicon:  article.Favicon,
		SiteName: article.SiteName,
	}, false
}
