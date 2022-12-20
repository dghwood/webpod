package t2s

import (
	"os"
	"testing"
)

func TestText2Speech(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/key.json")
	_, _, err := Text2Speech("hello")
	if err {
		t.Error("Text2Speech throws error")
	}
}
