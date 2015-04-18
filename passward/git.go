package passward

import (
	"errors"
	"regexp"
	"time"

	git2go "github.com/libgit2/git2go"

	"github.com/jandre/passward/util"
	// git2go "gopkg.in/libgit2/git2go.v22"
)

func detectGitName(url string) string {
	// find last instance of *.git
	re, err := regexp.Compile(".+/(.*).git")
	if err != nil {
		panic(err)
	}
	matches := re.FindStringSubmatch(url)
	if matches != nil && len(matches) == 2 {
		return matches[1]
	}
	return ""
}

//
// Contains git helpers
//
type Git struct {
	path        string
	credentials *Credentials
	repo        *git2go.Repository
}

var instance *Git

func credentialsCallback(url string, username string, allowedTypes git2go.CredType) (git2go.ErrorCode, *git2go.Cred) {
	return instance.getGitCredentials()
}

func certificateCheckCallback(cert *git2go.Certificate, valid bool, hostname string) git2go.ErrorCode {
	// Made this one just return 0 during troubleshooting...
	return 0
}

func (git *Git) CloneRepository(url string) error {
	instance = git
	opts := git2go.CloneOptions{}
	opts.RemoteCallbacks = &git2go.RemoteCallbacks{
		CredentialsCallback:      credentialsCallback,
		CertificateCheckCallback: certificateCheckCallback,
	}

	repo, err := git2go.Clone(url, git.path, &opts)
	if err != nil {
		return err
	}
	git.repo = repo
	return nil
}

func NewGit(path string, credentials *Credentials) *Git {
	return &Git{path: path, credentials: credentials}
}

func (git *Git) Push() error {

	remote, err := git.repo.LookupRemote("origin")

	if err != nil {
		return err
	}

	if remote == nil {
		return errors.New("No remote found, did you call `SetRemote`?")
	}

	instance = git
	cbs := &git2go.RemoteCallbacks{
		CredentialsCallback:      credentialsCallback,
		CertificateCheckCallback: certificateCheckCallback,
	}

	remote.SetCallbacks(cbs)

	return remote.Push([]string{"refs/heads/master"}, nil)

	// TODO: handle pull and sync etc
}

func (git *Git) SetRemote(remote string) error {
	gitRemote, err := git.repo.CreateRemote("origin", remote)

	if err != nil {
		return err
	}

	expected := []string{
		"+refs/heads/*:refs/remotes/origin/*",
	}

	if err := gitRemote.SetFetchRefspecs(expected); err != nil {
		return err
	}

	if err := gitRemote.Save(); err != nil {
		return err
	}

	// branch, err := git.repo.LookupBranch("master", git2go.BranchLocal)

	// if err != nil {
	// return err
	// }

	// branch.

	// if err := branch.SetUpstream("master"); err != nil {
	// log.Println("XXXX", err)
	// return err
	// }

	return nil
}

func (git *Git) Initialize() error {
	var err error

	if util.DirectoryExists(git.path) {
		if git.repo, err = git2go.OpenRepository(git.path); err != nil {
			return err
		}
	} else {
		if git.repo, err = git2go.InitRepository(git.path, false); err != nil {
			return err
		}
	}

	return nil
}

func (git *Git) getGitCredentials() (git2go.ErrorCode, *git2go.Cred) {
	err, cred := git2go.NewCredSshKey("git", git.credentials.PublicKeyPath,
		git.credentials.PrivateKeyPath, git.credentials.Passphrase())
	return git2go.ErrorCode(err), &cred
}

func (git *Git) makeSignature() *git2go.Signature {

	return &git2go.Signature{
		Name:  git.credentials.Name,
		Email: git.credentials.Email,
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

	err = idx.AddAll([]string{"**"}, git2go.IndexAddDefault, nil)

	if err != nil {
		return err
	}

	oid, err := idx.WriteTree()
	if err != nil {
		return err
	}

	sig := g.makeSignature()

	tree, err := g.repo.LookupTree(oid)

	if err != nil {
		return err
	}

	if branch, err := g.repo.Head(); branch != nil {
		tip, err = g.repo.LookupCommit(branch.Target())
		if err != nil {
			return err
		}
	}

	if tip != nil {
		commit, err = g.repo.CreateCommit("HEAD", sig, sig, msg, tree, tip)
	} else {
		commit, err = g.repo.CreateCommit("HEAD", sig, sig, msg, tree)
	}

	if err != nil {
		return err
	}

	// writes the index to disk
	err = idx.Write()

	debug("returned commit: %s", commit)
	return err
}
