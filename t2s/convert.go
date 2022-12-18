package t2s

import (
	"context"
	//"fmt"
	//"io/ioutil"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

func Text2Speech(text string) ([]byte, bool) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/key.json")
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
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

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	// The resp's AudioContent is binary.
	//filename := "output.mp3"
	//err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Audio content written to file: %v\n", filename)
	return resp.AudioContent, false
}
