package utils

import (
	"log"
	"net/http"
	"os"

	"bytes"
)

var (
	RecordingUrl string
)

func PushRecording(data string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}

	recording := []byte(`{"spendy_host":"` + hostname + `","data":"` + data + `"}`)

	r, err := http.NewRequest("POST", RecordingUrl, bytes.NewBuffer(recording))
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
