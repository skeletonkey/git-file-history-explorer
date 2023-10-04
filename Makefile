build:
	CGO_ENABLED=1 go build -o bin/git-file-history-explorer

deps:
	go get fyne.io/fyne/v2@latest
	go get github.com/go-git/go-git/v5@latest
	go mod tidy
