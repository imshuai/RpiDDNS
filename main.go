package main

import "fmt"

import "os"

var (
	c = &config{}
)

func init() {
	err := c.init()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	tPublic = Public{
		LoginToken:   c.Token,
		Format:       "json",
		Lang:         "en",
		ErrorOnEmpty: "no",
	}
}

func main() {
	record := getRecord(c.DomainID, c.SubDomainRecordID)
	fmt.Println("Old Record:", record.Value)
	ips := getIPv6()
	fmt.Println("ipv6:", ips)
	value := ""
	for _, ip := range ips {
		if ip == record.Value {
			continue
		}
		value = ip
	}
	if value == "" {
		fmt.Println("Same Record, abort!")
		return
	}
	fmt.Println("New Value", value)
	newRecord, err := updateRecord(c.DomainID, c.SubDomainRecordID, c.SubDomain, "AAAA", value)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Actually new record:", newRecord.Value)
}
