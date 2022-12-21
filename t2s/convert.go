package t2s

import (
	"context"
	"errors"
	"fmt"
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

func splitText(text string, limit int) ([]string, error) {
	lines := strings.Split(text, "\n")
	returnStrings := make([]string, 0)
	buffer := ""
	for _, line := range lines {
		line += "\n"
		if len(buffer)+len(line) < limit {
			buffer += line
		} else {
			returnStrings = append(returnStrings, buffer)
			buffer = line
		}
	}
	if len(buffer) > limit {
		return nil, errors.New("text unable to be split into requests")
	}
	returnStrings = append(returnStrings, buffer)
	return returnStrings, nil
}

func Text2SpeechLong(text string) ([]byte, string, bool) {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	texts, _ := splitText(text, 5000-1)
	fmt.Println(len(texts))

	files := make([]wav.File, len(texts))

	// Run each text piece via T2S
	for i := 0; i < len(texts); i++ {
		req := buildWAVRequest(texts[i])
		resp, err := client.SynthesizeSpeech(ctx, &req)
		if err != nil {
			log.Fatal(err)
		}
		wav.Unmarshal(resp.AudioContent, &files[i])
	}
	// Then concat the files into one
	content, _ := wav.New(files[0].SamplesPerSec(), files[0].BitsPerSample(), files[0].Channels())
	for i := 0; i < len(texts); i++ {
		io.Copy(content, &files[i])
	}
	audioContent, _ := wav.Marshal(content)

	return audioContent, ".wav", false
}
