package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	gh "github.com/google/go-github/v24/github"
	"github.com/gosimple/slug"
	"github.com/mavimo/git-helper-gh/pkg/cli"
	"github.com/mavimo/git-helper-gh/pkg/git"
	"github.com/mavimo/git-helper-gh/pkg/github"
	"github.com/mavimo/git-helper-gh/pkg/path"
	"gopkg.in/gookit/color.v1"
)

func main() {
	repoPath, err := path.GetGitRepoPath()
	if err != nil {
		os.Exit(1)
	}

	config, err := git.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	err = cli.CheckConfig(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	values := strings.Split(config["gh.project"], "/")
	client, err := github.NewClient(values[0], values[1], config["gh.token"])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		color.Error.Println("You should specify an issue")
		os.Exit(1)
	}

	issueNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	issue, err := client.GetIssue(issueNumber)
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	err = validateIssue(issue, config["gh.username"])
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	printIssueInformations(issue)

	gitPath, err := exec.LookPath("git")
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}

	branchName := fmt.Sprintf("feature/%d-%s", *issue.Number, slug.Make(*issue.Title))

	cmd := exec.Command(gitPath, "-C", repoPath, "checkout", "-b", branchName)
	cmd.Dir = repoPath
	if err = cmd.Run(); err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}
}

func validateIssue(issue *gh.Issue, username string) error {
	if issue.IsPullRequest() {
		return fmt.Errorf("You should not create a new feature branch from pull-request, you need an issue")
	}

	if *issue.State != "open" {
		return fmt.Errorf("You should not create a new feature branch from a closed issue")
	}

	isAssignedToYou := false
	for _, assignee := range issue.Assignees {
		if *assignee.Login == username {
			isAssignedToYou = true
			break
		}
	}

	if isAssignedToYou == false {
		return fmt.Errorf("You should not create a new feature branch from an issue that's not assigned to you")
	}

	for _, label := range issue.Labels {
		if *label.Name == "status/ON-HOLD" {
			return fmt.Errorf("You should not create a new feature branch from an issue that's ON-HOLD")
		}
	}

	return nil
}

func printIssueInformations(issue *gh.Issue) {
	issueArea := []string{}
	issueKind := []string{}
	issueStatus := []string{}
	issueTribe := []string{}
	for _, label := range issue.Labels {
		labelItems := strings.Split(*label.Name, "/")
		if len(labelItems) != 2 {
			continue
		}
		switch labelItems[0] {
		case "area":
			issueArea = append(issueArea, strings.TrimSpace(labelItems[1]))
		case "kind":
			issueKind = append(issueKind, strings.TrimSpace(labelItems[1]))
		case "status":
			issueStatus = append(issueStatus, strings.TrimSpace(labelItems[1]))
		case "tribe":
			tribeItems := strings.Split(labelItems[1], ":")
			issueTribe = append(issueTribe, strings.TrimSpace(tribeItems[0]))
		}
	}

	color.Printf(`
 <op=underscore;>                                                               </>
|
| <op=bold;fg=yellow;>Area</>     : <fg=green;>%s</>
| <op=bold;fg=yellow;>Kind</>     : <fg=green;>%s</>
| <op=bold;fg=yellow;>Status</>   : <fg=green;>%s</>
| <op=bold;fg=yellow;>Tribe</>    : <fg=green;>%s</>
|<op=underscore;>                                                               </>
|
%s
|<op=underscore;>                                                               </>
|
| <op=bold;fg=yellow;>Assigee</>  : <fg=green;>%s</>
| <op=bold;fg=yellow;>Reporter</> : <fg=green;>%s</>
| <op=bold;fg=yellow;>Link</>     : <op=underscore;>%s</>
|<op=underscore;>                                                               </>

`,
		strings.Join(issueArea, ", "),
		strings.Join(issueKind, ", "),
		strings.Join(issueStatus, ", "),
		strings.Join(issueTribe, ", "),
		prefixBody(issue),
		getAssignees(issue),
		*issue.User.Login,
		*issue.HTMLURL,
	)
}

func getAssignees(issue *gh.Issue) string {
	if len(issue.Assignees) == 0 {
		return "N/A"
	}

	assignees := []string{}
	for _, assignee := range issue.Assignees {
		assignees = append(assignees, *assignee.Login)
	}

	return strings.Join(assignees, ", ")
}

func prefixBody(issue *gh.Issue) string {
	return "| " + strings.Join(strings.Split(*issue.Body, "\n"), "\n| ")
}
