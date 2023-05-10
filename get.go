package main

import (
	"errors"
	"fmt"

	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var errInvalidURLFormat = errors.New("invalid repository URL format")

const usage = `Usage: get <repository URL>
Example repository URLs:
 get github.com/username/repo
 get github.com/username/repo.git
 get http://github.com/username/repo
 get https://github.com/username/repo
`

func strip(url string) string {
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	fmt.Println(url)
	return url
}

func parts(url string) (domain, account, repo string, err error) {
	ps := strings.Split(url, "/")
	if len(ps) < 3 {
		return "", "", "", errInvalidURLFormat
	}
	return ps[0], ps[1], strings.Join(ps[2:], "/"), nil
}

type Getwder interface {
	Getwd() (dir string, err error)
}

type GetwdFunc func() (string, error)

func (fn GetwdFunc) Getwd() (dir string, err error) {
	return fn()
}

func rooted(gopath string, pwd func() (dir string, err error)) string {
	dir := fmt.Sprintf("%s/%s", gopath, "src")
	if gopath == "" {
		fmt.Println("GOPATH environment variable is not set.")
		var err error
		dir, err = pwd()
		fmt.Printf("Using current directory %s as root dir. error: %v\n", dir, err)
	}
	return dir
}

func updateRepository(path string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func cloneRepository(repoURL, path string) error {
	cmd := exec.Command("git", "clone", repoURL, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		return
	}

	url := strip(os.Args[1])
	domain, account, repo, err := parts(url)
	if err != nil {
		fmt.Println("Invalid repository URL format.")
		return
	}

	dir := rooted(os.Getenv("GOPATH"), os.Getwd)
	path := filepath.Join(dir, domain, account, repo)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Repository already exists at %s, updating...\n", path)
		if err := updateRepository(path); err != nil {
			fmt.Printf("Error updating repository: %s\n", err)
			return
		}
		fmt.Println("Repository updated successfully!")
		return
	}

	repoURL := fmt.Sprintf("https://%s.git", url)
	if err := cloneRepository(repoURL, path); err != nil {
		fmt.Printf("Error cloning repository: %s\n", err)
		return
	}
	fmt.Println("Repository cloned successfully!")
}
