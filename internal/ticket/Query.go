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

//QueryLeftTicket 查询指定信息的余票信息，若有余票则返回有余票的车次
func QueryLeftTicket(trainTicket *model.TrainTicket) (tti *domain.TrainTicketInfo, err error) {
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
	var infoes []*domain.TrainTicketInfo
	for _, t := range ti {
		trainTicketInfo := transferToTrainTicketInfo(t.String())
		if !trainTicketInfo.CanBuy {
			continue
		}
		if trainTicketInfo.SecondClassSeat == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "O"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.HardSeat == "有" || trainTicketInfo.NoSeat == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "1"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.HardSleeper == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "3"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.MCSleeper == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "F"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.SoftSleeper == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "4"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.SoftSeat == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "2"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.FirstClassSeat == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "M"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.BusinessClassSeat == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "9"
			infoes = append(infoes, trainTicketInfo)
		} else if trainTicketInfo.SeniorSoftSleeper == "有" {
			trainTicketInfo.PassengerTicketSeatCode = "6"
			infoes = append(infoes, trainTicketInfo)
			// } else if trainTicketInfo.OtherSeat == "有" {
			// 	infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SecondClassSeat); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "O"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.HardSeat); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "1"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.HardSleeper); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "3"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.MCSleeper); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "F"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.SoftSleeper); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "4"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.SoftSeat); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "2"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.NoSeat); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "1"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.FirstClassSeat); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "M"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.BusinessClassSeat); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "9"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ = strconv.Atoi(trainTicketInfo.SeniorSoftSleeper); c >= trainTicket.TicketNum {
			trainTicketInfo.PassengerTicketSeatCode = "6"
			infoes = append(infoes, trainTicketInfo)
			// } else if c, _ := strconv.Atoi(trainTicketInfo.OtherSeat); c >= passengerCount {
			// 	infoes = append(infoes, trainTicketInfo)
		}
	}
	var selected []*domain.TrainTicketInfo
	wantSeats := strings.Split(trainTicket.WantedSeats, ",")
	wantTrans := strings.Split(trainTicket.WantedTrains, ",")
	isSeatSelect := len(wantSeats) > 0
	isTrainSelect := len(wantTrans) > 0
	if trainTicket.SeatFirst == 1 {
		if isSeatSelect {
			selected = selectSeatTrainTicketInfo(infoes, wantSeats, trainTicket.TicketNum)
		}
		if isTrainSelect {
			if isSeatSelect {
				selected = selectTrainTrainTicketInfo(selected, wantTrans, trainTicket.TicketNum)
			} else {
				selected = selectTrainTrainTicketInfo(infoes, wantTrans, trainTicket.TicketNum)
			}
		}
	} else { // 车次优先
		if isTrainSelect {
			selected = selectTrainTrainTicketInfo(infoes, wantTrans, trainTicket.TicketNum)
		}
		if isSeatSelect {
			if isTrainSelect {
				selected = selectSeatTrainTicketInfo(selected, wantSeats, trainTicket.TicketNum)
			} else {
				selected = selectSeatTrainTicketInfo(infoes, wantSeats, trainTicket.TicketNum)
			}
		}
	}
	if isSeatSelect || isTrainSelect {
		if len(selected) > 0 {
			tti = selected[0]
		}
	} else {
		if len(infoes) > 0 {
			tti = infoes[0]
		}
	}
	return
}

func selectTrainTrainTicketInfo(all []*domain.TrainTicketInfo, wantTrains []string, passengerCount int) (trainOrdered []*domain.TrainTicketInfo) {
	for _, train := range wantTrains {
		for _, tti := range all {
			if tti.TrainCode == train {
				if tti.SecondClassSeat == "有" || tti.HardSeat == "有" || tti.HardSleeper == "有" || tti.MCSleeper == "有" || tti.SoftSleeper == "有" || tti.SoftSeat == "有" || tti.NoSeat == "有" || tti.FirstClassSeat == "有" || tti.BusinessClassSeat == "有" || tti.SeniorSoftSleeper == "有" || tti.OtherSeat == "有" {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.SecondClassSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.HardSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.HardSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.MCSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.SoftSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.SoftSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.NoSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.FirstClassSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.BusinessClassSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.SeniorSoftSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ = strconv.Atoi(tti.OtherSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				}
				break
			}
		}
	}
	return
}

func selectSeatTrainTicketInfo(all []*domain.TrainTicketInfo, wantSeats []string, passengerCount int) (seatsOrdered []*domain.TrainTicketInfo) {
	// notSelected := all
	// TODO 测试一下这种写法会不会影响原数组
	for _, seat := range wantSeats {
		for _, tti := range all {
			if tti.Empty {
				continue
			}
			switch seat {
			case "二等座":
				if tti.SecondClassSeat == "" || tti.SecondClassSeat == "无" {
					continue
				}
				if tti.SecondClassSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "O"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SecondClassSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "O"
					tti.Empty = true
				}
			case "硬座":
				if tti.HardSeat == "" || tti.HardSeat == "无" {
					continue
				}
				if tti.HardSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "1"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.HardSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "1"
					tti.Empty = true
				}
			case "硬卧":
				if tti.HardSleeper == "" || tti.HardSleeper == "无" {
					continue
				}
				if tti.HardSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "3"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.HardSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "3"
					tti.Empty = true
				}
			case "动卧": // 动卧
				if tti.MCSleeper == "" || tti.MCSleeper == "无" {
					continue
				}
				if tti.MCSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "F"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.MCSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "F"
					tti.Empty = true
				}
			case "软卧":
				if tti.SoftSleeper == "" || tti.SoftSleeper == "无" {
					continue
				}
				if tti.SoftSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "4"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SoftSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "4"
					tti.Empty = true
				}
			case "软座":
				if tti.SoftSeat == "" || tti.SoftSeat == "无" {
					continue
				}
				if tti.SoftSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "2"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SoftSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "2"
					tti.Empty = true
				}
			case "无座":
				if tti.NoSeat == "" || tti.NoSeat == "无" {
					continue
				}
				if tti.NoSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "1"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.NoSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "1"
					tti.Empty = true
				}
			case "一等座":
				if tti.FirstClassSeat == "" || tti.FirstClassSeat == "无" {
					continue
				}
				if tti.FirstClassSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "M"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.FirstClassSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "M"
					tti.Empty = true
				}
			case "特等座":
				if tti.SpecialClassSeat == "" || tti.SpecialClassSeat == "无" {
					continue
				}
				if tti.SpecialClassSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "9"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SpecialClassSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "9"
					tti.Empty = true
				}
			case "商务座":
				if tti.BusinessClassSeat == "" || tti.BusinessClassSeat == "无" {
					continue
				}
				if tti.BusinessClassSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "P"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.BusinessClassSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "P"
					tti.Empty = true
				}
			case "高级软卧":
				if tti.SeniorSoftSleeper == "" || tti.SeniorSoftSleeper == "无" {
					continue
				}
				if tti.SeniorSoftSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "6"
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SeniorSoftSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.PassengerTicketSeatCode = "6"
					tti.Empty = true
				}
				// case "X": // 其他 - 暂不支持吧
				// 	if tti.OtherSeat == "" || tti.OtherSeat == "无" {
				// 		continue
				// 	}
				// 	if tti.OtherSeat == "有" {
				// 		seatsOrdered = append(seatsOrdered, tti)
				// 		tti.Empty = true
				// 	} else if c, _ := strconv.Atoi(tti.OtherSeat); c >= passengerCount {
				// 		seatsOrdered = append(seatsOrdered, tti)
				// 		tti.Empty = true
				// 	}
			}
		}
	}
	return
}

func transferToTrainTicketInfo(result string) (tti *domain.TrainTicketInfo) {
	sa := strings.Split(result, "|")
	tti = &domain.TrainTicketInfo{}
	tti.Secret = sa[0]
	// tti.Secret, _ = url.QueryUnescape(sa[0])
	tti.TrainNo = sa[2]
	tti.TrainCode = sa[3]
	tti.StartStationCode = sa[4]
	tti.EndStationCode = sa[5]
	tti.FromStationCode = sa[6]
	tti.ToStationCode = sa[7]
	tti.StartTime = sa[8]
	tti.ArriveTime = sa[9]
	tti.Duration = sa[10]
	tti.CanBuy = sa[11] == "Y"
	tti.LeftTicketSecret = sa[12]
	tti.StartTrainDate = sa[13]
	tti.TrainLocationCode = sa[15]
	tti.SeniorSoftSleeper = sa[21]
	tti.OtherSeat = sa[22]
	tti.SoftSleeper = sa[23]
	tti.SoftSeat = sa[24]
	tti.SpecialClassSeat = sa[25]
	tti.NoSeat = sa[26]
	tti.HardSleeper = sa[28]
	tti.HardSeat = sa[29]
	tti.SecondClassSeat = sa[30]
	tti.FirstClassSeat = sa[31]
	tti.BusinessClassSeat = sa[32]
	tti.MCSleeper = sa[33]

	return
}
