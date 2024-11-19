package main

import "fmt"

var (
	Version   string
	BuildTime string
	GitBranch string
	GitCommit string
)

func PrintVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("BuildTime: %s\n", BuildTime)
	fmt.Printf("GitBranch: %s\n", GitBranch)
	fmt.Printf("GitCommit: %s\n", GitCommit)
}
