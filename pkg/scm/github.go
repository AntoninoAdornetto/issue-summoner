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

