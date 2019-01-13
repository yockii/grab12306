package major

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"strconv"

	"github.com/yockii/grab12306/captcha"

	netutil "github.com/yockii/grab12306/utils/net"

	"github.com/yockii/grab12306/config"
)

func getLoginCaptchaResult() (captchaResult string, err error) {
	imgContent, err := getLoginCaptcha()
	// captcha.RequestRuokuaiCaptchaResult("captcha.jpg")
	// captcha.RequestRuokuaiCaptchaResult(imgContent)
	captchaResult = captcha.RuokuaiCaptchaResult(imgContent)
	if captchaResult == "" {
		err = errors.New("验证码识别失败")
	}
	return
}

func getLoginCaptcha() (imgContent []byte, err error) {
	var urlBuffer bytes.Buffer
	urlBuffer.WriteString(config.Urls["loginCaptcha"])
	urlBuffer.WriteString("?login_site=E&module=login&rand=sjrand&")
	urlBuffer.WriteString(strconv.FormatFloat(rand.Float64(), 'f', 16, 64))

	client := netutil.GetMajorClient()

	res, err := client.Get(urlBuffer.String())
	if err != nil {
		return
	}
	defer res.Body.Close()
	imgContent, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	// fh, _ := os.Create("captcha.jpg")
	// defer fh.Close()
	// io.Copy(fh, bytes.NewReader(imgContent))
	return
}
