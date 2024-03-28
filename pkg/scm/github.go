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
	ACCESS_TOKEN  = "/access_token"
	REPO_SCOPE    = "repo"
	ACCEPT_HEADER = "application/json"
)

type GitHubManager struct {
	AccessToken     string
	RefreshToken    string
	ExpiresAt       string
	IsAuthenticated bool
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
	deviceChan := make(chan verifyDeviceResponse)
	tokenChan := make(chan createTokenResponse)
	errChan := make(chan error)

	go initDeviceFlow(deviceChan, errChan)
	device := <-deviceChan
	go pollTokenService(tokenChan, device, errChan, &once)

	select {
	case token := <-tokenChan:
		return WriteToken(token.TokenType, GH)
	case err := <-errChan:
		return err
	}
}

// @TODO implement IsAuthorized
// read config file and check for access token presence
func (gh *GitHubManager) IsAuthorized() bool {
	return false
}

func initDeviceFlow(vd chan verifyDeviceResponse, ec chan error) {
	resp, err := verifyDevice()
	if err != nil {
		ec <- err
		return
	}

	fmt.Printf("User Code: %s\n", resp.UserCode)
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
// in the **verifyDeviceResponse** struct. GitHub will respond with a 200 status code and
// an error response
func pollTokenService(
	tc chan createTokenResponse,
	device verifyDeviceResponse,
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
		if err != nil {
			fmt.Println(err)
		} else {
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
	headers.Add("Accept", ACCEPT_HEADER)

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
