package utility

import (
	"WebPathScanner/core"
	"github.com/panjf2000/ants"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Run map every url from GlobalConfig.TargetList to pool
func Run() {
	log.Infoln("[+] Coroutine Initializing")
	defer ants.Release()

	runTimes := len(core.GlobalConfig.TargetList)
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(core.GlobalConfig.ThreadNum, func(url interface{}) {
		Bruter(url.(string))
		wg.Done()
	})
	defer p.Release()
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(core.GlobalConfig.TargetList[i])
	}
	wg.Wait()
}

// SubRun Map payloads from GlobalConfig.PayloadList to the url
func SubRun(urlWithSlash string) {
	defer ants.Release()
	client := GenerateHTTPClient()
	runTimes := len(core.GlobalConfig.PayloadList)
	var wg sync.WaitGroup
	type args struct {
		urlWithSlash string
		client       *http.Client
		payloadIndex int
	}
	p, _ := ants.NewPoolWithFunc(int(core.Vipe.Get("General.RequestLimit").(int64)), func(inputArg interface{}) {

		inputArgs := inputArg.(args)
		Worker(inputArgs.urlWithSlash, inputArgs.client, inputArgs.payloadIndex)
		wg.Done()
	})
	defer p.Release()
	delayTime := int(core.Vipe.Get("General.DelayTime").(int64))
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		if delayTime > 0 {
			time.Sleep(time.Duration(rand.Intn(delayTime)) * time.Second)
		}
		inputArg := args{urlWithSlash: urlWithSlash, client: &client, payloadIndex: i}
		_ = p.Invoke(inputArg)
	}
	wg.Wait()
}

func GenerateHTTPClient() http.Client {
	client := http.Client{Timeout: time.Second * 5}
	if core.GlobalConfig.Proxy != "" {
		proxyUrl, err := url.Parse(core.GlobalConfig.Proxy)
		if err != nil {
			log.Fatal("Illegal Proxy Address")
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	return client
}
