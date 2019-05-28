package main

import (
	"github.com/colt3k/utils/updater"
	"github.com/colt3k/utils/updater/artifactory"
)

func main() {
	// test to see if working validly

	c := updater.Connection{
		Name:        "main",
		User:        "username",
		PassOrToken: "user token goes here",
		URLPrefix:   "http://localhost:8081/artifactory/",
		Repository:  "go-release-local/",
		Path:        "tunler/",
		OnAvailable: "http://localhost:8081",
		OnAvailableViaHTTP:true,
	}

	v := updater.Version{
		Version:   "v0.0.9",
		BuildDate: "1542223120",
	}

	cons := make([]updater.Connection, 0)
	cons = append(cons, c)
	artifactory.PerformUpdate("myappname", cons, v, true)
}
