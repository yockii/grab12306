package constant

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// Station 站点信息
type Station struct {
	Name     string
	Code     string
	Pinyin   string
	JianPin3 string
	JianPin  string
	ID       string
}

// Stations 站点编码对应的站点信息
var Stations = make(map[string]*Station)

// StationCodes 站点名称对应编码 如 北京北: VAP
var StationCodes = make(map[string]string)

// FetchStationNamesFrom12306 从12306官方获取站点名称代码信息
func FetchStationNamesFrom12306() {
	res, err := http.Get(Urls["stationNames"])
	if err != nil {
		fmt.Println("获取站点信息出错!", err)
		return
	}
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	reg := regexp.MustCompile("var\\s+station_names\\s+=\\s*'(.*)';")
	stationInfoResult := reg.FindStringSubmatch(string(content[:]))
	stationInfoes := stationInfoResult[1]
	stations := strings.Split(stationInfoes, "@")

	for _, s := range stations {
		if s != "" {
			ss := strings.Split(s, "|")
			if len(ss) > 0 {
				Stations[ss[2]] = &Station{
					Name:     ss[1],
					Code:     ss[2],
					Pinyin:   ss[3],
					JianPin3: ss[0],
					JianPin:  ss[4],
					ID:       ss[5],
				}
				StationCodes[ss[1]] = ss[2]
			}
		}
	}
}
