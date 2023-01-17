package networking

import (
	"fmt"
	"testing"
)

func TestNetworking(*testing.T) {

	found, err := New().GetNetworkByIface("wlp0s20f3")
	fmt.Println(found, err)

}
