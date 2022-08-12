package utility

import (
	"WebPathScanner/core"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
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
	core.Mutex.Lock()
	log.Infof("[+] Current target: %s", strings.TrimRight(urlInput, "/"))
	core.Mutex.Unlock()
	//doing prepend work
	if core.GlobalConfig.AutoCheck404 {
		core.Mutex.Lock()
		log.Infoln("[*] Launching auto check 404")
		core.Mutex.Unlock()
		i := Inspector{target: urlInput}
		result, notfoundType := i.CheckThis()
		if notfoundType == Test404Md5 || notfoundType == Test404Ok {
			core.GlobalConfig.AutoMd5.Add(result)
		}
	}
	SubRun(urlInput)

}

func Worker(urlWithSlash string, client *http.Client, payloadIndex int) {
	craftedUrl := genTarget(urlWithSlash, payloadIndex) // Concat url with payload
	req, _ := http.NewRequest("GET", craftedUrl, nil)
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
		if core.GlobalConfig.Save {
			reqUrl, _ := url.Parse(urlWithSlash)
			core.GlobalConfig.Result[reqUrl.Host] = append(core.GlobalConfig.Result[reqUrl.Host], msg)
		}
	}
	core.ProgressBar.Add(1)
	core.Mutex.Unlock()

}

func genTarget(urlWithSlash string, payloadIndex int) string { // To make future feature ez to add since final payload will be generated here
	switch core.GlobalConfig.Mode {
	case 0:
		{
			return urlWithSlash + core.GlobalConfig.PayloadList[payloadIndex]
		}
	case 1:
		{
			return strings.ReplaceAll(strings.TrimRight(urlWithSlash, "/"), core.Vipe.Get("FuzzDictConfig.flag").(string), core.GlobalConfig.PayloadList[payloadIndex])
		}
	}
	return ""
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
