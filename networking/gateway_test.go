package networking

import (
	"fmt"
	pprint "github.com/NubeIO/lib-networking/print"
	"testing"
)

func Test_nets_GetGatewayIP(t *testing.T) {
	nets := &nets{}
	GetGatewayIP, err := nets.GetGatewayIP("wlp3s0")
	fmt.Println(err)
	pprint.PrintJOSN(GetGatewayIP)

	GetInterfacesNames, err := nets.GetInterfacesNames()
	fmt.Println(err)
	pprint.PrintJOSN(GetInterfacesNames)
	GetValidNetInterfaces, err := nets.GetValidNetInterfaces()
	fmt.Println(err)
	pprint.PrintJOSN(GetValidNetInterfaces)

	GetInternetIP, err := nets.GetInternetIP()
	fmt.Println(err)
	pprint.PrintJOSN(GetInternetIP)

}
