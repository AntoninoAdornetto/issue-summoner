package git

import (
	"bytes"
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
	github_api_version  = "2022-11-28"
	github_client_id    = "ca711ca70149e4948032"
	github_grant_type   = "urn:ietf:params:oauth:grant-type:device_code"
)

type githubManager struct {
	config       map[string]utils.IssueSummonerConfig
	repo         *Repository
	device       deviceVerificationResponse
	reportParams *createIssueRequest
}

func (ghub *githubManager) Report(req ReportRequest, res chan ReportResponse) {
	result := ReportResponse{Index: req.Index}

	if ghub.reportParams == nil {
		params, err := ghub.prepareReportRequest()
		if err != nil {
			result.Err = err
			res <- result
			return
		}
		ghub.reportParams = &params
	}

	data, err := json.Marshal(req)
	if err != nil {
		result.Err = err
		res <- result
		return
	}

	url, headers := ghub.reportParams.url, ghub.reportParams.headers
	body := bytes.NewBuffer(data)

	resp, data, err := utils.MakeRequest("POST", url, body, headers)
	if err != nil {
		result.Err = err
		res <- result
		return
	}

	if resp.StatusCode != 201 {
		err := handleCreateIssueErr(data, resp.StatusCode, req.Title)
		result.Err = err
		res <- result
		return
	}

	if err := json.Unmarshal(data, &result); err != nil {
		result.Err = err
		res <- result
		return
	}

	res <- result
}

type createIssueRequest struct {
	url     string
	headers http.Header
}

func (ghub *githubManager) prepareReportRequest() (createIssueRequest, error) {
	req := createIssueRequest{headers: http.Header{}}

	accessToken := ghub.config[GITHUB].Auth.AccessToken
	paths := []string{"repos", ghub.repo.UserName, ghub.repo.RepoName, "issues"}

	url, err := utils.BuildURL(github_api_base_url, nil, paths...)
	if err != nil {
		return req, err
	}

	req.headers.Add("Accept", "application/vnd.github+json")
	req.headers.Add("Authorization", "Bearer "+accessToken)
	req.headers.Add("X-GitHub-Api-Version", github_api_version)
	req.url = url
	return req, nil
}

type createIssueErrorResponse struct {
	Message string `json:"message"`
}

func handleCreateIssueErr(data []byte, statusCode int, title string) error {
	var res createIssueErrorResponse
	err := json.Unmarshal(data, &res)
	if err != nil {
		return fmt.Errorf(
			"failed to create issue <%s> with status code: %d\terror: %s",
			title,
			statusCode,
			err.Error(),
		)
	}

	if statusCode == 404 {
		return fmt.Errorf(
			"failed to create issue <%s> with status code: %d\terror: check that you are authorized to report issues. <issue-summoner authorize> command",
			title,
			statusCode,
		)
	}

	return fmt.Errorf(
		"failed to create issue <%s> with status code: %d\terror: %s",
		title,
		statusCode,
		res.Message,
	)
}

func (ghub *githubManager) Authorize() error {
	var err error

	if ghub.device, err = requestDeviceVerification(); err != nil {
		return err
	}

	if err = ghub.startDeviceFlow(); err != nil {
		return err
	}

	tokenChan, errChan := make(chan createTokenResponse), make(chan error)
	go ghub.pollAuthService(tokenChan, errChan)

	select {
	case token := <-tokenChan:
		ghub.config[GITHUB] = utils.IssueSummonerConfig{
			Auth: utils.AuthConfig{
				AccessToken: token.AccessToken,
				CreatedAt:   time.Now(),
			},
		}

		// @TODO move config marshaling into config util file
		data, err := json.Marshal(ghub.config)
		if err != nil {
			return err
		}

		return utils.WriteIssueSummonerConfig(data)
	case err = <-errChan:
		return err
	}
}

// @TODO IsAuthorized should be more intelligent in the way it validates an access token.
func (ghub *githubManager) IsAuthorized() bool {
	if entry, ok := ghub.config[GITHUB]; ok {
		return entry.Auth.AccessToken != ""
	}
	return false
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

	url, err := utils.BuildURL(
		github_base_url,
		device_verification_params,
		device_verification_uris...)
	if err != nil {
		return res, err
	}

	_, data, err := utils.MakeRequest("POST", url, nil, headers)
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return res, err
	}

	return res, nil
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

func (ghub *githubManager) createToken() (createTokenResponse, error) {
	var res createTokenResponse
	headers := http.Header{}
	headers.Add("Accept", "application/json")
	create_token_params["device_code"] = ghub.device.DeviceCode

	url, err := utils.BuildURL(github_base_url, create_token_params, create_token_uris...)
	if err != nil {
		return res, err
	}

	_, data, err := utils.MakeRequest("POST", url, nil, headers)
	if err != nil {
		return res, err
	}

	if createTokenErr := onCreateTokenError(data); createTokenErr.Error != "" {
		return res, errors.New(createTokenErr.ErrorDesc)
	}

	if err = json.Unmarshal(data, &res); err != nil {
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
func (ghub *githubManager) pollAuthService(res chan createTokenResponse, err chan error) {
	expireTime := time.Now().Add(time.Duration(ghub.device.ExpiresIn) * time.Second)
	ticker := time.NewTicker(time.Duration(ghub.device.Interval+1) * time.Second)
	defer ticker.Stop()

	for {
		if time.Now().After(expireTime) {
			err <- errors.New(
				"User Code has expired. Please re-run 'issue-summoner authorize' command to generate a new user code",
			)
			return
		}

		data, err := ghub.createToken()
		if err != nil {
			<-ticker.C
		} else {
			res <- data
		}
	}
}

// after requesting a device verification code, it is now time to instruct
// the user to open their browser and enter the user code. The function will
// attempt to open the users default browser. If for some reason we cannot open
// the browser, we will print the url in the terminal the user can manually visit
// device verification page.
func (ghub *githubManager) startDeviceFlow() error {
	if ghub.device.UserCode == "" {
		return errors.New("expected a device user code but got empty string")
	}

	if ghub.device.VerificationUri == "" {
		return errors.New("expected a valid verification url but got empty string")
	}

	fmt.Printf(
		"Enter User Code: %s at %s\n",
		ghub.device.UserCode,
		ghub.device.VerificationUri,
	)

	if err := utils.OpenBrowser(ghub.device.VerificationUri); err != nil {
		fmt.Printf(
			"Failed to open default browser. please visit %s and enter your User Code\n",
			ghub.device.VerificationUri,
		)
	}

	return nil
}
