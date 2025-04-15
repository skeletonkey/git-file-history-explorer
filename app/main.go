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
	version = "0.1.1"

	firstCommitIndex widget.ListItemID = 0

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

	var lastCommitId widget.ListItemID = firstCommitIndex

	repo := repository.NewRepo(filename)

	fileContents, err := repo.GetFileLogs(firstCommitIndex)
	report.PanicOnError(err)

	fileContentsLabel := widget.NewLabel(fileContents)
	commitDetailsLabel := widget.NewLabel(repo.Commits[firstCommitIndex].FullCommit)

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
	commitList.Select(firstCommitIndex)
	commitList.OnSelected = func(id widget.ListItemID) {
		fileContents, err := repo.GetFileLogs(id)
		if err != nil {
			report.ErrorPopUp(
				fmt.Errorf("unable to get logs for commit %s: %s", repo.Commits[id].Hash, err),
				w,
				func() { commitList.Select(lastCommitId) },
			)

			return
		}
		fileContentsLabel.SetText(fileContents)
		commitDetailsLabel.SetText(repo.Commits[id].FullCommit)
	}
	commitList.OnUnselected = func(id widget.ListItemID) {
		lastCommitId = id
	}

	leftSplit := container.NewVSplit(commitList, container.NewScroll(commitDetailsLabel))
	leftSplit.Offset = commitListToDetails

	screenSplit := container.NewHSplit(leftSplit, container.NewScroll(fileContentsLabel))
	screenSplit.Offset = commitInfoToFile
	w.SetContent(container.NewGridWithColumns(1, screenSplit))

	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.CenterOnScreen()
	w.SetTitle(fmt.Sprintf("%s [ver. %s]", repo.GetTitle(), version))
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
