package parse

import (
	readability "github.com/go-shiori/go-readability"
	"log"
	"time"
)

type Article struct {
	Title    string
	Text     string
	ImageURL string
}

func ParseArticle(url string) (Article, bool) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
		return Article{}, false
	}
	return Article{
		Title:    article.Title,
		Text:     article.TextContent,
		ImageURL: article.Image,
	}, false
}
