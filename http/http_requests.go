package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// NewHTTPClient -
func NewHTTPClient() *http.Client {
	timeout := time.Duration(30 * time.Second)
	return &http.Client{
		Timeout: timeout,
	}
}

var httpClient = NewHTTPClient()

func makeRequest(method, endpoint string, headers map[string]interface{}, body []byte) (*http.Request, error) {
	request, err := http.NewRequest(method, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, value := range headers {
			if value != nil {
				request.Header.Add(key, value.(string))
			}
		}
	}

	request.Header.Add("Content-Type", "application/json")
	return request, nil
}

// MakePutRequest -
func MakePutRequest(endpoint string, body []byte) (int, error) {
	request, err := makeRequest("PUT", endpoint, nil, body)
	if err != nil {
		log.Printf("Error Creating Request: %#v\n", err)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Printf("Got An Error: %#v\n", err.Error())
	}

	if IsSuccessful(resp) {
		return 200, nil
	}
	switch resp.StatusCode {
	case 404:
		return 404, nil
	case 401:
		return 401, fmt.Errorf("unauthorized error")
	default:
		return resp.StatusCode, fmt.Errorf(resp.Status)
	}
}

// MakePostRequest -
func MakePostRequest(endpoint string, body []byte) (int, error) {
	request, err := makeRequest("POST", endpoint, nil, body)
	if err != nil {
		fmt.Printf("Error Creating Request: %#v\n", err)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Printf("Got An Error: %#v", err.Error())
	}

	if IsSuccessful(resp) {
		return 200, nil
	}
	switch resp.StatusCode {
	case 404:
		return 404, nil
	case 401:
		return 401, fmt.Errorf("unauthorized error")
	default:
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		resMessage := string(data)
		return resp.StatusCode, fmt.Errorf("%s => %s", resp.Status, resMessage)
	}
}

// MakePostForResponse -
func MakePostForResponse(endpoint string, body []byte) (*http.Response, error) {
	request, err := makeRequest("POST", endpoint, nil, body)
	if err != nil {
		fmt.Printf("Error Creating Request: %#v\n", err)
	}
	return httpClient.Do(request)
}

// MakePostWithHeadersForResponse -
func MakePostWithHeadersForResponse(endpoint string, reqHeaders map[string]interface{}, body []byte) (*http.Response, error) {
	request, err := makeRequest("POST", endpoint, reqHeaders, body)
	if err != nil {
		fmt.Printf("Error Creating Request: %#v\n", err)
	}
	return httpClient.Do(request)
}

// IsSuccessful - returns if an http response was successful.
func IsSuccessful(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}

// MakePostRequestWithResponse -
func MakePostRequestWithResponse(endpoint string, body []byte, response interface{}) (int, error) {
	request, err := makeRequest("POST", endpoint, nil, body)
	if err != nil {
		fmt.Printf("Error Creating Request: %#v\n", err)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Printf("Got An Error: %#v", err.Error())
	}

	if IsSuccessful(resp) {
		defer resp.Body.Close()
		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 200, nil
		}
		if rb != nil {
			err = json.Unmarshal(rb, &response)
			if err != nil {
				return 200, err
			}
		}
		return 200, nil
	}
	switch resp.StatusCode {
	case 404:
		return 404, nil
	case 401:
		return 401, fmt.Errorf("unauthorized error")
	default:
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		resMessage := string(data)
		return resp.StatusCode, fmt.Errorf("%d => %s", resp.StatusCode, resMessage)
	}
}

// MakeGetRequest -
func MakeGetRequest(endpoint string) (string, int, error) {
	request, err := makeRequest("GET", endpoint, nil, nil)
	if err != nil {
		fmt.Printf("Error Creating Request: %v\n", err)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Printf("Error Making GetRequest: %v\n", err.Error())
	}

	if IsSuccessful(resp) {
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("unable to read body of content")
		}

		return string(content), 200, nil
	}
	switch resp.StatusCode {
	case 404:
		return "", 404, nil
	case 401:
		return "", 401, fmt.Errorf("unauthorized error")
	default:
		return "", resp.StatusCode, fmt.Errorf("unknown error")
	}
}
