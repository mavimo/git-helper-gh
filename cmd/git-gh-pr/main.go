package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	gh "github.com/google/go-github/v24/github"
	"github.com/mavimo/git-helper-gh/pkg/cli"
	"github.com/mavimo/git-helper-gh/pkg/git"
	"github.com/mavimo/git-helper-gh/pkg/github"
	"github.com/mavimo/git-helper-gh/pkg/path"
	"github.com/pkg/browser"
	"gopkg.in/gookit/color.v1"
)

func main() {
	repoPath, err := path.GetGitRepoPath()
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	config, err := git.GetConfig()
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	err = cli.CheckConfig(config)
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(gitPath, "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	branchName := strings.TrimSpace(string(output))

	branchKind, err := getBranchKind(branchName)

	url := ""
	if branchKind == "TECH" {
		url = getTechIssueURL(config, "develop", branchName)
	} else {
		issueID, err := strconv.Atoi(branchKind)
		if err != nil {
			color.Error.Println(err)
			os.Exit(1)
		}

		url = getStandardIssueURL(config, "develop", branchName, issueID)
	}

	browser.OpenURL(url)
}

func getBranchKind(branchName string) (string, error) {
	branchSplit := strings.Split(branchName, "/")
	if len(branchSplit) != 2 {
		color.Error.Println("branch should be feature/ISSUE_ID-slug-name or fix/ISSUE_ID-slug-name")
		os.Exit(1)
	}

	branchIssueSplit := strings.Split(branchSplit[1], "-")
	if len(branchIssueSplit) < 2 {
		color.Error.Println("branch should be feature/ISSUE_ID-slug-name or fix/ISSUE_ID-slug-name")
		os.Exit(1)
	}

	if branchIssueSplit[0] == "TECH" {
		return "TECH", nil
	}

	_, err := strconv.Atoi(branchIssueSplit[0])
	if err != nil {
		return "", err
	}

	return branchIssueSplit[0], nil
}

func getAssignees(issue *gh.Issue) string {
	if len(issue.Assignees) == 0 {
		return "N/A"
	}

	assignees := []string{}
	for _, assignee := range issue.Assignees {
		assignees = append(assignees, *assignee.Login)
	}

	return strings.Join(assignees, ",")
}

func getLabels(issue *gh.Issue) string {
	labels := []string{}
	for _, label := range issue.Labels {
		labels = append(labels, *label.Name)
	}
	labels = append(labels, "status/CR-NEEDED")

	return strings.Join(labels, ",")
}

func getTechIssueURL(config map[string]string, sourceBranchName string, destinationBranchName string) string {
	prTitle := "refs #TECH:"
	prBody := `
<!-- ARE YOU SURE WE DON'T HAVE AN ISSUE FOR IT -->

### Relevant commits/breaking changes

  -

### Database migrations

  - no

### Integration notes

  - no

### Deploy instructions

  - standard

### Checks

  - [ ] I manually tested the feature
  - [ ] I wrote automatic tests
  - [ ] I wrote fixtures for it <!-- remove if not needed -->
  - [ ] I wrote API doc for it  <!-- remove if not needed -->
  - [ ] I documented how to test it for QA in issue
  - [ ] I do a commit cleanup before ask for PRs
`

	url := fmt.Sprintf(
		"https://github.com/%s/compare/%s...%s?expand=1&pull_request[title]=%s&pull_request[body]=%s&milestone=%s&labels=%s&assignee=%s",
		config["gh.project"],
		sourceBranchName,
		destinationBranchName,
		prTitle,
		prBody,
		"next-release",
		strings.Join([]string{"area/TECH", "status/CR-NEEDED"}, ","),
		config["gh.username"],
	)

	return url
}

func getStandardIssueURL(config map[string]string, sourceBranchName string, destinationBranchName string, issueID int) string {
	values := strings.Split(config["gh.project"], "/")
	client, err := github.NewClient(values[0], values[1], config["gh.token"])
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	issue, err := client.GetIssue(issueID)
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	prTitle := fmt.Sprintf("refs #%d: %s", *issue.Number, *issue.Title)
	prBody := fmt.Sprintf(`
Closes #%d

### Relevant commits/breaking changes

  - %s

### Database migrations

  - no

### Integration notes

  - no

### Deploy instructions

  - standard

### Checks

  - [ ] I manually tested the feature
  - [ ] I wrote automatic tests
  - [ ] I wrote fixtures for it <!-- remove if not needed -->
  - [ ] I wrote API doc for it  <!-- remove if not needed -->
  - [ ] I documented how to test it for QA in issue
  - [ ] I do a commit cleanup before ask for PR
`,
		*issue.Number,
		*issue.Title,
	)

	url := fmt.Sprintf(
		"https://github.com/%s/compare/%s...%s?expand=1&pull_request[title]=%s&pull_request[body]=%s&milestone=%s&labels=%s&assignee=%s",
		config["gh.project"],
		sourceBranchName,
		destinationBranchName,
		prTitle,
		prBody,
		*issue.Milestone.Title,
		getLabels(issue),
		getAssignees(issue),
	)

	return url
}
