package cdnutil

import (
	"sync"

	"github.com/yockii/grab12306/config"
)

var cdnList []Cdn
var dealingList []Cdn

// InitCdn 获取并校验可用CDN列表
func InitCdn() {
	list := GetCdnList()
	total := len(list)
	dealingList = make([]Cdn, 0, total)
	i := 0
	var wg sync.WaitGroup
	for ; i < total/config.CdnDealCountPerTime; i++ {
		wg.Add(1)
		dealList := list[i*config.CdnDealCountPerTime : (i+1)*config.CdnDealCountPerTime]
		VerifyCdnList(dealList)
		go func() {
			defer wg.Done()
			for _, v := range dealList {
				cdn := v
				if cdn.Available {
					dealingList = append(dealingList, cdn)
				}
			}
		}()
	}
	wg.Wait()

	cdnList = dealingList
	dealingList = nil
}

// GetDealingAvailableCount 获取正在处理中的可用cdn数量
func GetDealingAvailableCount() int {
	if dealingList != nil {
		return len(dealingList)
	}
	return 0
}
