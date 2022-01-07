package util

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/yockii/qscore/pkg/httputil"
)

var client *http.Client
var once sync.Once

func GetHttpClient() *http.Client {
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

func GetCookieHttpClient(jar httputil.CookieJar) *http.Client {
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(10 * time.Second),
		Jar:       jar,
	}
	return c
}
