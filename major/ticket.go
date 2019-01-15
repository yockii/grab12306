package major

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yockii/grab12306/domain"

	"github.com/yockii/grab12306/constant"
	netutil "github.com/yockii/grab12306/utils/net"
)

/*
	订单提交步骤：
	1、submitOrder
	2、进入initDc (getDcHTML)
	3、解析html (ParseDcHTML)
	4、CheckOrderInfo
	5、GetQueueCount
	6、ConfirmQueue 提交订单
	7、轮询订票结果(入队不一定有票)
*/
// MakeOrder
// seatType 一等座 二等座 商务座。。。
// passengerType ADULT 表示成人票
func MakeOrder(tti *domain.TrainTicketInfo, passengers []*domain.Passenger, seatType, passengerType string) (success bool, err error) {
	startDate, err := time.Parse("20060102", tti.StartTrainDate)
	if err != nil {
		return
	}
	err = SubmitOrder(tti.Secret, startDate.Format("2006-01-02"), passengerType, constant.Stations[tti.FromStationCode].Name, constant.Stations[tti.ToStationCode].Name)
	if err != nil {
		return
	}
	initDcHTML, err := getDcHTML()
	if err != nil {
		return
	}
	repeatSubmitToken, ticketInfo, err := ParseDcHTML(initDcHTML)
	if err != nil {
		return
	}
	checkOrderInfoResponse, err := CheckOrderInfo(repeatSubmitToken, seatType, passengers, ticketInfo)
	if err != nil {
		return
	}
	fmt.Printf("订单信息确认, 选座%s, 选铺%s \n", checkOrderInfoResponse.Data.CanChooseSeats, checkOrderInfoResponse.Data.CanChooseBeds)
	queue, left, err := GetQueueCount(repeatSubmitToken, seatType, tti, ticketInfo)
	if err != nil {
		return
	}
	fmt.Printf("排队人数%d, 余票%d \n", queue, left)
	inQueue, err := ConfirmQueue(repeatSubmitToken, seatType, passengers, tti, ticketInfo)
	if err != nil {
		return
	}
	if !inQueue {
		fmt.Println("排队失败")
		return
	}

	for {
		waitCount, orderID, err := QueryOrder(repeatSubmitToken)
		if err != nil {
			break
		}
		if orderID != "" {
			success = true
			break
		}
		fmt.Printf("等待人数: %d", waitCount)
		time.Sleep(3 * time.Second)
	}
	return
}

// SubmitOrder 单程票提交订单
// secretStr
// startDate 2019-01-22
// passengerType ADULT-成人票
// fromStationName 出发站名 杭州东
// toStationName 到达站名 上海虹桥
func SubmitOrder(secretStr, startDate, passengerType, fromStationName, toStationName string) (err error) {
	data := make(url.Values)
	data["secretStr"] = []string{secretStr}
	data["train_date"] = []string{startDate}
	data["back_train_date"] = []string{time.Now().Format("2006-01-02")}
	data["tour_flag"] = []string{"dc"}
	data["purpose_codes"] = []string{passengerType}
	data["query_from_station_name"] = []string{fromStationName}
	data["query_to_station_name"] = []string{toStationName}
	data["undefined"] = []string{""}

	client := netutil.GetMajorClient()
	res, err := client.PostForm(constant.Urls["submitOrder"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	var soRes domain.SubmitOrderResponse
	err = json.Unmarshal(content, &soRes)
	if !soRes.Status {
		err = errors.New(soRes.Messages)
		return
	}
	return
}

// getDcHTML 获取订单页面html
func getDcHTML() (htmlContent []byte, err error) {
	data := make(url.Values)
	data["_json_att"] = []string{""}
	res, err := netutil.GetMajorClient().PostForm(constant.Urls["initDcPage"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()

	htmlContent, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}

// ParseDcHTML 解析单程订单页面元素获取需要的数据
func ParseDcHTML(htmlContent []byte) (repeatSubmitToken string, ticketInfo *domain.TicketInfoForPassengerForm, err error) {
	htmlString := string(htmlContent[:])
	globalRepeatSubmitTokenReg := regexp.MustCompile("var\\s+globalRepeatSubmitToken\\s*=\\s*'(.*)';")
	params := globalRepeatSubmitTokenReg.FindStringSubmatch(htmlString)
	if len(params) > 1 {
		repeatSubmitToken = params[1]
	}
	// ticketInfoForPassengerForm
	ticketInfoForPassengerFormReg := regexp.MustCompile("var\\s+ticketInfoForPassengerForm\\s*=\\s*(.*);")
	ticketInfoForPassengerFormSingle := ticketInfoForPassengerFormReg.FindStringSubmatch(htmlString)

	s := ticketInfoForPassengerFormSingle[1]
	s = strings.Replace(s, "'", "\"", -1)
	var ti domain.TicketInfoForPassengerForm
	err = json.Unmarshal([]byte(s), &ti)
	ticketInfo = &ti
	return
}

// CheckOrderInfo 检查订单信息
// seatType 一等座、二等座。。。。
func CheckOrderInfo(repeatSubmitToken, seatType string, passengers []*domain.Passenger, ticketInfo *domain.TicketInfoForPassengerForm) (coiRes domain.CheckOrderInfoResponse, err error) {
	data := make(url.Values)
	cancelFlag := "2"
	if ticketInfo.OrderRequestDTO.CancelFlag != "" {
		cancelFlag = string(ticketInfo.OrderRequestDTO.CancelFlag)
	}
	data["cancel_flag"] = []string{cancelFlag}
	bedLevelOrderNum := "000000000000000000000000000000"
	if ticketInfo.OrderRequestDTO.BedLevelOrderNum != "" {
		bedLevelOrderNum = string(ticketInfo.OrderRequestDTO.BedLevelOrderNum)
	}
	data["bed_level_order_num"] = []string{bedLevelOrderNum}
	data["passengerTicketStr"] = []string{generatePassengerTicket(seatType, passengers)}
	data["oldPassengerStr"] = []string{generateOldPassenger(passengers)}
	data["tour_flag"] = []string{"dc"}
	data["randCode"] = []string{""}     // 验证码(需要的话)
	data["whatsSelect"] = []string{"1"} // 是否常用联系人中选择
	data["_json_att"] = []string{""}
	data["REPEAT_SUBMIT_TOKEN"] = []string{repeatSubmitToken}

	res, err := netutil.GetMajorClient().PostForm(constant.Urls["checkOrderInfo"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &coiRes)
	return
}

// GetQueueCount 获取余票及排队信息
func GetQueueCount(repeatSubmitToken, seatType string, tti *domain.TrainTicketInfo, ticketInfo *domain.TicketInfoForPassengerForm) (queueCount, leftTicket int, err error) {
	// queueCount
	data := make(url.Values)
	trainDateStr := tti.StartTrainDate
	trainDate, err := time.ParseInLocation("20060102", trainDateStr, time.Local)
	data["train_date"] = []string{trainDate.Format("Mon Jan 2 2006 15:04:05 GMT+0800 (中国标准时间)")}
	data["train_no"] = []string{tti.TrainNo}
	data["stationTrainCode"] = []string{tti.TrainCode}
	data["seatType"] = []string{constant.PassengerTicketSeatTypes[seatType]}
	data["fromStationTelecode"] = []string{tti.FromStationCode}
	data["toStationTelecode"] = []string{tti.ToStationCode}
	data["leftTicket"] = []string{tti.LeftTicketSecret}
	data["purpose_codes"] = []string{string(ticketInfo.PurposeCodes)}
	data["train_location"] = []string{tti.TrainLocationCode}
	data["_json_att"] = []string{""}
	data["REPEAT_SUBMIT_TOKEN"] = []string{repeatSubmitToken}

	res, err := netutil.GetMajorClient().PostForm(constant.Urls["queueCount"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var gqcRes domain.GetQueueCountResponse
	err = json.Unmarshal(content, &gqcRes)
	if err != nil {
		return
	}
	if !gqcRes.Status {
		err = errors.New(strings.Join(gqcRes.Messages, ","))
		return
	}
	fmt.Printf("当前余票%s, 当前排队人数%s\n", gqcRes.Data.Ticket, gqcRes.Data.Count)
	queueCount, _ = strconv.Atoi(gqcRes.Data.Count)
	leftTicket, _ = strconv.Atoi(gqcRes.Data.Ticket)
	return
}

// ConfirmQueue 确认队列信息
func ConfirmQueue(repeatSubmitToken, seatType string, passengers []*domain.Passenger, tti *domain.TrainTicketInfo, ticketInfo *domain.TicketInfoForPassengerForm) (inQueue bool, err error) {
	// confirmQueue
	data := make(url.Values)
	data["passengerTicketStr"] = []string{generatePassengerTicket(seatType, passengers)}
	data["oldPassengerStr"] = []string{generateOldPassenger(passengers)}
	data["randCode"] = []string{""}
	data["purpose_codes"] = []string{string(ticketInfo.PurposeCodes)}
	data["key_check_isChange"] = []string{string(ticketInfo.KeyCheckIsChange)}
	data["leftTicketStr"] = []string{string(ticketInfo.LeftTicketStr)}
	data["train_location"] = []string{tti.TrainLocationCode}
	data["choose_seats"] = []string{""} // 选座
	seatDetailType := "000"
	if ticketInfo.OrderRequestDTO.SeatDetailTypeCode != "" {
		seatDetailType = string(ticketInfo.OrderRequestDTO.SeatDetailTypeCode)
	}
	data["seatDetailType"] = []string{seatDetailType}
	data["whatsSelect"] = []string{"1"} // 是否常用联系人中选择
	data["dwAll"] = []string{"N"}       // 动卧？
	data["_json_att"] = []string{""}
	data["REPEAT_SUBMIT_TOKEN"] = []string{repeatSubmitToken}

	res, err := netutil.GetMajorClient().PostForm(constant.Urls["confirmQueue"], data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var cqRes domain.ConfirmQueueResponse
	err = json.Unmarshal(content, &cqRes)
	inQueue = cqRes.Data.SubmitStatus
	return
}

func generatePassengerTicket(seatType string, passengers []*domain.Passenger) string {
	var bf bytes.Buffer
	for _, p := range passengers {
		bf.WriteString(constant.PassengerTicketSeatTypes[seatType])
		bf.WriteString(",0,")
		bf.WriteString(p.PassengerType)
		bf.WriteString(",")
		bf.WriteString(p.PassengerName)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDTypeCode)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDNo)
		bf.WriteString(",")
		if p.PhoneNo != "" {
			bf.WriteString(p.PhoneNo)
		}
		bf.WriteString(",N_")
	}
	s := bf.String()
	return s[:len(s)-1]
}

func generateOldPassenger(passengers []*domain.Passenger) string {
	var bf bytes.Buffer
	for _, p := range passengers {
		bf.WriteString(p.PassengerName)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDTypeCode)
		bf.WriteString(",")
		bf.WriteString(p.PassengerIDNo)
		bf.WriteString(",")
		bf.WriteString(p.PassengerType)
		bf.WriteString("_")
	}
	return bf.String()
}

// QueryOrder 查询订单信息
func QueryOrder(repeatSubmitToken string) (waitCount int, orderId string, err error) {
	// queryOrder
	var bs bytes.Buffer
	bs.WriteString(constant.Urls["queryOrder"])
	bs.WriteString("?random=")
	bs.WriteString(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	bs.WriteString("&tourFlag=dc&_json_att=&REPEAT_SUBMIT_TOKEN=")
	bs.WriteString(repeatSubmitToken)

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
	if qoRes.Data.OrderID != "" {
		orderId = string(qoRes.Data.OrderID)
	}
	return
}
