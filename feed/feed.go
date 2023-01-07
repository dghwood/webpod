/*
	 Fetch Feeds of Articles

		* Currently only supports RSS
*/
package feed

import (
	"time"
)

type Feed struct {
	Items []FeedItem
	Url   string
	Title string
}

type FeedItem struct {
	Title       string
	Url         string
	Description string
	Date        time.Time
}

func FetchFeed(url string) (feed Feed, err error) {
	feed, err = fetchRSSFeed(url)
	return
}
