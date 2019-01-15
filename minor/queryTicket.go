package minor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/yockii/grab12306/constant"
	"github.com/yockii/grab12306/domain"
	netutil "github.com/yockii/grab12306/utils/net"
)

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
	tti.NoSeat = sa[26]
	tti.HardSleeper = sa[28]
	tti.HardSeat = sa[29]
	tti.SecondClassSeat = sa[30]
	tti.FirstClassSeat = sa[31]
	tti.BusinessClassSeat = sa[32]
	tti.MCSleeper = sa[33]

	return
}

// QueryTicket 查询余票，如果查到余票返回对应结果,
// cdn: 要使用的cdn的ip,
// fromStationCode: 出发站代码,
// toStationCode: 到达站代码,
// date: 购票日期,
// passengerTypeCode: 乘客类型代码,
// trainCodeStr: 要购买的车次(用,分割),
// seatStr: 要购买的坐席(用,分割),
// seatFirst: 坐席优先
func QueryTicket(cdn, fromStationCode, toStationCode, date, passengerTypeCode, trainCodeStr, seatStr string, seatFirst bool, passengerCount int) (tti *domain.TrainTicketInfo, err error) {
	wantSeats := strings.Split(seatStr, ",")
	wantTrains := strings.Split(trainCodeStr, ",")

	var buffer bytes.Buffer
	if cdn == "" {
		buffer.WriteString(constant.Urls["leftTicket"])
	} else {
		buffer.WriteString(fmt.Sprintf(constant.Urls["leftTicketCDN"], cdn))
	}
	buffer.WriteString("?leftTicketDTO.train_date=")
	buffer.WriteString(date)
	buffer.WriteString("&leftTicketDTO.from_station=")
	buffer.WriteString(fromStationCode)
	buffer.WriteString("&leftTicketDTO.to_station=")
	buffer.WriteString(toStationCode)
	buffer.WriteString("&purpose_codes=")
	buffer.WriteString(passengerTypeCode)

	content, err := netutil.Get(buffer.String())
	if err != nil {
		return
	}

	var ltRes domain.LeftTicketResponse
	err = json.Unmarshal(content, &ltRes)
	if err != nil {
		return
	}

	if !ltRes.Status {
		err = errors.New(ltRes.Messages)
		return
	}

	if len(ltRes.Data.Result) == 0 {
		err = errors.New("车次查询信息有误，或返回数据有异常")
		return
	}
	var infoes []*domain.TrainTicketInfo
	for _, s := range ltRes.Data.Result {
		trainTicketInfo := transferToTrainTicketInfo(s)
		if !trainTicketInfo.CanBuy {
			continue
		}
		if trainTicketInfo.SecondClassSeat == "有" || trainTicketInfo.HardSeat == "有" || trainTicketInfo.HardSleeper == "有" || trainTicketInfo.MCSleeper == "有" || trainTicketInfo.SoftSleeper == "有" || trainTicketInfo.SoftSeat == "有" || trainTicketInfo.NoSeat == "有" || trainTicketInfo.FirstClassSeat == "有" || trainTicketInfo.BusinessClassSeat == "有" || trainTicketInfo.SeniorSoftSleeper == "有" || trainTicketInfo.OtherSeat == "有" {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SecondClassSeat); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.HardSeat); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.HardSleeper); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.MCSleeper); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SoftSleeper); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SoftSeat); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.NoSeat); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.FirstClassSeat); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.BusinessClassSeat); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SeniorSoftSleeper); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.OtherSeat); c >= passengerCount {
			infoes = append(infoes, trainTicketInfo)
		}
	}
	var selected []*domain.TrainTicketInfo
	isSeatSelect := wantSeats[0] != ""
	isTrainSelect := wantTrains[0] != ""
	if seatFirst { // 坐席优先
		if isSeatSelect {
			selected = selectSeatTrainTicketInfo(infoes, wantSeats, passengerCount)
		}
		if isTrainSelect {
			if isSeatSelect {
				selected = selectTrainTrainTicketInfo(selected, wantTrains, passengerCount)
			} else {
				selected = selectTrainTrainTicketInfo(infoes, wantTrains, passengerCount)
			}
		}
	} else { // 车次优先
		if isTrainSelect {
			selected = selectTrainTrainTicketInfo(infoes, wantTrains, passengerCount)
		}
		if isSeatSelect {
			if isTrainSelect {
				selected = selectSeatTrainTicketInfo(selected, wantSeats, passengerCount)
			} else {
				selected = selectSeatTrainTicketInfo(infoes, wantSeats, passengerCount)
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

func selectTrainTrainTicketInfo(all []*domain.TrainTicketInfo, wantedTrains []string, passengerCount int) (trainOrdered []*domain.TrainTicketInfo) {

	for _, train := range wantedTrains {
		for _, tti := range all {
			if tti.TrainCode == train {
				if tti.SecondClassSeat == "有" || tti.HardSeat == "有" || tti.HardSleeper == "有" || tti.MCSleeper == "有" || tti.SoftSleeper == "有" || tti.SoftSeat == "有" || tti.NoSeat == "有" || tti.FirstClassSeat == "有" || tti.BusinessClassSeat == "有" || tti.SeniorSoftSleeper == "有" || tti.OtherSeat == "有" {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.SecondClassSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.HardSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.HardSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.MCSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.SoftSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.SoftSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.NoSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.FirstClassSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.BusinessClassSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.SeniorSoftSleeper); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				} else if c, _ := strconv.Atoi(tti.OtherSeat); c >= passengerCount {
					trainOrdered = append(trainOrdered, tti)
				}
				break
			}
		}
	}
	return
}

func selectSeatTrainTicketInfo(all []*domain.TrainTicketInfo, wantedSeats []string, passengerCount int) (seatsOrdered []*domain.TrainTicketInfo) {
	// notSelected := all
	// TODO 测试一下这种写法会不会影响原数组

	for _, seat := range wantedSeats {
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
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SecondClassSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "硬座":
				if tti.HardSeat == "" || tti.HardSeat == "无" {
					continue
				}
				if tti.HardSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.HardSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "硬卧":
				if tti.HardSleeper == "" || tti.HardSleeper == "无" {
					continue
				}
				if tti.HardSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.HardSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "动卧":
				if tti.MCSleeper == "" || tti.MCSleeper == "无" {
					continue
				}
				if tti.MCSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.MCSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "软卧":
				if tti.SoftSleeper == "" || tti.SoftSleeper == "无" {
					continue
				}
				if tti.SoftSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SoftSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "软座":
				if tti.SoftSeat == "" || tti.SoftSeat == "无" {
					continue
				}
				if tti.SoftSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SoftSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "无座":
				if tti.NoSeat == "" || tti.NoSeat == "无" {
					continue
				}
				if tti.NoSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.NoSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "一等座":
				if tti.FirstClassSeat == "" || tti.FirstClassSeat == "无" {
					continue
				}
				if tti.FirstClassSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.FirstClassSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "商务座":
				fallthrough
			case "特等座":
				if tti.BusinessClassSeat == "" || tti.BusinessClassSeat == "无" {
					continue
				}
				if tti.BusinessClassSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.BusinessClassSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "高级软卧":
				if tti.SeniorSoftSleeper == "" || tti.SeniorSoftSleeper == "无" {
					continue
				}
				if tti.SeniorSoftSleeper == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.SeniorSoftSleeper); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			case "其他":
				if tti.OtherSeat == "" || tti.OtherSeat == "无" {
					continue
				}
				if tti.OtherSeat == "有" {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				} else if c, _ := strconv.Atoi(tti.OtherSeat); c >= passengerCount {
					seatsOrdered = append(seatsOrdered, tti)
					tti.Empty = true
				}
			}
		}
	}
	return
}

// secretStr := sa[0]      // 购票用
// buttonTextInfo := sa[1] // 备注(预订字样)
// trainNo := sa[2]
// trainCode := sa[3]            // 车次
// startStationTelecode := sa[4] //起始站代号
// endStationTelecode := sa[5]   //终点站代号
// fromStationTelecode := sa[6]  //出发站代号
// toStationTelecode := sa[7]    //到达站代号
// startTime := sa[8]            // 出发时间
// arriveTime := sa[9]           // 到达时间
// duration := sa[10]            // 历时
// canWebBuy := sa[11]           //是否能购买：Y 可以
// ypInfo := sa[12]	// leftTicket - getQueueCount用到
// startTrainDate := sa[13] // 出发日期
// trainSeatFeature := sa[14]
// locationCode := sa[15]
// fromStationNo := sa[16]
// toStationNo := sa[17]
// isSupportCard := sa[18]
// controlledTrainFlag := sa[19]
// ggNum := sa[20]
// grNum := sa[21] // 高级软卧
// qtNum := sa[22] // 其他
// rwNum := sa[23] // 软卧
// rzNum := sa[24] //软座
// tzNum := sa[25]
// wzNum := sa[26] // 无座
// ybNum := sa[27]
// ywNum := sa[28]   // 硬卧
// yzNum := sa[29]   // 硬座
// zeNum := sa[30]   // 二等座
// zyNum := sa[31]   // 一等座
// swzNum := sa[32]  // 商务特等座
// srrbNum := sa[33] // 动卧
// ypEx := sa[34]
// seatTypes := sa[35]
// exchangeTrainFlag := sa[36]
