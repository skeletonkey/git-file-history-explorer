package main

//go:generate go run app-init.go

import instanceGen "github.com/skeletonkey/lib-instance-gen-go/app"

func main() {
	app := instanceGen.NewApp("git-file-history-explorer", "app")
	app.SetupApp(
		app.WithCGOEnabled(),
		app.WithDependencies("fyne.io/fyne/v2", "github.com/go-git/go-git/v5"),
		app.WithGoVersion("1.23"),
		app.WithMakefile(),
	).Generate()
}
