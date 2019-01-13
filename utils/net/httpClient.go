package netutil

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

var client *http.Client
var once sync.Once
var majorClient *http.Client

// GetMajorClient 获取主服务客户端(拥有cookies，主要负责登录账户，获取信息及订票)
func GetMajorClient() *http.Client {
	once.Do(func() {
		//跳过证书验证
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// Dial: (&net.Dialer{
			// 	Timeout: 30 * time.Second,
			// }).Dial,
			// TLSHandshakeTimeout: 30 * time.Second,
		}
		jar, _ := cookiejar.New(nil)

		majorClient = &http.Client{
			Transport: tr,
			Timeout:   time.Duration(30 * time.Second),
			Jar:       jar,
		}
	})
	return majorClient
}

// GetClient 获取http客户端直接处理请求
func GetClient() *http.Client {
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

// Get 通过Get请求数据
func Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Host = "kyfw.12306.cn"

	res, err := GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Post 通过post请求获取数据
func Post(url string, payload string) ([]byte, error) {
	return nil, errors.New("Not Implement")
}
