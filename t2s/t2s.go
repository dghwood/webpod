package t2s

import ()

func Text2Speech(req Text2SpeechRequest) (resp Text2SpeechResponse, err error) {
	resp, err = GCloudT2S(req)
	return
}
