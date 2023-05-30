package util

import (
	"encoding/csv"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net"
	"os"
)

func GetExcelDataByPath(filename string) ([][]string, error) {
	fileHandler, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	rows := fileHandler.GetRows("Sheet1")
	return rows, err
}

func GetCSVDataByPath(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 将CSV文件的字节数据转换为io.Reader接口
	f := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())

	reader := csv.NewReader(f)
	result := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, record)
	}
	return result, nil
}

func GetLocalIP() (ipv4, ipv6 net.IP, err error) {
	// 获取所有网络接口列表
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	// 遍历每个网络接口
	for _, iface := range ifaces {
		// 排除无效的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, nil, err
		}

		// 遍历每个地址，找到 IP 地址
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 排除非 IPv4 和 IPv6 地址
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				ipv4 = ip
			} else {
				ipv6 = ip
			}
		}
	}

	return ipv4, ipv6, nil
}

func GenerateSnowflakeInt64() int64 {
	snow, _ := GenerateSnowflake()

	return snow.Int64()
}

func GenerateSnowflake() (snowflake.ID, error) {
	node, errNode := snowflake.NewNode(1)
	if errNode != nil {
		return 0, errNode
	}
	return node.Generate(), nil
}
