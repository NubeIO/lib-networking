package address

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/types"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type IPv4 [4]int

func New() *IPv4 {
	return &IPv4{}
}

// ToIPv4 converts a string to a IPv4.
func (inst *IPv4) ToIPv4(ip string) IPv4 {
	var newIP IPv4
	ipS := strings.Split(ip, ".")
	for i, v := range ipS {
		newIP[i], _ = strconv.Atoi(v)
	}
	return newIP
}

// ToString converts an IP from IPv4 type to string.
func (inst *IPv4) ToString() string {
	ipStringed := strconv.Itoa(inst[0])
	for i := 1; i < 4; i++ {
		strI := strconv.Itoa(inst[i])
		ipStringed += "." + strI
	}
	return ipStringed
}

// IsValid checks an IP address as valid or not.
func (inst *IPv4) IsValid() bool {
	for i, oct := range inst {
		if i == 0 || i == 3 {
			if oct < 1 || oct > 254 {
				return false
			}
		} else {
			if oct < 0 || oct > 255 {
				return false
			}
		}
	}
	return true
}

// PlusPlus increments an IPv4 value.
func (inst *IPv4) PlusPlus() *IPv4 {
	if inst[3] < 254 {
		inst[3] = inst[3] + 1
	} else {
		if inst[2] < 255 {
			inst[2] = inst[2] + 1
			inst[3] = 1
		} else {
			if inst[1] < 255 {
				inst[1] = inst[1] + 1
				inst[2] = 1
				inst[3] = 1
			} else {
				if inst[0] < 255 {
					inst[0] = inst[0] + 1
					inst[1] = 1
					inst[2] = 1
					inst[3] = 1
				}
			}
		}
	}
	return inst
}

// ParseIPSequence gets a sequence of IP addresses correspondent from an
// "init-end" entry.
func (inst *IPv4) ParseIPSequence(ipSequence string) []IPv4 {
	var arrayIps []IPv4
	series, _ := regexp.Compile("([0-9]+)")
	// For sequence ips, using '-'
	lSeries := series.FindAllStringSubmatch(ipSequence, -1)
	for i := types.ToInt(lSeries[3][0]); i <= types.ToInt(lSeries[4][0]); i++ {
		arrayIps = append(arrayIps, IPv4{
			types.ToInt(lSeries[0][0]),
			types.ToInt(lSeries[1][0]),
			types.ToInt(lSeries[2][0]),
			i})
	}
	return arrayIps
}

//GetIpList  ipsSequence := []string{"192.168.15.1-2"}  => [[192 168 15 1] [192 168 15 2]]
func (inst *IPv4) GetIpList(ips []string) []IPv4 {
	var ipList []IPv4
	for _, i := range ips {
		if strings.Contains(i, "-") {
			ipList = append(ipList, inst.ParseIPSequence(i)...)
		} else {
			ip := inst.ToIPv4(i)
			if ip.IsValid() {
				ipList = append(ipList, ip)
			}
		}
	}
	return ipList
}

// GetIPSubnet GetIPSubnet("192.168.15.1", "255.255.255.0")  =>  192.168.15.0/24 <nil>
func GetIPSubnet(ip, netmask string) (ipPrefix, prefix string, err error) {
	// Check ip
	if net.ParseIP(ip) == nil {
		return "", "", fmt.Errorf("invalid IP address %s", ip)
	}
	// Check netmask
	maskIP := net.ParseIP(netmask).To4()
	if maskIP == nil {
		return "", "", fmt.Errorf("invalid Netmask %s", netmask)
	}
	// Get prefix
	mask := net.IPMask(maskIP)
	p, _ := mask.Size()

	// Get network
	sPrefix := strconv.Itoa(p)
	_, _, err = net.ParseCIDR(ip + "/" + sPrefix)
	if err != nil {
		return "", "", err
	}
	return ip + "/" + sPrefix, sPrefix, nil
}

// SubnetString SubnetString("24") => "255.255.255.0"
func SubnetString(cidr int) (string, error) {
	cid := types.ToString(cidr)
	var maskList []string
	var netMask string
	cidrInt, err := strconv.ParseUint(cid, 10, 8)
	if err != nil {
		return "", err
	}
	for i := 0; i < 4; i++ {
		tmp := ""
		for ii := 0; ii < 8; ii++ {
			if cidrInt > 0 {
				tmp = tmp + "1"
				cidrInt--
			} else {
				tmp = tmp + "0"
			}
		}
		n, err := strconv.ParseUint(tmp, 2, 64)
		if err != nil {
			return "", err
		}
		maskList = append(maskList, strconv.FormatUint(n, 10))
	}
	netMask = strings.Join(maskList, ".")
	return netMask, nil
}
