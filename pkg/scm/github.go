package scm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/utils"
)

const (
	BASE_URL           = "https://github.com"
	GITHUB_BASE_URL    = "https://api.github.com"
	CLIENT_ID          = "ca711ca70149e4948032"
	GRANT_TYPE         = "urn:ietf:params:oauth:grant-type:device_code"
	ACCESS_TOKEN       = "/access_token"
	SCOPES             = "repo"
	ACCEPT_JSON        = "application/json"
	ACCEPT_VDN         = "application/vnd.github+json"
	GITHUB_API_VERSION = "2022-11-28"
)

type GitHubManager struct{}

func (gh *GitHubManager) Report(issues []Issue) <-chan int64 {
	idChan := make(chan int64)
	wg := sync.WaitGroup{}
	wg.Add(len(issues))

	for _, issue := range issues {
		go func(is Issue) {
			defer wg.Done()
			resp, err := createIssue(is)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			idChan <- resp.ID
		}(issue)
	}

	go func() {
		wg.Wait()
		close(idChan)
	}()

	return idChan
}

type createIssueResponse struct {
	URL           string `json:"url"`
	RepositoryURL string `json:"repository_url"`
	LabelsURL     string `json:"labels_url"`
	CommentsURL   string `json:"comments_url"`
	EventsURL     string `json:"events_url"`
	HTMLURL       string `json:"html_url"`
	ID            int64  `json:"id"`
	NodeID        string `json:"node_id"`
	Number        int    `json:"number"`
	Title         string `json:"title"`
}

func createIssue(issue Issue, accessToken string) (createIssueResponse, error) {
	var res createIssueResponse

	repoName, err := RepoName()
	if err != nil {
		return res, err
	}

	userName, err := GlobalUserName()
	if err != nil {
		return res, err
	}

	paths := []string{"repos", userName, repoName, "issues"}
	url, err := utils.BuildURL(GITHUB_BASE_URL, paths, nil)
	if err != nil {
		return res, err
	}

	body, err := json.Marshal(issue)
	if err != nil {
		return res, err
	}

	headers := http.Header{}
	headers.Add("Accept", ACCEPT_VDN)
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	headers.Add("X-GitHub-Api-Version", GITHUB_API_VERSION)

	resp, err := utils.SubmitPostRequest(url, bytes.NewReader(body), headers)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Authorize satisfies the GitManager interface. Each source code management
// platform will have their own version of how to authorize so that
// the program can submit issues on the users behalf. This implementation
// uses GitHubs device oauth flow. See their docs for more detailed information.
// First, a user code is created and a browser opens to GitHubs verification url.
// While the program is waiting for the user to enter the code, we poll an endpoint
// and check if the user has authorized the app. Once they have done so, an access token
// is returned from the service and is then written to ~/.config/issue-summoner/config.json
func (gh *GitHubManager) Authorize() error {
	var once sync.Once
	deviceChan := make(chan requestDeviceVerificationResponse)
	tokenChan := make(chan createTokenResponse)
	errChan := make(chan error)

	go initDeviceFlow(deviceChan, errChan)
	device := <-deviceChan
	go pollTokenService(tokenChan, device, errChan, &once)

	select {
	case token := <-tokenChan:
		return WriteToken(token.AccessToken, GH)
	case err := <-errChan:
		return err
	}
}

// @TODO do we still need the IsAuthorized func?
// there is a check in git.go that validates if the user
// has an access token
func (gh *GitHubManager) IsAuthorized() (bool, error) {
	return CheckForAccess(GH)
}

func initDeviceFlow(vd chan requestDeviceVerificationResponse, ec chan error) {
	resp, err := requestDeviceVerification()
	if err != nil {
		ec <- err
		return
	}

	fmt.Printf(
		"User Code: %s - Please visit %s if you have any isues",
		resp.UserCode,
		resp.VerificationUri,
	)

	err = utils.OpenBrowser(resp.VerificationUri)
	if err != nil {
		fmt.Printf(
			"failed to open default browser. Please visit %s and enter your User Code",
			resp.VerificationUri,
		)
	}
	vd <- resp
}

// pollTokenService will make an http POST request to check if the user has successfully
// authorized the app by entering the user_code into the browser. The function will not
// poll the endpoint at a higher frequency than the frequency indicated by **interval**
// in the **requestDeviceVerificationResponse** struct. GitHub will respond with a 200 status code and
// an error response
func pollTokenService(
	tc chan createTokenResponse,
	device requestDeviceVerificationResponse,
	ec chan error,
	once *sync.Once,
) {
	expireTime := time.Now().Add(time.Duration(device.ExpiresIn) * time.Second)
	ticker := time.NewTicker(time.Duration(device.Interval+1) * time.Second)

	defer ticker.Stop()

	for {
		if time.Now().After(expireTime) {
			once.Do(func() {
				ec <- errors.New("User Code has expired, please re-run the 'issue-summoner report' command to generate a new user code")
				close(ec)
			})
			break
		}

		<-ticker.C
		resp, err := createToken(device.DeviceCode)
		if err == nil {
			tc <- resp
		}
	}
}

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
	headers.Add("Accept", ACCEPT_JSON)

	resp, err := utils.SubmitPostRequest(url, nil, headers)
	if err != nil {
		return res, err
	}

	tokenErr := handleCreateTokenErr(resp)
	if tokenErr.Error != "" {
		return res, errors.New(tokenErr.ErrorDesc)
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return res, err
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

type requestDeviceVerificationResponse struct {
	DeviceCode      string `json:"device_code"`      // device verification code used to verify the device
	UserCode        string `json:"user_code"`        // displayed on the device, user will enter code into browser (verificationUri)
	VerificationUri string `json:"verification_uri"` // the url where the user needs to enter the user_code
	ExpiresIn       int    `json:"expires_in"`       // seconds (default is 900 seconds = 15 minutes)
	Interval        int    `json:"interval"`         // min num of seconds that must pass before we make a new access token req
}

// requestDeviceVerification will make a POST request to GitHubs device and user verification
// code service. It returns a struct containing information that is needed
// to create an access token. This is step 1 of the device flow.
// See -> https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
func requestDeviceVerification() (requestDeviceVerificationResponse, error) {
	var res requestDeviceVerificationResponse
	paths := []string{"login", "device", "code"}
	params := map[string]string{"client_id": CLIENT_ID, "scope": SCOPES}

	url, err := utils.BuildURL(BASE_URL, paths, params)
	if err != nil {
		return res, err
	}

	headers := http.Header{}
	headers.Add("Accept", ACCEPT_JSON)

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
