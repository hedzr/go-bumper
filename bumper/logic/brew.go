package logic

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log"
	"github.com/hedzr/log/dir"
)

func brewIt(cmd *cmdr.Command, remainArgs []string) (err error) {
	//str := dumpIt(cmd, m)
	//outputWithFormat(str, "Address")
	prefix := cmd.GetDottedNamePath()
	tap := cmdr.GetStringRP(prefix, "tap")
	formula := cmdr.GetStringRP(prefix, "formula")
	ver := cmdr.GetStringRP(prefix, "release-version")
	actor := cmdr.GetStringRP(prefix, "actor")
	actorMail := cmdr.GetStringRP(prefix, "actor-mail")
	ref := cmdr.GetStringRP(prefix, "ref")
	if ref != "" {
		ver = strings.TrimLeft(ref, "refs/tags/")
	}
	token := cmdr.GetStringRP(prefix, "token")
	verS := strings.TrimLeft(ver, "v")
	fmt.Printf(`
> Bumping '%v' to '%v' into Tap '%v' ...

         Version: %v, %v
             Ref: %v
           Actor: %v

> Cloneing the brew tap repository ...
`, formula, ver, tap, ver, verS, ref, actor)

	var sha256table map[string]string
	log.Debugf("> loadShaFile ...")
	sha256table, err = loadShaFile(cmdr.GetStringRP(prefix, "sha256file"))
	if err != nil {
		return
	}

	var formulaFile string
	var repo *git.Repository
	log.Debugf("> cloneToLocal ...")
	formulaFile, repo, err = cloneToLocal(tap, actor, token, formula)
	if err != nil {
		return
	}
	log.Tracef("repo cloned: %v", repo)

	var count int
	log.Debugf("> updateFormulaFile ...")
	count, err = updateFormulaFile(formulaFile, ver, verS, sha256table)
	if err != nil {
		return
	}

	if count > 0 && repo != nil {
		fmt.Printf("> committing ...\n")
		err = commitToTap(repo, formula, ver, actor, actorMail, token, prefix)
	}
	return
}

func commitToTap(repo *git.Repository, formula, ver, actor, actorMail, token, prefix string) (err error) {
	var w *git.Worktree
	w, err = repo.Worktree()
	if err == nil {
		_, err = w.Add("Formula/" + formula + ".rb")

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
					if cmdr.GetBoolRP(prefix, "push") {
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

func updateFormulaFile(formulaFile, ver, verS string, sha256table map[string]string) (count int, err error) {
	var file, nf *os.File
	file, err = os.Open(formulaFile)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer file.Close()

	log.Debugf("> creating %v.new ...", formulaFile)
	nf, err = os.Create(formulaFile + ".new")
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer nf.Close()

	lastLineIsMatched, lastLine := false, ""
	verRE := regexp.MustCompile(`\d+\.\d+\.\d+(-[a-z0-9]+)?`)
	urlRE := regexp.MustCompile(`(url[ \t]+['"])(.+?)(['"])`)
	shaRE := regexp.MustCompile(`(sha256[ \t]+['"])(.+?)(['"])`)

	log.Debugf("> scanning %v ...", formulaFile)
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		line := scanner.Text()
		if lastLineIsMatched {
			lastLineIsMatched = false
			if strings.Contains(line, "sha256") {
				um := urlRE.FindSubmatch([]byte(lastLine))
				u := string(um[2])
				bu := path.Base(u)
				if sha, ok := sha256table[bu]; ok {
					line = shaRE.ReplaceAllString(line, "${1}"+sha+"${3}")
					count++
				}
			}
		} else if verRE.Match([]byte(line)) {
			if strings.Contains(line, "url") {
				lastLineIsMatched, lastLine = true, line
			}
			line = verRE.ReplaceAllString(line, verS)
			count++
		}

		_, err = nf.WriteString(line)
		_, err = nf.WriteString("\n")
	}

	if err = scanner.Err(); err != nil {
		//log.Fatal(err)
		return
	}

	fmt.Printf("\n%v replaced.\n\n\nNEW CONTENTS:\n", count)
	catFile(formulaFile + ".new")

	fmt.Println("\n\nbinaries.asc CONTENTS:")
	for k, v := range sha256table {
		fmt.Printf("%v => %v\n", k, v)
	}

	err = dir.DeleteFile(formulaFile)
	if err == nil {
		err = os.Rename(formulaFile+".new", formulaFile)
	}
	return
}

func loadShaFile(shaFile string) (sha256table map[string]string, err error) {
	sha256table = make(map[string]string)

	var file *os.File
	file, err = os.Open(shaFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		filename, sha := fields[len(fields)-1], fields[0]
		sha256table[path.Base(filename)] = sha
	}
	return
}

func cloneToLocal(tap, actor, token, formula string) (formulaFile string, repo *git.Repository, err error) {
	tgtDir := fmt.Sprintf("/tmp/%v", tap)
	url := tapToRepoUrl(tap, actor, token)
	formulaFile = path.Join(tgtDir, "Formula", formula+".rb")
	if dir.FileExists(tgtDir) {
		if dir.FileExists(path.Join(tgtDir, ".git")) && dir.FileExists(formulaFile) {
			repo, err = git.PlainOpen(tgtDir)
			return
		}
	}

	log.Debugf("       actor: %v", actor)
	log.Debugf("       token: %v", token)
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
	url, ok := tap, false
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
		//url = fmt.Sprintf("https://%v:%v@github.com/%v.git", actor, token, tap)
		url = fmt.Sprintf("https://github.com/%v.git", tap)
	}
	return
}

func catFile(filename string) {
	var file *os.File
	var err error
	file, err = os.Open(filename)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		line := scanner.Text()
		fmt.Printf("%v\n", line)
	}
}
