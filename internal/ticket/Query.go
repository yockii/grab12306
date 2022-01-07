package ticket

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/yockii/grab12306/internal/constant"
	"github.com/yockii/grab12306/internal/domain"
	"github.com/yockii/grab12306/internal/model"
)

var seatIndex = map[string]int{
	"特等座":  25,
	"一等座":  31,
	"二等座":  30,
	"动卧":   33,
	"硬座":   29,
	"无座":   26,
	"软座":   24,
	"硬卧":   28,
	"软卧":   23,
	"高级软卧": 21,
	"商务座":  32,
}

var seatType = map[string]string{
	"特等座":  "P",
	"一等座":  "M",
	"二等座":  "O",
	"动卧":   "F",
	"硬座":   "1",
	"无座":   "1",
	"软座":   "2",
	"硬卧":   "3",
	"软卧":   "4",
	"高级软卧": "6",
	"商务座":  "9",
}

func QueryLeftTicket(trainTicket *model.TrainTicket) (ttis []*domain.TrainTicketInfo, err error) {
	params := url.Values{}
	params.Set("leftTicketDTO.train_date", trainTicket.Date)
	params.Set("leftTicketDTO.from_station", trainTicket.From)
	params.Set("leftTicketDTO.to_station", trainTicket.To)
	params.Set("purpose_codes", trainTicket.PassengerType)

	resp, err := http.DefaultClient.Get(constant.UrlQueryLeftTicket + "?" + params.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respJson := gjson.ParseBytes(response)
	ti := respJson.Get("data.result").Array()
	for _, t := range ti {
		tInfo := strings.Split(t.String(), "|")
		if tInfo[11] != "Y" {
			continue
		}
		if trainTicket.TrainNo != "" && tInfo[3] != trainTicket.TrainNo {
			continue
		}
		if trainTicket.Seat != "" {
			if index, ok := seatIndex[trainTicket.Seat]; ok {
				seatCount := tInfo[index]
				if seatCount == "有" {
					ttis = append(ttis, &domain.TrainTicketInfo{
						Secret:              tInfo[0],
						TrainNo:             tInfo[2],
						StationTrainCode:    tInfo[3],
						SeatType:            seatType[trainTicket.Seat],
						FromStationTelecode: tInfo[6],
						ToStationTelecode:   tInfo[7],
						LeftTicket:          tInfo[12],
						PurposeCodes:        "", // ?
						TrainLocation:       tInfo[15],
						DepartDate:          tInfo[8],
						ArriveDate:          tInfo[9],
						Interval:            tInfo[10],
						RepeatSubmitToken:   "", // ?
					})
				} else {
					c, err0 := strconv.ParseInt(seatCount, 10, 64)
					if err0 != nil {
						continue
					}
					if int(c) > trainTicket.TicketNum {
						ttis = append(ttis, &domain.TrainTicketInfo{
							Secret:              tInfo[0],
							TrainNo:             tInfo[2],
							StationTrainCode:    tInfo[3],
							SeatType:            seatType[trainTicket.Seat],
							FromStationTelecode: tInfo[6],
							ToStationTelecode:   tInfo[7],
							LeftTicket:          tInfo[12],
							PurposeCodes:        "", // ?
							TrainLocation:       tInfo[15],
							DepartDate:          tInfo[8],
							ArriveDate:          tInfo[9],
							Interval:            tInfo[10],
							RepeatSubmitToken:   "", // ?
						})
					}
				}
			}
			continue
		}
		//TODO 不指定座位，有就行
	}
	return
}
