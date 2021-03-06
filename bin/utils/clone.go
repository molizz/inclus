package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Repository struct {
	Url string
	Dir string // clone的目录

	out io.Writer // 输出
}

func NewClone(out io.Writer, url, dir string) *Repository {

	if stat, _ := os.Stat(dir); stat != nil && !stat.IsDir() {
		panic("clone的初始化目录不存在.请在inclus.yaml中的 clone_dir 字段中设置")
	}

	return &Repository{
		Url: url,
		Dir: dir,
		out: out,
	}
}

func (r *Repository) Clone() error {
	fmt.Println("git clone ", r.Url)
	if !r.existCloned() {
		gitclone := exec.Command("git", "clone", r.Url)
		gitclone.Stderr = r.out
		gitclone.Dir = r.Dir
		err := gitclone.Run()
		if err != nil {
			return err
		}
	}

	err := r.fetch()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) fetch() error {
	fmt.Println("git fetch ", r.Url)
	projectDir, err := r.ProjectDir()
	if err != nil {
		return err
	}

	fetch := exec.Command("git", "fetch")
	fetch.Stderr = r.out
	fetch.Dir = projectDir
	err = fetch.Run()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ProjectDir() (string, error) {
	projectName, err := r.projectName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", r.Dir, projectName), nil
}

func (r *Repository) RepoDir() (string, error) {
	projectDir, err := r.ProjectDir()
	if err != nil {
		return "", err
	}
	repoDir := fmt.Sprintf("%s/.git", projectDir)
	return repoDir, nil
}

func (r *Repository) existCloned() bool {
	repoDir, err := r.RepoDir()
	if err != nil {
		return false
	}

	if stat, _ := os.Stat(repoDir); stat != nil && stat.IsDir() {
		return true
	}
	return false
}

// git@github.com:xxx/abc-123.git
// result abc-123
func (r *Repository) projectName() (string, error) {
	cloneURI, err := ParseCloneURL(r.Url)
	if err != nil {
		return "", err
	}

	return cloneURI.RepoName, nil
}

//git rev-list HEAD --count
//
func (r *Repository) CommitCountByTagName(tagName string) (uint, error) {
	projectDir, err := r.ProjectDir()
	if err != nil {
		return 0, err
	}

	var stdOutBuff bytes.Buffer

	total := exec.Command("git", "rev-list", tagName, "--count")
	fmt.Println(total.Args)

	total.Dir = projectDir
	total.Stdout = &stdOutBuff
	err = total.Run()
	if err != nil {
		r.out.Write(stdOutBuff.Bytes())
		return 0, err
	}

	countStr := strings.TrimSpace(string(stdOutBuff.Bytes()))
	fmt.Println("Count: ", countStr)

	count, err := strconv.ParseUint(countStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}
