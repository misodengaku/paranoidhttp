package paranoidhttp

import (
	"net"
	"testing"
)

func TestParanoidGet(t *testing.T) {
	_, err := ParanoidGet("http://oomurosakura.co") // IPv4 only
	if err != nil {
		t.Error("The request with an ordinal url should be successful: ", err)
	}

	// _, err = ParanoidGet("http://ipv6.google.com")
	// if err != nil {
	// 	t.Log("Warning: If your network supports IPv6, this request should be successful", err)
	// }

	_, err = ParanoidGet("http://redirect.test.oomurosakura.co") // redirect to 127.0.0.1
	if err == nil {
		t.Error("The request for bad host should be fail")
	}

	_, err = ParanoidGet("http://localv4.test.oomurosakura.co") // redirect to 127.0.0.1
	if err == nil {
		t.Error("The request for bad host should be fail")
	}

	_, err = ParanoidGet("http://localhost")
	if err == nil {
		t.Errorf("The request for localhost should be fail")
	}
	_, err = ParanoidGet("http://127.0.0.1")
	if err == nil {
		t.Errorf("The request for localhost should be fail")
	}
}

func TestIsBadHost(t *testing.T) {
	badHosts := []string{
		"localhost",
		"host has space",
		"isp-shared.test.oomurosakura.co", // => 100.64.0.1
		"localv4.test.oomurosakura.co",    // => 127.0.0.1
		"localv6.test.oomurosakura.co",    // => ::1
	}

	for _, h := range badHosts {
		isbad, _ := IsBadHost(h)
		if !isbad {
			t.Errorf("%s should be bad", h)
		}
	}

	notBadHosts := []string{
		"www.hatena.ne.jp",
		"ipv4.google.com",
		"ipv6.google.com",
		"xn--t8jx73hngb.jp",
	}

	for _, h := range notBadHosts {
		isbad, _ := IsBadHost(h)
		if isbad {
			t.Errorf("%s should not be bad", h)
		}
	}
}

func TestIsBadIPAddress(t *testing.T) {
	badIPs := []string{
		"0.0.0.0",                    // Unspecified
		"10.0.0.0", "10.255.255.255", // Private A
		"100.64.0.0", "100.127.255.255", // ISP Shared
		"127.0.0.0", "127.255.255.255", // Loopback
		"169.254.0.0", "169.254.255.255", // Link-local
		"172.16.0.0", "172.31.255.255", // Private B
		"192.0.2.0", "192.0.2.255", // Test-Net
		"192.88.99.0", "192.88.99.255", // 6to4 relay
		"192.168.0.0", "192.168.255.255", // Private C
		"198.18.0.0", "198.19.255.255", // Benchmark test address
		"198.51.100.0", "198.51.100.255", // Test-Net 2
		"203.0.113.0", "203.0.113.255", // Test-Net 3
		"224.0.0.0", "239.255.255.255", // Multicast
		"255.255.255.255",                 // Broadcast
		"::1", "::ffff:0.0.0.0", "fe80::", // IPv6
	}

	for _, ip := range badIPs {
		isbad, ip := isBadIPAddress(net.ParseIP(ip))
		if !isbad {
			t.Errorf("%s should be bad", ip)
		}
	}

	notBadIPs := []string{
		"0.0.0.1", "8.8.8.8",
		"126.255.255.255", "128.0.0.0",
		"9.255.255.255", "11.0.0.0",
		"172.15.255.255", "172.32.0.0",
		"192.167.255.255", "192.169.0.0",
		"192.88.98.255", "192.88.100.0",
		"223.255.255.255", "240.0.0.0",
		"169.253.255.255", "169.255.0.0",
		"100.63.255.255", "100.128.0.0",
		"198.52.0.1", "255.255.255.254",
	}

	for _, ip := range notBadIPs {
		isbad, ip := isBadIPAddress(net.ParseIP(ip))
		if isbad {
			t.Errorf("%s should not be bad", ip)
		}
	}
}
