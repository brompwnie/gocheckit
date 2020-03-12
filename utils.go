package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/jdxcode/netrc"
)

type Result struct {
	TopicTest  string
	ReadMeTest string
}

//Structs below are for parsing go.mod.json files
type Topics struct {
	Names []string
}

// https://golang.org/src/cmd/go/internal/modcmd/edit.go
type Module struct {
	Path    string
	Version string
}

type GoMod struct {
	Module  Module
	Go      string
	Require []Require
	Exclude []Module
	Replace []Replace
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

func doGithubApiRequest(url string) (string, error) {
	if *verbose {
		fmt.Println("[*] Making Github API Request", url)
	}
	client := &http.Client{Timeout: time.Second * 7}
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/vnd.github.mercy-preview+json")

	if *netrcFile != "nil" {
		if *verbose {
			fmt.Println("[*] Loading netrc creds from:", *netrcFile)
		}
		usr, err := user.Current()
		if err != nil {
			fmt.Println("[ERROR]", err)
		}
		n, err := netrc.Parse(filepath.Join(usr.HomeDir, *netrcFile))
		if err != nil {
			fmt.Println("[ERROR]", err)
		}
		req.SetBasicAuth(n.Machine("api.github.com").Get("user"), n.Machine("api.github.com").Get("password"))
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readBytesFromFile(sourcefile string) ([]byte, error) {
	file, err := os.Open(filepath.Clean(sourcefile))
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func getRepoReadme(url string) ([]byte, error) {
	if *verbose {
		fmt.Println("[*] Downloading file from:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Error: Return code was not 200 for:%s\n", url))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
