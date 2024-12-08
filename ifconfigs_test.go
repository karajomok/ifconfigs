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
	ifconfig, err := GetIfconfigData()
	if err != nil {
		t.Error(err)
	}

	res, err := freebsdParser(ifconfig)
	if err != nil {
		t.Error(err)
	}

	iface := res[0]
	if iface.Name != "vlan10" {
		t.Errorf("NetworkInterface.Name must be <vlan10>, you got <%v>", iface.Name)
	} else if iface.Status != "active" {
		t.Errorf("NetworkInterface.Status must be <active>, you got <%v>", iface.Status)
	} else if iface.Type != "vlan" {
		t.Errorf("NetworkInterface.Type must be <vlan>, you got <%v>", iface.Type)
	} else if iface.VlanID != "10" {
		t.Errorf("NetworkInterface.VlanID must be <10>, you got <%v>", iface.Type)
	} else if iface.Parent != "igb1" {
		t.Errorf("NetworkInterface.Parent must be <igb1>, you got <%v>", iface.Parent)
	} else if iface.MAC != "b4:2e:99:69:71:1d" {
		t.Errorf("NetworkInterface.MAC must be <b4:2e:99:69:71:1d>, you got <%v>", iface.MAC)
	} else if iface.MTU != "1500" {
		t.Errorf("NetworkInterface.MTU must be <1500>, you got <%v>", iface.MTU)
	} else if iface.Description != "mgt_switches" {
		t.Errorf("NetworkInterface.Description must be <mgt_switches>, you got <%v>", iface.Description)
	}

	// Check IPv4
	ipv4 := iface.IpAddresses[0]
	if ipv4.IP != "198.51.100.30/27" {
		t.Errorf("ipv4.IP must be <198.51.100.30/27>, you got <%v>", ipv4.IP)
	} else if ipv4.CARP != "MASTER" {
		t.Errorf("ipv4.CARP must be <MASTER>, you got <%v>", ipv4.CARP)
	} else if ipv4.VHID != "51" {
		t.Errorf("ipv4.Vhid must be <51>, you got <%v>", ipv4.VHID)
	}

	// Check IPv6
	ipv6 := iface.IpAddresses[1]
	if ipv6.IP != "2001:db8:beef:305::1/64" {
		t.Errorf("ipv6.IP must be <2001:db8:beef:305::1/64>, you got <%v>", ipv6.IP)
	} else if ipv6.CARP != "BACKUP" {
		t.Errorf("ipv6.CARP must be <BACKUP>, you got <%v>", ipv6.CARP)
	} else if ipv6.VHID != "181" {
		t.Errorf("ipv6.VHID must be <181>, you got <%v>", ipv6.VHID)
	}
}
