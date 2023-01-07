package parse

import (
	"testing"
)

func TestArticle(t *testing.T) {
	url := "https://www.economist.com/leaders/2022/12/15/the-french-exception"
	article, err := ParseArticle(url)
	if err != nil {
		t.Error("failed to parse URL")
	}
	if article.Title != "The French exception" {
		t.Error("title failed to parse")
	}
	if article.Text[0:18] != "As the world turns" {
		t.Error("text failed to parse")
	}
}

func TestNonURL(t *testing.T) {
	url := "asdasda"
	article, err := ParseArticle(url)
	if err == nil {
		t.Error("Should error out", article)
	}
}

func TestURLToB64(t *testing.T) {
	_, err := urlToDataURL("https://www.economist.com/engassets/ico/touch-icon-180x180.f1ea908894.png")
	if err != nil {
		t.Error(err)
	}
	//t.Error(dataURL)
}
