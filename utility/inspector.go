package utility

import (
	"WebPathScanner/core"
	"crypto/md5"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	Test404Ok     = 0
	Test404Md5    = 1
	Test404String = 2
	Test404Url    = 3
	Test404None   = 4
	DictMode      = 0
	FuzzMode      = 1
	CrawlMode     = 2
)

type Inspector struct {
	target string
	client http.Client
}

func (inspector Inspector) tryOnce() map[string]string {
	var randomPath string
	for i := 0; i < 42; i++ {
		randChar := string(rune(rand.Intn(25) + 97))
		randomPath += randChar
	}
	inspector.target += randomPath
	log.Infof("[+] Checking with: %s\n", inspector.target)

	inspector.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	// configure client
	if core.GlobalConfig.Proxy != "" {
		proxyUrl, err := url.Parse(core.GlobalConfig.Proxy)
		if err != nil {
			log.Fatal("Illegal Proxy Address")
		}
		inspector.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	inspector.client.Timeout = time.Second * 5
	req, err := http.NewRequest("GET", inspector.target, nil)
	if err != nil {
		log.Fatal(err)
	}
	if core.Vipe.Get("General.UserAgent").(string) != "" {
		req.Header.Set("User-Agent", core.Vipe.Get("General.UserAgent").(string))
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	}

	resp, err := inspector.client.Do(req)
	//https://stackoverflow.com/questions/53992694/what-does-netloc-mean
	/*From RFC 1808, Section 2.1, every URL should follow a specific format:

	<scheme>://<netloc>/<path>;<params>?<query>#<fragment>
	Lets break this format down syntactically:

	scheme: The protocol name, usually http/https
	netloc: Contains the network location - which includes the domain itself (and subdomain if present), the port number, along with an optional credentials in form of username:password. Together it may take form of username:password@domain.com:80.
	path: Contains information on how the specified resource needs to be accessed.
	params: Element which adds fine tuning to path. (optional)
	query: Another element adding fine grained access to the path in consideration. (optional)
	fragment: Contains bits of information of the resource being accessed within the path. (optional)
	Lets take a very simple example to understand the above clearly:

	https://cat.com/list;meow?breed=siberian#pawsize
	In the above example:

	https is the scheme (first element of a URL)
	cat.com is the netloc (sits between the scheme and path)
	/list is the path (between the netloc and params)
	meow is the param (sits between path and query)
	breed=siberian is the query (between the fragment and params)
	pawsize is the fragment (last element of a URL)
	*/

	if err != http.ErrUseLastResponse && err != nil {
		result := map[string]string{
			"target":   req.Host,
			"code":     "",
			"size":     "",
			"md5":      "",
			"content":  "",
			"location": "None"}
		return result
	} else if err == http.ErrUseLastResponse {
		body, _ := ioutil.ReadAll(resp.Body)
		md5 := md5.New()
		md5.Write(body)
		result := map[string]string{
			"target":   req.Host,
			"code":     strconv.Itoa(resp.StatusCode),
			"size":     strconv.Itoa(len(body)),
			"md5":      hex.EncodeToString(md5.Sum(nil)),
			"content":  string(body),
			"location": resp.Header.Get("Location")}
		return result
	} else { // err = nil
		body, _ := ioutil.ReadAll(resp.Body)
		md5 := md5.New()
		md5.Write(body)
		result := map[string]string{
			"target":   req.Host,
			"code":     strconv.Itoa(resp.StatusCode),
			"size":     strconv.Itoa(len(body)),
			"md5":      hex.EncodeToString(md5.Sum(nil)),
			"content":  string(body),
			"location": "None"}
		return result
	}
}

func (inspector Inspector) CheckThis() (string, int) {
	//Get the a request and decide what to do
	firstResult := inspector.tryOnce()

	if firstResult["code"] == "404" {
		return "", Test404Ok
	} else if firstResult["code"] == "302" || firstResult["location"] != "None" { // It seems that dirmap author just ignore this output maybe a 302 redirect just means the none-existence of the uri
		location := firstResult["location"]
		return location, Test404Url
	} else { //may be a fake 200 ok
		return firstResult["md5"], Test404Md5
	}

	return "", Test404None
}
