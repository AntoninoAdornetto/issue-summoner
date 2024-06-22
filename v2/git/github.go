package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/utils"
)

const (
	github_base_url     = "https://github.com"
	github_api_base_url = "https://api.github.com"
	github_client_id    = "ca711ca70149e4948032"
	github_grant_type   = "urn:ietf:params:oauth:grant-type:device_code"
)

type githubManager struct {
	config map[string]IssueSummonerConfig
	repo   *Repository
	device deviceVerificationResponse
}

// @TODO check for other services that can properly validate a bearer/access token versus checking for empty string
func (github *githubManager) IsAuthorized() bool {
	if entry, ok := github.config[GITHUB]; ok {
		return entry.Auth.AccessToken != ""
	}
	return false
}

func (github *githubManager) Authorize() error {
	var err error

	if github.device, err = requestDeviceVerification(); err != nil {
		return err
	}

	if err = github.print(); err != nil {
		return err
	}

	res, errChan := make(chan createTokenResponse), make(chan error)
	go github.pollAuthService(res, errChan)

	select {
	case token := <-res:
		github.config[GITHUB] = IssueSummonerConfig{
			Auth: authConfig{
				AccessToken: token.AccessToken,
				CreatedAt:   time.Now(),
			},
		}

		data, err := json.Marshal(github.config)
		if err != nil {
			return err
		}

		return utils.WriteIssueSummonerConfig(data)
	case err = <-errChan:
		return err
	}
}

type deviceVerificationResponse struct {
	DeviceCode      string `json:"device_code"`      // device verification code used to verify the device
	UserCode        string `json:"user_code"`        // displayed on the device, user will enter code into browser (verificationUri)
	VerificationUri string `json:"verification_uri"` // the url where the user needs to enter the user_code
	ExpiresIn       int    `json:"expires_in"`       // seconds (default is 900 seconds = 15 minutes)
	Interval        int    `json:"interval"`         // min num of seconds that must pass before we make a new access token req
}

var (
	device_verification_uris   = []string{"login", "device", "code"}
	device_verification_params = map[string]string{"client_id": github_client_id, "scope": "repo"}
)

// requestDeviceVerification sends a POST request to Githubs device and user verification code
// service. It returns a struct containing information that is needed to create an access token.
// https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
func requestDeviceVerification() (deviceVerificationResponse, error) {
	var res deviceVerificationResponse
	headers := http.Header{}
	headers.Add("Accept", "application/json")

	url, err := buildURL(github_base_url, device_verification_params, device_verification_uris...)
	if err != nil {
		return res, err
	}

	resp, err := makeRequest("POST", url, nil, headers)
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(resp, &res); err != nil {
		return res, err
	}

	return res, nil
}

func (github *githubManager) print() error {
	if github.device.UserCode == "" {
		return errors.New("expected a device user code but got empty string")
	}

	if github.device.VerificationUri == "" {
		return errors.New("expected a valid verification url but got empty string")
	}

	fmt.Printf(
		"User Code: %s\n",
		github.device.UserCode,
	)

	if err := utils.OpenBrowser(github.device.VerificationUri); err != nil {
		fmt.Printf(
			"Failed to open default browser. please visit %s and enter your User Code\n",
			github.device.VerificationUri,
		)
	}

	return nil
}

type createTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

var (
	create_token_uris   = []string{"login", "oauth", "access_token"}
	create_token_params = map[string]string{
		"client_id":   github_client_id,
		"device_code": "",
		"grant_type":  github_grant_type,
	}
)

func (github *githubManager) createToken() (createTokenResponse, error) {
	var res createTokenResponse
	headers := http.Header{}
	headers.Add("Accept", "application/json")
	create_token_params["device_code"] = github.device.DeviceCode

	url, err := buildURL(github_base_url, create_token_params, create_token_uris...)
	if err != nil {
		return res, err
	}

	resp, err := makeRequest("POST", url, nil, headers)
	if err != nil {
		return res, err
	}

	if createTokenErr := onCreateTokenError(resp); createTokenErr.Error != "" {
		return res, errors.New(createTokenErr.ErrorDesc)
	}

	if err = json.Unmarshal(resp, &res); err != nil {
		return res, err
	}

	return res, nil
}

type createTokenErrorResposne struct {
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
}

func onCreateTokenError(data []byte) createTokenErrorResposne {
	var res createTokenErrorResposne

	if err := json.Unmarshal(data, &res); err != nil {
		res.Error = err.Error()
	}

	return res
}

// pollAuthService makes an http POST request to githubs oauth token creation service to
// see if the user has successfully authorized the app by entering the user_code into
// the verification_uri. We do not poll the endpoint at a higher frequency than the
// interval specified at `interval` in `deviceVerificationResponse`. An error is returned
// if the user does not enter the user code within the alloted time frame
func (github *githubManager) pollAuthService(res chan createTokenResponse, err chan error) {

	expireTime := time.Now().Add(time.Duration(github.device.ExpiresIn) * time.Second)
	ticker := time.NewTicker(time.Duration(github.device.Interval+1) * time.Second)
	defer ticker.Stop()

	fmt.Println(expireTime)
	for {
		if time.Now().After(expireTime) {
			err <- errors.New(
				"User Code has expired. Please re-run 'issue-summoner authorize' command to generate a new user code",
			)
			return
		}

		data, err := github.createToken()
		if err != nil {
			<-ticker.C
		} else {
			res <- data
		}
	}
}
