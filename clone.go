package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

var errInvalidURLFormat = errors.New("invalid repository URL format")

const usage = `Usage: clone <repository URL>
Example repository URLs:
 clone github.com/username/repo
 clone github.com/username/repo.git
 clone http://github.com/username/repo
 clone https://github.com/username/repo
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

func rooted(gopath string, pwd func() (string, error)) string {
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

func get(args string) (string, string) {
	url := strip(args)
	domain, account, repo, err := parts(url)
	if err != nil {
		return "Invalid repository URL format.", ""
	}

	dir := rooted(os.Getenv("GOPATH"), os.Getwd)
	path := filepath.Join(dir, domain, account, repo)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Repository already exists at %s, updating...\n", path)
		if err := updateRepository(path); err != nil {
			return fmt.Sprintf("Error updating repository: %s\n", err), path
		}
		return "Repository updated successfully!", path
	}

	repoURL := fmt.Sprintf("https://%s.git", url)
	if err := cloneRepository(repoURL, path); err != nil {
		return fmt.Sprintf("Error cloning repository: %#v\n", err.Error()), path
	}
	return "Repository cloned successfully!", path
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("get command is version:", version)
		fmt.Print(usage)
		return
	}

	args := os.Args[1]
	result, path := get(args)
	fmt.Println(result)
	clipboard.WriteAll(path)
}
