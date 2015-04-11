package passward

import (
	"github.com/jandre/passward/util"
	git2go "gopkg.in/libgit2/git2go.v22"
)

//
// Contains git helpers
//
type Git struct {
	Path string
	repo *git2go.Repository
}

// TODO: creds?
func NewGit(path string) (*Git, error) {
	git := Git{Path: path}
	var err error

	if util.DirectoryExists(path) {
		if git.repo, err = git2go.OpenRepository(path); err != nil {
			return nil, err
		}
	} else {
		if git.repo, err = git2go.InitRepository(path, false); err != nil {
			return nil, err
		}
	}
	return &git, nil
}
