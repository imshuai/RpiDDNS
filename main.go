package main

import (
	"RpiDDNS/config"
	"RpiDDNS/dnspod"
	"RpiDDNS/utils"
	"fmt"
	"os"
	"time"
)

var (
	c = &config.Config{}
)

func init() {
	var err error
	if len(os.Args) > 1 {
		err = c.Init(os.Args[1], os.Args[2])
	} else {
		err = c.Init("", "")
	}
	if err != nil {
		c.Log.Fatal(err.Error())
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	ipv4, err := utils.GetIPv4()
	if err != nil {
		c.Log.Error("get ipv4 address fail with error:", err.Error())
		return
	}
	c.Log.Info("ipv4 address:", ipv4)
	ipv6, err := utils.GetIPv6()
	if err != nil {
		c.Log.Error("get ipv6 address fail with error:", err.Error())
		return
	}
	c.Log.Info("ipv6 address:", ipv6)
	c.IPv4 = ipv4
	c.IPv6 = ipv6
	for _, domain := range c.Domains {
		var e error
		switch domain.Provider.Name {
		case "dnspod":
			dnspod.SetToken(domain.Provider.LoginToken)
			domain.DomainID, e = dnspod.GetDomainID(domain.DomainName)
			if err != nil {
				c.Log.Error(fmt.Sprintf("[dnspod]get %s's id fail with error:%v", domain.DomainName, e))
				continue
			}
			for _, record := range domain.Records {
				record.SubDomainRecordID, err = dnspod.GetRecordID(record.SubDomain, domain.DomainID)
				if err != nil {
					c.Log.Error(fmt.Sprintf("[dnspod]get %s's %s record id fail with error:%v", domain.DomainName, record.SubDomain, err))
					continue
				}
				oldValue, err := dnspod.GetRecordValue(domain.DomainID, record.SubDomainRecordID)
				if err != nil {
					c.Log.Error(fmt.Sprintf("[dnspod]get %s's id fail with error:%v", domain.DomainName, e))
					continue
				}
				c.Log.Info(fmt.Sprintf("old value for %s.%s:%s", record.SubDomain, domain.DomainName, oldValue))
				var newValue string
				switch record.Type {
				case "A":
					if c.IPv4 == oldValue {
						c.Log.Info("same value, abort")
						continue
					}
					newValue = c.IPv4
				case "AAAA":
					if c.IPv6 == oldValue {
						c.Log.Info("same value, abort")
						continue
					}
					newValue = c.IPv6
				default:
					c.Log.Error("incorrect record type, want: A or AAAA, get:", record.Type)
					continue
				}
				err = dnspod.UpdateRecord(domain.DomainID, record.SubDomainRecordID, record.SubDomain, record.Type, newValue)
				if err != nil {
					c.Log.Error("update record fail with error:", err.Error())
					continue
				}
				c.Log.Info("update record success")
			}
		default:
			continue
		}
	}
	c.Log.Info("-----> update all domain's record complete <-----")
	time.Sleep(1 * time.Second)
}
