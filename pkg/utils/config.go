/*
Copyright © 2024 AntoninoAdornetto

READ/WRITE UTILITIES FOR ISSUE SUMMONER. THE MAIN PURPOSE OF THE
CONFIG FILE IS TO STORE TOKENS, SUCH AS AN ACCESS TOKEN, THAT WILL BE USED
IN HTTP REQUEST TO SOURCE CODE MANAGEMENT PLATFORMS WHEN REPORTING ISSUES
*/

package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type IssueSummonerConfig struct {
	Auth AuthConfig `json:"auth"`
}

type AuthConfig struct {
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

var (
	Config = map[string]IssueSummonerConfig{
		"github":    {},
		"gitlab":    {},
		"bitbucket": {},
	}
)

func ReadIssueSummonerConfig() ([]byte, error) {
	conf, err := getConfigFile()
	if err != nil {
		return nil, err
	}

	defer conf.Close()
	return io.ReadAll(conf)
}

func WriteIssueSummonerConfig(data []byte) error {
	conf, err := getConfigFile()
	if err != nil {
		return err
	}

	defer conf.Close()

	if _, err = conf.Seek(0, 0); err != nil {
		return err
	}

	if err = conf.Truncate(0); err != nil {
		return err
	}

	_, err = conf.Write(data)
	return err
}

func getConfigFile() (*os.File, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "issue-summoner")

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(filepath.Join(path, "config.json"), os.O_RDWR|os.O_CREATE, 0666)
}
