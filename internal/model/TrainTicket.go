package model

type TrainTicket struct {
	Id            string
	AccountId     string
	From          string
	To            string
	Date          string
	SeatFirst     int    // 坐席优先=1
	WantedTrains  string // 指定车次，逗号分隔
	WantedSeats   string // 指定座位，逗号分隔
	TicketNum     int    // 预定票数
	PassengerType string // 查询普通票还是学生票
	Status        int    //状态： 1-执行中 2-预定中 3-已预定 4-已关闭
}

type TrainTicketPassenger struct {
	Id            string
	TrainTicketId string
	PassengerId   string
}
