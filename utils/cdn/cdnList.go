package cdnutil

import (
	"math/rand"
	"sync"
	"time"

	"github.com/yockii/grab12306/constant"
)

// CdnList 可用的cdn列表
var CdnList []Cdn
var dealingList []Cdn

// InitCdn 获取并校验可用CDN列表
func InitCdn() {
	list := GetCdnList()
	total := len(list)
	dealingList = make([]Cdn, 0, total)
	i := 0
	var wg sync.WaitGroup
	for ; i < total/constant.CdnDealCountPerTime; i++ {
		wg.Add(1)
		dealList := list[i*constant.CdnDealCountPerTime : (i+1)*constant.CdnDealCountPerTime]
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

	CdnList = dealingList
	dealingList = nil
}

// GetDealingAvailableCount 获取正在处理中的可用cdn数量
func GetDealingAvailableCount() int {
	if dealingList != nil {
		return len(dealingList)
	}
	return 0
}

// GetRandomCDN 获取随机的CDN，获取的同时会对该cdn设置下次可执行的时间
func GetRandomCDN() *Cdn {
	size := len(CdnList)
	randomIndex := rand.Intn(size)
	for CdnList[randomIndex].NextAvailableTime > time.Now().UnixNano() {
		randomIndex = rand.Intn(size)
		time.Sleep(time.Second)
	}
	cdn := CdnList[randomIndex]
	cdn.NextAvailableTime = time.Now().UnixNano() + rand.Int63n(4*1e6) + 1e6 // 随机增加1-5s
	return &cdn
}
