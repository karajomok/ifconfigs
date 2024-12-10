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
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Type        string   `json:"type"`
	Parent      string   `json:"parent"`
	VlanID      string   `json:"vlanid"`
	MAC         string   `json:"mac"`
	MTU         string   `json:"mtu"`
	Description string   `json:"description"`
	IpAddresses []IPinfo `json:"ipaddresses"`
}

type IPinfo struct {
	IP      string `json:"ip"`
	CARP    string `json:"carp"`
	VHID    string `json:"vhid"`
	Advbase string `json:"advbase"`
	Advskew string `json:"advskew"`
}

func GetOutput() (string, error) {
	// temp
	out, err := os.ReadFile("testdata/test-ifconfig-output.txt")
	// out, err := exec.Command("ifconfig").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Parse(ifconfig string) ([]NetworkInterface, error) {
	// if runtime.GOOS == "freebsd" {
	// 	res, err := freebsdParser(ifconfig)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return res, err
	// 	}
	// 	return res, nil
	// } else {
	// 	return []NetworkInterface{}, fmt.Errorf("OS: %v not supported yet", runtime.GOOS)
	// }
	res, err := freebsdParser(ifconfig)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	return res, nil
}

func freebsdParser(ifconfig string) ([]NetworkInterface, error) {
	var interfaces []NetworkInterface

	// vlan10: flags=8943<UP,BROADCAST,RUNNING,PROMISC,SIMPLEX,MULTICAST> metric 0 mtu 1500
	reName := regexp.MustCompile(`^(\w+): flags=`)
	// status: active
	reStatus := regexp.MustCompile(`status: (active|no carrier)`)
	// groups: vlan
	reType := regexp.MustCompile(`groups: (\w+)`)
	// parent interface: igb1
	reParent := regexp.MustCompile(`parent interface: (\w+)`)
	// vlan: 10
	reVlanID := regexp.MustCompile(`vlan: (\d+) vlanpcp`)
	// ether b4:2e:99:69:71:1d
	reMAC := regexp.MustCompile(`ether ((?:\w{1,2}:){5}\w{1,2})`)
	// description: mgt_switches
	reDescription := regexp.MustCompile(`description: (.*)`)
	// vlan10: flags=8943<UP,BROADCAST,RUNNING,PROMISC,SIMPLEX,MULTICAST> metric 0 mtu 1500
	reMTU := regexp.MustCompile(`mtu (\d+)`)
	// inet 198.51.100.30 netmask 0xffffffe0 broadcast 198.51.100.31 vhid 51
	reIPv4 := regexp.MustCompile(`inet (?P<ipv4>(?:\d{1,3}.){3}\d{1,3}) netmask (?P<mask>\S+) broadcast (?P<ipv4>(?:\d{1,3}.){3}\d{1,3}) ?(?:vhid (?P<vhid>\d+))?`)
	// inet6 2001:db8:beef:305::1 prefixlen 64 vhid 181
	reIPv6 := regexp.MustCompile(`inet6 (?P<ipv6>(?:\w{1,4}::?)+\w{1,4}) prefixlen (?P<prefix>\d+)( vhid (?P<vhid>\d+))?`)
	// carp: MASTER vhid 51 advbase 2 advskew 50
	reCARP := regexp.MustCompile(`carp: (?P<status>\S+) vhid (?P<vhid>\d+) advbase (?P<advbase>\d+) advskew (?P<advskew>\d+)`)

	// Add a new line separator "[<tOsMo>]" between interfaces, so it's easier to separate them
	tempReplace := strings.ReplaceAll(ifconfig, "\n\t", "\t")
	lineBreak := strings.ReplaceAll(tempReplace, "\n", "[<tOsMo>]")
	lines := strings.Split(strings.ReplaceAll(lineBreak, "\t", "\n"), "[<tOsMo>]")

	for _, line := range lines {
		var currentInterface NetworkInterface
		carpstat := map[string]map[string]string{}

		if match := reName.FindStringSubmatch(line); match != nil {
			currentInterface.Name = strings.TrimSpace(match[1])
		}

		if match := reStatus.FindStringSubmatch(line); match != nil {
			if strings.TrimSpace(match[1]) == "active" {
				currentInterface.Status = "up"
			} else {
				currentInterface.Status = "down"
			}

		}

		if match := reType.FindStringSubmatch(line); match != nil {
			currentInterface.Type = strings.TrimSpace(match[1])
		}

		if match := reParent.FindStringSubmatch(line); match != nil {
			currentInterface.Parent = strings.TrimSpace(match[1])
		}

		if match := reVlanID.FindStringSubmatch(line); match != nil {
			currentInterface.VlanID = strings.TrimSpace(match[1])
		}

		if match := reMAC.FindStringSubmatch(line); match != nil {
			currentInterface.MAC = strings.TrimSpace(match[1])
		}

		if match := reDescription.FindStringSubmatch(line); match != nil {
			currentInterface.Description = strings.TrimSpace(match[1])
		}

		if match := reMTU.FindStringSubmatch(line); match != nil {
			currentInterface.MTU = strings.TrimSpace(match[1])
		}

		if match := reCARP.FindAllStringSubmatch(line, -1); match != nil {
			for _, v := range match {
				carpVhid := v[reCARP.SubexpIndex("vhid")]
				carpstat[carpVhid] = map[string]string{
					"status":  v[reCARP.SubexpIndex("status")],
					"advbase": v[reCARP.SubexpIndex("advbase")],
					"advskew": v[reCARP.SubexpIndex("advskew")],
				}
			}
		}

		if match := reIPv4.FindAllStringSubmatch(line, -1); match != nil {
			for _, ip := range match {
				ipaddress := IPinfo{}
				ipv4 := ip[reIPv4.SubexpIndex("ipv4")]
				mask := ip[reIPv4.SubexpIndex("mask")]
				prefix, err := NetmaskHexToPrefix(mask)
				if err != nil {
					prefix = "unknown"
				}
				vhid := ip[reIPv4.SubexpIndex("vhid")]
				ipaddress.IP = ipv4 + "/" + prefix
				ipaddress.VHID = vhid
				ipaddress.CARP = carpstat[vhid]["status"]
				ipaddress.Advbase = carpstat[vhid]["advbase"]
				ipaddress.Advskew = carpstat[vhid]["advskew"]
				currentInterface.IpAddresses = append(currentInterface.IpAddresses, ipaddress)
			}
		}

		if match := reIPv6.FindAllStringSubmatch(line, -1); match != nil {
			for _, ip := range match {
				ipaddress := IPinfo{}
				ipv6 := ip[reIPv6.SubexpIndex("ipv6")]
				prefix := ip[reIPv6.SubexpIndex("prefix")]
				vhid := ip[reIPv6.SubexpIndex("vhid")]

				ipaddress.IP = ipv6 + "/" + prefix
				ipaddress.VHID = vhid
				ipaddress.CARP = carpstat[vhid]["status"]
				ipaddress.Advbase = carpstat[vhid]["advbase"]
				ipaddress.Advskew = carpstat[vhid]["advskew"]
				currentInterface.IpAddresses = append(currentInterface.IpAddresses, ipaddress)
			}
		}

		// Add an interface to a slice, only if it has a name
		if currentInterface.Name != "" {
			interfaces = append(interfaces, currentInterface)
		}
	}

	if len(interfaces) == 0 {
		return []NetworkInterface{}, fmt.Errorf("parse failed, there is nothing to parse")
	}

	return interfaces, nil
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
