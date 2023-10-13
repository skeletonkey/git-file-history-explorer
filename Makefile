// lib-instance-gen-go: File auto generated -- DO NOT EDIT!!!
.DEFAULT_GOAL=build

build:
	go fmt ./...
	go vet ./...
	CGO_ENABLED=1 go build -o bin/git-file-history-explorer app/*.go

install:
	cp bin/git-file-history-explorer /usr/local/sbin/git-file-history-explorer

golib-latest:
	go get -u fyne.io/fyne/v2@latest
	go get -u github.com/go-git/go-git/v5@latest
	go get -u github.com/skeletonkey/lib-core-go@latest
	go get -u github.com/skeletonkey/lib-instance-gen-go@latest

	go mod tidy

app-init:
	go generate
