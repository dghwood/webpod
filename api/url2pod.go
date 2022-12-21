package api

import (
	b64 "encoding/base64"
	parse "github.com/dghwood/webpod/parse"
	storage "github.com/dghwood/webpod/storage"
	t2s "github.com/dghwood/webpod/t2s"
	"log"
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
	log.Println("HELLO", len(article.Text))
	audio, fileExtension, _ := t2s.Text2SpeechLong(article.Text)
	fileName := b64.StdEncoding.EncodeToString([]byte(request.URL)) + fileExtension
	audioURL := storage.Store(audio, fileName)
	resp.AudioURL = audioURL

	return resp, false
}
