package paranoidhttp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)

// DefaultClient is the default Client whose setting is the same as http.DefaultClient.
var DefaultClient *http.Client

func mustParseCIDR(addr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(addr)
	if err != nil {
		log.Fatalf("%s must be parsed", addr)
	}
	return ipnet
}

var (
	netPrivateClassA = mustParseCIDR("10.0.0.0/8")
	netPrivateClassB = mustParseCIDR("172.16.0.0/12")
	netPrivateClassC = mustParseCIDR("192.168.0.0/16")
	netBenchmark     = mustParseCIDR("198.18.0.0/15")
	netTestNet       = mustParseCIDR("192.0.2.0/24")
	netTestNet2      = mustParseCIDR("198.51.100.0/24")
	netTestNet3      = mustParseCIDR("203.0.113.0/24")
	net6To4Relay     = mustParseCIDR("192.88.99.0/24")
	netISPShared     = mustParseCIDR("100.64.0.0/10")
)

func init() {
	DefaultClient, _, _ = NewClient()
}

func safeAddr(ctx context.Context, resolver *net.Resolver, hostport string) (string, error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(host)
	if ip != nil {
		if isBadIPAddress(ip) {
			return "", fmt.Errorf("bad ip is detected: %v", ip)
		}
		return net.JoinHostPort(ip.String(), port), nil
	}

	if isBadHost(host) {
		return "", fmt.Errorf("bad host is detected: %v", host)
	}

	r := resolver
	if r == nil {
		r = net.DefaultResolver
	}
	addrs, err := r.LookupIPAddr(ctx, host)
	if err != nil || len(addrs) <= 0 {
		return "", err
	}
	for _, addr := range addrs {
		if isBadIPAddress(addr.IP) {
			return "", fmt.Errorf("bad ip is detected: %v", addr.IP)
		}
	}
	return net.JoinHostPort(addrs[0].IP.String(), port), nil
}

// NewDialer returns a dialer function which only allows IPv4 connections.
//
// This is used to create a new paranoid http.Client,
// because I'm not sure about a paranoid behavior for IPv6 connections :(
func NewDialer(dialer *net.Dialer) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, hostport string) (net.Conn, error) {
		switch network {
		case "tcp", "tcp4":
			addr, err := safeAddr(ctx, dialer.Resolver, hostport)
			if err != nil {
				return nil, err
			}
			return dialer.DialContext(ctx, "tcp4", addr)
		default:
			return nil, errors.New("does not support any networks except tcp4")
		}
	}
}

// NewClient returns a new http.Client configured to be paranoid for attackers.
//
// This also returns http.Tranport and net.Dialer so that you can customize those behavior.
func NewClient() (*http.Client, *http.Transport, *net.Dialer) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         NewDialer(dialer),
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}, transport, dialer
}

var regLocalhost = regexp.MustCompile("(?i)^localhost$")
var regHasSpace = regexp.MustCompile("(?i)\\s+")

func isBadHost(host string) bool {
	if regLocalhost.MatchString(host) {
		return true
	}
	if regHasSpace.MatchString(host) {
		return true
	}
	ipList, _ := net.LookupIP(host)
	// if err != nil {
	// 	return true
	// }

	for i := 0; i < len(ipList); i++ {
		if isBadIPAddress(ipList[i]) {
			return true
		}
	}

	return false
}

func isBadIPAddress(ip net.IP) bool {
	isbadv4, _ := isBadIPv4(ip)
	isbadv6, _ := isBadIPv6(ip)
	if isbadv4 || isbadv6 {
		return true
	}
	return false
}

func isBadIPv4(ip net.IP) (bool, error) {
	if ip.To4() == nil {
		return false, fmt.Errorf("not IPv4 address")
	}

	if ip.Equal(net.IPv4bcast) || !ip.IsGlobalUnicast() ||
		netPrivateClassA.Contains(ip) || netPrivateClassB.Contains(ip) || netPrivateClassC.Contains(ip) ||
		netTestNet.Contains(ip) || netTestNet2.Contains(ip) || netTestNet3.Contains(ip) ||
		net6To4Relay.Contains(ip) || netISPShared.Contains(ip) || netBenchmark.Contains(ip) {
		return true, nil
	}

	return false, nil
}

func isBadIPv6(ip net.IP) (bool, error) {
	if ip.To16() == nil {
		return false, fmt.Errorf("not IPv6 address")
	}

	if !ip.IsGlobalUnicast() {
		return true, nil
	}

	return false, nil
}
