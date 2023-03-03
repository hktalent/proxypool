package storage

import (
	"context"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
)

var szCheckUrl = "https://httpbin.org/get?show_env=1"

//var szCheckUrl = "https://www.google.com/manifest?pwa=webhp"

func CheckSocks5(szIp string) bool {
	dialer, err := proxy.SOCKS5("tcp", szIp, nil, proxy.Direct)
	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}
	transport := &http.Transport{DialContext: dialContext,
		DisableKeepAlives: true}
	cl := &http.Client{Transport: transport}

	resp, err := cl.Get(szCheckUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
