package repository

import "github.com/go-git/go-git/v5/plumbing/object"

const (
	commitLabelMaxLength = 34 // characters - zero-based # - appears to be limited by list widget
)

type commitData struct {
	author     string
	committer  string
	FullCommit string
	Hash       string
	message    string
	shortHash  string
}

func newCommitData(c *object.Commit) (d commitData) {
	d.author = c.Author.String()
	d.committer = c.Committer.String()
	d.FullCommit = c.String()
	d.Hash = c.Hash.String()
	d.message = c.Message
	d.shortHash = c.Hash.String()[:8]

	return d
}

func (c commitData) Label() string {
	msg := c.shortHash + " - " + c.message
	if len(msg) > commitLabelMaxLength {
		msg = msg[:commitLabelMaxLength-3] + "..."
	}
	return msg
}
