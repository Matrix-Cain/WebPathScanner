package core

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type globalConfig struct {
	Target         string     // required
	ConfPath       string     //optional
	Proxy          string     //optional
	Save           bool       //optional default false
	FileName       string     //optional
	AutoCheck404   bool       //optional default true
	AutoMd5        mapset.Set //not required
	ConfigFilePath string     //optional
	RandomSleep    int        //optional default 0
	Mode           int        //default dict mode
	ThreadNum      int        //default 10
	TargetList     []string   //auto filled
	PayloadList    []string   //auto filled
	TargetType     string     //auto filled
}

var GlobalConfig = globalConfig{}
var Vipe *viper.Viper
var ProgressBar *progressbar.ProgressBar
var Mutex = &sync.Mutex{}

func LoadAllConfig() {
	loadConfigFromFile() // Load config from persistent config.toml
	globalConfigRegister()
	engineRegister()   // Set pool size for multi-threading
	targetRegister()   // Handle target param from input
	payloadRegister()  // Load payload
	progressRegister() // Determine the tasks number using [target(s)Number * payloadsNumber]
}

func loadConfigFromFile() {

	var confPath string
	if GlobalConfig.ConfPath != "" {
		confPath = GlobalConfig.ConfPath
	} else {
		confPath = "./config.toml"
	}
	content, err := ioutil.ReadFile(confPath)

	if err != nil {
		log.Fatal(err)
	}
	Vipe = viper.GetViper()
	Vipe.SetConfigType("toml")
	err = Vipe.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		log.Fatal(err)
	}

}

func globalConfigRegister() {
	GlobalConfig.AutoMd5 = mapset.NewSet()
	if Vipe.IsSet("General.AutoCheck404") {
		GlobalConfig.AutoCheck404 = Vipe.Get("General.AutoCheck404").(bool)
	} else {
		GlobalConfig.AutoCheck404 = true
	}

}

func targetRegister() { // All the target will be regarded as legal target after this process
	log.Infoln("[*] Initialize targets...")
	if GlobalConfig.TargetType == "url" {
		err := ParseTarget()
		if err != nil {
			log.Errorln(err) // Output Error Details
			log.Fatal("Invalid input in [-i], Example: -i [http://]target.com or 192.168.1.1[/24] or 192.168.1.1-192.168.1.100")
		}
	} else {
		file, err := os.Open(GlobalConfig.Target)
		defer file.Close()
		if err != nil {
			log.Fatalf("Error when opening file: %s", err)
		}
		fileScanner := bufio.NewScanner(file)
		// read line by line
		for fileScanner.Scan() {
			if fileScanner.Text() != "" {
				GlobalConfig.Target = fileScanner.Text()
				err := ParseTarget()
				if err != nil {
					log.Errorln(err) // Output Error Details
					log.Fatal("Invalid input in [-i], Example: -i [http://]target.com or 192.168.1.1[/24] or 192.168.1.1-192.168.1.100")
				}
			}
		}
		// handle first encountered error while reading
		if err := fileScanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
		}
	}
	// now filter repeated target
	targetFilter()

	//check number of target(s) in utility.GlobalConfig.TargetList
	if len(GlobalConfig.TargetList) == 0 {
		log.Fatal("[!] No targets found.Please load targets with [-i|-iF]")
	}
	log.Infoln("[√] Targets Initialized")
}

func engineRegister() {
	if GlobalConfig.ThreadNum > 200 || GlobalConfig.ThreadNum < 1 {
		log.Warnln("[*] Invalid input in [-t](range: 1 to 200), has changed to default(10)")
		GlobalConfig.ThreadNum = 10
	}
}

//payloadRegister Handle Scan Mode to apply different processing to the given urls and payloads
func payloadRegister() {
	switch GlobalConfig.Mode {
	case 0:
		{ // default mode 0 as dictionary iteration

			file, err := os.Open(viper.Get("VintageDictConfig.path").(string))
			defer file.Close()
			if err != nil {
				log.Fatalf("Error when opening file: %s", err)
			}
			fileScanner := bufio.NewScanner(file)
			// read line by line
			for fileScanner.Scan() {
				if fileScanner.Text() != "" {
					GlobalConfig.PayloadList = append(GlobalConfig.PayloadList, fileScanner.Text())
				}
			}
			// handle first encountered error while reading
			if err := fileScanner.Err(); err != nil {
				log.Fatalf("Error while reading file: %s", err)
			}
		}
	case 1:
		{ // mode 1 as fuzz mode
			flag := Vipe.Get("FuzzDictConfig.flag")
			if flag.(string) != "" {
				fuzzModeCheck(flag.(string))
			} else {
				log.Fatal("[!]Fuzz flag not present in config")
			}

			file, err := os.Open(viper.Get("FuzzDictConfig.path").(string))
			defer file.Close()
			if err != nil {
				log.Fatalf("Error when opening file: %s", err)
			}
			fileScanner := bufio.NewScanner(file)
			// read line by line
			for fileScanner.Scan() {
				if fileScanner.Text() != "" {
					GlobalConfig.PayloadList = append(GlobalConfig.PayloadList, fileScanner.Text())
				}
			}
			// handle first encountered error while reading
			if err := fileScanner.Err(); err != nil {
				log.Fatalf("Error while reading file: %s", err)
			}

		}
	}
}

func progressRegister() {
	msg := fmt.Sprintf("[!]Total %v Targets with %v Payloads in total", len(GlobalConfig.TargetList), len(GlobalConfig.PayloadList))
	log.Infoln(msg)
	ProgressBar = progressbar.NewOptions(len(GlobalConfig.TargetList)*len(GlobalConfig.PayloadList),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan]Progress🚀[reset] Scanning..."),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
			log.Infoln("[√]All Task Done")
		}),
		//progressbar.OptionShowIts(), //this option will cause flicker
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]━━[reset]",
			SaucerHead:    "[green]━━[reset]",
			SaucerPadding: " ",
			BarStart:      "",
			BarEnd:        "",
		}))

}

//TargetFilter filter the same target And do
func targetFilter() {
	tmp := make([]interface{}, len(GlobalConfig.TargetList))
	for i := range GlobalConfig.TargetList {
		tmp[i] = GlobalConfig.TargetList[i]
	}
	setTarget := mapset.NewSetFromSlice(tmp).ToSlice()
	filteredTargetList := make([]string, len(setTarget))
	for i := range setTarget {
		filteredTargetList[i] = setTarget[i].(string)
	}
	GlobalConfig.TargetList = filteredTargetList
}

func fuzzModeCheck(flag string) {
	illegal := false
	var illegalList []string
	for _, v := range GlobalConfig.TargetList {
		if !strings.Contains(v, flag) {
			illegal = true
			illegalList = append(illegalList, v)
		}
	}
	if illegal {
		for _, v := range illegalList {
			log.Warnln(v)
		}
		log.Fatal("[!]Illegal Input that not correspond to the fuzz flag")
	}
}
