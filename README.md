# Gocheckit
Gocheckit is a Go tool that can be used to help identify Go modules that are potentially Deprecated.

# What does it do?
Gocheckit uses your go.mod in JSON format and analyzes the modules to determine if they might be deprecated. It does the following:
- Check a Go modules Github repo for any Deprecated Topics
- Analyze the contents of a modules README for the presence of certain keywords i.e deprecated
- Utilise your .netrc Github creds so you aren't throttled and can access Private Repos
- It's multithreaded, if you're in a rush, use those cores!

# Installation

## Binaries
For installation instructions from binaries please visit the [Releases Page](https://github.com/brompwnie/gocheckit/releases).

## Via Go
```
go get github.com/brompwnie/gocheckit
```

# Building from source

Building Gocheckit via Go:
```
go build
```
Building Gocheckit via Make:
```
make
```

# Usage
Gocheckit can be compiled into a binary for the targeted platform and supports the following usage
```
Usage of ./gocheckit:
  -gomod string
        go.mod JSON file to analyze. Run 'go mod edit -json' to get the JSON for your go.mod (default "go.m
od.json")
  -netrc string
        Makes Gocheckit scan the provided '.netrc' file in the user's home directory for login name and pas
sword. (default "nil")
  -repo string
        Repo to analyze (default "nil")
  -threads int
        Amount of threads to spawn. (default 5)
  -verbose
        Verbosity Level
```

# Example

Simplest Usage

```
# go mod edit -json > go.mod.json
# ./gocheckit
[+] Go Checkit
[*] Loading Modules from: go.mod.json
[*] 11 Modules Loaded
[!] Deprecated README Identified: github.com/bsm/sarama-cluster
[!] Deprecated Topic Identified: github.com/bsm/sarama-cluster
[!] Deprecated README Identified: github.com/uudashr/go-module
```

Throw some threads at the problem
```
# go mod edit -json > go.mod.json
# ./gocheckit -threads=10
[+] Go Checkit
[*] Loading Modules from: go.mod.json
[*] 11 Modules Loaded
[!] Deprecated README Identified: github.com/bsm/sarama-cluster
[!] Deprecated Topic Identified: github.com/bsm/sarama-cluster
[!] Deprecated README Identified: github.com/uudashr/go-module
```

User your Github creds from .netrc for no throttling
```
# go mod edit -json > go.mod.json
# ./gocheckit -threads=10 -netrc=.netrc
[+] Go Checkit
[*] Loading Modules from: go.mod.json
[*] 11 Modules Loaded
[!] Deprecated README Identified: github.com/bsm/sarama-cluster
[!] Deprecated Topic Identified: github.com/bsm/sarama-cluster
[!] Deprecated README Identified: github.com/uudashr/go-module
```

# Issues, Bugs and Improvements
For any bugs, please submit an issue. There is a long list of improvements but please submit an Issue if there is something you want to see added to Gocheckit.

 # License
 Gocheckit is licensed under a Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License (http://creativecommons.org/licenses/by-nc-sa/4.0).
