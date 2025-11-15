package helper

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"strings"
)

// IPLocation IP地理位置信息
type IPLocation struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
	Raw      string `json:"raw"` // 原始字符串
}

var (
	dbPath = "./runtime/ip2region_v4.xdb" // 数据库文件路径
)

// 根据IP地址获取地理位置信息
func GetIPLocation(ip string) (*IPLocation, error) {
	if ip == "" {
		return nil, fmt.Errorf("IP地址不能为空")
	}
	// 创建查询对象
	searcher, err := xdb.NewWithFileOnly(xdb.IPv4, dbPath)
	if err != nil {
		return nil, fmt.Errorf("创建查询对象失败: %w", err)
	}
	defer searcher.Close()
	// 查询IP
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		return nil, fmt.Errorf("查询IP失败: %w", err)
	}
	// 解析结果
	location := parseRegion(ip, region)
	return location, nil
}

// 解析地理位置字符串
// 格式：国家|省份|城市|ISP
func parseRegion(ip, region string) *IPLocation {
	parts := strings.Split(region, "|")

	location := &IPLocation{
		IP:  ip,
		Raw: region,
	}
	if len(parts) >= 1 {
		location.Country = parts[0]
	}
	if len(parts) >= 2 {
		location.Province = parts[1]
	}
	if len(parts) >= 3 {
		location.City = parts[2]
	}
	if len(parts) >= 4 {
		location.ISP = parts[3]
	}
	return location
}
