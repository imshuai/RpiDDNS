package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//IsIPv4 判断地址是否是IPv4地址
func IsIPv4(address string) bool {
	return strings.Count(address, ":") <= 1
}

//IsIPv6 判断地址是否是IPv6地址
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

//GetIPv6 获取本机全球单播IPv6地址
func GetIPv6() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("net.Interfaces failed, err: %s", err.Error())
	}
	t := []string{}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To16() != nil && IsIPv6(ipnet.IP.String()) && ipnet.IP.IsGlobalUnicast() {
						t = append(t, ipnet.IP.String())
					}
				}
			}
		}
	}
	if len(t) == 0 {
		return "", errors.New("cannot find a ipv6 address")
	}
	return t[0], nil
}

//GetIPv4 获取本机外网IPv4地址
func GetIPv4() (string, error) {
	resp, err := http.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	byts, _ := ioutil.ReadAll(resp.Body)
	ip := string(byts)
	ip = strings.Trim(ip, "\n")
	return ip, nil
}

//PostData 推送数据
func PostData(url string, v url.Values, h http.Header) ([]byte, error) {
	client := http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("User-Agent", fmt.Sprintf("RpiDDNS/0.1 (%s)", c.Email))
	req.Header = h
	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("post data fail with error:%s", err.Error())
	}

	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)

	return resp, nil
}

//GetCurPath 获取当前文件执行的路径
func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])
	fmt.Println(1, file)
	path, _ := filepath.Abs(file)
	fmt.Println(2, path)
	rst := filepath.Dir(path)
	fmt.Println(3, rst)
	return rst
}
