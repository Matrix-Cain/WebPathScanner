package controller

import (
	"WebPathScanner/core"
	"WebPathScanner/utility"
	log "github.com/sirupsen/logrus"
)

func PipeLine() {
	core.LoadAllConfig() // Doing pre-processing staff for later scanning
	utility.Run()        // Running scan task
	// do post-processing staff
	log.Infoln("[*]Doing Post-Processing staff")
	utility.SaveToFile()
}
