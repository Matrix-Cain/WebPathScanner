package utility

import (
	"WebPathScanner/core"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Bruter(urlInput string) {
	// url initialize
	urlStructure, err := url.Parse(urlInput)
	if err != nil {
		log.Fatal("Invalid URL Input")
	}
	if urlStructure.Scheme == "" {
		urlInput = "http://" + urlInput
		urlStructure, err = url.ParseRequestURI(urlInput)
		if err != nil {
			log.Fatal("Invalid URL Input")
		}
	}
	if urlStructure.Scheme != "http" && urlStructure.Scheme != "https" {
		log.Fatal("Unsupported Scheme")
	}
	// Fix the trailing slash if not exist
	if !strings.HasSuffix(urlInput, "/") {
		urlInput += "/"
	}
	// Print Current Target
	log.Infof("[+] Current target: %s", urlInput)

	//doing prepend work
	if core.GlobalConfig.AutoCheck404 {
		log.Infoln("[*] Launching auto check 404")
		i := Inspector{target: urlInput}
		result, notfoundType := i.CheckThis()
		if notfoundType == Test404Md5 || notfoundType == Test404Ok {
			core.GlobalConfig.AutoMd5.Add(result)
		}
	}
	SubRun(urlInput)

}

func Worker(urlWithSlash string, client *http.Client, payloadIndex int) {
	req, _ := http.NewRequest("GET", urlWithSlash+core.GlobalConfig.PayloadList[payloadIndex], nil)
	if core.Vipe.Get("General.UserAgent").(string) != "" {
		req.Header.Set("User-Agent", core.Vipe.Get("General.UserAgent").(string))
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	}

	resp, err := client.Do(req)
	if err != nil {
		core.Mutex.Lock()
		core.ProgressBar.Add(1)
		core.Mutex.Unlock()
		return
	}
	defer resp.Body.Close()
	msg := responseHandler(resp, req)
	core.Mutex.Lock()
	if msg != "" {
		_ = core.ProgressBar.Clear()
		log.Infoln(msg)
	}
	core.ProgressBar.Add(1)
	core.Mutex.Unlock()

}

func loadDict(filePath string) []string {
	var payloadList []string
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	// read line by line
	for fileScanner.Scan() {
		if fileScanner.Text() != "" {
			payloadList = append(payloadList, strings.Trim(fileScanner.Text(), "\t"))
		}
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	return payloadList
}

func responseHandler(resp *http.Response, req *http.Request) string {
	if core.GlobalConfig.AutoCheck404 {
		body, _ := ioutil.ReadAll(resp.Body)
		md5 := md5.New()
		md5.Write(body)
		if core.GlobalConfig.AutoMd5.Contains(hex.EncodeToString(md5.Sum(nil))) {
			return ""
		}
	}
	if mapset.NewSetFromSlice(core.Vipe.Get("General.ResponseStatusCode").([]interface{})).Contains(strconv.Itoa(resp.StatusCode)) {
		msg := fmt.Sprintf("[%v] ", resp.StatusCode)
		msg += req.URL.String()
		//save result stripe the same
		return msg
	}
	return ""
}
