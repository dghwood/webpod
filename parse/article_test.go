package parse

import (
	//"fmt"
	"testing"
)

func TestArticle(t *testing.T) {
	//url := "https://www.nytimes.com/2022/12/16/us/politics/justice-eastman-trump-lawyers-fake-electors.html"
	url := "https://www.economist.com/leaders/2022/12/15/the-french-exception"
	//url = "https://www.nytimes.com/2022/12/16/us/politics/justice-eastman-trump-lawyers-fake-electors.amp.html"
	article, err := ParseArticle(url)
	if err {
		t.Error("failed to parse URL")
	}
	if article.Title != "The French exception" {
		t.Error("title failed to parse")
	}
	if article.Text[0:18] != "As the world turns" {
		t.Error("text failed to parse")
	}
	if article.ImageURL != "https://www.economist.com/img/b/1280/720/90/media-assets/image/20221217_LDP503.jpg" {
		t.Error("imageURL failed to parse")
	}
}
