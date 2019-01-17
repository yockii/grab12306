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
	BusinessClassSeat string // 商务座
	SpecialClassSeat  string // 特等座
	MCSleeper         string // 动卧
	TrainLocationCode string // 余票及排队信息查询时用到的TrainLocation
	LeftTicketSecret  string // 余票及排队信息查询时用到的LeftTicket

	PassengerTicketSeatCode string // 乘客本车次应购买的坐席代码

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
	PurposeCodes          Nstring `json:"purpose_codes"`
	QueryLeftNewDetailDTO struct {
		BXRZNum                string  `json:"BXRZ_num"`
		BXRZPrice              string  `json:"BXRZ_price"`
		BXYWNum                string  `json:"BXYW_num"`
		BXYWPrice              string  `json:"BXYW_price"`
		EDRZNum                string  `json:"EDRZ_num"`
		EDRZPrice              string  `json:"EDRZ_price"`
		EDSRNum                string  `json:"EDSR_num"`
		EDSRPrice              string  `json:"EDSR_price"`
		ERRBNum                string  `json:"ERRB_num"`
		ERRBPrice              string  `json:"ERRB_price"`
		GGNum                  string  `json:"GG_num"`
		GGPrice                string  `json:"GG_price"`
		GRNum                  string  `json:"GR_num"`
		GRPrice                string  `json:"GR_price"`
		HBRWNum                string  `json:"HBRW_num"`
		HBRWPrice              string  `json:"HBRW_price"`
		HBRZNum                string  `json:"HBRZ_num"`
		HBRZPrice              string  `json:"HBRZ_price"`
		HBYWNum                string  `json:"HBYW_num"`
		HBYWPrice              string  `json:"HBYW_price"`
		HBYZNum                string  `json:"HBYZ_num"`
		HBYZPrice              string  `json:"HBYZ_price"`
		RWNum                  string  `json:"RW_num"`
		RWPrice                string  `json:"RW_price"`
		RZNum                  string  `json:"RZ_num"`
		RZPrice                string  `json:"RZ_price"`
		SRRBNum                string  `json:"SRRB_num"`
		SRRBPrice              string  `json:"SRRB_price"`
		SWZNum                 string  `json:"SWZ_num"` // 商务座
		SWZPrice               string  `json:"SWZ_price"`
		TDRZNum                string  `json:"TDRZ_num"`
		TDRZPrice              string  `json:"TDRZ_price"`
		TZNum                  string  `json:"TZ_num"`
		TZPrice                string  `json:"TZ_price"`
		WZNum                  string  `json:"WZ_num"`
		WZPrice                string  `json:"WZ_price"`
		WZSeatTypeCode         string  `json:"WZ_seat_type_code"`
		YBNum                  string  `json:"YB_num"`
		YBPrice                string  `json:"YB_price"`
		YDRZNum                string  `json:"YDRZ_num"`
		YDRZPrice              string  `json:"YDRZ_price"`
		YDSRNum                string  `json:"YDSR_num"`
		YDSRPrice              string  `json:"YDSR_price"`
		YRRBNum                string  `json:"YRRB_num"`
		YRRBPrice              string  `json:"YRRB_price"`
		YWNum                  string  `json:"YW_num"`
		YWPrice                string  `json:"YW_price"`
		YYRWNum                string  `json:"YYRW_num"`
		YYRWPrice              string  `json:"YYRW_price"`
		YZNum                  string  `json:"YZ_num"`
		YZPrice                string  `json:"YZ_price"`
		ZENum                  string  `json:"ZE_num"`
		ZEPrice                string  `json:"ZE_price"`
		ZYNum                  string  `json:"ZY_num"`
		ZYPrice                string  `json:"ZY_price"`
		ArriveTime             string  `json:"arrive_time"`
		ControlTrainDay        string  `json:"control_train_day"`
		ControlledTrainFlag    Nstring `json:"controlled_train_flag"`
		ControlledTrainMessage Nstring `json:"controlled_train_message"`
		DayDifference          Nstring `json:"day_difference"`
		EndStationName         Nstring `json:"end_station_name"`
		EndStationTelecode     Nstring `json:"end_station_telecode"`
		FromStationName        string  `json:"from_station_name"`
		FromStationTelecode    string  `json:"from_station_telecode"`
		IsSupportCard          Nstring `json:"is_support_card"`
		Lishi                  string  `json:"lishi"`
		SeatFeature            string  `json:"seat_feature"`
		StartStationName       Nstring `json:"start_station_name"`
		StartStationTelecode   Nstring `json:"start_station_telecode"`
		StartTime              string  `json:"start_time"`
		StartTrainDate         string  `json:"start_train_date"`
		StationTrainCode       string  `json:"station_train_code"`
		ToStationName          string  `json:"to_station_name"`
		ToStationTelecode      string  `json:"to_station_telecode"`
		TrainClassName         Nstring `json:"train_class_name"`
		TrainNo                string  `json:"train_no"`
		TrainSeatFeature       string  `json:"train_seat_feature"`
		YpEx                   string  `json:"yp_ex"`
	} `json:"queryLeftNewDetailDTO"`
	QueryLeftTicketRequestDTO struct {
		ArriveTime        string  `json:"arrive_time"`
		Bigger20          string  `json:"bigger20"`
		ExchangeTrainFlag string  `json:"exchange_train_flag"`
		FromStation       string  `json:"from_station"`
		FromStationName   string  `json:"from_station_name"`
		FromStationNo     string  `json:"from_station_no"`
		Lishi             string  `json:"lishi"`
		LoginID           Nstring `json:"login_id"`
		LoginMode         Nstring `json:"login_mode"`
		LoginSite         Nstring `json:"login_site"`
		PurposeCodes      string  `json:"purpose_codes"`
		QueryType         Nstring `json:"query_type"`
		SeatTypeAndNum    Nstring `json:"seatTypeAndNum"`
		SeatTypes         string  `json:"seat_types"`
		StartTime         string  `json:"start_time"`
		StartTimeBegin    Nstring `json:"start_time_begin"`
		StartTimeEnd      Nstring `json:"start_time_end"`
		StationTrainCode  string  `json:"station_train_code"`
		TicketType        Nstring `json:"ticket_type"`
		ToStation         string  `json:"to_station"`
		ToStationName     string  `json:"to_station_name"`
		ToStationNo       string  `json:"to_station_no"`
		TrainDate         string  `json:"train_date"`
		TrainFlag         Nstring `json:"train_flag"`
		TrainHeaders      Nstring `json:"train_headers"`
		TrainNo           string  `json:"train_no"`
		UseMasterPool     bool    `json:"useMasterPool"`
		UseWB10LimitTime  bool    `json:"useWB10LimitTime"`
		UsingGemfireCache bool    `json:"usingGemfireCache"`
		YpInfoDetail      string  `json:"ypInfoDetail"`
	} `json:"queryLeftTicketRequestDTO"`
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
