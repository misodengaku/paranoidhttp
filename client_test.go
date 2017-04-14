package paranoidhttp

import (
	"fmt"
	"net"
	"testing"
)

func TestRequest(t *testing.T) {
	_, err := DefaultClient.Get("http://ipv4.google.com")
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		t.Error("The request with an ordinal url should be successful")
	}

	_, err = DefaultClient.Get("http://localhost")
	if err == nil {
		t.Errorf("The request for localhost should be fail")
	}
}

func TestIsBadHost(t *testing.T) {
	badHosts := []string{
		"localhost",
		"host has space",
		"isp-shared.test.oomurosakura.co", // => 100.64.0.1
		"local.test.oomurosakura.co",      // => 127.0.0.1
		"localv6.test.oomurosakura.co",    // => ::1
	}

	for _, h := range badHosts {
		if !isBadHost(h) {
			t.Errorf("%s should be bad", h)
		}
	}

	notBadHosts := []string{
		"www.hatena.ne.jp",
		"www.google.com",
		"xn--t8jx73hngb.jp",
	}

	for _, h := range notBadHosts {
		if isBadHost(h) {
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
		if !isBadIPAddress(net.ParseIP(ip)) {
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
		if isBadIPAddress(net.ParseIP(ip)) {
			t.Errorf("%s should not be bad", ip)
		}
	}
}
