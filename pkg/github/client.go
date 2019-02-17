package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

type Client struct {
	organization string
	project      string
	client       *gh.Client
	ctx          context.Context
}

// NewClient create a new github client
func NewClient(organization string, project string, token string) (*Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)

	return &Client{
		organization: organization,
		project:      project,
		client:       client,
		ctx:          ctx,
	}, nil
}

// GetIssue return the github issue with the specified ID
func (c *Client) GetIssue(number int) (*gh.Issue, error) {
	issue, _, err := c.client.Issues.Get(c.ctx, c.organization, c.project, number)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

// GetMilestoneByTitle return the github milestone with the specified title
func (c *Client) GetMilestoneByTitle(title string) (*gh.Milestone, error) {
	milestones, _, err := c.client.Issues.ListMilestones(c.ctx, c.organization, c.project, &gh.MilestoneListOptions{
		State: "open", // We keep only open milestone or we should add a loop to fetch all milestones since open are not in the first batch
	})

	if err != nil {
		return nil, err
	}
	for _, m := range milestones {
		if *m.Title == title {
			return m, nil
		}
	}

	return nil, fmt.Errorf("milestone %q not fund", title)
}

// GetIssuesFromMilestone return the list of issues, excluding PRs from a specific milestone
func (c *Client) GetIssuesFromMilestone(milestone *gh.Milestone) ([]*gh.Issue, error) {
	issues := []*gh.Issue{}

	morePage := true
	page := 0
	for morePage {
		tempIssues, resp, err := c.client.Issues.ListByRepo(c.ctx, c.organization, c.project, &gh.IssueListByRepoOptions{
			State:     "all",
			Milestone: fmt.Sprintf("%d", *milestone.Number),
			ListOptions: gh.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, err
		}
		page = resp.NextPage
		if resp.NextPage == 0 {
			morePage = false
		}

		for _, i := range tempIssues {
			if i.IsPullRequest() {
				continue
			}
			issues = append(issues, i)
		}
	}

	return issues, nil
}
