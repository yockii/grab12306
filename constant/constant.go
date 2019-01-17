package constant

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
	"loginPage":                 "https://kyfw.12306.cn/otn/login/init",
	"loginCaptcha":              "https://kyfw.12306.cn/passport/captcha/captcha-image",
	"checkCaptcha":              "https://kyfw.12306.cn/passport/captcha/captcha-check",
	"baseLogin":                 "https://kyfw.12306.cn/passport/web/login",
	"auth":                      "https://kyfw.12306.cn/passport/web/auth/uamtk",
	"uamauthclient":             "https://kyfw.12306.cn/otn/uamauthclient",
	"passengers":                "https://kyfw.12306.cn/otn/passengers/query",
	"leftTicket":                "https://kyfw.12306.cn/otn/leftTicket/queryZ",
	"leftTicketCDN":             "https://%s/otn/leftTicket/queryZ",
	"submitOrder":               "https://kyfw.12306.cn/otn/leftTicket/submitOrderRequest",
	"initDcPage":                "https://kyfw.12306.cn/otn/confirmPassenger/initDc",
	"checkOrderInfo":            "https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo",
	"stationNames":              "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9090",
	"queueCount":                "https://kyfw.12306.cn/otn/confirmPassenger/getQueueCount",
	"confirmQueue":              "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueue",
	"queryOrder":                "https://kyfw.12306.cn/otn/confirmPassenger/queryOrderWaitTime",
	"autoSubmitOrder":           "https://kyfw.12306.cn/otn/confirmPassenger/autoSubmitOrderRequest",
	"queueCountAsync":           "https://kyfw.12306.cn/otn/confirmPassenger/getQueueCountAsync",
	"confirmSingleForQueueAsys": "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueueAsys",
	// 用于快速分析一系列url
	"quickAnalysis": "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc",
	"otnUrlPrefix":  "https://kyfw.12306.cn/otn/",
}

// PassengerTicketSeatTypes 乘客座位类型代码
var PassengerTicketSeatTypes = map[string]string{
	"特等座":  "P",
	"一等座":  "M",
	"二等座":  "O",
	"动卧":   "F",
	"硬座":   "1",
	"无座":   "1",
	"软座":   "2",
	"硬卧":   "3",
	"软卧":   "4",
	"高级软卧": "6",
	"商务座":  "9",
}

// SeatTypeCodes 坐席代码
var SeatTypeCodes = map[string]string{
	"商务座":  "SWZ",
	"特等座":  "TZ",
	"一等座":  "ZY",
	"二等座":  "ZE",
	"高级软卧": "GR",
	"软卧":   "RW",
	"硬卧":   "YW",
	"动卧":   "SRRB",
	"高级动卧": "YYRW",
	"软座":   "RZ",
	"硬座":   "YZ",
	"无座":   "WZ",
}

// PassengerTypes 乘客类型代码
var PassengerTypes = map[string]string{
	"成人": "ADULT",
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
