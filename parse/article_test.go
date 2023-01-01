package parse

import (
	"os"
	"testing"
)

func TestArticle(t *testing.T) {
	//url := "https://www.nytimes.com/2022/12/16/us/politics/justice-eastman-trump-lawyers-fake-electors.html"
	url := "https://www.economist.com/leaders/2022/12/15/the-french-exception"
	//url = "https://www.nytimes.com/2022/12/16/us/politics/justice-eastman-trump-lawyers-fake-electors.amp.html"
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
	if article.ImageURL != "https://www.economist.com/img/b/1280/720/90/media-assets/image/20221217_LDP503.jpg" {
		t.Error("imageURL failed to parse")
	}
}

func TestArticleLong(t *testing.T) {
	url := "https://www.adexchanger.com/commerce/more-performance-less-transparency-inside-metas-advantage-shopping-black-box/"
	_, err := ParseArticle(url)
	if err != nil {
		t.Error("failed to parse URL ", url)
	}
}

func TestArticleNoNewLines(t *testing.T) {
	url := "https://www.theguardian.com/commentisfree/2022/dec/21/politicians-candidates-tax-returns-mandatory-congress-stocks"
	article, err := ParseArticle(url)
	if err != nil {
		t.Error("failed to parse URL ", url)
	}
	os.WriteFile("../samples/test_article_no_new_lines.txt", []byte(article.Text), 0644)
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
