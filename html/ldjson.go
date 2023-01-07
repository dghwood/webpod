package html

import (
	"encoding/json"
	"golang.org/x/net/html"
	"time"
)

func linkedDataCond(n *html.Node) (found bool) {
	a, found := getAttribute(n, "type")
	if !found {
		return false
	}
	return a.Val == "application/ld+json"
}

type LinkedData struct {
	Type          string      `json:"@type"`
	Context       string      `json:"@context"`
	Id            string      `json:"@id"`
	Image         []string    `json:"image"`
	Headline      string      `json:"headline"`
	DatePublished time.Time   `json:"datePublished"`
	Publisher     LDPublisher `json:"publisher"`
}
type LDPublisher struct {
	Name string `json:"name"`
	Logo LDLogo `json:"logo"`
}
type LDLogo struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func GetLinkedData(doc *html.Node) (nodes []LinkedData) {
	n, found := findNode(doc, linkedDataCond)
	data := make([]LinkedData, 0)
	if !found {
		return data
	}

	err := json.Unmarshal([]byte(n.FirstChild.Data), &data)
	if err != nil {
		// Try single entry
		singleLD := LinkedData{}
		err = json.Unmarshal([]byte(n.FirstChild.Data), &singleLD)
		if err != nil {
			return data
		}
		data = append(data, singleLD)
	}
	return data
}
