package major

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/yockii/grab12306/config"
	netutil "github.com/yockii/grab12306/utils/net"
)

func openLoginPage() (err error) {
	client := netutil.GetMajorClient()
	req, err := http.NewRequest("GET", config.Urls["loginPage"], nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return
}

// LoginSuit 登录的一系列处理
func LoginSuit(username, password string) (isPass bool, err error) {
	err = openLoginPage()
	if err != nil {
		return
	}
	platformCaptchaAnswer, err := getLoginCaptchaResult()
	fmt.Println("验证码识别完成，识别结果:", platformCaptchaAnswer)
	selected := strings.Split(platformCaptchaAnswer, "")
	var buffer bytes.Buffer
	for i, v := range selected {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(config.CaptchaPoints[strings.TrimSpace(v)])
	}
	captchaAnswer := buffer.String()
	fmt.Println("识别结果处理完毕，提交的验证码答案坐标:", captchaAnswer)
	err = captchaCheck(captchaAnswer)
	if err != nil {
		return
	}
	// 验证码通过验证
	fmt.Println("验证码校验通过，开始登录")
	// 登录
	success := login(username, password)
	if !success {
		err = errors.New("baseLogin登录失败")
		return
	}
	fmt.Println("登录成功，开始获取token")
	success = auth()
	if !success {
		err = errors.New("登录验证失败")
		return
	}
	// 获取到最新的app token后，获取用户姓名
	fmt.Println("token获取成功，开始获取用户名")
	isPass = getUsername()
	return
}

func getUsername() (success bool) {
	client := netutil.GetMajorClient()
	data := make(url.Values)
	data["tk"] = []string{MyAccount.AppToken}
	res, err := client.PostForm(config.Urls["uamauthclient"], data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	resultCode := jsoniter.Get(content, "result_code").ToInt()
	if resultCode == 0 {
		MyAccount.Username = jsoniter.Get(content, "username").ToString()
		success = true
		return
	}
	return
}

func auth() (success bool) {
	client := netutil.GetMajorClient()
	data := make(url.Values)
	data["appid"] = []string{"otn"}
	res, err := client.PostForm(config.Urls["auth"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	resultCode := jsoniter.Get(content, "result_code").ToInt()
	if resultCode == 0 {
		fmt.Println("校验成功")
		success = true
		newapptk := jsoniter.Get(content, "newapptk").ToString()
		MyAccount.AppToken = newapptk
		MyAccount.Logined = true
		return
	}
	return
}

func login(username, password string) (success bool) {
	client := netutil.GetMajorClient()

	data := make(url.Values)
	data["username"] = []string{username}
	data["password"] = []string{password}
	data["appid"] = []string{"otn"}

	res, err := client.PostForm(config.Urls["baseLogin"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	resultCode := jsoniter.Get(content, "result_code").ToInt()
	if resultCode == 0 {
		fmt.Println("登录成功")
		success = true
		// uamtk = jsoniter.Get(content, "uamtk").ToString()
		return
	}
	return
}

func captchaCheck(captchaAnswer string) (err error) {
	// client := netutil.GetMajorClient()
	// u, err := url.Parse()
	// if err != nil {
	// 	return
	// }
	// q := u.Query()
	// // q.Set("callback", "jQuery19105351335444165837_1547224186345")
	// q.Set("answer", captchaAnswer)
	// q.Set("rand", "sjrand")
	// q.Set("login_site", "E")

	// req, err := http.NewRequest("POST", config.Urls["checkCaptcha"], nil)
	// if err != nil {
	// 	return
	// }
	// req.Header.Set("Referer", "https://kyfw.12306.cn/otn/login/init")
	// res, err := client.Do(req)

	client := netutil.GetMajorClient()

	data := make(url.Values)
	data["answer"] = []string{captchaAnswer}
	data["rand"] = []string{"sjrand"}
	data["login_site"] = []string{"E"}

	res, err := client.PostForm(config.Urls["checkCaptcha"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return
	}
	resultCode := jsoniter.Get(content, "result_code").ToInt()

	if resultCode == 4 {
		fmt.Println("验证码校验通过!")
		return
	}
	resultMsg := jsoniter.Get(content, "result_message").ToString()
	fmt.Println("验证码校验失败!", resultMsg)
	err = errors.New(resultMsg)
	return
}