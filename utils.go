package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

func getIPv6() []string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return []string{}
	}
	t := []string{}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To16() != nil && isIPv6(ipnet.IP.String()) && ipnet.IP.IsGlobalUnicast() {
						t = append(t, ipnet.IP.String())
					}
				}
			}
		}
	}
	return t
}

func postData(url string, v url.Values) []byte {
	client := http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("RpiDDNS/0.1 (%s)", c.Email))
	response, err := client.Do(req)

	if err != nil {
		log.Println("Post failed...")
		log.Println(err)
		return nil
	}

	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)

	return resp
}
