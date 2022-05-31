package scanner

import (
	"fmt"
	"testing"
)

func TestPortScanner(t *testing.T) {

	ports := []string{"22", "1414", "1883", "1660", "502", "80", "7"}

	address, err := New().ResoleAddress("er432", 50, "")
	if err != nil {
		fmt.Println("err msg", err)
		return
	}

	results := New().IPScanner(address, ports, true)
	fmt.Println(results)
	fmt.Println("-------------HOSTS------------------")

}
