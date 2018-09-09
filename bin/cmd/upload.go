package cmd

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/molizz/inclus/bin/utils"
	"github.com/spf13/cobra"
)

type Upload struct {
	PushURL string `yaml:"push_url"`
	Branch  string `yaml:"branch"`
	Path    string `yaml:"push_to"`

	GithubToken string
}

const (
	clientTimeout = 10 * time.Second
)

var (
	ctx = context.Background()
)

// 目前只支持github

func init() {
	RootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "u",
	Short: "传入版本号",
	Long:  "传入配置文件, 默认使用当前目录的 " + ConfigFile,
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

func push(c *Upload) error {
	owner, repo, err := ownerAndPath(c.PushURL)
	if err != nil {
		return err
	}

	client := http.DefaultClient
	client.Timeout = clientTimeout
	client.Transport = &github.BasicAuthTransport{
		Username: *owner,
		Password: c.GithubToken,
	}

	var message = fmt.Sprintf("%s from inclus.", time.Now().Format(time.RFC3339))
	content, err := c.commitContent()
	if err != nil {
		return err
	}
	// contentEncode, err := encodeBase64(content)
	// if err != nil {
	// 	return err
	// }

	var shaBuff bytes.Buffer
	shaBuff.Write([]byte("blob"))
	shaBuff.Write([]byte(fmt.Sprintf(" %d", len(content))))
	shaBuff.Write([]byte{0})
	shaBuff.Write(content)

	sha := sha1.New()
	sha.Write(shaBuff.Bytes())
	shaString := fmt.Sprintf("%x", sha.Sum(nil))

	// github.CommitAuthor{}

	opt := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: content,
		SHA:     &shaString,
		Branch:  &c.Branch,
	}

	fmt.Println(*opt.Message, *opt.SHA)

	gh := github.NewClient(client)
	_, resp, err := gh.Repositories.UpdateFile(ctx, *owner, *repo, c.Path, opt)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 == 2 {
		fmt.Println("更新成功.")
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

func (c *Upload) commitContent() ([]byte, error) {
	content, err := ioutil.ReadFile(c.Path)
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

func prepareCommit(args ...string) (*Upload, error) {
	c := &Upload{}
	if len(args) > 0 {
		c.Path = args[0]
	} else {
		c.Path = ConfigFile
	}

	c.GithubToken = os.Getenv("TOKEN")

	v, err := utils.GetViper(c.Path)
	if err != nil {
		return nil, err
	}

	var ok bool
	repository := v.Get("upload").(map[string]interface{})
	c.PushURL, ok = repository["push_url"].(string)
	if !ok {
		return nil, errors.New("not found github repository url")
	}

	c.Branch, ok = repository["branch"].(string)
	if !ok {
		return nil, errors.New("not found github repository branch")
	}

	c.Path, ok = repository["path"].(string)
	if !ok {
		return nil, errors.New("not found github repository gitpath")
	}

	ok = len(c.GithubToken) > 0
	if !ok {
		return nil, errors.New("not found github repository token")
	}

	return c, nil
}

func verifyCommit(c *Upload) error {
	if !utils.FileExist(c.Path) {
		return ErrConfigNotFount
	}

	return nil
}
