package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hedzr/log"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/logex/build"
	"golang.org/x/sys/unix"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var pushAllowed bool

func init() {
	build.New(build.NewLoggerConfigWith(true, "logrus", "debug"))
	pushAllowed = true
}

func main() {
	tap := "hedzr/test1"
	actor := "hedzr"
	token := os.Getenv("GITHUB_TOKEN")
	formula := ""

	log.Debugf("> cloneToLocal")
	ff, repo, err := cloneToLocal(tap, actor, token, formula)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Debugf("ff: %v", ff)

	log.Debugf("> updateFile")
	err = updateFile(ff)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Debugf("> commitToTap")
	err = commitToTap(repo, "README.md", formula, "ver", actor, "hedzrz@gmail.com", token, "prefix")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func commitToTap(repo *git.Repository, ff, formula, ver, actor, actorMail, token, prefix string) (err error) {
	var w *git.Worktree
	w, err = repo.Worktree()
	if err == nil {
		_, err = w.Add(ff)

		if err == nil {
			var commit plumbing.Hash
			var obj *object.Commit
			commit, err = w.Commit(fmt.Sprintf("bump %v to %v", formula, ver), &git.CommitOptions{
				Author: &object.Signature{
					Name:  actor,
					Email: actorMail,
					When:  time.Now(),
				},
			})

			if err == nil {
				obj, err = repo.CommitObject(commit)
				if err == nil {
					log.Debugf("committed: %v", obj)
					if pushAllowed {
						log.Debugf("pushing as %q...", actor)
						err = repo.Push(&git.PushOptions{
							//Auth: &http.TokenAuth{
							//	Token: token,
							//},
							Auth: &http.BasicAuth{
								Username: actor, // yes, this can be anything except an empty string
								Password: token,
							},
							RemoteName: "origin",
							Progress:   os.Stdout,
						})
						if err == nil {
							log.Debugf("pushed ...")
						} else {
							log.Errorf("push failed: %v", err)
						}
					} else {
						//patch := obj.Patch(previousCommit)
					}
				}
			}
		}
	}
	return
}

func updateFile(ff string) (err error) {
	var file *os.File
	file, err = os.OpenFile(ff, os.O_RDWR, 0664)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer file.Close()

	_, err = file.Seek(0, unix.SEEK_END)
	_, err = file.WriteString(time.Now().String())
	_, err = file.WriteString("\n\n")
	return
}

func cloneToLocal(tap, actor, token, formula string) (formulaFile string, repo *git.Repository, err error) {
	tgtDir := fmt.Sprintf("/tmp/%v", tap)
	url := tapToRepoUrl(tap, actor, token)
	//formulaFile = path.Join(tgtDir, "Formula", formula+".rb")
	formulaFile = path.Join(tgtDir, "README.md")
	if dir.FileExists(tgtDir) {
		if dir.FileExists(path.Join(tgtDir, ".git")) && dir.FileExists(formulaFile) {
			repo, err = git.PlainOpen(tgtDir)
			return
		}
	}

	log.Debugf("       actor: %v", actor)
	log.Debugf("       study: %v", token)
	log.Debugf("  clone from: %v", url)
	repo, err = git.PlainClone(tgtDir, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: actor, // yes, this can be anything except an empty string
			Password: token,
		},
		// cannot work:
		//Auth: &http.TokenAuth{
		//	Token: token,
		//},
		URL:      url,
		Progress: os.Stdout,
	})
	if err == nil {
		fmt.Printf("repo cloned: %v\n", repo)
		rs, _ := repo.Remotes()
		for _, r := range rs {
			fmt.Printf("  remote: %v", r.Config().Name)
		}
	}
	return
}

func tapToRepoUrl(tap, actor, token string) (url string) {
	url, ok := tap, true
	if !strings.Contains(url, "://") {
		ok = false
	} else {
		matched, err := regexp.Match("@.+:", []byte(url))
		if err != nil {
			ok = false
		} else {
			ok = matched
		}
	}
	if !ok {
		url = fmt.Sprintf("https://%v:%v@github.com/%v.git", actor, token, tap)
		//url = fmt.Sprintf("https://github.com/%v.git", tap)
	}
	return
}
