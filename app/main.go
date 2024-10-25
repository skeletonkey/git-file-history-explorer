package main

import (
	"fmt"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/skeletonkey/git-file-history-explorer/pkg/report"
	"github.com/skeletonkey/git-file-history-explorer/pkg/repository"
)

// Notes:
// goland ~/.gvm/pkgsets/go1.20/global/pkg/mod/fyne.io/fyne/v2\@v2.3.5/
//    cmd/fyne_demo/tutorials
// fyne_demo &

const (
	// window sizes
	windowHeight float32 = 700
	windowWidth  float32 = 1200

	// splits
	commitListToDetails float64 = .8
	commitInfoToFile    float64 = .35
)

func main() {
	filename := getFileName()

	a := app.New()
	w := a.NewWindow(filename)

	repo := repository.NewRepo(filename)
	fileContentsLabel := widget.NewLabel(repo.GetFileLogs(0))
	commitDetailsLabel := widget.NewLabel(repo.Commits[0].FullCommit)

	commitList := widget.NewList(
		func() int {
			return len(repo.Commits)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(repo.Commits[i].Label())
		})
	commitList.Select(0)
	commitList.OnSelected = func(id widget.ListItemID) {
		fileContentsLabel.SetText(repo.GetFileLogs(id))
		commitDetailsLabel.SetText(repo.Commits[id].FullCommit)
	}

	leftSplit := container.NewVSplit(commitList, container.NewScroll(commitDetailsLabel))
	leftSplit.Offset = commitListToDetails

	screenSplit := container.NewHSplit(leftSplit, container.NewScroll(fileContentsLabel))
	screenSplit.Offset = commitInfoToFile
	w.SetContent(container.NewGridWithColumns(1, screenSplit))

	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.CenterOnScreen()
	w.SetTitle(repo.GetTitle())
	w.ShowAndRun()
}

// getFileName checks that a filename is provided and that the file exists
func getFileName() string {
	if len(os.Args) == 1 {
		_, filename := path.Split(os.Args[0])
		report.PanicOnError(fmt.Errorf("%s requires a filename as an argument", filename))
	}
	fileInfo, err := os.Stat(os.Args[1])
	if err != nil {
		report.PanicOnError(fmt.Errorf("error attempting to get FileInfo for %s: %s", os.Args[1], err))
	}
	if fileInfo.IsDir() {
		report.PanicOnError(fmt.Errorf("filename (%s) provided is a directory", fileInfo.Name()))
	}

	return os.Args[1]
}
