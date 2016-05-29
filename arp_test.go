package whoshome

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var arpFixture = `IP address       HW type     Flags       HW address            Mask     Device
10.10.10.1      0x1         0x2         00:01:02:03:04:05     *        br0
10.10.11.1      0x1         0x2         00:01:02:03:04:06	  *        br1
10.10.12.1      0x1         0x0         00:01:02:03:04:07	  *        br2
10.10.13.1      0x1         0x2         00:01:02:03:04:08	  *        br2
`

func TestARPProvider(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	err = ioutil.WriteFile(f.Name(), []byte(arpFixture), 0666)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	p := NewARPProvider(f.Name(), map[string]string{"00:01:02:03:04:05": "user1", "00:01:02:03:04:08": "user2"})
	l, err := p.Present()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	ls := []string{"user1", "user2"}
	if !reflect.DeepEqual(ls, l) {
		t.Errorf("Expected:\n%s\nGot:\n%s", ls, l)
	}
}
