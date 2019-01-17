package major

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/yockii/grab12306/constant"
	netutil "github.com/yockii/grab12306/utils/net"

	"github.com/yockii/grab12306/domain"
)

// QuickSubmitTicketSuit 快速提交订单
func QuickSubmitTicketSuit(orderInfo *domain.OrderInfo, tti *domain.TrainTicketInfo) (orderID string, err error) {
	// 1、 autoSubmitOrder
	needCaptcha, trainLocation, keyCheckIsChange, leftTicketStr, err := autoSubmitOrder(orderInfo, tti)
	// 2、 queueCountAsync
	queueCountAsync(orderInfo, tti)
	// 3、 confirmSingleForQueueAsys
	confirmSingleForQueueAsys(needCaptcha, trainLocation, keyCheckIsChange, leftTicketStr, orderInfo, tti.PassengerTicketSeatCode)
	// 4、不断queryOrderWaitTime
	for i := 0; i < 20; i++ {
		waitCount, waitTime, orderID, _ := queryOrderWaitTime()
		if orderID != "" || waitCount == 0 || waitTime == -1 {
			break
		}
		time.Sleep(12 * time.Second)
	}
	return
}

func autoSubmitOrder(orderInfo *domain.OrderInfo, tti *domain.TrainTicketInfo) (needCaptcha bool, trainLocation, keyCheckIsChange, leftTicketStr string, err error) {
	data := make(url.Values)
	data["secretStr"] = []string{tti.Secret}
	data["train_date"] = []string{orderInfo.TravelDate.Format("2006-01-02")}
	data["tour_flag"] = []string{"dc"}
	data["purpose_codes"] = []string{orderInfo.QueryPassengerType}
	data["query_from_station_name"] = []string{constant.Stations[tti.FromStationCode].Name}
	data["query_to_station_name"] = []string{constant.Stations[tti.ToStationCode].Name}
	cancelFlag := "2"
	data["cancel_flag"] = []string{cancelFlag}
	bedLevelOrderNum := "000000000000000000000000000000"
	data["bed_level_order_num"] = []string{bedLevelOrderNum}
	data["passengerTicketStr"] = []string{generatePassengerTicket(tti.PassengerTicketSeatCode, orderInfo.Passengers)}
	data["oldPassengerStr"] = []string{generateOldPassenger(orderInfo.Passengers)}

	res, err := netutil.GetMajorClient().PostForm(constant.Urls["autoSubmitOrder"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var asoRes domain.AutoSubmitOrderResponse
	err = json.Unmarshal(content, &asoRes)
	if err != nil {
		return
	}
	needCaptcha = asoRes.Data.IfShowPassCode == "Y"
	// "result":
	// "H2#989E53A408B49D8370E33FB48CD0C790AE931FAEEE24172941CE2D7F#GceWlk%2FKqedREUUil0P4qgHQ%2BzoMb7TXgcPhUM55aELa9gWV#1",
	resultStrs := strings.Split(asoRes.Data.Result, "#")
	trainLocation = resultStrs[0]
	keyCheckIsChange = resultStrs[1]
	leftTicketStr = resultStrs[2]
	return
}

func queueCountAsync(orderInfo *domain.OrderInfo, tti *domain.TrainTicketInfo) {
	data := make(url.Values)
	data["train_date"] = []string{orderInfo.TravelDate.Format("Mon Jan 2 2006 15:04:05 GMT+0800 (中国标准时间)")}
	data["train_no"] = []string{tti.TrainNo}
	data["stationTrainCode"] = []string{tti.TrainCode}
	data["seatType"] = []string{tti.PassengerTicketSeatCode}
	data["fromStationTelecode"] = []string{tti.FromStationCode}
	data["toStationTelecode"] = []string{tti.ToStationCode}
	data["leftTicket"] = []string{tti.LeftTicketSecret}
	data["purpose_codes"] = []string{orderInfo.QueryPassengerType}
	data["_json_att"] = []string{""}

	res, err := netutil.GetMajorClient().PostForm(constant.Urls["queueCountAsync"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var qcaRes domain.GetQueueCountResponse
	err = json.Unmarshal(content, &qcaRes)
	if err != nil {
		return
	}
	return
}
func confirmSingleForQueueAsys(needCaptcha bool, trainLocation, keyCheckIsChange, leftTicketStr string, orderInfo *domain.OrderInfo, passengerTicketSeatType string) (err error) {
	data := make(url.Values)
	data["passengerTicketStr"] = []string{generatePassengerTicket(passengerTicketSeatType, orderInfo.Passengers)}
	data["oldPassengerStr"] = []string{generateOldPassenger(orderInfo.Passengers)}
	captchaResult := ""
	if needCaptcha {
		// 验证码打码
	}
	data["randCode"] = []string{captchaResult}
	data["purpose_codes"] = []string{orderInfo.QueryPassengerType}
	data["key_check_isChange"] = []string{keyCheckIsChange}
	data["leftTicketStr"] = []string{leftTicketStr}
	data["train_location"] = []string{trainLocation}
	data["choose_seats"] = []string{""}
	data["seatDetailType"] = []string{""}
	data["_json_att"] = []string{""}

	res, err := netutil.GetMajorClient().PostForm(constant.Urls["confirmSingleForQueueAsys"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var csfqaRes domain.ConfirmQueueResponse
	err = json.Unmarshal(content, &csfqaRes)
	if err != nil {
		return
	}
	return
}

func queryOrderWaitTime() (waitCount, waitTime int, orderID string, err error) {
	// queryOrder
	var bs bytes.Buffer
	bs.WriteString(constant.Urls["queryOrder"])
	bs.WriteString("?random=")
	bs.WriteString(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	bs.WriteString("&tourFlag=dc&_json_att=")

	res, err := netutil.GetMajorClient().Get(bs.String())
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	var qoRes domain.QueryOrderResponse
	err = json.Unmarshal(content, &qoRes)
	if err != nil {
		return
	}
	waitCount = qoRes.Data.WaitCount
	waitTime = qoRes.Data.WaitTime
	if qoRes.Data.OrderID != "" {
		orderID = string(qoRes.Data.OrderID)
	}
	return
}
