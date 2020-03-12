package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var githubRepo, gomodjson, netrcFile *string
var verbose *bool
var resultsChan chan Result
var threadCount *int

func main() {

	githubRepo = flag.String("repo", "nil", "Repo to analyze")
	gomodjson = flag.String("gomod", "go.mod.json", "go.mod JSON file to analyze. Run 'go mod edit -json' to get the JSON for your go.mod")
	verbose = flag.Bool("verbose", false, "Verbosity Level")
	netrcFile = flag.String("netrc", "nil", "Makes Gocheckit scan the provided '.netrc' file in the user's home directory for login name and password.")
	threadCount = flag.Int("threads", 5, "Amount of threads to spawn.")
	flag.Parse()

	resultsChan = make(chan Result)

	fmt.Println("[+] Go Checkit")

	if *githubRepo != "nil" {
		checkSingleRepo(*githubRepo)
	}

	analyzeGoMod(*gomodjson, *threadCount)
}

func performDeprecationTests(module Require) (Result, error) {
	var result Result

	stringSlice := strings.Split(module.Path, "/")
	if len(stringSlice) != 3 {
		if *verbose {
			fmt.Println("[*] Encountered Github repo error", module.Path)
		}
		return result, nil
	}
	if strings.ToLower(stringSlice[0]) != "github.com" {
		if *verbose {
			fmt.Println("[*] Ignoring golang module:", module.Path)
		}
		return result, nil
	}

	topicsUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/topics", stringSlice[1], stringSlice[2])
	check, err := checkRepoForDeprecatedTopic(topicsUrl, module)
	if err != nil {
		if *verbose {
			fmt.Println("[ERROR]", err)
		}
	}
	if check {
		resultMessage := fmt.Sprintf("[!] Deprecated Topic Identified: %s", module.Path)
		if *verbose {
			fmt.Println(resultMessage)
		}
		result.TopicTest = resultMessage
	}

	readmeUrl := fmt.Sprintf("https://github.com/%s/%s/raw/master/README.md", stringSlice[1], stringSlice[2])
	resp, err := getRepoReadme(readmeUrl)
	if err != nil {
		if *verbose {
			fmt.Println("[ERROR]", err)
		}
		return result, err
	}

	if checkReadmeForDeprecated(string(resp)) {
		resultMessage := fmt.Sprintf("[!] Deprecated README Identified: %s", module.Path)
		if *verbose {
			fmt.Println(resultMessage)
		}
		result.ReadMeTest = resultMessage
	}
	return result, nil
}

func analyzeGoMod(filename string, threadCount int) {
	if !fileExists(filename) {
		fmt.Printf("[ERROR] %s does not exist\n", filename)
	}
	fmt.Println("[*] Loading Modules from:", filename)
	data, err := readBytesFromFile(filename)
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(1)
	}

	var goModJson GoMod
	err = json.Unmarshal(data, &goModJson)
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(1)
	}

	modules := goModJson.Require
	fmt.Printf("[*] %d Modules Loaded\n", len(modules))
	jobs := make(chan Require, len(modules))
	results := make(chan Result, len(modules))

	for w := 1; w <= threadCount; w++ {
		go worker(jobs, results)
	}

	for _, module := range modules {
		jobs <- module
	}
	close(jobs)
	for _, _ = range modules {
		result := <-results
		if len(result.ReadMeTest) > 0 {
			fmt.Println(result.ReadMeTest)
		}

		if len(result.TopicTest) > 0 {
			fmt.Println(result.TopicTest)
		}
	}
}

func worker(jobs <-chan Require, results chan<- Result) {
	for module := range jobs {
		result, _ := performDeprecationTests(module)
		results <- result
	}
}

// TODO y'all
func checkSingleRepo(url string) {
	resp, err := getRepoReadme(*githubRepo)
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	if checkReadmeForDeprecated(string(resp)) {
		fmt.Println("[!] Repo README.md contains DEPRECATED:", *githubRepo)
	}
}

// Example deprecated repo -> https://github.com/bsm/sarama-cluster
func checkRepoForDeprecatedTopic(url string, module Require) (bool, error) {
	if *verbose {
		fmt.Printf("[+] Checking Module '%s' for DEPRECATED topic\n", module.Path)
	}

	containsDeprecatedTopic := false

	resp, err := doGithubApiRequest(url)
	if err != nil {
		return containsDeprecatedTopic, err
	}

	if strings.Contains(strings.ToLower(resp), "deprecated") {
		containsDeprecatedTopic = true
	}
	return containsDeprecatedTopic, nil
}

func checkReadmeForDeprecated(readme string) bool {
	// fmt.Println("[+] Checking a README.md for DEPRECATED")
	containsDeprecated := false

	if strings.Contains(strings.ToLower(readme), "deprecated") {
		containsDeprecated = true
	}
	return containsDeprecated
}
