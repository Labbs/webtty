package utils

import (
	"log"
	"net/http"
	"os"

	"bytes"
)

var (
	RecordingUrl     string
	RecordingEnabled bool
)

func PushRecording(data string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}

	recording := []byte(data)
	var url string = "https://" + RecordingUrl + "/recording/" + hostname

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(recording))
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
