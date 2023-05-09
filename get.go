package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "strings"
)

func main() {
    repoURL := os.Args[1]
    repoParts := strings.Split(repoURL, "/")

    var gitHost, owner, repo string

    switch {
    case strings.Contains(repoURL, "github.com"):
        gitHost = "github.com"
        owner = repoParts[1]
        repo = repoParts[2]
    case strings.Contains(repoURL, "gitlab.com"):
        gitHost = "gitlab.com"
        owner = repoParts[1]
        if strings.Contains(repoURL, "/tree/") {
            repo = repoParts[3]
        } else {
            repo = strings.TrimSuffix(repoParts[2], ".git")
        }
    case strings.Contains(repoURL, "bitbucket.org"):
        gitHost = "bitbucket.org"
        owner = repoParts[1]
        repo = strings.TrimSuffix(repoParts[2], ".git")
    }

    if gitHost == "" || owner == "" || repo == "" {
        fmt.Println("Invalid repository URL")
        return
    }

    cloneURL := fmt.Sprintf("https://%s/%s/%s.git", gitHost, owner, repo)
    cmd := exec.Command("git", "clone", cloneURL)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Run()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Repository cloned successfully")
}

