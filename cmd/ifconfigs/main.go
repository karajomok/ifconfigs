package main

import (
	"fmt"

	"github.com/karajomok/ifconfigs"
)

func main() {
	ifconfig, err := ifconfigs.GetIfconfigData()
	if err != nil {
		fmt.Println(err)
		return
	}
	interfaces := ifconfigs.ParseIfconfig(ifconfig)
	fmt.Println(interfaces)
}
