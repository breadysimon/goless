package git

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/breadysimon/goless/crypto"
	"github.com/breadysimon/goless/logging"
	"github.com/breadysimon/goless/util"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var log *logging.Logger = logging.GetLogger()

type GitRepository struct {
	gogit.Repository
}

func GetChanges(repo *gogit.Repository, hashOld, hashNew plumbing.Hash) *object.Changes {
	oldEncodeObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, hashOld)
	if err != nil {
		log.Fatal("can't get old encode object: ", err)

	}

	newEncodeObject, err := repo.Storer.EncodedObject(plumbing.CommitObject, hashNew)
	if err != nil {
		log.Fatal("can'tget new encode object: ", err)
	}

	old, _ := object.DecodeCommit(repo.Storer, oldEncodeObject)
	new, _ := object.DecodeCommit(repo.Storer, newEncodeObject)

	oldTree, _ := old.Tree()
	newTree, _ := new.Tree()

	changes, err := oldTree.Diff(newTree)
	if err != nil {
		log.Fatal(err)
	}
	return &changes
}

func CloneGitRepository(gitUrl string, key transport.AuthMethod, depth int) (r *gogit.Repository, err error) {
	var auth transport.AuthMethod
	if key == nil && util.StartWith(gitUrl, "ssh://") {
		var me *user.User
		if me, err = user.Current(); err == nil {
			auth, err = ssh.NewPublicKeysFromFile("git", filepath.Join(me.HomeDir, ".ssh/id_rsa"), "")
		}
	} else {
		auth = key
	}
	if err == nil {
		if r, err = gogit.Clone(memory.NewStorage(), nil, &gogit.CloneOptions{URL: gitUrl, Auth: auth, Depth: depth}); err == nil {
			return
		}
	}
	return
}

type Credential struct {
	name, secret, createdBy string
}

var (
	SECRET_KEY = ""
)

func (o *Credential) GetUserPassword() (username, password string, err error) {
	s, err := crypto.Decrypt(SECRET_KEY, o.secret)
	if err == nil {
		ss := strings.Split(s, "\n")
		if len(ss) == 2 {
			username, password = ss[0], ss[1]
		} else {
			err = errors.New("malformed secret")
		}
	}
	return
}
func (o *Credential) GetPem() (pem []byte, err error) {
	p, err := crypto.Decrypt(SECRET_KEY, o.secret)
	return []byte(p), err
}
func ConfigKey(k string) {
	SECRET_KEY = k
}

type GitRepo struct {
	repo   *gogit.Repository
	gitUrl string
	auth   transport.AuthMethod
	depth  int
	err    error
}

func NewRepo(gitUrl string) (r *GitRepo) {
	r = &GitRepo{}
	r.gitUrl = gitUrl
	return
}
func (r *GitRepo) SetSshKey(pem []byte) *GitRepo {
	if r.err == nil {
		r.auth, r.err = ssh.NewPublicKeys("git", pem, "")
	}
	return r
}
func (r *GitRepo) SetBasicAuth(username, password string) *GitRepo {
	if r.err == nil {
		r.auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}
	return r
}
func (r *GitRepo) Error() error {
	return r.err
}

func (r *GitRepo) Depth(n int) *GitRepo {
	r.depth = n
	return r
}
func removeAny(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return os.RemoveAll(path)
	}
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (r *GitRepo) CloneToDir(path string) *GitRepo {
	if r.err == nil {
		r.err = removeAny(path)
		if r.err == nil {
			r.repo, r.err = gogit.PlainClone(path, false, &gogit.CloneOptions{
				URL:               r.gitUrl,
				Auth:              r.auth,
				Depth:             r.depth,
				RecurseSubmodules: gogit.DefaultSubmoduleRecursionDepth,
			})
		}
	}
	return r
}
func (r *GitRepo) CloneToMemory() *GitRepo {
	if r.err == nil {
		r.repo, r.err = gogit.Clone(memory.NewStorage(), nil, &gogit.CloneOptions{
			URL:               r.gitUrl,
			Auth:              r.auth,
			Depth:             r.depth,
			RecurseSubmodules: gogit.DefaultSubmoduleRecursionDepth,
		})
	}
	return r
}
func OpenRepo(path string) *GitRepo {
	r := &GitRepo{}
	r.repo, r.err = gogit.PlainOpen(path)
	// lst, _ := r.repo.Remotes()
	// for l := range lst {
	// 	fmt.Println(l)
	// }
	return r
}

func (r *GitRepo) Checkout(hash interface{}) *GitRepo {
	if r.err == nil {
		if w, err := r.repo.Worktree(); err != nil {
			r.err = err
		} else {
			r.err = w.Checkout(&gogit.CheckoutOptions{
				Hash: r.GetHash(hash),
			})
		}
	}
	return r
}
func (r *GitRepo) CheckoutBranch(name string) *GitRepo {
	if r.err == nil {
		if w, err := r.repo.Worktree(); err != nil {
			r.err = fmt.Errorf("checkout branch: %v", err)
		} else {
			branch := plumbing.ReferenceName("refs/heads/" + name)
			r.err = w.Checkout(&gogit.CheckoutOptions{
				Branch: branch,
			})
		}
	}
	return r
}
func (r *GitRepo) Branch(name string, hash interface{}) *GitRepo {
	if r.err == nil {
		var h plumbing.Hash
		if c, ok := hash.(string); ok {
			h = plumbing.NewHash(c)
		} else {
			h = hash.(plumbing.Hash)
		}
		branch := plumbing.ReferenceName("refs/heads/" + name)
		ref := plumbing.NewHashReference(branch, h)
		r.err = r.repo.Storer.SetReference(ref)
	}
	return r
}
func (r *GitRepo) GetHash(hash interface{}) (h plumbing.Hash) {
	switch v := hash.(type) {
	case nil:
		href, _ := r.repo.Head()
		h = href.Hash()
	case string:
		h = plumbing.NewHash(v)
	case plumbing.Hash:
		h = v
	}
	return
}
func (r *GitRepo) FindTagHash(tag string) (hash plumbing.Hash, err error) {
	var tagrefs storer.ReferenceIter
	if tagrefs, err = r.repo.Tags(); err == nil {
		err = tagrefs.ForEach(func(t *plumbing.Reference) error {
			tt := t.Name().String()
			if tag == tt[len("refs/tags/"):] {
				hash = t.Hash()
				return storer.ErrStop
			}
			return nil
		})
	}
	return
}
func (r *GitRepo) CheckoutTag(tag string) *GitRepo {
	if r.err == nil {
		var hash plumbing.Hash
		if hash, r.err = r.FindTagHash(tag); r.err == nil {
			if hash.IsZero() {
				r.err = fmt.Errorf("can not find tag: %s", tag)
			} else {
				// return r.Branch(tag, hash).CheckoutBranch(tag)
				return r.Checkout(hash)
			}
		}
	}
	return r
}
func (r *GitRepo) Tag(tag string, hash interface{}) *GitRepo {
	if r.err == nil {
		tagName := plumbing.ReferenceName("refs/tags/" + tag)
		ref := plumbing.NewHashReference(tagName, r.GetHash(hash))
		r.err = r.repo.Storer.SetReference(ref)
	}
	return r
}

func (r *GitRepo) Pull() *GitRepo {
	if r.err == nil {
		var w *gogit.Worktree
		if w, r.err = r.repo.Worktree(); r.err == nil {
			r.err = w.Pull(&gogit.PullOptions{
				RemoteName: "origin",
				Auth:       r.auth,
			})
			if r.err == gogit.NoErrAlreadyUpToDate {
				r.err = nil
			}
		}
	}
	return r
}

func ShowChanges(r *gogit.Repository) error {
	ref, _ := r.Head()

	cIter, _ := r.Log(&gogit.LogOptions{From: ref.Hash()})
	head, _ := cIter.Next()
	pre, _ := cIter.Next()

	cc := GetChanges(r, head.Hash, pre.Hash)
	fmt.Println(cc)
	return nil
}
