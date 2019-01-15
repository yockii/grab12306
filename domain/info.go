package domain

import "encoding/json"

// TrainTicketInfo 每个车次的余票情况
type TrainTicketInfo struct {
	Secret            string // 密钥
	TrainNo           string // 车号
	TrainCode         string // 车次
	StartStationCode  string // 起始站
	EndStationCode    string // 终点站
	FromStationCode   string // 出发站
	ToStationCode     string // 到达站
	StartTime         string // 出发时间
	ArriveTime        string // 到达时间
	Duration          string // 历时
	CanBuy            bool   // 是否允许购买
	StartTrainDate    string // 出发日期 20190201
	SeniorSoftSleeper string // 高级软卧
	OtherSeat         string // 其他席位
	SoftSleeper       string // 软卧
	SoftSeat          string // 软座
	NoSeat            string // 无座
	HardSleeper       string // 硬卧
	HardSeat          string // 硬座
	SecondClassSeat   string // 二等座
	FirstClassSeat    string // 一等座
	BusinessClassSeat string // 商务特等座
	MCSleeper         string // 动卧
	TrainLocationCode string // 余票及排队信息查询时用到的TrainLocation
	LeftTicketSecret  string // 余票及排队信息查询时用到的LeftTicket

	Empty bool // 空数据
}

// TicketInfoForPassengerForm 从initDc解析出来的ticketInfoForPassengerForm变量json
type TicketInfoForPassengerForm struct {
	CardTypes             []CardOrSeatType `json:"cardTypes"`
	IsAsync               Nstring          `json:"isAsync"`
	KeyCheckIsChange      Nstring          `json:"key_check_isChange"`
	LeftDetails           []Nstring        `json:"leftDetails"`
	LeftTicketStr         Nstring          `json:"leftTicketStr"`
	LimitBuySeatTicketDTO struct {
		SeatTypeCodes     []CardOrSeatType `json:"seat_type_codes"`
		TicketSeatCodeMap struct {
			K1 []CardOrSeatType `json:"1"`
			K2 []CardOrSeatType `json:"2"`
			K3 []CardOrSeatType `json:"3"`
			K4 []CardOrSeatType `json:"4"`
		} `json:"ticket_seat_codeMap"`
		TicketTypeCodes []CardOrSeatType `json:"ticket_type_codes"`
	} `json:"limitBuySeatTicketDTO"`
	MaxTicketNum    Nstring `json:"maxTicketNum"`
	OrderRequestDTO struct {
		AdultNum            int      `json:"adult_num"`
		ApplyOrderNo        Nstring  `json:"apply_order_no"`
		BedLevelOrderNum    Nstring  `json:"bed_level_order_num"`
		BureauCode          Nstring  `json:"bureau_code"`
		CancelFlag          Nstring  `json:"cancel_flag"`
		CardNum             Nstring  `json:"card_num"`
		Channel             Nstring  `json:"channel"`
		ChildNum            int      `json:"child_num"`
		ChooseSeat          Nstring  `json:"choose_seat"`
		DisabilityNum       int      `json:"disability_num"`
		EndTime             InfoTime `json:"end_time"`
		ExchangeTrainFlag   Nstring  `json:"exchange_train_flag"`
		FromStationName     Nstring  `json:"from_station_name"`
		FromStationTelecode Nstring  `json:"from_station_telecode"`
		TetTicketPass       Nstring  `json:"get_ticket_pass"`
		IDMode              Nstring  `json:"id_mode"`
		IsShowPassCode      Nstring  `json:"isShowPassCode"`
		LeftTicketGenTime   Nstring  `json:"leftTicketGenTime"`
		OrderDate           Nstring  `json:"order_date"`
		PassengerFlag       Nstring  `json:"passengerFlag"`
		RealleftTicket      Nstring  `json:"realleftTicket"`
		ReqIPAddress        Nstring  `json:"reqIpAddress"`
		ReqTimeLeftStr      Nstring  `json:"reqTimeLeftStr"`
		ReserveFlag         Nstring  `json:"reserve_flag"`
		SeatDetailTypeCode  Nstring  `json:"seat_detail_type_code"`
		SeatTypeCode        Nstring  `json:"seat_type_code"`
		SequenceNo          Nstring  `json:"sequence_no"`
		StartTime           InfoTime `json:"start_time"`
		StartTimeStr        Nstring  `json:"start_time_str"`
		StationTrainCode    Nstring  `json:"station_train_code"`
		StudentNum          int      `json:"student_num"`
		TicketNum           int      `json:"ticket_num"`
		TicketTypeOrderNum  Nstring  `json:"ticket_type_order_num"`
		ToStationName       Nstring  `json:"to_station_name"`
		ToStationTelecode   Nstring  `json:"to_station_telecode"`
		TourFlag            Nstring  `json:"tour_flag"`
		TrainCodeText       Nstring  `json:"trainCodeText"`
		TrainDate           InfoTime `json:"train_date"`
		TrainDateStr        Nstring  `json:"train_date_str"`
		TrainLocation       Nstring  `json:"train_location"`
		TrainNo             Nstring  `json:"train_no"`
		TrainOrder          Nstring  `json:"train_order"`
		VarStr              Nstring  `json:"varStr"`
	} `json:"orderRequestDTO"`
	PurposeCodes Nstring `json:"purpose_codes"`
	// 省略。。。。queryLeftNewDetailDTO（余票信息）/queryLeftTicketRequestDTO
	// 'queryLeftNewDetailDTO': {
	// 	'BXRZ_num': '-1',
	// 	'BXRZ_price': '0',
	// 	'BXYW_num': '-1',
	// 	'BXYW_price': '0',
	// 	'EDRZ_num': '-1',
	// 	'EDRZ_price': '0',
	// 	'EDSR_num': '-1',
	// 	'EDSR_price': '0',
	// 	'ERRB_num': '-1',
	// 	'ERRB_price': '0',
	// 	'GG_num': '-1',
	// 	'GG_price': '0',
	// 	'GR_num': '-1',
	// 	'GR_price': '0',
	// 	'HBRW_num': '-1',
	// 	'HBRW_price': '0',
	// 	'HBRZ_num': '-1',
	// 	'HBRZ_price': '0',
	// 	'HBYW_num': '-1',
	// 	'HBYW_price': '0',
	// 	'HBYZ_num': '-1',
	// 	'HBYZ_price': '0',
	// 	'RW_num': '-1',
	// 	'RW_price': '0',
	// 	'RZ_num': '-1',
	// 	'RZ_price': '0',
	// 	'SRRB_num': '-1',
	// 	'SRRB_price': '0',
	// 	'SWZ_num': '26',
	// 	'SWZ_price': '02195',
	// 	'TDRZ_num': '-1',
	// 	'TDRZ_price': '0',
	// 	'TZ_num': '-1',
	// 	'TZ_price': '0',
	// 	'WZ_num': '0',
	// 	'WZ_price': '00730',
	// 	'WZ_seat_type_code': 'O',
	// 	'YB_num': '-1',
	// 	'YB_price': '0',
	// 	'YDRZ_num': '-1',
	// 	'YDRZ_price': '0',
	// 	'YDSR_num': '-1',
	// 	'YDSR_price': '0',
	// 	'YRRB_num': '-1',
	// 	'YRRB_price': '0',
	// 	'YW_num': '-1',
	// 	'YW_price': '0',
	// 	'YYRW_num': '-1',
	// 	'YYRW_price': '0',
	// 	'YZ_num': '-1',
	// 	'YZ_price': '0',
	// 	'ZE_num': '777',
	// 	'ZE_price': '00730',
	// 	'ZY_num': '115',
	// 	'ZY_price': '01170',
	// 	'arrive_time': '0954',
	// 	'control_train_day': '',
	// 	'controlled_train_flag': null,
	// 	'controlled_train_message': null,
	// 	'day_difference': null,
	// 	'end_station_name': null,
	// 	'end_station_telecode': null,
	// 	'from_station_name': '\u4E0A\u6D77\u8679\u6865',
	// 	'from_station_telecode': 'AOH',
	// 	'is_support_card': null,
	// 	'lishi': '00:54',
	// 	'seat_feature': '',
	// 	'start_station_name': null,
	// 	'start_station_telecode': null,
	// 	'start_time': '0900',
	// 	'start_train_date': '',
	// 	'station_train_code': 'G7505',
	// 	'to_station_name': '\u676D\u5DDE\u4E1C',
	// 	'to_station_telecode': 'HGH',
	// 	'train_class_name': null,
	// 	'train_no': '5l000G750502',
	// 	'train_seat_feature': '',
	// 	'yp_ex': ''
	// },
	// 'queryLeftTicketRequestDTO': {
	// 	'arrive_time': '09:54',
	// 	'bigger20': 'Y',
	// 	'exchange_train_flag': '1',
	// 	'from_station': 'AOH',
	// 	'from_station_name': '\u4E0A\u6D77\u8679\u6865',
	// 	'from_station_no': '07',
	// 	'lishi': '00:54',
	// 	'login_id': null,
	// 	'login_mode': null,
	// 	'login_site': null,
	// 	'purpose_codes': '00',
	// 	'query_type': null,
	// 	'seatTypeAndNum': null,
	// 	'seat_types': 'O9MO',
	// 	'start_time': '09:00',
	// 	'start_time_begin': null,
	// 	'start_time_end': null,
	// 	'station_train_code': 'G7505',
	// 	'ticket_type': null,
	// 	'to_station': 'HGH',
	// 	'to_station_name': '\u676D\u5DDE\u4E1C',
	// 	'to_station_no': '09',
	// 	'train_date': '20190212',
	// 	'train_flag': null,
	// 	'train_headers': null,
	// 	'train_no': '5l000G750502',
	// 	'useMasterPool': true,
	// 	'useWB10LimitTime': true,
	// 	'usingGemfireCache': false,
	// 	'ypInfoDetail': 'DSjv1392B2dI5dHj5ASZugvagKUDNNR%2FYV0asz6Gq4CtlF9QjviPIsA1wko%3D'
	// },

	TourFlag      Nstring `json:"tour_flag"`
	TrainLocation Nstring `json:"train_location"`
}

// InfoTime ticketInfo中的结构
type InfoTime struct {
	Date           int `json:"date"`
	Day            int `json:"day"`
	Hours          int `json:"hours"`
	Minutes        int `json:"minutes"`
	Month          int `json:"month"`
	Seconds        int `json:"seconds"`
	Time           int `json:"time"`
	TimezoneOffset int `json:"timezoneOffset"`
	Year           int `json:"year"`
}

// CardOrSeatType ticketInfo中的结构
type CardOrSeatType struct {
	EndStationName   Nstring `json:"end_station_name"`
	EndTime          Nstring `json:"end_time"`
	ID               Nstring `json:"id"`
	StartStationName Nstring `json:"start_station_name"`
	StartTime        Nstring `json:"start_time"`
	Value            Nstring `json:"value"`
}

// Nstring 允许null转化
type Nstring string

// UnmarshalJSON 自定义解析json函数
func (n *Nstring) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}
