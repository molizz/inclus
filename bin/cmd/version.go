package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	DefConfigName = "includes.yml"
)

var (
	definitions = make(map[string]interface{})
	versions    = make(map[string]interface{})
)

type Download struct {
	name string
	url  string
}

type Command struct {
	version string
	config  string
}

func initConfig(c *Command) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	viper.SetConfigType("yaml")
	body, err := ioutil.ReadFile(c.config)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	definitions = viper.Get("definitions").(map[string]interface{})
	versions = viper.Get("versions").(map[string]interface{})
	return nil
}

func buildDownloads(ver string) ([]*Download, error) {
	version, ok := versions[ver]
	if !ok {
		return nil, errors.New("not found " + ver)
	}

	// 该版本下的所有的依赖
	versionsMap := version.(map[string]interface{})

	// build url
	downloads := make([]*Download, 0)

	for defName, verName := range versionsMap {
		defKeyMap, ok := definitions[defName].(map[string]interface{})
		if !ok {
			return nil, errors.New("wrong format")
		}
		durl := defKeyMap["url"].(string)
		durl = strings.Replace(durl, "{{version}}", verName.(string), -1)

		name, ok := defKeyMap["filename"].(string)
		if !ok {
			urls := strings.Split(durl, "/")
			if len(urls) > 0 {
				name = urls[len(urls)-1]
			}
		}

		downloads = append(downloads, &Download{
			url:  durl,
			name: name,
		})
	}

	return downloads, nil
}

func download(urls []*Download) error {
	if len(urls) == 0 {
		return errors.New("urls is null.")
	}

	downloadFunc := func(download *Download) error {
		httpClient := http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := httpClient.Get(download.url)
		if err != nil {
			fmt.Println("url: ", download.url)
			return err
		}
		defer resp.Body.Close()

		file, err := os.OpenFile(download.name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}

		return nil
	}

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(durl *Download) {
			defer wg.Done()
			err := downloadFunc(durl)
			if err != nil {
				fmt.Println("downloading error:", err, " url: ", durl.url)
			}
		}(u)
	}

	fmt.Println("downloading...")

	wg.Wait()

	fmt.Println("done... bye.")

	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version VERSION_NAME CONFIG",
	Short: "传入版本号",
	Long:  "传入配置文件, 默认使用当前目录的includes.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := verify(args...)
		if err != nil {
			return err
		}

		c := prepare(args...)
		err = initConfig(c)
		if err != nil {
			return err
		}

		urls, err := buildDownloads(c.version)
		if err != nil {
			return err
		}

		if len(urls) == 0 {
			return errors.New("Not match download url.")
		}

		err = download(urls)
		if err != nil {
			return err
		}

		return nil
	},
}

func prepare(args ...string) *Command {
	c := &Command{}
	c.version = args[0]
	if len(args) >= 2 {
		c.config = args[1]
	} else {
		c.config = DefConfigName
	}
	return c
}

func verify(args ...string) error {
	if len(args) == 0 {
		return errors.New("version is not set.")
	}

	existFunc := func(path string) error {
		configInfo, err := os.Stat(path)
		if os.IsNotExist(err) || configInfo.IsDir() {
			return errors.New("config file is not exist.")
		}
		return nil
	}

	if len(args) < 2 {
		err := existFunc(DefConfigName)
		return err
	}

	err := existFunc(args[1])
	return err
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
