package major

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/yockii/grab12306/constant"

	netutil "github.com/yockii/grab12306/utils/net"
)

func openLoginPage() (err error) {
	client := netutil.GetMajorClient()
	req, err := http.NewRequest("GET", constant.Urls["loginPage"], nil)
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
		buffer.WriteString(constant.CaptchaPoints[strings.TrimSpace(v)])
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

// UamauthClientResponse 获取校验后的用户名返回信息
type UamauthClientResponse struct {
	ResultCode int    `json:"result_code"`
	Username   string `json:"username"`
}

func getUsername() (success bool) {
	client := netutil.GetMajorClient()
	data := make(url.Values)
	data["tk"] = []string{MyAccount.AppToken}
	res, err := client.PostForm(constant.Urls["uamauthclient"], data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	var uRes UamauthClientResponse
	err = json.Unmarshal(content, &uRes)
	if uRes.ResultCode == 0 {
		MyAccount.Username = uRes.Username
		success = true
		return
	}
	return
}

// AuthResponse 校验请求返回信息，主要获取token
type AuthResponse struct {
	ResultCode int    `json:"result_code"`
	AppToken   string `json:"newapptk"`
}

func auth() (success bool) {
	client := netutil.GetMajorClient()
	data := make(url.Values)
	data["appid"] = []string{"otn"}
	res, err := client.PostForm(constant.Urls["auth"], data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	var aRes AuthResponse
	err = json.Unmarshal(content, &aRes)
	if err != nil {
		fmt.Println(err)
		return
	}

	if aRes.ResultCode == 0 {
		fmt.Println("校验成功")
		success = true
		MyAccount.AppToken = aRes.AppToken
		MyAccount.Logined = true
		return
	}
	return
}

// BaseLoginResponse 登录返回
type BaseLoginResponse struct {
	ResultCode int    `json:"result_code"`
	Uamtk      string `json:"uamtk"`
}

func login(username, password string) (success bool) {
	client := netutil.GetMajorClient()

	data := make(url.Values)
	data["username"] = []string{username}
	data["password"] = []string{password}
	data["appid"] = []string{"otn"}

	res, err := client.PostForm(constant.Urls["baseLogin"], data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	var blRes BaseLoginResponse
	err = json.Unmarshal(content, &blRes)
	if err != nil {
		fmt.Println(err)
		return
	}

	if blRes.ResultCode == 0 {
		fmt.Println("登录成功")
		success = true
		return
	}
	return
}

// CheckCaptchaResponse 登录返回
type CheckCaptchaResponse struct {
	ResultCode    int    `json:"result_code"`
	ResultMessage string `json:"result_message"`
}

func captchaCheck(captchaAnswer string) (err error) {
	client := netutil.GetMajorClient()

	data := make(url.Values)
	data["answer"] = []string{captchaAnswer}
	data["rand"] = []string{"sjrand"}
	data["login_site"] = []string{"E"}

	res, err := client.PostForm(constant.Urls["checkCaptcha"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return
	}

	var ccRes CheckCaptchaResponse
	err = json.Unmarshal(content, &ccRes)

	if ccRes.ResultCode == 4 {
		fmt.Println("验证码校验通过!")
		return
	}
	err = errors.New(ccRes.ResultMessage)
	return
}
