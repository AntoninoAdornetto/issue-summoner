package scm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/utils"
)

const (
	BASE_URL      = "https://github.com"
	CLIENT_ID     = "Iv1.587a6b18c40684ba"
	GRANT_TYPE    = "urn:ietf:params:oauth:grant-type:device_code"
	CONFIG_PATH   = "~/.config/issue-summoner/scm.json"
	ACCESS_TOKEN  = "/access_token"
	REPO_SCOPE    = "repo"
	ACCEPT_HEADER = "application/json"
)


type createTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"` // bearer
	Scope       string `json:"scope"`      // "repo, gist, ..."
}

func createToken(deviceCode string) (createTokenResponse, error) {
	var res createTokenResponse
	paths := []string{"login", "oauth", "access_token"}
	params := map[string]string{
		"client_id":   CLIENT_ID,
		"device_code": deviceCode,
		"grant_type":  GRANT_TYPE,
	}

	url, err := utils.BuildURL(BASE_URL, paths, params)
	if err != nil {
		return res, err
	}

	headers := http.Header{}
	headers.Add("Accept", ACCEPT_HEADER)

	resp, err := utils.SubmitPostRequest(url, nil, headers)
	if err != nil {
		return res, err
	}

	tokenErr := handleCreateTokenErr(resp)
	if tokenErr.Error != "" {
		return res, errors.New(tokenErr.ErrorDesc)
	}

	return res, nil
}

type createTokenError struct {
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
}

func handleCreateTokenErr(data []byte) createTokenError {
	var tokenErr createTokenError
	err := json.Unmarshal(data, &tokenErr)
	if err != nil {
		tokenErr.Error = err.Error()
		return tokenErr
	}
	return tokenErr
}

type verifyDeviceResponse struct {
	DeviceCode      string `json:"device_code"`      // device verification code used to verify the device
	UserCode        string `json:"user_code"`        // displayed on the device, user will enter code into browser (verificationUri)
	VerificationUri string `json:"verification_uri"` // the url where the user needs to enter the user_code
	ExpiresIn       int    `json:"expires_in"`       // seconds (default is 900 seconds = 15 minutes)
	Interval        int    `json:"interval"`         // min num of seconds that must pass before we make a new access token req
}

// verifyDevice will make a POST request to GitHubs device and user verification
// code service. It returns a struct containing information that is needed
// to create an access token. This is step 1 of the device flow.
// See -> https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
func verifyDevice() (verifyDeviceResponse, error) {
	var res verifyDeviceResponse
	paths := []string{"login", "device", "code"}
	params := map[string]string{"client_id": CLIENT_ID, "scope": REPO_SCOPE}

	url, err := utils.BuildURL(BASE_URL, paths, params)
	if err != nil {
		return res, err
	}

	headers := http.Header{}
	headers.Add("Accept", ACCEPT_HEADER)

	resp, err := utils.SubmitPostRequest(url, nil, headers)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
