package constant

import "github.com/yockii/grab12306/internal/domain"

// Domain 站点主域名
const Domain = ""

// Verify12306URI 12306验证用的地址
const Verify12306URI = "otn/zwdch/init"

// Server12306URL 12306客运服务地址
const Server12306URL = "https://kyfw.12306.cn"

// CdnDealCountPerTime 每次处理校验cdn的数量
const CdnDealCountPerTime = 100

const (
	UrlCheckLogin      = "https://kyfw.12306.cn/otn/login/conf"
	UrlLoginPage       = "https://kyfw.12306.cn/otn/resources/login.html"
	UrlGetQrCode       = "https://kyfw.12306.cn/passport/web/create-qr64"
	UrlCheckQrCode     = "https://kyfw.12306.cn/passport/web/checkqr"
	UrlUserAuth        = "https://kyfw.12306.cn/passport/web/auth/uamtk"
	UrlUamAuth         = "https://exservice.12306.cn/excater/uamauthclient"
	UrlGetPassengers   = "https://kyfw.12306.cn/otn/passengers/query"
	UrlQueryLeftTicket = "https://kyfw.12306.cn/otn/leftTicket/query"
)

// Urls 所有12306处理请求的的Url信息
var Urls_Declared = map[string]string{

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

var (
	SeatTypeTeDengZuo   = &domain.SeatType{Name: "特等座", Index: 25, Code: "P"}
	SeatTypeYiDengZuo   = &domain.SeatType{Name: "一等座", Index: 31, Code: "M"}
	SeatTypeErDengZuo   = &domain.SeatType{Name: "二等座", Index: 30, Code: "O"}
	SeatTypeDongWo      = &domain.SeatType{Name: "动卧", Index: 33, Code: "F"}
	SeatTypeYingZuo     = &domain.SeatType{Name: "硬座", Index: 29, Code: "1"}
	SeatTypeWuZuo       = &domain.SeatType{Name: "无座", Index: 26, Code: "1"}
	SeatTypeRuanZuo     = &domain.SeatType{Name: "软座", Index: 24, Code: "2"}
	SeatTypeYingWo      = &domain.SeatType{Name: "硬卧", Index: 28, Code: "3"}
	SeatTypeRuanWo      = &domain.SeatType{Name: "软卧", Index: 23, Code: "4"}
	SeatTypeGaoJiRuanWo = &domain.SeatType{Name: "高级软卧", Index: 21, Code: "6"}
	SeatTypeShangWuZuo  = &domain.SeatType{Name: "商务座", Index: 32, Code: "9"}
)

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

// PassengerTypes 乘客类型代码
var PassengerTypes = map[string]string{
	"成人": "ADULT",
	"学生": "0X00",
}
