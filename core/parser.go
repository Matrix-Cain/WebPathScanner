package core

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseTarget() error {

	ipv4withmaskRe, _ := regexp.Compile(`^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])/(3[0-2]|[1-2]?[0-9])$`)
	ipv4rangeRe, _ := regexp.Compile(`^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])-(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`)

	//# 尝试解析url
	parsedUrl, err := url.Parse(GlobalConfig.Target)
	if err != nil {
		log.Errorln("Invalid Target")
		return err
	}
	//# 判断不带http
	if parsedUrl.Scheme != "http" {
		//# 判断IP/Mask格式
		if ipv4withmaskRe.MatchString(parsedUrl.Path) {
			//# 是子网还是网址 e.g. 192.168.1.1/24 or http://192.168.1.1/24
			log.Warnf("[*] %s is IP/Mask[Y] or URL(http://)[n]? [Y/n]", GlobalConfig.Target)
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)
			if strings.Compare(strings.ToLower(text), "n") == 0 || strings.Compare(strings.ToLower(text), "no") == 0 {
				//# 按链接处理 e.g. http://192.168.1.1/24
				GlobalConfig.TargetList = append(GlobalConfig.TargetList, GlobalConfig.Target)
			} else { //Treat input target as CIDR expression
				slice, err := cidrToIPSlice(GlobalConfig.Target)
				if err != nil {
					return err
				}
				//https://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go
				GlobalConfig.TargetList = append(GlobalConfig.TargetList, slice...)

			}

		} else if ipv4rangeRe.MatchString(GlobalConfig.Target) { //判断网络范围格式 e.g. 192.168.1.1-192.168.1.100

			cidr, err := iPv4RangeToCIDRRange(strings.Split(GlobalConfig.Target, "-")[0], strings.Split(GlobalConfig.Target, "-")[1])
			if err != nil {
				return err
			}
			slice, err := cidrToIPSlice(cidr[0])
			if err != nil {
				return err
			}
			GlobalConfig.TargetList = append(GlobalConfig.TargetList, slice...)

		} else { //# 按照链接处理
			GlobalConfig.TargetList = append(GlobalConfig.TargetList, GlobalConfig.Target)
		}
		//# 为http://格式
	} else {
		GlobalConfig.TargetList = append(GlobalConfig.TargetList, GlobalConfig.Target)
	}
	return nil
}

// Convert IPv4 range into CIDR
//https://gist.github.com/P-A-R-U-S/a090dd90c5104ce85a29c32669dac107
func iPv4RangeToCIDRRange(ipStart string, ipEnd string) (cidrs []string, err error) {

	cidr2mask := []uint32{
		0x00000000, 0x80000000, 0xC0000000,
		0xE0000000, 0xF0000000, 0xF8000000,
		0xFC000000, 0xFE000000, 0xFF000000,
		0xFF800000, 0xFFC00000, 0xFFE00000,
		0xFFF00000, 0xFFF80000, 0xFFFC0000,
		0xFFFE0000, 0xFFFF0000, 0xFFFF8000,
		0xFFFFC000, 0xFFFFE000, 0xFFFFF000,
		0xFFFFF800, 0xFFFFFC00, 0xFFFFFE00,
		0xFFFFFF00, 0xFFFFFF80, 0xFFFFFFC0,
		0xFFFFFFE0, 0xFFFFFFF0, 0xFFFFFFF8,
		0xFFFFFFFC, 0xFFFFFFFE, 0xFFFFFFFF,
	}

	ipStartUint32 := iPv4ToUint32(ipStart)
	ipEndUint32 := iPv4ToUint32(ipEnd)

	if ipStartUint32 > ipEndUint32 {
		msg := fmt.Sprintf("start IP:%s must be less than end IP:%s", ipStart, ipEnd)
		return nil, errors.New(msg)
	}

	for ipEndUint32 >= ipStartUint32 {
		maxSize := 32
		for maxSize > 0 {

			maskedBase := ipStartUint32 & cidr2mask[maxSize-1]

			if maskedBase != ipStartUint32 {
				break
			}
			maxSize--

		}

		x := math.Log(float64(ipEndUint32-ipStartUint32+1)) / math.Log(2)
		maxDiff := 32 - int(math.Floor(x))
		if maxSize < maxDiff {
			maxSize = maxDiff
		}

		cidrs = append(cidrs, uInt32ToIPv4(ipStartUint32)+"/"+strconv.Itoa(maxSize))

		ipStartUint32 += uint32(math.Exp2(float64(32 - maxSize)))
	}

	return cidrs, err
}

//Convert IPv4 to uint32
func iPv4ToUint32(iPv4 string) uint32 {

	ipOctets := [4]uint64{}

	for i, v := range strings.SplitN(iPv4, ".", 4) {
		ipOctets[i], _ = strconv.ParseUint(v, 10, 32)
	}

	result := (ipOctets[0] << 24) | (ipOctets[1] << 16) | (ipOctets[2] << 8) | ipOctets[3]

	return uint32(result)
}

//Convert uint32 to IP
func uInt32ToIPv4(iPuInt32 uint32) (iP string) {
	iP = fmt.Sprintf("%d.%d.%d.%d",
		iPuInt32>>24,
		(iPuInt32&0x00FFFFFF)>>16,
		(iPuInt32&0x0000FFFF)>>8,
		iPuInt32&0x000000FF)
	return iP
}

func cidrToIPSlice(cidrString string) ([]string, error) {
	/*
		https://stackoverflow.com/questions/60540465/how-to-list-all-ips-in-a-network
	*/
	var ipSlice []string
	_, ipv4Net, err := net.ParseCIDR(cidrString)
	if err != nil {
		return nil, err
	}
	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		GlobalConfig.TargetList = append(ipSlice, ip.String())
	}
	return ipSlice, nil
}
