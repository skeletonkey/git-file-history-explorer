package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"os/exec"
	"strings"
)

// Notes:
// goland ~/.gvm/pkgsets/go1.20/global/pkg/mod/fyne.io/fyne/v2\@v2.3.5/
//    cmd/fyne_demo/tutorials
// fyne_demo &

func main() {
	filename := getFileName()

	a := app.New()
	w := a.NewWindow(filename)

	commits := getCommits(".", filename)
	fileContentsLabel := widget.NewLabel(commits[0].getFile(filename))
	commitInfoLabel := widget.NewLabel(commits[0].fullCommit)

	listWidget := widget.NewList(
		func() int {
			return len(commits)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(commits[i].label())
		})
	listWidget.OnSelected = func(id widget.ListItemID) {
		fileContentsLabel.SetText(commits[id].getFile(filename))
		commitInfoLabel.SetText(commits[id].fullCommit)
	}

	w.SetContent(container.NewBorder(
		nil,
		nil,
		container.NewVSplit(listWidget, container.NewVScroll(commitInfoLabel)),
		container.NewVScroll(fileContentsLabel),
	))

	w.ShowAndRun()
}

func getCommits(dir string, filename string) []commitData {
	var commits []commitData

	repo, err := git.PlainOpen(dir)
	panicOnError(err)
	logOptions := git.LogOptions{
		FileName: &filename,
	}
	cIter, err := repo.Log(&logOptions)
	panicOnError(err)
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, newCommitData(c))
		return nil
	})
	panicOnError(err)

	return commits
}

type commitData struct {
	author     string
	committer  string
	hash       string
	message    string
	shortHash  string
	fullCommit string
}

func newCommitData(c *object.Commit) (d commitData) {
	d.author = c.Author.String()
	d.committer = c.Committer.String()
	d.hash = c.Hash.String()
	d.message = c.Message
	d.shortHash = c.Hash.String()[:8]
	d.fullCommit = c.String()

	return d
}

func (c commitData) label() string {
	msg := c.message
	if len(msg) > 20 {
		msg = msg[:20]
	}
	return c.shortHash + " - " + msg
}

func (c commitData) getFile(filename string) string {
	var out strings.Builder
	cmd := exec.Command("git", "show", c.hash+":"+filename)
	cmd.Stdout = &out
	err := cmd.Run()
	panicOnError(err)
	return out.String()
}

func getFileName() string {
	if len(os.Args) == 1 {
		panicOnError(fmt.Errorf("%s requires a filename as an argument", os.Args[0]))
	}
	if _, err := os.Stat(os.Args[1]); err != nil {
		panicOnError(err)
	}
	return os.Args[1]
}

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
