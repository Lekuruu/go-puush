package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

const GitHubUsername = "Lekuruu"
const GitHubRepo = "go-puush"
const GitHubBranch = "main"

var client = github.NewClient(nil)
var requiredFolders = []string{"web/static", "web/templates"}

func ListDirectoryContents(path string) ([]*github.RepositoryContent, error) {
	_, contents, _, err := client.Repositories.GetContents(
		context.Background(),
		GitHubUsername,
		GitHubRepo,
		path,
		&github.RepositoryContentGetOptions{
			Ref: GitHubBranch,
		},
	)
	return contents, err
}

func DownloadDirectory(path string) error {
	contents, err := ListDirectoryContents(path)
	if err != nil {
		return err
	}
	os.MkdirAll(path, os.ModePerm)

	for _, content := range contents {
		switch *content.Type {
		case "dir":
			err := DownloadDirectory(*content.Path)
			if err != nil {
				return err
			}
		case "file":
			err := DownloadFile(*content.Path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DownloadFile(path string) error {
	u := &url.URL{
		Scheme: "https",
		Host:   "raw.githubusercontent.com",
		Path:   fmt.Sprintf("%s/%s/%s/%s", GitHubUsername, GitHubRepo, GitHubBranch, path),
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: %s", u.String(), resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
