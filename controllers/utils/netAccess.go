package utils

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HttpRequester interface {
	//Verify that the provided URL returns http status code 200
	VerifyHealthyResponse(url string) bool
}

type HttpRequesterImpl struct{}

func NewHttpRequester() HttpRequester {
	return &HttpRequesterImpl{}
}

func (hri *HttpRequesterImpl) VerifyHealthyResponse(url string) bool {
	success := false
	tr := &http.Transport{
		// skipping TLS verification for now in order to allow for dev clusters with self-signed certs
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 1,
		Transport: tr,
	}
	resp, err := netClient.Get(url + "health")
	if err == nil {
		// log --> fmt.Printf(url + "health = %d\n", resp.StatusCode)
		if resp.StatusCode == 200 {
			success = true
		}
	}
	return success
}
