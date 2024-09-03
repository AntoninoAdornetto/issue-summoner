package git

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/common"
)

const (
	githubBaseUrl    = "https://github.com"
	githubApiBaseUrl = "https://api.github.com"
	githubApiVersion = "2022-11-28"
	githubClientId   = "ca711ca70149e4948032"
	githubGrantType  = "urn:ietf:params:oauth:grant-type:device_code"
)

type githubManager struct {
	conf        common.Config
	repo        *Repository
	device      requestDeviceResponse
	reportReq   *http.Request
	getIssueReq *http.Request
}

// Authorize will request device infromation details from githubs device flow,
// i.e., [device.UserCode] [device.VerificationUri], open the users default browser,
// and poll an authentication endpoint for an access token.
func (ghub *githubManager) Authorize() error {
	var err error

	if ghub.device, err = requestDevice(); err != nil {
		return err
	}

	if err = ghub.startDeviceFlow(); err != nil {
		return err
	}

	tokenChan, errChan := make(chan createTokenResponse), make(chan error)
	go ghub.pollAuthService(tokenChan, errChan)

	select {
	case token := <-tokenChan:
		if entry, ok := ghub.conf[Github]; ok {
			entry.Auth = common.AuthConfig{
				AccessToken: token.AccessToken,
				CreatedAt:   time.Now(),
			}
			ghub.conf[Github] = entry
		} else {
			return fmt.Errorf("No entry for %s in config map", Github)
		}

		return common.WriteToConfig(ghub.conf)
	case err = <-errChan:
		return err
	}
}

type requestDeviceResponse struct {
	DeviceCode      string `json:"device_code"`      // device verification code used to verify the device
	ExpiresIn       int    `json:"expires_in"`       // seconds (default is 900 seconds = 15 minutes)
	Interval        int    `json:"interval"`         // min num of seconds that must pass before we make a new access token req
	UserCode        string `json:"user_code"`        // displayed on the device, user will enter code into browser (verificationUri)
	VerificationUri string `json:"verification_uri"` // the url where the user needs to enter the user_code
}

var (
	verificationUris   = []string{"login", "device", "code"}
	verificationParams = map[string]string{"client_id": githubClientId, "scope": "repo"}
)

// requestDevice sends a POST request to Githubs device and user verification code
// service. It returns a struct containing information that is needed to create an access token.
// https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
func requestDevice() (requestDeviceResponse, error) {
	var res requestDeviceResponse
	headers := http.Header{}
	headers.Add("Accept", "application/json")

	url, err := common.ConstructURL(githubBaseUrl, verificationParams, verificationUris...)
	if err != nil {
		return res, err
	}

	_, data, err := common.Request("POST", url, nil, headers)
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return res, err
	}

	return res, err
}

// after requesting a device verification code, it is now time to instruct
// the user to open their browser and enter the user code. The function will
// attempt to open the users default browser. If for some reason we cannot open
// the browser, we will print the url in the terminal the user can manually visit
// device verification page.
func (ghub *githubManager) startDeviceFlow() error {
	if ghub.device.UserCode == "" {
		return errors.New("expected a device user code but got an empty string")
	}

	if ghub.device.VerificationUri == "" {
		return errors.New("expected a verification url but got an empty string")
	}

	fmt.Printf("Enter User Code: %s at %s\n", ghub.device.UserCode, ghub.device.VerificationUri)

	if err := common.OpenBrowser(ghub.device.VerificationUri); err != nil {
		fmt.Printf(
			"Failed to open default browser. Please open a browser, visit %s, and enter your User Code\n",
			ghub.device.VerificationUri,
		)
	}

	return nil
}

// pollAuthService makes an HTTP POST request to githubs oauth token creation service to
// determine if the user has successfully authorized the github app by entering the [device.UserCode]
// into the browser at [device.VerificationUri]. The function does not poll the endpoint at a higher
// frequency than the interval [device.Interval]. An error is returned if the user does not enter their
// user code within the alloted time [device.ExpiresIn]
func (ghub *githubManager) pollAuthService(res chan createTokenResponse, err chan error) {
	expireTime := time.Now().Add(time.Duration(ghub.device.ExpiresIn) * time.Second)
	ticker := time.NewTicker(time.Duration(ghub.device.Interval+1) * time.Second)
	defer ticker.Stop()

	for {
		if time.Now().After(expireTime) {
			err <- errors.New("User Code has expired. Please re-run <issue-summoner authorize> command to generate a new user code")
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

type createTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

var (
	createTokenUris   = []string{"login", "oauth", "access_token"}
	createTokenParams = map[string]string{
		"client_id":  githubClientId,
		"grant_type": githubGrantType,
	}
)

func (ghub *githubManager) createToken() (createTokenResponse, error) {
	var res createTokenResponse
	headers := http.Header{}
	headers.Add("Accept", "application/json")
	createTokenParams["device_code"] = ghub.device.DeviceCode

	url, err := common.ConstructURL(githubBaseUrl, createTokenParams, createTokenUris...)
	if err != nil {
		return res, err
	}

	_, data, err := common.Request("POST", url, nil, headers)
	if err != nil {
		return res, err
	}

	if tokenErr := onCreateTokenError(data); tokenErr.Error != "" {
		return res, errors.New(tokenErr.ErrorDesc)
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return res, err
	}

	return res, nil
}

type createTokenErr struct {
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
}

func onCreateTokenError(data []byte) createTokenErr {
	var res createTokenErr

	if err := json.Unmarshal(data, &res); err != nil {
		res.Error = err.Error()
	}

	return res
}

// @TODO Authenticated func should have a more intelligent way to validate an access token
func (ghub *githubManager) Authenticated() bool {
	if entry, ok := ghub.conf[Github]; ok {
		return entry.Auth.AccessToken != ""
	}
	return false
}

type githubReportResponse struct {
	URL         string `json:"url"`
	ID          int64  `json:"id"`
	IssueNumber int    `json:"number"`
}

var (
	errReport = "Failed to create issue <%s> with error: %s"
)

func (ghub *githubManager) Report(issue ReportRequest, res chan ReportResponse) {
	result := ReportResponse{Index: issue.Index}
	if ghub.reportReq == nil {
		if err := ghub.prepareReportRequest(); err != nil {
			result.Err = fmt.Errorf(errReport, issue.Title, err)
			res <- result
			return
		}
	}

	data, err := json.Marshal(issue)
	if err != nil {
		result.Err = fmt.Errorf(errReport, issue.Title, err)
		res <- result
		return
	}

	body, url := bytes.NewBuffer(data), ghub.reportReq.URL.String()
	resp, data, err := common.Request("POST", url, body, ghub.reportReq.Header)
	if err != nil {
		result.Err = fmt.Errorf(errReport, issue.Title, err)
		res <- result
		return
	}

	if resp.StatusCode != http.StatusCreated {
		result.Err = createIssueErr(data, resp.StatusCode, issue.Title)
		res <- result
		return
	}

	createIssueRes := githubReportResponse{}
	if err := json.Unmarshal(data, &createIssueRes); err != nil {
		result.Err = fmt.Errorf(errReport, issue.Title, err)
		res <- result
		return
	}

	result.ID = createIssueRes.IssueNumber
	res <- result
}

func (ghub *githubManager) prepareReportRequest() error {
	accessToken := ghub.conf[Github].Auth.AccessToken
	ghub.reportReq = &http.Request{
		Header: make(http.Header),
	}

	paths := []string{"repos", ghub.repo.UserName, ghub.repo.RepoName, "issues"}
	u, err := common.ConstructURL(githubApiBaseUrl, nil, paths...)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	ghub.reportReq.Header.Add("Accept", "application/vnd.github+json")
	ghub.reportReq.Header.Add("Authorization", "Bearer "+accessToken)
	ghub.reportReq.Header.Add("X-GitHub-Api-Version", githubApiVersion)
	ghub.reportReq.URL = parsedURL
	return nil
}

type createIssueErrorResponse struct {
	Message string `json:"message"`
}

var (
	errCreateIssue = "Failed to create issue <%s> with status code: %d\terror :%s"
)

func createIssueErr(data []byte, statusCode int, title string) error {
	var res createIssueErrorResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf(errCreateIssue, title, statusCode, err.Error())
	}

	if statusCode == http.StatusNotFound {
		return fmt.Errorf(
			errCreateIssue,
			title,
			statusCode,
			"check that you are authorized to report issues. <issue-summoner authorize> command",
		)
	}

	return fmt.Errorf(errCreateIssue, title, statusCode, res.Message)
}

type githubIssueStatusResponse struct {
	State string `json:"state"`
}

func (ghub *githubManager) GetStatus(issueNum, index int, status chan StatusResponse) {
	res := StatusResponse{Index: index, Resolved: false}
	if ghub.getIssueReq == nil {
		if err := ghub.prepareGetStatusRequest(issueNum); err != nil {
			res.Err = err
			status <- res
			return
		}
	}

	url, headers := ghub.getIssueReq.URL.String(), ghub.getIssueReq.Header
	resp, data, err := common.Request("GET", url, nil, headers)
	if err != nil {
		res.Err = err
		status <- res
		return
	}

	if resp.StatusCode != http.StatusOK {
		errRes := onGetIssueError(data)
		res.Err = errors.New(
			fmt.Sprintf("%s with status code: %d", errRes.Message, resp.StatusCode),
		)
		status <- res
		return
	}

	val := githubIssueStatusResponse{}
	if err := json.Unmarshal(data, &val); err != nil {
		status <- StatusResponse{Err: err}
		return
	}

	if val.State == "closed" {
		res.Resolved = true
	}

	status <- res
}

func (ghub *githubManager) prepareGetStatusRequest(issueNum int) error {
	accessToken := ghub.conf[Github].Auth.AccessToken

	ghub.getIssueReq = &http.Request{
		Header: make(http.Header),
	}

	path := fmt.Sprintf("/repos/%s/%s/issues/%d", ghub.repo.UserName, ghub.repo.RepoName, issueNum)

	u, err := common.ConstructURL(githubApiBaseUrl, nil, path)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	ghub.getIssueReq.Header.Add("Accept", "application/vnd.github+json")
	ghub.getIssueReq.Header.Add("Authorization", "Bearer "+accessToken)
	ghub.getIssueReq.Header.Add("X-GitHub-Api-Version", githubApiVersion)
	ghub.getIssueReq.URL = parsedURL
	return nil
}

type getIssueErr struct {
	Message string `json:"message"`
}

func onGetIssueError(data []byte) getIssueErr {
	var res getIssueErr
	if err := json.Unmarshal(data, &res); err != nil {
		res.Message = err.Error()
	}

	return res
}
