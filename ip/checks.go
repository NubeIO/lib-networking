package address

import (
	"errors"
	"net"
	"strings"
)

// IsIPAddrErr return an error if not a vaild ip
func (inst *IPv4) IsIPAddrErr(ip string) error {
	ipaddr := inst.IsIPAddr(ip)
	if ipaddr {
		return nil
	}
	return errors.New("invalid ip address")
}

// IsIPAddr return true if string ip contains a valid representation of an IPv4 or IPv6 address
func (inst *IPv4) IsIPAddr(ip string) bool {
	ipaddr := net.ParseIP(inst.normaliseIPAddr(ip))
	return ipaddr != nil
}

// NormaliseIPAddr return ip address without /32 (IPv4 or /128 (IPv6)
func (inst *IPv4) normaliseIPAddr(ip string) string {
	if strings.HasSuffix(ip, "/32") && strings.Contains(ip, ".") { // single host (IPv4)
		ip = strings.TrimSuffix(ip, "/32")
	} else {
		if strings.HasSuffix(ip, "/128") { // single host (IPv6)
			ip = strings.TrimSuffix(ip, "/128")
		}
	}

	return ip
}

func (inst *IPv4) IsIPSubnet(address string) (bool, error) {
	mask := net.IPMask(net.ParseIP(address).To4()) // If you have the mask as a string
	prefixSize, _ := mask.Size()
	if prefixSize == 0 {
		return false, errors.New("invalid subnet address, 255.255.255.0")
	}
	return true, nil
}
