package common

import (
	"encoding/json"
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

type Config = map[string]IssueSummonerConfig

var (
	// emptyConfig is utilized when the config.json file does not exist within the root dir for user specific
	// configuration data. The location depends on the operating system. Additionally, emptyConfig defines
	// the allowed source code hosting configurations.
	emptyConfig = Config{
		"github":    {},
		"gitlab":    {},
		"bitbucket": {},
	}
)

func ReadConfig() (Config, error) {
	f, err := getConfigFile()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return emptyConfig, nil
	}

	var res Config
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func WriteToConfig(conf Config) error {
	f, err := getConfigFile()
	if err != nil {
		return err
	}

	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Seek(0, 0); err != nil {
		return err
	}

	if err = f.Truncate(0); err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func getConfigFile() (*os.File, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "issue-summoner")
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(filepath.Join(path, "config.json"), os.O_RDWR, 0666)
}
