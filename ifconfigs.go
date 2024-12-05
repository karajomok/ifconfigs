package ifconfigs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	// "os/exec"

	"strconv"
)

type NetworkInterface struct {
	Name        string
	Status      string
	Parent      string
	MAC         string
	MTU         string
	Description string
	IP          string
	Vhid        string
	CARP        string
}

func GetIfconfigData() (string, error) {
	// out, err := os.ReadFile("raw-ifconfig.txt")
	out, err := os.ReadFile("ifconfig-output.raw.txt")
	// out, err := exec.Command("ifconfig", "-v").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func ParseIfconfig(ifconfig string) []NetworkInterface {
	var interfaces []NetworkInterface

	// vlan2600: flags=8943<UP,BROADCAST,RUNNING,PROMISC,SIMPLEX,MULTICAST> metric 0 mtu 1500
	reName := regexp.MustCompile(`^(\w+): flags=`)
	// status: no carrier
	reStatus := regexp.MustCompile(`status: (active|no carrier)`)
	// ether b4:2e:99:58:54:df
	reMAC := regexp.MustCompile(`ether ((?:\w+:?)+)`)
	// description: -- TECH DEP --
	reDescription := regexp.MustCompile(`description: (\S+)`)
	// vlan2600: flags=8943<UP,BROADCAST,RUNNING,PROMISC,SIMPLEX,MULTICAST> metric 0 mtu 1500
	reMTU := regexp.MustCompile(`mtu (\d+)`)
	// inet 10.212.226.254 netmask 0xffffff00 broadcast 10.212.226.255 vhid 226
	reIPv4 := regexp.MustCompile(`inet (?P<ipv4>(?:\d+.?)+) netmask (?P<mask>\S+)(?:.*vhid (?P<vhid>\d+))?`)
	// inet6 fc00:1:1:1::1 prefixlen 64 vhid 108
	reIPv6 := regexp.MustCompile(`inet6 (?P<ipv6>(?:\w+::?)+\w+) prefixlen (?P<prefix>\d+)( vhid (?P<vhid>\d+))?`)
	// carp: INIT vhid 108 advbase 2 advskew 100
	reCARP := regexp.MustCompile(`carp: (?P<status>\S+) vhid (?P<vhid>\d+)`)

	// Need for os.exec
	lines := strings.Split(strings.ReplaceAll(ifconfig, "\n\t", "\t"), "\n")
	// lines := strings.Split(ifconfig, "<blank>")

	for _, line := range lines {
		fmt.Println(line)
		var currentInterface NetworkInterface
		var subInterfaces []map[string]string
		carpstat := map[string]string{}

		if match := reName.FindStringSubmatch(line); match != nil {
			currentInterface.Name = match[1]
		}

		if match := reStatus.FindStringSubmatch(line); match != nil {
			currentInterface.Status = match[1]
		}

		if match := reMAC.FindStringSubmatch(line); match != nil {
			currentInterface.MAC = match[1]
		}

		if match := reDescription.FindStringSubmatch(line); match != nil {
			currentInterface.Description = match[1]
		}

		if match := reIPv4.FindAllStringSubmatch(line, -1); match != nil {
			for _, ip := range match {
				ipv4 := ip[reIPv4.SubexpIndex("ipv4")]
				mask := ip[reIPv4.SubexpIndex("mask")]
				prefix, err := NetmaskHexToPrefix(mask)
				if err != nil {
					prefix = "unknown"
				}
				vhid := ip[reIPv4.SubexpIndex("vhid")]

				sub := map[string]string{}
				sub["ip"] = ipv4 + "/" + prefix
				sub["vhid"] = vhid
				subInterfaces = append(subInterfaces, sub)
			}
		}

		if match := reIPv6.FindAllStringSubmatch(line, -1); match != nil {
			for _, ip := range match {
				ipv6 := ip[reIPv6.SubexpIndex("ipv6")]
				prefix := ip[reIPv6.SubexpIndex("prefix")]
				vhid := ip[reIPv6.SubexpIndex("vhid")]

				sub := map[string]string{}
				sub["ip"] = ipv6 + "/" + prefix
				sub["vhid"] = vhid
				subInterfaces = append(subInterfaces, sub)
			}
		}

		if match := reCARP.FindAllStringSubmatch(line, -1); match != nil {
			for _, v := range match {
				carpStatus := v[reCARP.SubexpIndex("status")]
				carpVhid := v[reCARP.SubexpIndex("vhid")]
				carpstat[carpVhid] = carpStatus
			}
		}

		if match := reMTU.FindStringSubmatch(line); match != nil {
			currentInterface.MTU = match[1]
		}

		// If currentInterface not empty - append
		if (NetworkInterface{}) != currentInterface {
			if len(subInterfaces) != 0 {
				currentInterface.IP = subInterfaces[0]["ip"]
				currentInterface.Vhid = subInterfaces[0]["vhid"]
				currentInterface.CARP = carpstat[subInterfaces[0]["vhid"]]
			}
			interfaces = append(interfaces, currentInterface)
		}
		fmt.Println(currentInterface)
		fmt.Println(subInterfaces)
		for i, m := range subInterfaces {
			if i == 0 {
				continue
			}
			var subIface NetworkInterface
			subIface.IP = m["ip"]
			subIface.Vhid = m["vhid"]
			subIface.CARP = carpstat[m["vhid"]]

			if (NetworkInterface{}) != subIface {
				// If sub interface not empty, append it
				interfaces = append(interfaces, subIface)
			}
		}
	}

	// fmt.Println(interfaces)
	return interfaces
}

// NetmaskHexToPrefix converts hex netmask to prefix lenght:
// 0xffffffc0 -> 4294967232 -> 11111111111111111111111111000000 => 26
func NetmaskHexToPrefix(netmaskHex string) (string, error) {
	// Hex string to uint64
	netmaskInt, err := strconv.ParseUint(netmaskHex[2:], 16, 64)
	if err != nil {
		return "", fmt.Errorf("can't convert [%v] to decimal - %v", netmaskHex, err)
	}

	// Converts decimal to binary and count how many 1 bits in netmask
	prefix := 0
	for _, v := range strconv.FormatInt(int64(netmaskInt), 2) {
		if string(v) == "1" {
			prefix += 1
		}
	}

	return strconv.Itoa(prefix), nil
}
