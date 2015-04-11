package passward

import (
	"errors"
	"time"

	"github.com/jandre/passward/util"
	git2go "gopkg.in/libgit2/git2go.v22"
)

//
// Contains git helpers
//
type Git struct {
	Name  string
	Email string
	Path  string
	repo  *git2go.Repository
}

func NewGit(path string, name string, email string) *Git {
	return &Git{Path: path, Name: name, Email: email}
}

func (git *Git) Initialize() error {
	var err error

	if util.DirectoryExists(git.Path) {
		if git.repo, err = git2go.OpenRepository(git.Path); err != nil {
			return err
		}
	} else {
		if git.repo, err = git2go.InitRepository(git.Path, false); err != nil {
			return err
		}
	}

	return nil
}

func (git *Git) makeSignature() *git2go.Signature {

	return &git2go.Signature{
		Name:  git.Name,
		Email: git.Email,
		When:  time.Now(),
	}
}

//
// equivalent of git add . ; git commit -a -m <msg>
//
func (g *Git) CommitAllChanges(msg string) error {
	var tip *git2go.Commit
	var commit *git2go.Oid

	if g.repo == nil {
		return errors.New("No repo - have you called Initialize()?")
	}

	idx, err := g.repo.Index()
	if err != nil {
		return err
	}

	err = idx.AddAll([]string{"*/**", "*"}, git2go.IndexAddDefault, nil)
	if err != nil {
		return err
	}
	oid, err := idx.WriteTree()
	if err != nil {
		return err
	}

	sig := g.makeSignature()

	if branch, err := g.repo.Head(); branch != nil {
		tip, err = g.repo.LookupCommit(branch.Target())
		if err != nil {
			return err
		}
	}

	tree, err := g.repo.LookupTree(oid)

	if err != nil {
		return err
	}

	if tip != nil {
		commit, err = g.repo.CreateCommit("HEAD", sig, sig, msg, tree, tip)
	} else {
		commit, err = g.repo.CreateCommit("HEAD", sig, sig, msg, tree)
	}

	debug("returned commit: %s", commit)
	return err
}
