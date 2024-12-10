package ifconfigs

import (
	"testing"
)

func Test_NetmaskHexToPrefix(t *testing.T) {
	exampleOne := "0xffffffc0"
	resOne, err := NetmaskHexToPrefix(exampleOne)
	if err != nil {
		t.Error(err)
	} else if resOne != "26" {
		t.Errorf("value <%v> must be converted to <26>, you got <%v>", exampleOne, resOne)
	}

	exampleTwo := "0xffffffe0"
	resTwo, err := NetmaskHexToPrefix(exampleTwo)
	if err != nil {
		t.Error(err)
	} else if resTwo != "27" {
		t.Errorf("value <%v> must be converted to <27>, you got <%v>", exampleTwo, resTwo)
	}

}

func Test_freebsdParser(t *testing.T) {
	ifconfig, err := GetOutput()
	if err != nil {
		t.Error(err)
		return
	}

	res, err := freebsdParser(ifconfig)
	if err != nil {
		t.Error(err)
		return
	}

	var iface NetworkInterface
	var ips IPinfo

	// vlan10
	iface = res[0]
	if res := iface.Name; res != "vlan10" {
		t.Errorf("value must be <vlan10>, you got <%v>", res)
	} else if res := iface.Status; res != "up" {
		t.Errorf("value must be <up>, you got <%v>", res)
	} else if res := iface.Type; res != "vlan" {
		t.Errorf("value must be <vlan>, you got <%v>", res)
	} else if res := iface.VlanID; res != "10" {
		t.Errorf("value must be <10>, you got <%v>", res)
	} else if res := iface.Parent; res != "igb1" {
		t.Errorf("value must be <igb1>, you got <%v>", res)
	} else if res := iface.MAC; res != "b4:2e:99:69:71:1d" {
		t.Errorf("value must be <b4:2e:99:69:71:1d>, you got <%v>", res)
	} else if res := iface.MTU; res != "1500" {
		t.Errorf("value must be <1500>, you got <%v>", res)
	} else if res := iface.Description; res != "mgt_switches" {
		t.Errorf("value must be <mgt_switches>, you got <%v>", res)
	}

	// Check IPv4
	if len(iface.IpAddresses) != 4 {
		t.Error("IP addresses for vlan10 not found, must be 4 addresses")
		return
	}
	ips = iface.IpAddresses[0]
	if res := ips.IP; res != "198.51.3.30/27" {
		t.Errorf("value must be <198.51.3.30/27>, you got <%v>", res)
	} else if res := ips.CARP; res != "MASTER" {
		t.Errorf("value must be <MASTER>, you got <%v>", res)
	} else if res := ips.VHID; res != "51" {
		t.Errorf("value must be <51>, you got <%v>", res)
	} else if res := ips.Advbase; res != "2" {
		t.Errorf("value must be <2>, you got <%v>", res)
	} else if res := ips.Advskew; res != "50" {
		t.Errorf("value must be <50>, you got <%v>", res)
	}
	// Check IPv4 without CARP
	ips = iface.IpAddresses[1]
	if res := ips.IP; res != "192.0.0.254/24" {
		t.Errorf("value must be <192.0.0.254/24>, you got <%v>", res)
	}

	// Check IPv6
	ips = iface.IpAddresses[2]
	if res := ips.IP; res != "2001:db8:beef:305::1/64" {
		t.Errorf("value must be <2001:db8:beef:305::1/64>, you got <%v>", res)
	} else if res := ips.CARP; res != "BACKUP" {
		t.Errorf("value must be <BACKUP>, you got <%v>", res)
	} else if res := ips.VHID; res != "181" {
		t.Errorf("value must be <181>, you got <%v>", res)
	} else if res := ips.Advbase; res != "2" {
		t.Errorf("value must be <2>, you got <%v>", res)
	} else if res := ips.Advskew; res != "155" {
		t.Errorf("value must be <155>, you got <%v>", res)
	}
	// Check IPv6 without CARP
	ips = iface.IpAddresses[3]
	if res := ips.IP; res != "2001:db8:bd:91::1/64" {
		t.Errorf("value must be <2001:db8:bd:91::1/64>, you got <%v>", res)
	}

	// ipsec99
	iface = res[1]
	if res := iface.Name; res != "ipsec99" {
		t.Errorf("value must be <ipsec99>, you got <%v>", res)
	} else if res := iface.Type; res != "ipsec" {
		t.Errorf("value must be <ipsec>, you got <%v>", res)
	} else if res := iface.MTU; res != "1500" {
		t.Errorf("value must be <1500>, you got <%v>", res)
	} else if res := iface.Description; res != "tunnel hostA -> hostB" {
		t.Errorf("value must be <tunnel hostA -> hostB>, you got <%v>", res)
	}

	// Check IPs
	if len(iface.IpAddresses) != 3 {
		t.Error("IP addresses for ipsec not found, must be 3 addresses")
		return
	}

	ips = iface.IpAddresses[0]
	if res := ips.IP; res != "tunnel 10.0.0.1 --> 172.16.0.18" {
		t.Errorf("value must be <tunnel 10.0.0.1 --> 172.16.0.18>, you got <%v>", res)
	}
	ips = iface.IpAddresses[1]
	if res := ips.IP; res != "198.51.100.0 --> 198.51.100.1" {
		t.Errorf("value must be <198.51.100.0 --> 198.51.100.1>, you got <%v>", res)
	}
	ips = iface.IpAddresses[2]
	if res := ips.IP; res != "2001:db8:ada:4::2 --> 2001:db8:ada:4::1" {
		t.Errorf("value must be <2001:db8:ada:4::2 --> 2001:db8:ada:4::1>, you got <%v>", res)
	}

}
