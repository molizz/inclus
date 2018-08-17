package cmd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/url"
	"strings"
)

const (
	clientTimeout = 10 * time.Second
)

var (
	ctx = context.Background()
)

// 目前只支持github

type Commit struct {
	gitUrl     string
	gitBranch  string
	gitToken   string
	gitPath    string
	configPath string
}

func init() {
	RootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit [CONFIG]",
	Short: "传入版本号",
	Long:  "传入配置文件, 默认使用当前目录的 " + ConfigName,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := prepareCommit(args...)
		if err != nil {
			return err
		}

		err = verifyCommit(config)
		if err != nil {
			return err
		}

		err = push(config)
		if err != nil {
			return err
		}
		return nil
	},
}

func push(c *Commit) error {
	client := http.DefaultClient
	client.Timeout = clientTimeout
	client.Transport = &github.BasicAuthTransport{
		Username: c.gitToken,
		Password: c.gitToken,
	}

	owner, repo, err := ownerAndPath(c.gitUrl)
	if err != nil {
		return err
	}

	var message = "from inclus"
	content, err := c.commitContent()
	if err != nil {
		return err
	}
	contentEncode, err := encodeBase64(content)
	if err != nil {
		return err
	}

	sha := sha1.New()
	sha.Write(content)
	shaString := fmt.Sprintf("%x", sha.Sum(nil))

	// github.CommitAuthor{}

	opt := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: contentEncode,
		SHA:     &shaString,
		Branch:  &c.gitBranch,
	}

	gh := github.NewClient(client)
	_, resp, err := gh.Repositories.UpdateFile(ctx, *owner, *repo, c.gitPath, opt)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	} else {
		return errors.New(resp.Status)
	}
}

func encodeBase64(src []byte) (dst []byte, err error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("encode base64 is error: src length is %d", length)
	}

	dst = make([]byte, base64.StdEncoding.EncodedLen(length))
	base64.StdEncoding.Encode(dst, src)
	return
}

func (c *Commit) commitContent() ([]byte, error) {
	content, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func ownerAndPath(gitUrl string) (owner *string, repo *string, err error) {
	uri, err := url.Parse(gitUrl)
	if err != nil {
		return
	}
	gitPath := strings.Trim(uri.Path, "/")

	ownerPaths := strings.Split(gitPath, "/")
	if len(ownerPaths) < 2 {
		return nil, nil, errors.New("not match owner and path")
	}

	owner = &ownerPaths[0]
	repo = &ownerPaths[1]
	return
}

func prepareCommit(args ...string) (*Commit, error) {
	c := &Commit{}
	if len(args) > 0 {
		c.configPath = args[0]
	} else {
		c.configPath = ConfigName
	}

	v, err := GetViper(c.configPath)
	if err != nil {
		return nil, err
	}

	var ok bool
	repository := v.Get("repository").(map[string]interface{})
	c.gitUrl, ok = repository["giturl"].(string)
	if !ok {
		return nil, errors.New("not found github repository url")
	}

	c.gitBranch, ok = repository["gitbranch"].(string)
	if !ok {
		return nil, errors.New("not found github repository branch")
	}

	c.gitPath, ok = repository["gitpath"].(string)
	if !ok {
		return nil, errors.New("not found github repository gitpath")
	}

	c.gitToken, ok = repository["token"].(string)
	if !ok {
		return nil, errors.New("not found github repository token")
	}

	return c, nil
}

func verifyCommit(c *Commit) error {
	if !fileExist(c.configPath) {
		return ErrConfigNotFount
	}

	return nil
}
