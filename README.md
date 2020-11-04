# GitHub CLI helper

CLI utils for GitHub

## Install

### Using homebrew

```
brew install https://raw.githubusercontent.com/mavimo/git-helper-gh/master/git-helper-gh.rb
```

### From source

1. Check that you installed `jq` and `sed` (GNU version) on your machine
1. Clone this repo
1. Include `bin` directory in your path, or create symlinks:
    - `ln -s $(pwd)/bin/git-gh-pr /usr/local/bin/git-gh-pr`
    - `ln -s $(pwd)/bin/git-gh-start /usr/local/bin/git-gh-start`
    - `ln -s $(pwd)/bin/git-gh-release /usr/local/bin/git-gh-release`
    - `chmod +x /usr/local/bin/git-gh-*`

## Configuration

1. Create a new GitHub token with "repo" permissions at [settings page](https://github.com/settings/tokens)
1. Configure each project using:
    - `git config --add gh.username YOURNAME` (replace `YOURNAME` with your GitHub username, eg. `mavimo`)
    - `git config --add gh.project PROJECT/NAME` (replace `PROJECT/NAME` with GitHub project name, eg.: `mavimo/git-helper-gh`. PAY ATTENTION as it's case-sensitive!)
    - `git config --add gh.token GITHUB_TOKEN` (replace `GITHUB_TOKEN` with your GitHub token generated above)

## Usage

1. Use `git checkout BRANCH` to switch to the branch you want to use as base (eg. `develop`)
1. Use `git gh-start ISSUE_ID` to start to work on a new feature (replace `ISSUE_ID` with the issue ID you are start to working)
1. After pushing the branch containing your desired changes, use `git gh-pr` to create a new pull request. A new tab in your browser will open with preconfigured PR label, author, title and content; you will only have to choose an assignee and add more information if needed.
    You can specify an optional parameter so set the base branch for the PR: `git gh-pr BASE-BRANCH-NAME`
1. Use `git gh-release MILESTONE_NAME` to create milestone report (replace `MILESTONE_NAME` with the name of the milestone)
