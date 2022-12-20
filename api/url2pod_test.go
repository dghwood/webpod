package api

import (
	"os"
	"testing"
)

func TestURL2Pod(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/key.json")
	resp, err := URL2Pod(URL2PodRequest{
		URL: "https://www.economist.com/leaders/2022/12/15/the-french-exception",
	})
	if err {
		t.Error("Thrown error")
	}
	t.Error(resp)
}
