package utils

import (
	"errors"
	"strings"
)

type CloneURLParser struct {
	rawURL string

	User     string
	Domain   string
	Owner    string
	RepoName string
}

func ParseCloneURL(url string) (*CloneURLParser, error) {

	c := &CloneURLParser{
		rawURL: url,
	}

	f := func(c rune) bool {
		switch c {
		case '@', ':', '/':
			return true
		default:
			return false
		}
	}
	rawFields := strings.FieldsFunc(c.rawURL, f)

	if len(rawFields) < 4 {
		return nil, errors.New("url不合法")
	}

	c.User = rawFields[0]
	c.Domain = rawFields[1]
	c.Owner = rawFields[2]
	c.RepoName = func() string {
		rawRepoName := rawFields[3]
		return rawRepoName[:len(rawRepoName)-len(".git")]
	}()

	return c, nil
}
