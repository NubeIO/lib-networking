package networking

import (
	"errors"
	"fmt"
	"net"
	"regexp"
)

/*
Interfaces
*/

type InterfaceNames struct {
	Names []string `json:"interface_names"`
}

func (inst *nets) GetValidNetInterfaces() (interfaces []net.Interface, err error) {
	iFaces, err := net.Interfaces()
	for i := range iFaces {
		interfaces = append(interfaces, iFaces[i])
	}
	return
}

func (inst *nets) GetInterfacesNames() (interfaces InterfaceNames, err error) {
	i, err := inst.GetValidNetInterfaces()
	if err != nil {
		return interfaces, errors.New("couldn't get interfaces")
	}
	for _, n := range i {
		interfaces.Names = append(interfaces.Names, n.Name)
	}
	return
}

func (inst *nets) CheckInterfacesName(iFaceName string) (bool, error) {
	names, err := inst.GetInterfacesNames()
	if err != nil {
		return false, err
	}
	for _, iface := range names.Names {
		matched, _ := regexp.MatchString(iface, iFaceName)
		if matched {
			return true, nil
		}
	}
	return false, fmt.Errorf("network interface not found")

}
