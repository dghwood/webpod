package storage

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestStorage(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/key.json")
	fileContents := []byte("HELLO")
	url, err := Store(fileContents, "test2.txt")
	if err != nil {
		t.Error(err)
	}
	resp, err := http.Get(url)
	if err != nil {
		t.Error("http failed", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error("body ready failed", err)
	}
	if string(body) != string(fileContents) {
		t.Error("file not correct")
	}
}
