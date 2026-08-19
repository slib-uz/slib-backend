package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type CHTTpClient struct {
	*http.Client
}

type BasicAuth struct {
	Username string
	Password string
}

func NewBasicAuth(username string, password string) *BasicAuth {
	return &BasicAuth{Username: username, Password: password}
}

// @inject
func NewCHTTpClient(client *http.Client) *CHTTpClient {
	return &CHTTpClient{Client: client}
}

func (this *CHTTpClient) super() *http.Client {
	return this.Client

}

func (this *CHTTpClient) Get(url string, auth *BasicAuth, header map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	resp, err := this.super().Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (this *CHTTpClient) Post(url string, body any, auth *BasicAuth, header map[string]string) (*http.Response, error) {

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}

	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	if header == nil || header["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := this.super().Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (this *CHTTpClient) Put(url string, body any, auth *BasicAuth, header map[string]string) (*http.Response, error) {

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}

	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := this.super().Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (this *CHTTpClient) PostFormData(url string, file []byte, fileName string, payload map[string]any, auth *BasicAuth, header map[string]string) (*http.Response, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}
	if _, err := part.Write(file); err != nil {
		return nil, fmt.Errorf("write file content error: %v", err)
	}

	for key, value := range payload {
		_ = writer.WriteField(key, fmt.Sprintf("%v", value))
	}

	// Multipartni yopish
	_ = writer.Close()

	// HTTP so‘rov
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, fmt.Errorf("request create error: %v", err)
	}

	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := this.super().Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request error: %v", err)
	}
	return resp, nil
}

func Closer(c io.Closer) func() {
	return func() {
		_ = c.Close()
	}
}
