package domain

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

type SeatType struct {
	Name  string
	Index int
	Code  string
}
