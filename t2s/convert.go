package t2s

import (
	"context"
	"errors"

	wav "github.com/moutend/go-wav"
	"io"
	"log"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

func buildRequest(text string) texttospeechpb.SynthesizeSpeechRequest {
	return texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}
}

func buildWAVRequest(text string) texttospeechpb.SynthesizeSpeechRequest {
	return texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
		},
	}
}

func Text2Speech(text string) ([]byte, string, bool) {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	req := buildRequest(text)
	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	return resp.AudioContent, ".mp3", false
}

/* Split text by , since that makes sense right? */
const splitChar = ","

func splitText(text string, limit int) ([]string, error) {
	lines := strings.Split(text, splitChar)
	returnStrings := make([]string, 0)
	buffer := ""
	for _, line := range lines {
		line += splitChar
		if len(line) >= limit {
			return nil, errors.New("text unable to be split into requests")
		} else if len(buffer)+len(line) < limit {
			buffer += line
		} else {
			returnStrings = append(returnStrings, buffer)
			buffer = line
		}
	}
	if len(buffer) > limit {
		log.Println("text unable to be split into requests")
		return nil, errors.New("text unable to be split into requests")
	}
	returnStrings = append(returnStrings, buffer)
	return returnStrings, nil
}

type Text2SpeechResponse struct {
	Duration      float32
	FileExtension string
	AudioBytes    []byte
}

func Text2SpeechLong(text string) (resp Text2SpeechResponse, err error) {
	resp = Text2SpeechResponse{}

	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return resp, err
	}

	defer client.Close()

	texts, _ := splitText(text, 5000-1)
	log.Println("TextSplit #", len(texts))

	files := make([]wav.File, len(texts))

	// Run each text piece via T2S
	// do this concurrently
	for i := 0; i < len(texts); i++ {
		log.Println("T2S for text:", len(texts[i]))
		req := buildWAVRequest(texts[i])
		r, err := client.SynthesizeSpeech(ctx, &req)
		if err != nil {
			return resp, err
		}
		wav.Unmarshal(r.AudioContent, &files[i])
	}
	// Then concat the files into one
	content, _ := wav.New(files[0].SamplesPerSec(), files[0].BitsPerSample(), files[0].Channels())
	for i := 0; i < len(texts); i++ {
		io.Copy(content, &files[i])
	}
	resp.FileExtension = "wav"

	audioContent, _ := wav.Marshal(content)
	resp.AudioBytes = audioContent
	// built in content.Duration is broken
	duration := float32(content.Length()) * 1. / float32(content.BlockAlign()) * 1. / float32(content.SamplesPerSec())
	resp.Duration = duration
	return resp, nil
}
