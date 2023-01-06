package api

import (
	"os"
	"testing"
)

func TestURL2Pod(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/webpods_key.json")
	_, err := URL2Pod(URL2PodRequest{
		URL: "https://www.economist.com/the-americas/2022/12/31/brazils-new-president-faces-a-fiscal-crunch-and-a-fickle-congress",
	})
	if err != nil {
		t.Error("Thrown error", err)
	}
}

func TestURL2Pod2(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/dgh_wood/webpods_key.json")
	_, err := URL2Pod(URL2PodRequest{
		URL: "https://www.adexchanger.com/commerce/more-performance-less-transparency-inside-metas-advantage-shopping-black-box/",
	})
	if err != nil {
		t.Error("Thrown error", err)
	}
}
