package utility

import (
	"WebPathScanner/core"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func SaveToFile() {
	filenameMid := time.Now().Format("2006-02-01")
	for target, results := range core.GlobalConfig.Result {
		if len(results) == 0 {
			continue
		}
		for num := 0; ; {
			if _, err := os.Stat(target + "_" + filenameMid + "_t" + strconv.Itoa(num) + ".log"); err == nil {
				// path exists
				num += 1
				continue
			} else if errors.Is(err, os.ErrNotExist) {
				// path not exist
				file, err := os.Create(target + "_" + filenameMid + "_t" + strconv.Itoa(num) + ".log")
				if err != nil {
					log.Fatal("[x]Unexpected error occurred while creating file")
				}
				for _, result := range results {
					file.WriteString(result + "\n")
				}
				break
			} else {
				log.Fatal("[x]Unexpected error occurred while saving to file")
			}
		}

	}
	log.Infoln("[√]Saving Completed Successfully!")
}

// TODO: Further Features in dev...
