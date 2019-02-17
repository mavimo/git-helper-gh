build:
	go build -o bin/git-gh-pr cmd/git-gh-pr/main.go
	go build -o bin/git-gh-start cmd/git-gh-start/main.go
	go build -o bin/git-gh-release cmd/git-gh-release/main.go
