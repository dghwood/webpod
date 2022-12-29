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
	if err != nil {
		t.Error("Thrown error", err)
	}
	t.Error(resp)
}

func TestURL2Pod2(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/key.json")
	resp, err := URL2Pod(URL2PodRequest{
		URL: "https://www.adexchanger.com/commerce/more-performance-less-transparency-inside-metas-advantage-shopping-black-box/",
	})
	if err != nil {
		t.Error("Thrown error", err)
	}
	t.Error(resp)
}
