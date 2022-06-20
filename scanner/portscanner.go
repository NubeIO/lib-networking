package scanner

import (
	"fmt"
	"github.com/NubeIO/lib-networking/ip"
	"github.com/NubeIO/lib-networking/networking"
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

type PortList struct {
	ServiceName string `json:"service"`
	Port        string `json:"port"`
}

type Host struct {
	IP    string     `json:"ip"`
	Ports []PortList `json:"ports"`
}

type Hosts struct {
	Hosts []Host `json:"hosts"`
}

func (inst *Scanner) ResoleAddress(ip string, count int, interfaceName string) (ipRange []string, err error) {
	ipUtil := address.New()
	networks := networking.New()
	if count > 254 {
		count = 254
	}
	if count < 0 {
		count = 1
	}
	if ip != "" {
		err := ipUtil.IsIPAddrErr(ip)
		if err != nil {
			return nil, err
		}

	} else {
		if interfaceName != "" {
			n, err := networks.GetNetworkByIface(interfaceName)
			if err != nil {
				return nil, err
			}
			ip = n.Gateway
		} else {
			gateway, err := networks.GetNetworksThatHaveGateway()
			if err != nil {
				return nil, err
			}
			for i, n := range gateway {
				if i == 0 {
					interfaceName = n.Interface
				}
			}
			n, err := networks.GetNetworkByIface(interfaceName)
			if err != nil {
				return nil, err
			}
			ip = n.Gateway
		}
	}
	ipsSequence := []string{fmt.Sprintf("%s-%d", ip, count)}
	return ipsSequence, nil

}

// IPScanner scans all IP addresses in ipList for every port in portList.
func (inst *Scanner) IPScanner(ips []string, portStr []string, printResults bool) (hostsFound *Hosts) {
	ipUtil := address.New()
	var ipList []address.IPv4
	var portList []string
	var wg sync.WaitGroup
	hostsFound = &Hosts{}
	if len(portStr) == 1 {
		portList = parsePortList(portStr[0])
	} else {
		portList = portStr
	}
	if len(ips) == 0 {
		ipList = append(ipList, address.IPv4{127, 0, 0, 1})
	} else {
		for _, i := range ips {
			if strings.Contains(i, "-") {
				ipList = append(ipList, ipUtil.ParseIPSequence(i)...)
			} else {
				ip := ipUtil.ToIPv4(i)
				if ipUtil.IsValid() {
					ipList = append(ipList, ip)
				}
			}
		}
	}
	for _, ip := range ipList {
		wg.Add(1)
		go func(ip address.IPv4) {
			defer wg.Done()
			ports := portScanner(ip, portList)
			if len(ports) > 0 {
				var pl []PortList
				for _, port := range ports {
					pl = append(pl, PortList{
						Port:        port,
						ServiceName: portShortList[port],
					})
				}
				hostsFound.Hosts = append(hostsFound.Hosts, Host{IP: ip.ToString(), Ports: pl})
				if printResults {
					presentResults(ip, ports)
				}
			}
		}(ip)
	}
	wg.Wait()
	return hostsFound
}

// ParsePortList gets a port list from its port entry in arguments.
func parsePortList(rawPorts string) []string {
	var ports []string
	individuals, _ := regexp.Compile("([0-9]+)[,]*")
	series, _ := regexp.Compile("([0-9]+)[-]([0-9]+)")

	// For individual ports, separated by ','
	lIndividuals := individuals.FindAllStringSubmatch(rawPorts, -1)

	// For sequence ports, using '-'
	lSeries := series.FindAllStringSubmatch(rawPorts, -1)

	if len(lSeries) > 0 {
		for _, s := range lSeries {
			init, _ := strconv.Atoi(s[1])
			end, _ := strconv.Atoi(s[2])
			for i := init + 1; i < end; i++ {
				ports = append(ports, strconv.Itoa(i))
			}
		}
	}
	for _, port := range lIndividuals {
		ports = append(ports, port[1])
	}
	sort.Strings(ports)
	return ports
}

// PortScanner scans IP:port pairs looking for open ports on IP addresses.
func portScanner(ip address.IPv4, portList []string) []string {
	var open []string
	for _, port := range portList {
		conn, err := net.DialTimeout("tcp",
			ip.ToString()+":"+port,
			300*time.Millisecond)
		if err == nil {
			conn.Close()
			open = append(open, port)
		}
	}
	return open
}

// PresentResults presents all results in console.
func presentResults(ip address.IPv4, ports []string) int {
	for _, port := range ports {
		log.Println("IP:", ip.ToString(), " PORT:"+port+"\t"+"Description: "+portShortList[port])
	}
	return 0
}
