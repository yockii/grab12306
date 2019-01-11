package cdnutil

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/yockii/grab12306/config"
)

// Cdn 结构
type Cdn struct {
	// IP cdn的ip地址
	IP string
	// cdn 被12306命名的名字
	Name string
	// Verified 是否已校验
	Verified bool
	// Available 是否可用
	Available bool
	// 响应速度
	Speed int64
}

var client *http.Client
var once sync.Once

func getClient() *http.Client {
	once.Do(func() {
		//跳过证书验证
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// Dial: (&net.Dialer{
			// 	Timeout: 30 * time.Second,
			// }).Dial,
			// TLSHandshakeTimeout: 30 * time.Second,
		}
		client = &http.Client{
			Transport: tr,
			Timeout:   time.Duration(30 * time.Second),
		}
	})
	return client
}

// GetCdnList 从服务器获取Cdn列表
func GetCdnList() []Cdn {
	var r []Cdn
	if config.FetchCdnFromGithub {
		resp, err := http.Get(config.CdnFetchURI)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Fatal error ", err.Error())
		}
		cs := jsoniter.Get(content, "content").ToString()
		dbs, err := base64.StdEncoding.DecodeString(cs)
		if err != nil {
			panic(err)
		}

		ipsBytes := bytes.Split(dbs, []byte("\n"))
		// r = make([]Cdn, len(ipsBytes))
		for _, v := range ipsBytes {
			r = append(r, Cdn{
				IP: string(v),
			})
		}
	}
	return r
}

// VerifyCdnList 校验Cdn是否可用
func VerifyCdnList(cdnList []Cdn) int {
	var wg sync.WaitGroup
	var availableCount int32
	for _, v := range cdnList {
		cdn := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			if verifyCdn(&cdn) {
				atomic.AddInt32(&availableCount, 1)
			}
		}()
	}
	wg.Wait()
	return int(atomic.LoadInt32(&availableCount))
}

func verifyCdn(cdn *Cdn) bool {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/%s", cdn.IP, config.Verify12306URI), nil)
	if err != nil {
		panic(err)
	}
	req.Host = "kyfw.12306.cn"
	// req.Header.Add("Content-Type", "text/html;charset=utf-8")
	// req.Header.Set("Host", "kyfw.12306.cn")
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	// req.Header.Add("Accept", "*/*")

	t0 := time.Now()

	res, err := getClient().Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		elapsed := time.Since(t0)
		cdn.Available = true
		cdn.Speed = elapsed.Nanoseconds()
		via := res.Header.Get("X-Via")
		if len(via) > 1 {
			cdn.Name = strings.Split(via, " ")[1]
		}
	} else {
		cdn.Available = false
	}
	cdn.Verified = true
	return cdn.Available
}
