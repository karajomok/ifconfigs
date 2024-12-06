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
