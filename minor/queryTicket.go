package minor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

// QueryTicket 查询余票，如果查到余票返回对应结果,
// cdn: 要使用的cdn的ip,
// fromStationCode: 出发站代码,
// toStationCode: 到达站代码,
// date: 购票日期,
// passengerTypeCode: 乘客类型代码,
// trainCodeStr: 要购买的车次(用,分割),
// seatStr: 要购买的坐席(用,分割),
// seatFirst: 坐席优先
func QueryTicket(cdn, fromStationCode, toStationCode, passengerTypeCode string, date time.Time, wantTrainCodes, wantSeatCodes []string, seatFirst bool, passengerCount int) (tti *domain.TrainTicketInfo, err error) {

	var buffer bytes.Buffer
	if cdn == "" {
		buffer.WriteString(constant.Urls["leftTicket"])
	} else {
		buffer.WriteString(fmt.Sprintf(constant.Urls["leftTicketCDN"], cdn))
	}
	buffer.WriteString("?leftTicketDTO.train_date=")
	buffer.WriteString(date.Format("2006-01-02"))
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
		} else if c, _ := strconv.Atoi(trainTicketInfo.SecondClassSeat); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "O"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.HardSeat); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "1"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.HardSleeper); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "3"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.MCSleeper); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "F"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SoftSleeper); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "4"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SoftSeat); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "2"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.NoSeat); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "1"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.FirstClassSeat); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "M"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.BusinessClassSeat); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "9"
			infoes = append(infoes, trainTicketInfo)
		} else if c, _ := strconv.Atoi(trainTicketInfo.SeniorSoftSleeper); c >= passengerCount {
			trainTicketInfo.PassengerTicketSeatCode = "6"
			infoes = append(infoes, trainTicketInfo)
			// } else if c, _ := strconv.Atoi(trainTicketInfo.OtherSeat); c >= passengerCount {
			// 	infoes = append(infoes, trainTicketInfo)
		}
	}
	var selected []*domain.TrainTicketInfo
	isSeatSelect := len(wantSeatCodes) > 0
	isTrainSelect := len(wantTrainCodes) > 0
	if seatFirst { // 坐席优先
		if isSeatSelect {
			selected = selectSeatTrainTicketInfo(infoes, wantSeatCodes, passengerCount)
		}
		if isTrainSelect {
			if isSeatSelect {
				selected = selectTrainTrainTicketInfo(selected, wantTrainCodes, passengerCount)
			} else {
				selected = selectTrainTrainTicketInfo(infoes, wantTrainCodes, passengerCount)
			}
		}
	} else { // 车次优先
		if isTrainSelect {
			selected = selectTrainTrainTicketInfo(infoes, wantTrainCodes, passengerCount)
		}
		if isSeatSelect {
			if isTrainSelect {
				selected = selectSeatTrainTicketInfo(selected, wantSeatCodes, passengerCount)
			} else {
				selected = selectSeatTrainTicketInfo(infoes, wantSeatCodes, passengerCount)
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

func selectTrainTrainTicketInfo(all []*domain.TrainTicketInfo, wantTrainCodes []string, passengerCount int) (trainOrdered []*domain.TrainTicketInfo) {
	for _, train := range wantTrainCodes {
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

func selectSeatTrainTicketInfo(all []*domain.TrainTicketInfo, wantSeatCodes []string, passengerCount int) (seatsOrdered []*domain.TrainTicketInfo) {
	// notSelected := all
	// TODO 测试一下这种写法会不会影响原数组
	for _, seat := range wantSeatCodes {
		for _, tti := range all {
			if tti.Empty {
				continue
			}
			switch seat {
			case "ZE":
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
			case "YZ":
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
			case "YW":
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
			case "SRRB": // 动卧
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
			case "RW":
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
			case "RZ":
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
			case "WZ":
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
			case "ZY":
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
			case "TZ":
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
			case "SWZ":
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
			case "GR":
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
// tzNum := sa[25] // 特等座
// wzNum := sa[26] // 无座
// ybNum := sa[27] //
// ywNum := sa[28]   // 硬卧
// yzNum := sa[29]   // 硬座
// zeNum := sa[30]   // 二等座
// zyNum := sa[31]   // 一等座
// swzNum := sa[32]  // 商务座
// srrbNum := sa[33] // 动卧
// ypEx := sa[34]
// seatTypes := sa[35]
// exchangeTrainFlag := sa[36]
