package path

import (
	"os"
	"path/filepath"
	"strings"
)

// GetGitRepoPath the the repo path
func GetGitRepoPath() (string, error) {
	dir, err := currentDirectory()
	if err != nil {
		return "", err
	}

	for {
		if dir == "" {
			return "", nil
		}

		_, err := gitConfigFile(dir)

		if err == nil {
			return dir, nil
		}

		dir, _ = filepath.Split(strings.TrimRight(dir, string(os.PathSeparator)))
	}
}

// GetGitConfigFile return the file or error if not fund of the git config
// starting from current directory and walking to the parent.
func GetGitConfigFile() (string, error) {
	dir, err := currentDirectory()
	if err != nil {
		return "", err
	}

	for {
		if dir == "" {
			return "", nil
		}

		configFile, err := gitConfigFile(dir)

		if err == nil {
			return configFile, nil
		}

		dir, _ = filepath.Split(strings.TrimRight(dir, string(os.PathSeparator)))
	}
}

func currentDirectory() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func gitConfigFile(folder string) (string, error) {
	filename := filepath.Join(folder, ".git", "config")

	fileInfo, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return "", err
	}

	if fileInfo != nil && fileInfo.IsDir() {
		return "", err
	}

	return filename, nil
}
