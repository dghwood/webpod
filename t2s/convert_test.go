package t2s

import (
	"testing"
)

func TestText2Speech(t *testing.T) {
	_, err := Text2Speech("hello")
	if err {
		t.Error("Text2Speech throws error")
	}
}
