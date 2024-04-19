package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func SubmitPostRequest(url string, body io.Reader, h http.Header) ([]byte, error) {
	var res []byte
	req, _ := http.NewRequest("POST", url, body)
	req.Header = h

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	res, err = io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return res, err
}

func BuildURL(baseURL string, paths []string, qParams map[string]string) (string, error) {
	u, err := url.JoinPath(baseURL, paths...)
	if err != nil {
		return "", err
	}

	size := 0
	params := url.Values{}
	for key, val := range qParams {
		params.Set(key, val)
		size++
	}

	if size > 0 {
		return fmt.Sprintf("%s?%s", u, params.Encode()), nil
	}

	return u, nil
}
