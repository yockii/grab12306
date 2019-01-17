package constant

import (
	"regexp"

	netutil "github.com/yockii/grab12306/utils/net"
)

func init() {
	content, err := netutil.Get(Urls["quickAnalysis"])
	if err != nil {
		return
	}
	htmlStr := string(content[:])

	// ticket查询url
	queryURLReg := regexp.MustCompile("var\\s+CLeftTicketUrl\\s*=\\s*'(.*)'")
	queryURLRes := queryURLReg.FindStringSubmatch(htmlStr)
	if len(queryURLRes) > 1 && queryURLRes[1] != "" {
		Urls["leftTicket"] = Urls["otnUrlPrefix"] + queryURLRes[1]
		Urls["leftTicketCDN"] = "https://%s/otn/" + queryURLRes[1]
	}

	// 登录验证码地址
	passportCaptchaURLReg := regexp.MustCompile("var\\s+passport_captcha\\s*=\\s*'(.*)'")
	passportCaptchaURLRes := passportCaptchaURLReg.FindStringSubmatch(htmlStr)
	if len(passportCaptchaURLRes) > 1 && passportCaptchaURLRes[1] != "" {
		Urls["loginCaptcha"] = passportCaptchaURLRes[1]
	}

	// 登录验证码校验地址
	passportCaptchaCheckURLReg := regexp.MustCompile("var\\s+passport_captcha_check\\s*=\\s*'(.*)'")
	passportCaptchaCheckURLRes := passportCaptchaCheckURLReg.FindStringSubmatch(htmlStr)
	if len(passportCaptchaCheckURLRes) > 1 && passportCaptchaCheckURLRes[1] != "" {
		Urls["checkCaptcha"] = passportCaptchaCheckURLRes[1]
	}

	// 登录接口地址
	passportLoginURLReg := regexp.MustCompile("var\\s+passport_login\\s*=\\s*'(.*)'")
	passportLoginURLRes := passportLoginURLReg.FindStringSubmatch(htmlStr)
	if len(passportLoginURLRes) > 1 && passportLoginURLRes[1] != "" {
		Urls["baseLogin"] = passportLoginURLRes[1]
	}

	// 登录验证接口地址
	passportAuthuamURLReg := regexp.MustCompile("var\\s+passport_authuam\\s*=\\s*'(.*)'")
	passportAuthuamURLRes := passportAuthuamURLReg.FindStringSubmatch(htmlStr)
	if len(passportAuthuamURLRes) > 1 && passportAuthuamURLRes[1] != "" {
		Urls["auth"] = passportAuthuamURLRes[1]
	}

	// 用户名接口地址uamauthclient
	passportAuthclientURLReg := regexp.MustCompile("var\\s+passport_authclient\\s*=\\s*'(.*)'")
	passportAuthclientURLRes := passportAuthclientURLReg.FindStringSubmatch(htmlStr)
	if len(passportAuthclientURLRes) > 0 && passportAuthclientURLRes[1] != "" {
		Urls["uamauthclient"] = Urls["otnUrlPrefix"] + passportAuthclientURLRes[1]
	}

}
