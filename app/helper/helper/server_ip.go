package helper

import (
	"net"
	"sync"
)

var (
	serverIPs    []string // 缓存 IP 列表
	serverIPOnce sync.Once
)

// 初始化时预加载 IP（线程安全）
func initserverIPs() {
	serverIPOnce.Do(func() {
		interfaces, _ := net.Interfaces()
		for _, iface := range interfaces {
			// 快速过滤：排除回环、未启用、Docker 虚拟接口等
			if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 || len(iface.HardwareAddr) == 0 {
				continue
			}

			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok || ipNet.IP.IsLoopback() {
					continue
				}

				ip := ipNet.IP.To4()
				if ip != nil && !ip.IsLinkLocalUnicast() {
					serverIPs = append(serverIPs, ip.String())
				}
			}
		}
	})
}

// 获取缓存的 IP 列表
func GetServerIPs() []string {
	initserverIPs()
	return serverIPs
}

// 获取第一个非回环 IPv4（常用场景）
func GetFirstServerIP() string {
	initserverIPs()
	if len(serverIPs) > 0 {
		return serverIPs[0]
	}
	return "127.0.0.1" // 默认回环地址
}
