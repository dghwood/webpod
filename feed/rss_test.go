package feed

import (
	"testing"
)

func TestRSS(t *testing.T) {
	_, err := fetchRSSFeed("https://www.theguardian.com/rss")
	if err != nil {
		t.Error(err)
	}
	//t.Error(rss)
}

func TestRSS2(t *testing.T) {
	_, err := fetchRSSFeed("https://www.adexchanger.com/rss")
	if err != nil {
		t.Error(err)
	}
	//t.Error(rss)
}

func TestFeed1(t *testing.T) {
	_, err := fetchRSSFeed("https://www.theguardian.com/rss")
	if err != nil {
		t.Error(err)
	}
	//t.Error(rss)
}

func TestFeed2(t *testing.T) {
	_, err := fetchRSSFeed("https://www.adexchanger.com/rss")
	if err != nil {
		t.Error(err)
	}
	//t.Error(rss)
}
