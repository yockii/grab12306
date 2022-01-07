package ticket

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/httputil"
	coreUtil "github.com/yockii/qscore/pkg/util"

	"github.com/yockii/grab12306/internal/constant"
	"github.com/yockii/grab12306/internal/model"
	"github.com/yockii/grab12306/internal/util"
)

func CheckLoginStatus(accountId string) (bool, error) {
	client, cookieJar, err := getHttpClientFromAccount(accountId)
	if err != nil {
		return false, err
	}
	if cookieJar == nil {
		return false, nil
	}
	if resp, err := client.Get(constant.UrlCheckLogin); err != nil {
		return false, err
	} else {
		defer resp.Body.Close()
		response, err0 := io.ReadAll(resp.Body)
		if err0 != nil {
			return false, err0
		}
		j := gjson.ParseBytes(response)
		if j.Get("data.is_login").String() == "Y" {
			return true, nil
		}
	}

	return true, nil
}

func getHttpClientFromAccount(accountId string) (*http.Client, *httputil.CookieJar, error) {
	account := new(model.Account)
	if exist, err := database.DB.ID(accountId).Get(account); err != nil {
		return nil, nil, err
	} else if !exist {
		return nil, nil, errors.New("Account not found ")
	}
	if account.CookieJar == nil {
		return util.GetHttpClient(), nil, nil
	}
	jar := httputil.Decode(account.CookieJar)
	if jar.IsZero() {
		return util.GetHttpClient(), nil, nil
	}
	client := util.GetCookieHttpClient(jar)
	return client, &jar, nil
}

func GetLoginQrCode(accountId string) (qrCode string, qrUuid string, err error) {
	jar := httputil.CookieJar{}
	client := util.GetCookieHttpClient(jar)
	client.Get(constant.UrlLoginPage)
	// 获取cookies完毕
	v := url.Values{}
	v.Set("appid", "otn")

	resp, err := client.PostForm(constant.UrlGetQrCode, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respJson := gjson.ParseBytes(response)
	qrUuid = respJson.Get("uuid").String()
	qrCode = respJson.Get("image").String()
	// 更新cookies到数据库
	database.DB.ID(accountId).Update(&model.Account{CookieJar: jar.Encode()})
	return
}

func CheckQrScan(qrUuid, accountId string) (login bool, uamtk string, err error) {
	client, cookieJar, err := getHttpClientFromAccount(accountId)
	if err != nil {
		return
	}
	if cookieJar == nil {
		err = errors.New("No Cookies Found ")
		return
	}
	u, _ := url.Parse(constant.Server12306URL)
	railDeviceId := ""
	railExpiration := ""
	cs := cookieJar.Cookies(u)
	for _, c := range cs {
		if c.Name == "RAIL_DEVICEID" {
			railDeviceId = c.Value
		} else if c.Name == "RAIL_EXPIRATION" {
			railExpiration = c.Value
		}
		if railDeviceId != "" && railExpiration != "" {
			break
		}
	}

	v := url.Values{}
	v.Set("appid", "otn")
	v.Set("uuid", qrUuid)
	v.Set("RAIL_DEVICEID", railDeviceId)
	v.Set("RAIL_EXPIRATION", railExpiration)

	resp, err := client.PostForm(constant.UrlCheckQrCode, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respJson := gjson.ParseBytes(response)
	if respJson.Get("result_code").String() == "2" {
		// 更新cookies到数据库
		database.DB.ID(accountId).Update(&model.Account{CookieJar: cookieJar.Encode()})
		return true, respJson.Get("uamtk").String(), nil
	}
	return
}

func UserAuth(uamtk, accountId string) (tk string, err error) {
	client, cookieJar, err := getHttpClientFromAccount(accountId)
	if err != nil {
		return
	}
	if cookieJar == nil {
		err = errors.New("No Cookies Found ")
		return
	}
	v := url.Values{}
	v.Set("appid", "excater")
	v.Set("uamtk", uamtk)
	resp, err := client.PostForm(constant.UrlUserAuth, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respJson := gjson.ParseBytes(response)
	tk = respJson.Get("newapptk").String()
	// 更新cookies到数据库
	database.DB.ID(accountId).Update(&model.Account{CookieJar: cookieJar.Encode()})
	return
}

func UamAuth(apptk, accountId string) (username string, err error) {
	client, cookieJar, err := getHttpClientFromAccount(accountId)
	if err != nil {
		return
	}
	if cookieJar == nil {
		err = errors.New("No Cookies Found ")
		return
	}
	v := url.Values{}
	v.Set("tk", apptk)
	resp, err := client.PostForm(constant.UrlUamAuth, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respJson := gjson.ParseBytes(response)
	username = respJson.Get("username").String()
	return
}

func CheckLogin(accountId string) (ok bool, err error) {
	client, cookieJar, err := getHttpClientFromAccount(accountId)
	if err != nil {
		return
	}
	if cookieJar == nil {
		err = errors.New("No Cookies Found ")
		return
	}
	resp, err := client.PostForm(constant.UrlUamAuth, url.Values{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respJson := gjson.ParseBytes(response)
	return respJson.Get("data.flag").String() == "1", nil
}

func GetPassengers(accountId string) (err error) {
	client, cookieJar, err := getHttpClientFromAccount(accountId)
	if err != nil {
		return
	}
	if cookieJar == nil {
		err = errors.New("No Cookies Found ")
		return
	}
	pageTotal := 1
	var ps []gjson.Result
	for pageIndex := 1; pageIndex <= pageTotal; pageIndex++ {
		var result []gjson.Result
		pageTotal, result, err = getPassengers(client, pageIndex)
		if err != nil {
			return
		}
		ps = append(ps, result...)
	}
	for _, p := range ps {
		passenger := &model.Passenger{
			AccountId:     accountId,
			PassengerUuid: p.Get("passenger_uuid").String(),
		}
		exist := false
		if exist, err = database.DB.Get(passenger); err != nil {
			return
		}
		if !exist {
			passenger.Id = coreUtil.GenerateDatabaseID()
			passenger.PassengerName = p.Get("passenger_name").String()
			passenger.PassengerType = p.Get("passenger_type").String()
			passenger.AllEncStr = p.Get("allEncStr").String()
			_, err = database.DB.Insert(passenger)
			if err != nil {
				return
			}
		}
	}
	return
}

func getPassengers(client *http.Client, pageIndex int) (pageTotal int, ps []gjson.Result, err error) {
	v := url.Values{}
	v.Set("pageIndex", fmt.Sprintf("%d", pageIndex))
	v.Set("pageSize", "10")

	var resp *http.Response
	resp, err = client.PostForm(constant.UrlGetPassengers, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respJson := gjson.ParseBytes(response)
	pageTotal = int(respJson.Get("data.pageTotal").Int())
	ps = respJson.Get("data.datas").Array()
	return

}
