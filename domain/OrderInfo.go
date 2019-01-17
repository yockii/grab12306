package domain

import (
	"time"
)

// OrderInfo 用户抢票信息
type OrderInfo struct {
	FromStationCode    string
	ToStationCode      string
	TravelDate         time.Time
	Passengers         []*Passenger
	WantedTrainCodes   []string
	WantedSeatCodes    []string // SWZ/TZ/YZ....
	SeatFirst          bool
	QueryPassengerType string // ADULT
	Stop               bool
}
