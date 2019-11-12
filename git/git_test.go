package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sshGitUrl, httpGitUrl string
	basicAuth, pemAuth    Credential
)

func loadSecrets() {
	s, _ := ioutil.ReadFile("secret.txt")
	ss := strings.Split(string(s), "\n")
	ConfigKey(ss[0])
	sshGitUrl, httpGitUrl = ss[1], ss[2]
	basicAuth = Credential{secret: ss[3]}
	pemAuth = Credential{secret: ss[4]}
	return
}

func TestGit(t *testing.T) {
	loadSecrets()

	pem, _ := pemAuth.GetPem()
	assert.Nil(t,
		NewRepo(sshGitUrl).
			SetSshKey(pem).
			Depth(2).
			CloneToMemory().
			Error())

	u, p, _ := basicAuth.GetUserPassword()
	localRepo := filepath.Join(os.TempDir(), "repo1")
	fmt.Println(localRepo)

	assert.Nil(t,
		NewRepo(httpGitUrl).
			SetBasicAuth(u, p).
			CloneToDir(localRepo).Pull().
			Error())

	assert.Nil(t, OpenRepo(localRepo).
		SetBasicAuth(u, p).
		Pull().
		Error())

	assert.Nil(t, OpenRepo(localRepo).
		SetBasicAuth(u, p).
		Pull().
		Checkout("2821118e33296d96181c046dd018222077ca43d0").
		Tag("test0", "5f8169f83de623f808d7ade198d101b0c29c0a06").
		Tag("xxx", nil).
		CheckoutTag("test0").
		Branch("test1", "2821118e33296d96181c046dd018222077ca43d0").
		CheckoutBranch("test1").
		CheckoutBranch("master").Pull().
		Error())
}

func TestRemoveAny(t *testing.T) {
	s, _ := ioutil.TempDir(os.TempDir(), "aaa")
	ioutil.WriteFile(filepath.Join(s, "xxx"), []byte("123"), os.ModePerm)
	ioutil.WriteFile(filepath.Join(s, "yyy"), []byte("123"), os.ModePerm)
	stat, _ := os.Stat(s)
	assert.True(t, stat.IsDir())
	stat, _ = os.Stat(filepath.Join(s, "xxx"))
	assert.False(t, stat.IsDir())
	err := removeAny(filepath.Join(s, "xxx"))
	assert.Nil(t, err)
	err = removeAny(s)
	assert.Nil(t, err)

}
