package address

import (
	"fmt"
	"testing"
)

func TestNetIP(*testing.T) {

	fmt.Println(GetIPSubnet("192.168.15.10", "255.255.255.0"))
	fmt.Println(SubnetString(6))

}
