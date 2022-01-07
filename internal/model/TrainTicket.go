package model

type TrainTicket struct {
	Id            string
	AccountId     string
	From          string
	To            string
	Date          string
	TrainNo       string // 指定车次
	Seat          string // 指定座位
	TicketNum     int    // 预定票数
	PassengerType string // 查询普通票还是学生票
	Status        int    //状态： 1-执行中 2-预定中 3-已预定 4-已关闭
}

type TrainTicketPassenger struct {
	Id            string
	TrainTicketId string
	PassengerId   string
}
