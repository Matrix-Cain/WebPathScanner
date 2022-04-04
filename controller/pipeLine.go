package controller

import (
	"WebPathScanner/core"
	"WebPathScanner/utility"
	log "github.com/sirupsen/logrus"
)

func PipeLine() {
	core.LoadAllConfig()
	utility.Run()
	// do post-processing staff
	log.Infoln("[*]Doing Post-Processing staff")
}
