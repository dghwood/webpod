package t2s

type Text2SpeechResponse struct {
	Duration      float64
	FileExtension string
	AudioBytes    []byte
}

type Text2SpeechRequest struct {
	Text         string
	Model        string
	LanguageCode string
	Voice        string
}
