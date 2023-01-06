package t2s

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"errors"
	wav "github.com/dghwood/goaudio/wav"
	"strings"
)

func buildWAVRequest(req Text2SpeechRequest) (pb texttospeechpb.SynthesizeSpeechRequest) {
	// Populate defaults
	if req.LanguageCode == "" {
		req.LanguageCode = "en-US"
	}

	pb = texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: req.Text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: req.LanguageCode,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16, // WAV Format
		},
	}
	// Set Voice if there is one
	if req.Voice != "" {
		pb.Voice.Name = req.Voice
	}
	return
}

/*
	 Given the Cloud T2S API has a hard limit of 5000 characters
		 This function splits the text into batches of below 5000 characters
		 TODO: Does "," make sense as a splitter character?
*/
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
		return nil, errors.New("text unable to be split into requests")
	}
	returnStrings = append(returnStrings, buffer)
	return returnStrings, nil
}

func GCloudT2S(req Text2SpeechRequest) (resp Text2SpeechResponse, err error) {
	resp = Text2SpeechResponse{}
	text := req.Text
	if len(text) > 5*5000 {
		return resp, errors.New("article longer than 15000 chars")
	}

	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return resp, err
	}

	defer client.Close()

	texts, err := splitText(text, 5000-1)
	if err != nil {
		return resp, err
	}

	files := make([][]byte, len(texts))

	// Run each text piece via T2S
	// TODO: Do this concurrently?
	for i := 0; i < len(texts); i++ {
		t2sreq := req
		t2sreq.Text = texts[i]
		req := buildWAVRequest(t2sreq)
		r, err := client.SynthesizeSpeech(ctx, &req)
		if err != nil {
			return resp, err
		}
		files[i] = r.AudioContent
	}

	// Then concat the files into one
	wavFile, err := wav.FromBytes(files[0])
	if err != nil {
		return resp, err
	}
	for i := 1; i < len(texts); i++ {
		wavFile.AppendBytes(files[i])
	}
	resp.FileExtension = "wav"

	audioContent, _ := wavFile.Bytes()
	resp.AudioBytes = audioContent
	duration := wavFile.Seconds()
	resp.Duration = duration
	return resp, nil
}
