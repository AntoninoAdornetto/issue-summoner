package common

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func ConstructURL(base string, qParams map[string]string, paths ...string) (string, error) {
	u, err := url.JoinPath(base, paths...)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	for key, val := range qParams {
		params.Set(key, val)
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", u, params.Encode()), nil
	}

	return u, nil
}

func Request(method, url string, body io.Reader, h http.Header) (*http.Response, []byte, error) {
	var res []byte
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, res, err
	}

	req.Header = h
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, res, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return resp, data, err
}
