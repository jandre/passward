package passward

import (
	"errors"
	"regexp"
	"time"

	pb "github.com/cheggaaa/pb"
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
// Contains git helpers used by the vaults
//
type Git struct {
	path        string
	credentials *Credentials
	repo        *git2go.Repository
	progressBar *pb.ProgressBar
}

var instance *Git

func credentialsCallback(url string, username string, allowedTypes git2go.CredType) (git2go.ErrorCode, *git2go.Cred) {
	return instance.getGitCredentials()
}

func certificateCheckCallback(cert *git2go.Certificate, valid bool, hostname string) git2go.ErrorCode {
	// Made this one just return 0 during troubleshooting...
	return 0
}

func transferProgressCallback(stats git2go.TransferProgress) git2go.ErrorCode {
	return instance.PrintTransferProgress(stats)
}

func pushTransferProgressCallback(current uint32, total uint32, bytes uint) git2go.ErrorCode {
	return instance.PrintPushTransferProgress(current, total, bytes)
}

//
// HasRemote is true if the repository has a remote
//
func (git *Git) HasRemote() bool {

	if git.repo != nil {
		remote, err := git.repo.LookupRemote("origin")
		if remote != nil && err == nil {
			return true
		}
	}

	return false
}

func (git *Git) PrintPushTransferProgress(current uint32, total uint32, bytes uint) git2go.ErrorCode {

	if total != 0 {
		if git.progressBar != nil {
			git.progressBar.Set(int(current))
		} else {
			git.progressBar = pb.StartNew(int(total))
		}
	}
	return git2go.ErrorCode(0)
}

func (git *Git) PrintTransferProgress(stats git2go.TransferProgress) git2go.ErrorCode {
	if git.progressBar != nil {
		git.progressBar.Set(int(stats.ReceivedObjects))
		//	fmt.Printf("-- Received %d/%d Objects...\n", stats.ReceivedObjects, stats.TotalObjects)
	} else {
		git.progressBar = pb.StartNew(int(stats.TotalObjects))
	}
	return git2go.ErrorCode(0)
}

func (git *Git) PrintTransferCompletion(complete git2go.CompletionCallback) git2go.ErrorCode {
	// fmt.Println("Transfer Complete!")
	return git2go.ErrorCode(0)
}

//
// Clone will clone a remote vault.
//
func (git *Git) Clone(url string) error {
	instance = git
	opts := git2go.CloneOptions{}

	defer func() {
		if git.progressBar != nil {
			git.progressBar.FinishPrint("Transfer complete!")
		}
		git.progressBar = nil
	}()

	opts.RemoteCallbacks = &git2go.RemoteCallbacks{
		CredentialsCallback:      credentialsCallback,
		CertificateCheckCallback: certificateCheckCallback,
		TransferProgressCallback: transferProgressCallback,
	}

	repo, err := git2go.Clone(url, git.path, &opts)
	if err != nil {
		return err
	}
	git.repo = repo
	return nil
}

//
// NewGit creates a new git.  the `path` is the path of
// the repository; the credentials contain the ssh credentials
// used to commit and push remote repositories.
//
// The same ssh key is used for encrypting/decrypting the keys.
//
func NewGit(path string, credentials *Credentials) *Git {
	return &Git{path: path, credentials: credentials}
}

//
// Push will sync the repository to the remote, much like `git push`.
// It does not handle merge conflicts currently.
//
func (git *Git) Push() error {

	remote, err := git.repo.LookupRemote("origin")

	if err != nil {
		debug("no remote repository found:", err)
		return err
	}

	if remote == nil {
		debug("no remote repository found")
		return errors.New("No remote found, did you call `SetRemote`?")
	}

	instance = git
	cbs := &git2go.RemoteCallbacks{
		CredentialsCallback:          credentialsCallback,
		CertificateCheckCallback:     certificateCheckCallback,
		PushTransferProgressCallback: pushTransferProgressCallback,
	}

	defer func() {
		if git.progressBar != nil {
			git.progressBar.FinishPrint("Transfer complete!")
		}
		git.progressBar = nil
	}()

	remote.SetCallbacks(cbs)

	return remote.Push([]string{"refs/heads/master"}, nil)

	// TODO: handle pull and sync etc
}

//
// SetRemote will set the remote url `remote`, e.g. git@github.com:jandre/work.git
//
// Currently it only supports ssh remotes, *not* passowrd or any other kind of
// authentication.
//
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
