package cmd

import (
	"github.com/google/go-github/github"
	"github.com/molizz/inclus/bin/utils"
	"github.com/spf13/cobra"

	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Download struct {
	Dir     string
	Version string
}

// 只支持github

func init() {
	RootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "d",
	Short: "传入版本号下载",
	Long:  "传入版本号, 下载对应的软件包",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("请传入版本号")
		}

		down := NewDownload(args[0])
		err := down.Process()
		if err != nil {
			return err
		}

		return nil
	},
}

func NewDownload(version string) *Download {
	dir := config.GetString("releases_dir")
	return &Download{
		Dir:     dir,
		Version: version,
	}
}

func (d *Download) Process() error {
	defedVersionMap := config.GetStringMap("definitions")
	for k, v := range defedVersionMap {
		repoCfg.Definitions[k] = v.(string)
	}

	downVersions, ok := utils.Ver.Versions[d.Version]
	if !ok {
		return errors.New("这个版本不在 " + VerionsFile + " 中")
	}

	var downGroup sync.WaitGroup

	for name, ver := range downVersions {
		url := defedVersionMap[name].(string)
		cloneURI, err := utils.ParseCloneURL(url)
		if err != nil {
			return err
		}

		go func(owner, repo, tag, storeDir string) {
			downGroup.Add(1)
			defer downGroup.Done()

			err := d.download(owner, repo, tag, storeDir)
			if err != nil {
				fmt.Println(owner, repo, tag)
				fmt.Println("download release is err: ", err)
			}
		}(cloneURI.Owner, cloneURI.RepoName, ver, d.Dir)
	}

	downGroup.Wait()
	time.Sleep(1 * time.Second) // 可能有打印无法输出

	return nil
}

func (d *Download) download(owner, repo string, tagName string, toDir string) error {
	// fmt.Println(repo, owner, tagName, toDir)
	client := http.DefaultClient
	client.Timeout = clientTimeout
	client.Transport = &github.BasicAuthTransport{
		Username: owner,
		Password: os.Getenv(GithubToken),
	}

	gh := github.NewClient(client)
	release, resp, err := gh.Repositories.GetReleaseByTag(ctx, owner, repo, tagName)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s/%s 获取tag信息错误, 错误码: %d", owner, repo, resp.StatusCode)
	}

	assets := release.Assets
	if len(assets) == 0 {
		return fmt.Errorf("%s/%s 的tag版本 %s 没有上传asset", owner, repo, tagName)
	}

	asset := assets[0]
	if asset.GetState() != "uploaded" {
		return fmt.Errorf("%s/%s 资源还未上传完成,当前状态:%s", owner, repo, asset.GetState())
	}

	downloadUrl := asset.BrowserDownloadURL
	if downloadUrl == nil {
		return fmt.Errorf("%s/%s 并没有找到下载地址", owner, repo)
	}

	filename := asset.GetName()
	filesize := asset.GetSize()

	downResp, err := client.Get(asset.GetBrowserDownloadURL())
	if err != nil {
		return err
	}
	defer downResp.Body.Close()

	path := filepath.Join(d.Dir, filename)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	dsize, err := io.Copy(file, downResp.Body)
	if err != nil {
		return err
	}

	if dsize != int64(filesize) {
		return fmt.Errorf("%s/%s 警告! 下载的文件大小异常. 源文件大小 %d, 下载的大小 %d", owner, repo, filesize, dsize)
	}

	return nil
}
