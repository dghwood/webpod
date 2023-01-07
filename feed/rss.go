package feed

import (
	"encoding/xml"
	"github.com/dghwood/webpod/fetch"
	"io"
	"time"
)

type RSS struct {
	Channel []RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title string    `xml:"title"`
	Item  []RSSItem `xml:"item"`
	Link  string    `xml:"link"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func fetchRSS(urlString string) (rss RSS, err error) {
	resp, _, err := fetch.Fetch(urlString)
	if err != nil {
		return rss, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = xml.Unmarshal(respBytes, &rss)
	return
}

func fetchRSSFeed(url string) (feed Feed, err error) {
	rss, err := fetchRSS(url)
	if err != nil {
		return
	}
	/* Parse RSS to Feed */
	feed.Title = rss.Channel[0].Title
	feed.Url = rss.Channel[0].Link
	feed.Items = make([]FeedItem, len(rss.Channel[0].Item))
	for i, item := range rss.Channel[0].Item {
		feed.Items[i].Description = item.Description
		feed.Items[i].Title = item.Title
		feed.Items[i].Url = item.Link
		date, err := time.Parse(time.RFC1123, item.PubDate)
		if err == nil {
			feed.Items[i].Date = date
		}
	}
	return
}
