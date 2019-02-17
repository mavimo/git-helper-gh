package cli

import (
	"fmt"
)

func CheckConfig(config map[string]string) error {
	missingParameter := false
	if config["gh.project"] == "" {
		formatError("project")
		missingParameter = true
	}
	if config["gh.token"] == "" {
		formatError("token")
		missingParameter = true
	}
	if config["gh.username"] == "" {
		formatError("username")
		missingParameter = true
	}

	if missingParameter {
		return fmt.Errorf("Some configuration is missing, please fix it before continue")
	}

	return nil
}

func formatError(token string) {
	fmt.Printf(`Configure your github %s with:

	git config --add gh.%s VALUE

`, token, token)
}
