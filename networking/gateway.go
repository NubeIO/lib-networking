package networking

import (
	"github.com/jackpal/gateway"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

//GetGatewayIP Get gateway IP address
func (inst *nets) GetGatewayIP(iFaceName string) (ip string, err error) {
	if runtime.GOOS == "linux" {
		return inst.GetGatewayIPLinux(iFaceName)
	} else {
		discoverInterface, err := gateway.DiscoverGateway()
		if err != nil {
			return "", err
		}
		return discoverInterface.String(), nil
	}
}

//GetGatewayIPLinux Get gateway IP address
func (inst *nets) GetGatewayIPLinux(iFaceName string) (ip string, err error) {
	if iFaceName == "" {
		discoverInterface, err := gateway.DiscoverGateway()
		if err != nil {
			return "", err
		}
		return discoverInterface.String(), nil
	}
	cmd := exec.Command("ip", "route", "show", "dev", iFaceName)
	d, err := cmd.Output()
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		return
	}
	fields := strings.Fields(string(d))
	if len(fields) < 3 || fields[0] != "default" {
		return
	}
	getIP := net.ParseIP(fields[2])
	if getIP == nil {
		return
	}

	return getIP.String(), nil
}
