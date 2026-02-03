package main

import (
	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/updater"
	"github.com/colt3k/utils/updater/artifactory"
)

func main() {
	// test to see if working validly
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca))

	c := updater.Connection{
		Name:               "main",
		User:               "username",
		PassOrToken:        "user token goes here",
		URLPrefix:          "http://localhost:8081/artifactory/",
		Repository:         "go-release-local/",
		Path:               "tunler/",
		OnAvailable:        "http://localhost:8081",
		OnAvailableTimeout: 10,
		OnAvailableViaHTTP: true,
	}

	v := updater.Version{
		Version:   "v0.0.9",
		BuildDate: "1542223120",
	}

	cons := make([]updater.Connection, 0)
	cons = append(cons, c)
	artifactory.PerformUpdate("myappname", cons, v, true, false)
}
