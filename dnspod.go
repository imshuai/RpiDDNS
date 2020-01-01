package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
)

type Public struct {
	LoginToken   string `form:"login_token"`
	Format       string `form:"format"`
	Lang         string `form:"lang"`
	ErrorOnEmpty string `form:"error_on_empty"`
}

type Status struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type Domain struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Record struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

var (
	tPublic Public
)

func getDomainID(tDomain string) int {
	apiHost := "https://dnsapi.cn/Domain.List"
	params := struct {
		Public
		Keyword string `form:"keyword"`
	}{
		tPublic,
		tDomain,
	}
	paramsValue := paramsMarshal(params)

	//fmt.Println(paramsValue.Encode())

	resp := postData(apiHost, paramsValue)
	ret := &struct {
		Status
		Domains []Domain `json:"domains"`
	}{}
	json.Unmarshal(resp, ret)
	for _, domain := range ret.Domains {
		if domain.Name == tDomain {
			return domain.ID
		}
	}
	return 0
}
func getRecords(subDomain string, domainID string) []Record {
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

	//fmt.Println(paramsValue.Encode())

	resp := postData(apiHost, paramsValue)
	ret := &struct {
		Status
		Records []Record `json:"records"`
	}{}
	json.Unmarshal(resp, ret)
	return ret.Records
}

func getRecord(domainID string, recordID string) Record {
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

	//fmt.Println(paramsValue.Encode())

	resp := postData(apiHost, paramsValue)
	ret := &struct {
		Status
		Record struct {
			ID        string `json:"id"`
			SubDomain string `json:"sub_domain"`
			Value     string `json:"value"`
		} `json:"record"`
	}{}

	//fmt.Println(string(resp))

	json.Unmarshal(resp, ret)
	record := Record{}
	record.ID = ret.Record.ID
	record.Name = ret.Record.SubDomain
	record.Value = ret.Record.Value
	return record
}

func updateRecord(domainID string, recordID, subDomain, recordType, value string) (Record, error) {
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

	//fmt.Println(paramsValue.Encode())

	resp := postData(apiHost, paramsValue)
	ret := &struct {
		Status
		Record Record `json:"record"`
	}{}
	json.Unmarshal(resp, ret)
	if ret.Status.Code != "1" {
		return Record{}, errors.New(ret.Status.Message)
	}
	return ret.Record, nil
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
