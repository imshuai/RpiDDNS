package main

import (
	"RpiDDNS/config"
	"RpiDDNS/dnspod"
	"RpiDDNS/utils"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

const (
	//NAME 程序名称
	NAME = "RpiDDNS"
	//VERSION 程序版本号
	VERSION = "v1.0.2"
	//AUTHOR 作者
	AUTHOR = "im帥"
	//EMAIL 作者邮箱
	EMAIL = "iris-me@live.com"
)

var (
	c = &config.Config{}
)

// func init() {
// 	var err error
// 	if len(os.Args) > 1 {
// 		err = c.Init(os.Args[1], os.Args[2])
// 	} else {
// 		err = c.Init("", "")
// 	}
// 	if err != nil {
// 		c.Log.Fatal(err.Error())
// 		fmt.Println(err)
// 		os.Exit(-1)
// 	}
// }

func main() {
	app := cli.NewApp()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  AUTHOR,
			Email: EMAIL,
		},
	}
	app.Name = NAME
	app.Version = VERSION
	app.Usage = "Update this computer's ipv4/ipv6 address to your dns record."
	app.Copyright = time.Now().Format("2006") + " © " + AUTHOR

	app.Commands = []cli.Command{
		cli.Command{
			Name: "run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config,c",
					Usage: "path to config file",
				},
			},
			Usage: "run with config file",
			Action: func(c *cli.Context) error {
				//TODO:完成初版功能

				return nil
			},
		},
		cli.Command{
			Name: "update",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain,d",
					Usage: "your domain name. like: baidu.com google.com home.yourdomain.com",
				},
				cli.StringFlag{
					Name:  "subdomain,s",
					Usage: "your subdomain name. like: pc router",
				},
				cli.StringFlag{
					Name:  "type,t",
					Usage: "your dns record type. like: A AAAA",
				},
				cli.StringFlag{
					Name:  "ip",
					Usage: "ip for update to dns record",
				},
				cli.StringFlag{
					Name:  "provider,pd",
					Value: "dnspod",
					Usage: "your dns provider name. for now only support \"dnspod\"",
				},
				cli.StringFlag{
					Name:  "username,u",
					Usage: "your username for login ",
				},
				cli.StringFlag{
					Name:  "password,p",
					Usage: "your password for login",
				},
				cli.StringFlag{
					Name:  "token",
					Usage: "your access token for login",
				},
				cli.StringFlag{
					Name:  "loglevel",
					Value: "info",
					Usage: "log level: none|fatal|error|warning|info|debug",
				},
			},
			Usage: "update record once",
			Action: func(c *cli.Context) error {
				//TODO 完成初版功能
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func e() {
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
