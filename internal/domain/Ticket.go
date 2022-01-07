package domain

type TrainTicketInfo struct {
	Secret              string
	TrainNo             string
	StationTrainCode    string // 可理解的车次号
	SeatType            string
	FromStationTelecode string
	ToStationTelecode   string
	LeftTicket          string
	PurposeCodes        string // 00?
	TrainLocation       string // H3
	DepartDate          string
	ArriveDate          string
	Interval            string
	RepeatSubmitToken   string
}

type SeatType struct {
	Name  string
	Index int
	Code  string
}
