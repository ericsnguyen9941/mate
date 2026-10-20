package dashboard

import (
	"net"
	"strconv"
	"strings"
	"testing"
)

// A loopback bind has no LAN URL; an unspecified bind has one exactly when
// the machine has a private IPv4 address, and it is never loopback.
func TestLANURL(t *testing.T) {
	lo, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer lo.Close()
	if got := LANURL(lo); got != "" {
		t.Errorf("loopback bind has LAN URL %q", got)
	}
	all, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Fatal(err)
	}
	defer all.Close()
	got := LANURL(all)
	if ip := firstPrivateIPv4(); ip == nil {
		if got != "" {
			t.Errorf("no private IPv4 but LAN URL %q", got)
		}
		return
	}
	port := all.Addr().(*net.TCPAddr).Port
	if !strings.HasPrefix(got, "http://") || strings.HasPrefix(got, "http://127.") || !strings.HasSuffix(got, ":"+strconv.Itoa(port)+"/") {
		t.Errorf("LAN URL = %q", got)
	}
}
