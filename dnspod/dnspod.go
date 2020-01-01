package dnspod

import (
	"RpiDDNS/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

//Public 公共参数结构体
type Public struct {
	LoginToken   string `form:"login_token"`
	Format       string `form:"format"`
	Lang         string `form:"lang"`
	ErrorOnEmpty string `form:"error_on_empty"`
}

//Status 公共返回结构体
type Status struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

//Domain 域名信息返回
type Domain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//Record 解析记录返回
type Record struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

var (
	tPublic Public
)

//SetToken 设置登陆密匙
func SetToken(token string) {
	tPublic.LoginToken = token
	tPublic.ErrorOnEmpty = "false"
	tPublic.Format = "json"
	tPublic.Lang = "en"
}

//NewDomain 创建域名实例
func NewDomain() (d *Domain, err error) {
	return nil, nil
}

//GetDomainID 获取域名ID
func GetDomainID(tDomain string) (string, error) {
	apiHost := "https://dnsapi.cn/Domain.Info"
	params := struct {
		Public
		Domain string `form:"domain"`
	}{
		tPublic,
		tDomain,
	}
	paramsValue := paramsMarshal(params)

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("User-Agent", fmt.Sprintf("RpiDDNS/0.1 (%s)", "iris-me@live.com"))
	resp, err := utils.PostData(apiHost, paramsValue, header)
	if err != nil {
		return "", err
	}
	ret := &struct {
		Status
		Domain Domain `json:"domain"`
	}{}
	err = json.Unmarshal(resp, ret)
	if err != nil {
		return "", err
	}
	if ret.Code == "1" {
		return ret.Domain.ID, nil
	}
	return "", errors.New(ret.Message)
}

//GetRecordID 获取subDomain的解析记录ID
func GetRecordID(subDomain string, domainID string) (string, error) {
	apiHost := "https://dnsapi.cn/Record.List"
	params := struct {
		Public
		DomainID  string `form:"domain_id"`
		SubDomain string `form:"sub_domain"`
	}{
		tPublic,
		domainID,
		subDomain,
	}
	paramsValue := paramsMarshal(params)

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("User-Agent", fmt.Sprintf("RpiDDNS/0.1 (%s)", "iris-me@live.com"))
	resp, err := utils.PostData(apiHost, paramsValue, header)
	if err != nil {
		return "", err
	}
	ret := &struct {
		Status
		Records []Record `json:"records"`
	}{}
	json.Unmarshal(resp, ret)
	if ret.Code != "1" {
		return "", errors.New(ret.Message)
	}
	return ret.Records[0].ID, nil
}

//GetRecordValue 获取解析记录信息
func GetRecordValue(domainID string, recordID string) (string, error) {
	apiHost := "https://dnsapi.cn/Record.Info"
	params := struct {
		Public
		DomainID string `form:"domain_id"`
		RecordID string `form:"record_id"`
	}{
		tPublic,
		domainID,
		recordID,
	}
	paramsValue := paramsMarshal(params)

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("User-Agent", fmt.Sprintf("RpiDDNS/0.1 (%s)", "iris-me@live.com"))
	resp, err := utils.PostData(apiHost, paramsValue, header)
	if err != nil {
		return "", err
	}
	ret := &struct {
		Status
		Record struct {
			ID        string `json:"id"`
			SubDomain string `json:"sub_domain"`
			Value     string `json:"value"`
		} `json:"record"`
	}{}
	json.Unmarshal(resp, ret)
	if ret.Code != "1" {
		return "", errors.New(ret.Message)
	}
	return ret.Record.Value, nil
}

//UpdateRecord 提交更新记录
func UpdateRecord(domainID string, recordID, subDomain, recordType, value string) error {
	apiHost := "https://dnsapi.cn/Record.Modify"
	params := struct {
		Public
		DomainID   string `form:"domain_id"`
		RecordID   string `form:"record_id"`
		SubDomain  string `form:"sub_domain"`
		RecordType string `form:"record_type"`
		RecordLine string `form:"record_line"`
		Value      string `form:"value"`
	}{
		tPublic,
		domainID,
		recordID,
		subDomain,
		recordType,
		"默认",
		value,
	}
	paramsValue := paramsMarshal(params)

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("User-Agent", fmt.Sprintf("RpiDDNS/0.1 (%s)", "iris-me@live.com"))
	resp, err := utils.PostData(apiHost, paramsValue, header)
	if err != nil {
		return err
	}
	ret := &struct {
		Status
	}{}
	json.Unmarshal(resp, ret)
	if ret.Status.Code != "1" {
		return errors.New(ret.Status.Message)
	}
	return nil
}

func paramsMarshal(p interface{}) url.Values {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	params := url.Values{}
	if v.Kind() != reflect.Struct {
		return params
	}
	params = getParams(t, v)
	return params
}
func getParams(t reflect.Type, v reflect.Value) url.Values {
	params := url.Values{}
	if v.Kind() != reflect.Struct {
		return params
	}
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			tParams := getParams(f.Type, v.Field(i))
			for kk, vv := range tParams {
				params.Set(kk, vv[0])
			}
		} else {
			name, ok := f.Tag.Lookup("form")
			if !ok {
				continue
			}
			vv := v.FieldByName(f.Name).String()
			if vv == "" {
				continue
			}
			params.Set(name, vv)
		}
	}
	return params
}
