package repository

import (
	"path"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/skeletonkey/git-file-history-explorer/pkg/report"
)

type repo struct {
	baseDir      string
	Commits      []commitData
	relativeFile string
}

func NewRepo(file string) repo {
	fileDir, _ := path.Split(file)
	if fileDir == "" {
		fileDir = "."
	}

	dir := executeCmd("git", "-C", fileDir, "rev-parse", "--show-toplevel")

	fullFilename, err := filepath.Abs(file)
	report.PanicOnError(err)

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
	report.PanicOnError(err)

	filename := r.relativeFile
	logOptions := git.LogOptions{
		FileName: &filename,
	}
	cIter, err := repo.Log(&logOptions)
	report.PanicOnError(err)
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, newCommitData(c))
		return nil
	})
	report.PanicOnError(err)

	r.Commits = commits
}

func (r *repo) GetFileLogs(commitID int) string {
	return executeCmd("git", "-C", r.baseDir, "show", r.Commits[commitID].hash+":"+r.relativeFile)
}

func (r *repo) GetTitle() string {
	return r.baseDir + ":" + r.relativeFile
}
