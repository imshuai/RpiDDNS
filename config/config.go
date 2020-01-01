package config

import (
	"RpiDDNS/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/imshuai/lightlog"
	"github.com/imshuai/sysutils"
)

//Provider 服务提供商
type Provider struct {
	Name       string
	Username   string
	Password   string
	LoginToken string
}

//Record 解析记录
type Record struct {
	SubDomain         string
	SubDomainRecordID string `json:"-"`
	Type              string
	Value             string    `json:"-"`
	LastUpdateTime    time.Time `json:"-"`
}

//Domain 域名信息
type Domain struct {
	Provider   Provider
	DomainName string
	DomainID   string `json:"-"`
	Records    []Record
}

//Config 整体配置文件
type Config struct {
	Domains  []Domain
	Log      *lightlog.Logger `json:"-"`
	LogLevel string
	IPv4     string `json:"-"`
	IPv6     string `json:"-"`
}

//Init 初始化配置
func (c *Config) Init(confPath, logPath string) error {
	c.Log = lightlog.NewLogger(10)
	c.Log.TimeFormat = "2006-01-02 15:04:05 -0700"
	c.Log.Prefix = "[RpiDDNS]"
	if confPath == "" {
		confPath = utils.GetCurPath() + sysutils.PathSeparator() + "settings.json"
		c.Log.Debug("use default settings file:", confPath)
	}
	byts, err := ioutil.ReadFile(confPath)
	if err != nil {
		c.Log.Error("read settings file fail with error:", err.Error())
		if os.IsNotExist(err) {
			tConfig := &Config{
				Domains: []Domain{
					{
						Provider:   Provider{},
						DomainName: "",
						Records: []Record{
							{
								SubDomain: "",
								Type:      "",
							},
						},
					},
				},
			}
			byts, _ := json.MarshalIndent(tConfig, "", "\t")
			ioutil.WriteFile(confPath, byts, os.ModePerm)
			c.Log.Info("create new settings file:", confPath)
		}
		return err
	}
	err = json.Unmarshal(byts, c)
	if err != nil {
		c.Log.Error("unmarshal config fail with error:", err.Error())
		return err
	}
	for _, domain := range c.Domains {
		if len(domain.Records) == 0 || domain.DomainName == "" || (domain.Provider.Name == "" && domain.Provider.LoginToken == "" && (domain.Provider.Username == "" || domain.Provider.Password == "")) {
			err := errors.New("incorrect settings, please check it again")
			c.Log.Error(err.Error())
			return err
		}
		for _, record := range domain.Records {
			if record.SubDomain == "" {
				err := errors.New("incorrect settings, please check it again")
				c.Log.Error(err.Error())
				return err
			}
			switch record.Type {
			case "A", "AAAA":
				continue
			default:
				err := errors.New("incorrect settings, please check it again")
				c.Log.Error(err.Error())
				return err
			}
		}
	}

	if logPath == "" {
		logPath = utils.GetCurPath() + sysutils.PathSeparator() + "ddns.log"
		c.Log.Debug("use default log file:", logPath)
	}
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		c.Log.Error("open log file fail with error:", err.Error())
		return err
	}
	c.Log.FileOut = f
	switch c.LogLevel {
	case "debug":
		c.Log.Level = lightlog.LevelDebug
	case "info":
		c.Log.Level = lightlog.LevelInfo
	case "warning":
		c.Log.Level = lightlog.LevelWarning
	case "error":
		c.Log.Level = lightlog.LevelError
	case "fatal":
		c.Log.Level = lightlog.LevelFatal
	case "none":
		c.Log.Level = lightlog.LevelNone
	default:
		c.Log.Level = lightlog.LevelError
	}
	c.Log.Info("initialize settings complete")
	return nil
}
