package utils

import (
	"net"
	"net/url"
	"strconv"
	"time"
)

// 生成唯一的sessionID，用于溯源每个请求日志
func CreateUniSessionId() string {
	unixNano := time.Now().UnixNano()           // 获取当前时间的Unix纳秒时间戳 1680067671341495000
	unixNano = (unixNano * 100000) & 0x7FFFFFFF //1365038848
	return strconv.FormatInt(unixNano, 10)      //1365038848
}

func GetLocalIP() ([]string, error) {
	var ipStr []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ipStr, err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					//获取IPv6
					/*if ipnet.IP.To16() != nil {
						mylog.Log.Info(ipnet.IP.String())
						ipStr = append(ipStr, ipnet.IP.String())
					}*/
					//获取IPv4
					if ipnet.IP.To4() != nil {
						// mylog.Log.Info(ipnet.IP.String())
						ipStr = append(ipStr, ipnet.IP.String())
					}
				}
			}
		}
	}
	return ipStr, nil
}

func GetDomainFromReferer(referer string) (string, error) {
	parsedURL, err := url.Parse(referer)
	if err != nil {
		return "", err
	}

	return parsedURL.Hostname(), nil
}

func GetRequestIp(reqIp string) string {
	if reqIp == "::1" {
		reqIp = "127.0.0.1"
	}
	return reqIp
}
