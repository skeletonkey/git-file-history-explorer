package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

	commitLabelMaxLength = 34 // characters - zero-based # - appears to be limited by list widget
)

func main() {
	filename := getFileName()

	a := app.New()
	w := a.NewWindow(filename)

	repo := newRepo(filename)
	fileContentsLabel := widget.NewLabel(repo.getFileLogs(0))
	commitDetailsLabel := widget.NewLabel(repo.commits[0].fullCommit)

	commitList := widget.NewList(
		func() int {
			return len(repo.commits)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(repo.commits[i].label())
		})
	commitList.Select(0)
	commitList.OnSelected = func(id widget.ListItemID) {
		fileContentsLabel.SetText(repo.getFileLogs(id))
		commitDetailsLabel.SetText(repo.commits[id].fullCommit)
	}

	leftSplit := container.NewVSplit(commitList, container.NewScroll(commitDetailsLabel))
	leftSplit.Offset = commitListToDetails

	screenSplit := container.NewHSplit(leftSplit, container.NewScroll(fileContentsLabel))
	screenSplit.Offset = commitInfoToFile
	w.SetContent(container.NewGridWithColumns(1, screenSplit))

	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.CenterOnScreen()
	w.SetTitle(repo.baseDir + ":" + repo.relativeFile)
	w.ShowAndRun()
}

type repo struct {
	baseDir      string
	commits      []commitData
	relativeFile string
}

func newRepo(file string) repo {
	fileDir, _ := path.Split(file)
	if fileDir == "" {
		fileDir = "."
	}

	dir := executeCmd("git", "-C", fileDir, "rev-parse", "--show-toplevel")

	fullFilename, err := filepath.Abs(file)
	panicOnError(err)

	r := repo{
		baseDir:      dir,
		relativeFile: fullFilename[len(dir)+1:],
	}
	r.setCommits()

	return r
}

func (r *repo) setCommits() {
	var commits []commitData

	repo, err := git.PlainOpen(r.baseDir)
	panicOnError(err)

	filename := r.relativeFile
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

	r.commits = commits
}

func (r *repo) getFileLogs(commitID int) string {
	return executeCmd("git", "-C", r.baseDir, "show", r.commits[commitID].hash+":"+r.relativeFile)
}

type commitData struct {
	author     string
	committer  string
	fullCommit string
	hash       string
	message    string
	shortHash  string
}

func newCommitData(c *object.Commit) (d commitData) {
	d.author = c.Author.String()
	d.committer = c.Committer.String()
	d.fullCommit = c.String()
	d.hash = c.Hash.String()
	d.message = c.Message
	d.shortHash = c.Hash.String()[:8]

	return d
}

func (c commitData) label() string {
	msg := c.shortHash + " - " + c.message
	if len(msg) > commitLabelMaxLength {
		msg = msg[:commitLabelMaxLength-3] + "..."
	}
	return msg
}

func executeCmd(cmdName string, args ...string) string {
	var out strings.Builder
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	panicOnError(err)
	return strings.TrimRight(out.String(), "\n")
}

// getFileName checks that a filename is provided and that the file exists
func getFileName() string {
	if len(os.Args) == 1 {
		_, filename := path.Split(os.Args[0])
		panicOnError(fmt.Errorf("%s requires a filename as an argument", filename))
	}
	fileInfo, err := os.Stat(os.Args[1])
	if err != nil {
		panicOnError(fmt.Errorf("error attempting to get FileInfo for %s: %s", os.Args[1], err))
	}
	if fileInfo.IsDir() {
		panicOnError(fmt.Errorf("filename (%s) provided is a directory", fileInfo.Name()))
	}

	return os.Args[1]
}

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
