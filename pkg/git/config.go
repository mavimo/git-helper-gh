package git

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/mavimo/git-helper-gh/pkg/path"
	"github.com/muja/goconfig"
)

func getParsedConfig(gitconfig string) (map[string]string, error) {
	bytes, err := ioutil.ReadFile(gitconfig)
	if err != nil {
		return nil, err
	}
	config, _, err := goconfig.Parse(bytes)

	return config, err
}

func GetConfig() (map[string]string, error) {
	configFileLocal, err := path.GetGitConfigFile()
	if err != nil {
		return nil, err
	}
	configLocal, err := getParsedConfig(configFileLocal)
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	configFileGlobal := filepath.Join(currentUser.HomeDir, ".gitconfig")
	configGlobal, err := getParsedConfig(configFileGlobal)

	for key, value := range configLocal {
		configGlobal[key] = value
	}

	return configGlobal, nil
}
