package main

import (
	"WebPathScanner/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})
	cmd.Execute()
}
