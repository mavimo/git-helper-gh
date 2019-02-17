package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mavimo/git-helper-gh/pkg/cli"
	"github.com/mavimo/git-helper-gh/pkg/git"
	"github.com/mavimo/git-helper-gh/pkg/github"
)

func main() {
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
		fmt.Println("You should specify a version")
		os.Exit(1)
	}
	milestoneTitle := os.Args[1]
	milestone, err := client.GetMilestoneByTitle(milestoneTitle)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	issues, err := client.GetIssuesFromMilestone(milestone)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, i := range issues {
		fmt.Printf(" - refs [#%d](%s): %s\n", *i.Number, *i.HTMLURL, *i.Title)
	}
}
