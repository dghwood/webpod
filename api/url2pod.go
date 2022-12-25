package api

import (
	b64 "encoding/base64"
	parse "github.com/dghwood/webpod/parse"
	storage "github.com/dghwood/webpod/storage"
	t2s "github.com/dghwood/webpod/t2s"
	"time"
)

type URL2PodRequest struct {
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

type URL2PodResponse struct {
	Timestamp  time.Time     `json:"timestamp"`
	ArticleURL string        `json:"article_url"`
	Article    parse.Article `json:"article"`
	AudioURL   string        `json:"audio_url"`
	Duration   float32       `json:"duration"`
	Error      bool          `json:"error"`
}

func URL2Pod(request URL2PodRequest) (URL2PodResponse, bool) {

	resp := URL2PodResponse{ArticleURL: request.URL}
	resp.Timestamp = request.Timestamp
	article, err := parse.ParseArticle(request.URL)
	if err {
		resp.Error = true
		return resp, true
	}
	resp.Article = article
	audioResp, err := t2s.Text2SpeechLong(article.Text)
	if err {
		resp.Error = true
		return resp, true
	}
	resp.Duration = audioResp.Duration
	fileName := b64.StdEncoding.EncodeToString([]byte(request.URL)) + "." + audioResp.FileExtension
	audioURL := storage.Store(audioResp.AudioBytes, fileName)
	resp.AudioURL = audioURL

	return resp, false
}
