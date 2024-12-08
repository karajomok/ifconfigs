package ifconfigs

import "testing"

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
	res, err := freebsdParser("testdata/test-ifconfig-output.txt")
	if err != nil {
		t.Error(err)
	} else if res[0].Name != "vlan10" {
		t.Errorf("NetworkInterface.Name must be <vlan10>, you got <%v>", res[0].Name)
	} else if res[0].Status != "active" {
		t.Errorf("NetworkInterface.Status must be <active>, you got <%v>", res[0].Status)
	} else if res[0].Parent != "igb1" {
		t.Errorf("NetworkInterface.Parent must be <igb1>, you got <%v>", res[0].Parent)
	} else if res[0].MAC != "b4:2e:99:69:71:1d" {
		t.Errorf("NetworkInterface.MAC must be <b4:2e:99:69:71:1d>, you got <%v>", res[0].MAC)
	} else if res[0].MTU != "1500" {
		t.Errorf("NetworkInterface.MTU must be <1500>, you got <%v>", res[0].MTU)
	} else if res[0].Description != "mgt_switches" {
		t.Errorf("NetworkInterface.Description must be <mgt_switches>, you got <%v>", res[0].Description)
	} else if res[0].IpAddresses[0].IP != "172.16.186.158/27" {
		t.Errorf("NetworkInterface.IpAddresses[0].IP must be <172.16.186.158/27>, you got <%v>", res[0].IpAddresses[0].IP)
	} else if res[0].IpAddresses[0].CARP != "MASTER" {
		t.Errorf("NetworkInterface.IpAddresses[0].CARP must be <MASTER>, you got <%v>", res[0].IpAddresses[0].CARP)
	} else if res[0].IpAddresses[0].VHID != "51" {
		t.Errorf("NetworkInterface.IpAddresses[0].Vhid must be <51>, you got <%v>", res[0].IpAddresses[0].VHID)
	}
}
