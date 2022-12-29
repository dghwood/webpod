package t2s

import (
	"os"
	"testing"
)

func TestText2SpeechLong(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/key.json")
	buffer := ""
	for i := 0; i < 100; i++ {
		buffer += "hello hello hello hello hello hello hello hello hello" + splitChar
	}
	if len(buffer) < 5000 {
		t.Error("buffer is too short to split", len(buffer))
	}
	resp, err := Text2SpeechLong(buffer)
	if err != nil {
		t.Error("Text2Speech throws error")
	}

	// download the file for testing
	os.WriteFile("testtext2speechlong.wav", resp.AudioBytes, 0644)
}

func TestSplitText(t *testing.T) {
	text := `one
two
three`
	result, err := splitText(text, 9)
	if err != nil {
		t.Error("splitText threw an error")
	}
	if len(result) != 2 {
		t.Error("splitText returned the wrong number of splits")
	}
}

func TestSplitTextLong(t *testing.T) {
	buffer := ""
	for i := 0; i < 100; i++ {
		buffer += "hello hello hello hello hello hello hello hello hello" + splitChar
	}
	result, err := splitText(buffer, 5000)
	if err != nil {
		t.Error(err)
	}
	if len(result) != 2 {
		t.Error("splitText returned wrong number of splits")
	}
}

func TestSplitTextNoNewLines(t *testing.T) {
	text, err := os.ReadFile("../samples/test_article_no_new_lines.txt")
	if err != nil {
		t.Error("failed to read file")
	}
	result, err := splitText(string(text), 5000-1)
	if err != nil {
		t.Error("failed to split", result)
	}
}
