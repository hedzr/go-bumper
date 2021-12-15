package logic

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log"
	"github.com/hedzr/log/dir"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
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
	sha256table, err = loadShaFile(cmdr.GetStringRP(prefix, "sha256file"))
	if err != nil {
		return
	}

	var formulaFile string
	var repo *git.Repository
	formulaFile, repo, err = cloneToLocal(tap, actor, token, formula)
	if err != nil {
		return
	}
	log.Tracef("repo cloned: %v", repo)

	var count int
	count, err = updateFormulaFile(formulaFile, ver, verS, sha256table)
	if err != nil {
		return
	}

	if count > 0 && repo != nil {
		fmt.Printf("> committing ...\n")
		err = commitToTap(repo, formula, ver, actor, actorMail, prefix)
	}
	return
}

func commitToTap(repo *git.Repository, formula, ver, actor, actorMail, prefix string) (err error) {
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
						err = repo.Push(&git.PushOptions{})
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

	nf, err = os.Create(formulaFile + ".new")
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer nf.Close()

	lastLineIsMatched, lastLine := false, ""
	verRE := regexp.MustCompile(`\d+\.\d+\.\d`)
	urlRE := regexp.MustCompile(`(url[ \t]+['"])(.+?)(['"])`)
	shaRE := regexp.MustCompile(`(sha256[ \t]+['"])(.+?)(['"])`)

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

	fmt.Printf("\n%v replaced.\n", count)
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
		sha256table[filename] = sha
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

	repo, err = git.PlainClone(tgtDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	fmt.Printf("repo cloned: %v\n", repo)
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
		url = fmt.Sprintf("https://%v:%v@github.com/%v", actor, token, tap)
	}
	return
}
