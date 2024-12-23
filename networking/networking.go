package networking

import (
	"errors"
	"net"
	"strings"
)

type Nets interface {
	GetGatewayIP(iFaceName string) (ip string, err error) //gateways

	GetNetworks() (interfaces []NetworkInterfaces, err error) //networks
	GetNetworkByIface(name string) (network NetworkInterfaces, err error)
	GetValidNetInterfaces() (interfaces []net.Interface, err error)
	GetNetworksThatHaveGateway() (interfaces []NetworkInterfaces, err error)
	CheckInterfacesName(iFaceName string) (bool, error)
	GetInterfacesNames() (InterfaceNames, error) //internet
	GetValidNetInterfacesForWeb() ([]NetInterface, error)
	CheckInternetByInterface(iface string) (connection Check, err error)
	GetInternetIP() (connection Check, err error)

	GetSubnet(iFaceName string) (subnet string, err error) //subnet
	GetSubnetCIDR(iFaceName string) (cidr int, err error)
}

// nets type
type nets struct{}

// New make an instance of nets
func New() Nets {
	return &nets{}
}

/*
Networks
*/

type NetworkInterfaces struct {
	Interface     string `json:"interface"`
	IP            string `json:"ip"`
	IPMask        string `json:"ip_and_mask"`
	NetMask       string `json:"netmask"`
	NetMaskLength int    `json:"net_mask_length"`
	Gateway       string `json:"gateway"`
	MacAddress    string `json:"mac_address"`
}

type CIDR struct {
	ip    net.IP
	ipnet *net.IPNet
}

func ParseCIDR(s string) (*CIDR, error) {
	i, n, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &CIDR{ip: i, ipnet: n}, nil
}

// GetNetworks fetches system network addresses
func (inst *nets) GetNetworks() (interfaces []NetworkInterfaces, err error) {
	var networkInterfaces NetworkInterfaces
	if ifaces, err := net.Interfaces(); err == nil {
		if err != nil {
			return nil, err
		}
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			if addrs, err := iface.Addrs(); err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
						networkInterfaces.NetMask = net.IP(v.Mask).String()
					case *net.IPAddr:
						ip = v.IP
						networkInterfaces.NetMask = ""
					}
					if ip == nil || ip.IsLoopback() {
						continue
					}
					ip = ip.To4()
					if ip == nil {
						continue
					}
					mask := strings.Split(addr.String(), "/")
					networkInterfaces.Interface = iface.Name
					networkInterfaces.IP = ip.String()
					networkInterfaces.IPMask = addr.String()
					if len(mask) >= 1 {
						networkInterfaces.NetMaskLength = ToInt(mask[1])
					}
					networkInterfaces.Gateway, err = inst.GetGatewayIP(iface.Name)
					networkInterfaces.MacAddress = iface.HardwareAddr.String()
					interfaces = append(interfaces, networkInterfaces)
				}
			}
		}
	}
	return interfaces, err
}

func (inst *nets) GetNetworksThatHaveGateway() (interfaces []NetworkInterfaces, err error) {
	var networkInterfaces NetworkInterfaces
	if ifaces, err := net.Interfaces(); err == nil {
		if err != nil {
			return nil, err
		}
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			if addrs, err := iface.Addrs(); err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
						networkInterfaces.NetMask = net.IP(v.Mask).String()
					case *net.IPAddr:
						ip = v.IP
						networkInterfaces.NetMask = ""
					}
					if ip == nil || ip.IsLoopback() {
						continue
					}
					ip = ip.To4()
					if ip == nil {
						continue
					}
					mask := strings.Split(addr.String(), "/")
					networkInterfaces.Interface = iface.Name
					networkInterfaces.IP = ip.String()
					networkInterfaces.IPMask = addr.String()
					if len(mask) >= 1 {
						networkInterfaces.NetMaskLength = ToInt(mask[1])
					}
					networkInterfaces.Gateway, err = inst.GetGatewayIP(iface.Name)
					networkInterfaces.MacAddress = iface.HardwareAddr.String()
					if networkInterfaces.Gateway != "" {
						interfaces = append(interfaces, networkInterfaces)
					}
				}
			}
		}
	}
	return interfaces, err
}

func (inst *nets) GetNetworkByIface(name string) (network NetworkInterfaces, err error) {
	all, err := inst.GetNetworks()
	if err != nil {
		return
	}

	for _, interfaces := range all {
		if name == interfaces.Interface {
			return interfaces, nil
		}
	}

	return NetworkInterfaces{}, errors.New("interface not found")
}
