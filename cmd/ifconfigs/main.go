package main

import (
	"encoding/json"
	"fmt"

	"github.com/karajomok/ifconfigs"
)

func main() {
	ifconfig, err := ifconfigs.GetOutput()
	if err != nil {
		fmt.Println(err)
		return
	}
	interfaces, err := ifconfigs.Parse(ifconfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := json.MarshalIndent(interfaces, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
}
