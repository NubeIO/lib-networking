package scanner

import (
	"fmt"
	"testing"
)

func TestPortScanner(t *testing.T) {

	ports := []string{"22", "1414", "1883", "1660", "502", "80", "7"}

	address, err := New().ResoleAddress("", 50, "wlp3s0")
	if err != nil {
		fmt.Println("err msg", err)
		return
	}

	results := New().IPScanner(address, ports, true)
	fmt.Println(results)
	fmt.Println("-------------HOSTS------------------")

}
