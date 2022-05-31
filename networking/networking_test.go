package networking

import (
	"fmt"
	"testing"
)

func TestNetworking(*testing.T) {

	found, err := New().CheckInterfacesName("www")
	fmt.Println(found, err)

}
