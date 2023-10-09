package main

//go:generate go run app-init.go

import instanceGen "github.com/skeletonkey/lib-instance-gen-go/app"

func main() {
	app := instanceGen.NewApp("git-file-history-explorer", "app")
	app.WithCGOEnabled().
		WithDependencies("fyne.io/fyne/v2", "github.com/go-git/go-git/v5").
		WithGoVersion("1.21").
		WithGithubWorkflows("linter", "test").
		WithMakefile()
}
