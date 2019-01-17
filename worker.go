package main

import (
	"fmt"
	"time"

	"github.com/yockii/grab12306/config"
	"github.com/yockii/grab12306/domain"
	"github.com/yockii/grab12306/major"
	"github.com/yockii/grab12306/minor"
	cdnutil "github.com/yockii/grab12306/utils/cdn"
)

// QueryAndOrderTicket 查询并下订单
func QueryAndOrderTicket(oi *domain.OrderInfo) (err error) {
	// 查询
	cdnList := cdnutil.CdnList
	size := len(cdnList)
	if size == 0 {
		cdnutil.InitCdn()
		cdnList = cdnutil.CdnList
		size = len(cdnList)
	}
	if size == 0 {
		fmt.Println("无法取得有效CDN, 请检查网络或检查公网IP")
		return
	}

	var ttiChan = make(chan *domain.TrainTicketInfo)
	go func() {
		for !oi.Stop {
			cdn := cdnutil.GetRandomCDN()
			if cdn != nil {
				tti, err := QueryTicket(cdn, oi.FromStationCode, oi.ToStationCode, oi.QueryPassengerType, oi.TravelDate, len(oi.Passengers), oi.WantedTrainCodes, oi.WantedSeatCodes, oi.SeatFirst)
				if err != nil {
					// 检查网络连接错误？
				}
				if tti != nil {
					oi.Stop = true
					ttiChan <- tti
					return
				}
			}
			time.Sleep(time.Millisecond * 100) // 100ms查询一次
		}
	}()

	trainTicketInfo := <-ttiChan
	orderID, err := major.QuickSubmitTicketSuit(oi, trainTicketInfo)
	if err != nil {
		// 检查错误?
	}
	if orderID != "" {
		// 订票成功!
	}
	return
}

// QueryTicket 查询余票，通过单个cdn
// travDate 2019-01-02
// passengerType ADULT
func QueryTicket(cdn *cdnutil.Cdn, fromStationCode, toStationCode, passengerType string, TravelDate time.Time, passengerCount int, wantedTrainCodes, WantedSeatCodes []string, seatFirst bool) (tti *domain.TrainTicketInfo, err error) {
	tti, err = minor.QueryTicket(cdn.IP, fromStationCode, toStationCode, passengerType, TravelDate, wantedTrainCodes, WantedSeatCodes, seatFirst, passengerCount)
	if err != nil {
		cdn.NextAvailableTime = time.Now().UnixNano() + (config.GetInstance().Basic.CdnBlackHouseTimeInSecond * 1e6)
		return
	}
	return
}
