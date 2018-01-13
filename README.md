# Git utils

Utility to help GitHub usage

## Install

1. Clone this repo
1. Include `bin` directory in your path, or create symlinks:
    - `ln -s $(pwd)/bin/git-gh-pr /usr/local/bin/git-gh-pr`
    - `ln -s $(pwd)/bin/git-gh-start /usr/local/bin/git-gh-start`

## Configure

1. Create a new GitHub token from [settings page](https://github.com/settings/tokens)
1. Configure each project using:
    - `git config --add gh.username YOURNAME` (replace `YOURNAME` with your GitHub username, eg. `mavimo`)
    - `git config --add gh.project PROJECT/NAME` (replace `PROJECT/NAME` with GitHub project name eg.: `mavimo/git-helper-gh`)
    - `git config --add gh.token GITHUB_TOKEN` (replace `GITHUB_TOKEN` with your GitHub token generated above)

## Usage

1. Start to work on a new feature `git gh-start ISSUE_ID` (replace `ISSUE_ID` with the issue ID you are start to working)
1. Create a pull request using `git gh-pr`, a new tab in your browser will be open with preconfigured PR label, author, title and content; please choise assegnee and add more info if needed.
