package api

import (
	b64 "encoding/base64"

	parse "github.com/dghwood/webpod/parse"
	storage "github.com/dghwood/webpod/storage"
	t2s "github.com/dghwood/webpod/t2s"
)

type URL2PodRequest struct {
	URL string `json:"url"`
}

type URL2PodResponse struct {
	ArticleURL string        `json:"article_url"`
	Article    parse.Article `json:"article"`
	AudioURL   string        `json:"audio_url"`
}

func URL2Pod(request URL2PodRequest) (URL2PodResponse, bool) {
	resp := URL2PodResponse{ArticleURL: request.URL}

	article, _ := parse.ParseArticle(request.URL)
	resp.Article = article

	audio, mimeType, _ := t2s.Text2Speech(article.Text)
	// TODO: Need to sort out mimeType vs fileExtension
	// Seems like I need .mp3 otherwise Chrome won't play it.
	fileName := b64.StdEncoding.EncodeToString([]byte(request.URL)) + ".mp3"
	audioURL := storage.Store(audio, fileName, mimeType)
	resp.AudioURL = audioURL

	return resp, false
}
