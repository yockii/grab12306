package config

// Domain 站点主域名
const Domain = ""

// FetchCdnFromGithub 从Github上获取cdn列表
const FetchCdnFromGithub = true

// CdnFetchURI 获取cdn的地址
const CdnFetchURI = "https://api.github.com/repos/testerSunshine/12306/contents/cdn_list"

// Verify12306URI 12306验证用的地址
const Verify12306URI = "otn/zwdch/init"

// Server12306URL 12306客运服务地址
const Server12306URL = "kyfw.12306.cn"

// CdnDealCountPerTime 每次处理校验cdn的数量
const CdnDealCountPerTime = 100

// Urls 所有12306处理请求的的Url信息
var Urls = map[string]string{
	"loginPage":    "https://kyfw.12306.cn/otn/login/init",
	"loginCaptcha": "https://kyfw.12306.cn/passport/captcha/captcha-image",
	// "loginCaptcha": "http://localhost:8000/captcha-image.jpg",
	"checkCaptcha":  "https://kyfw.12306.cn/passport/captcha/captcha-check",
	"baseLogin":     "https://kyfw.12306.cn/passport/web/login",
	"auth":          "https://kyfw.12306.cn/passport/web/auth/uamtk",
	"uamauthclient": "https://kyfw.12306.cn/otn/uamauthclient",
	"passengers":    "https://kyfw.12306.cn/otn/passengers/query",
}

// CaptchaPoints 验证码坐标，后期改为通过随机数计算
var CaptchaPoints = map[string]string{
	"1": "41,36",
	"2": "107,41",
	"3": "196,46",
	"4": "247,46",
	"5": "43,102",
	"6": "123,123",
	"7": "181,117",
	"8": "273,113",
}
